package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/SergeyDarn/git-repos-manager/utils"

	"github.com/charmbracelet/log"
)

type JSONConfig struct {
	Dependencies    map[string]string
	DevDependencies map[string]string
}

var reposDirPermissions os.FileMode = 0o777

func getReposDir() string {
	return utils.GetReposDir()
}

func ReadConfigJSON() JSONConfig {
	configFile := utils.GetConfigFile()
	return utils.ReadJSON[JSONConfig](configFile)
}

func CreateReposDir() {
	err := os.MkdirAll(getReposDir(), reposDirPermissions)
	utils.CheckError(err, "error when creating repos folder")
}

func RemoveReposDir() {
	err := os.RemoveAll(getReposDir())
	utils.CheckError(err, "error while trying to delete repos folder")
	log.Debug(utils.ColorDanger("repos folder deleted before cloning"))
}

func CloneRepos(repos map[string]string, parallelCloning bool) {
	log.Infof("cloning repos into %s", getReposDir())
	fmt.Println()

	start := time.Now()

	if !parallelCloning {
		for name, url := range repos {
			cloneRepo(name, url)
		}
	} else {
		var waitGroup sync.WaitGroup

		for name, url := range repos {
			waitGroup.Add(1)

			go func() {
				cloneRepo(name, url)
				defer waitGroup.Done()
			}()
		}

		waitGroup.Wait()
	}

	log.Infof(
		utils.ColorSuccess("cloning of %d repos took %s"),
		len(repos),
		time.Since(start),
	)
}

func cloneRepo(name string, url string) {
	if !utils.IsGitURL(url) {
		return
	}

	repoDir := getRepoDir(name)
	err, isRepoNotCloned := checkRepoDirStatus(repoDir)

	if isRepoNotCloned {
		utils.GitClone(name, url, repoDir)
		return
	}

	utils.CheckErrorf(err, "error while reading repo %s folder", name)

	gitStatus := utils.GitDirStatus(repoDir)
	if gitStatus.String() != "" {
		log.Warnf(
			utils.ColorWarning(
				"\nthere are unsaved changes for repo \"%s\" - skipping it\n"+
					"\n%s",
			),
			name,
			gitStatus.String(),
		)
		return
	}

	err = os.RemoveAll(repoDir)
	utils.CheckErrorf(err, "error while trying to delete repo folder for %s", name)

	utils.GitClone(name, url, repoDir)
}

func PullRepos() {
	repos, err := os.ReadDir(getReposDir())
	utils.CheckError(err, "error reading repos folder")

	for _, repo := range repos {
		utils.GitPull(
			getRepoDir(repo.Name()),
		)
	}
}

func ShowReposStatus() {
	repos, err := os.ReadDir(getReposDir())
	utils.CheckError(err, "error reading repos folder")
	reposWithChanges := []string{}

	for _, repo := range repos {
		gitStatus := utils.GitDirStatus(
			getRepoDir(repo.Name()),
		)

		if gitStatus.String() != "" {
			reposWithChanges = append(reposWithChanges, repo.Name())
		}
	}

	if len(reposWithChanges) == 0 {
		utils.LogSuccess("you have no unsaved repos")
		return
	}

	log.Infof(
		utils.ColorWarning("\nunsaved repos (%d):\n\n%s"),
		len(reposWithChanges),
		strings.Join(reposWithChanges, "\n"),
	)
}

func checkRepoDirStatus(repoDir string) (error, bool) {
	_, err := os.Stat(repoDir)
	return err, os.IsNotExist(err)
}

func getRepoDir(repoName string) string {
	return filepath.Join(getReposDir(), repoName)
}
