package manager

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/SergeyDarn/git-repos-manager/utils"

	"github.com/charmbracelet/log"
)

type cloneRepoTest struct {
	name     string
	repoName string
	repoURL  string
	want     cloneRepoWant
}

type cloneRepoWant struct {
	noDir         bool
	testGitStatus bool
}

func TestCloneRepo(t *testing.T) {
	tests := []cloneRepoTest{
		{"Not a git repo", "not_a_git_repo", "5.0.0", cloneRepoWant{
			noDir: true,
		}},
		{"Simple clone", "simple_clone", "https://github.com/SergeyDarn/test-module-js.git", cloneRepoWant{}},
		{"Git status", "git_status", "https://github.com/SergeyDarn/test-module-js.git", cloneRepoWant{
			testGitStatus: true,
		}},
	}

	utils.InitEnv()
	log.SetLevel(log.ErrorLevel)

	rootReposDir := utils.GetPathRootDir(getReposDir())
	os.RemoveAll(rootReposDir)

	t.Run("CloneRepos Group", func(t *testing.T) {
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()
				testCloneRepo(t, test)
			})
		}
	})

	os.RemoveAll(rootReposDir)
}

func testCloneRepo(t *testing.T, test cloneRepoTest) {
	repoDir := getRepoDir(test.repoName)

	var initialGitStatus string
	if test.want.testGitStatus {
		cloneRepoWithChanges(t, test, repoDir)

		initialGitStatus = utils.GitDirStatus(repoDir).String()
	}

	cloneRepo(test.repoName, test.repoURL)

	err, _ := checkRepoDirStatus(repoDir)
	if !test.want.noDir && err != nil {
		t.Fatalf("Expected repo to exist after cloning, but got error: %s", err.Error())
	}

	if test.want.noDir || !test.want.testGitStatus {
		return
	}

	gitStatus := utils.GitDirStatus(repoDir).String()

	if (initialGitStatus != "") && (gitStatus != initialGitStatus) {
		t.Fatalf("Expected repo %s that has git changes to preserve them after cloneRepo call", test.repoName)
	}
}

func cloneRepoWithChanges(t *testing.T, test cloneRepoTest, repoDir string) {
	testFile := "test.txt"

	cloneRepo(test.repoName, test.repoURL)

	err, _ := checkRepoDirStatus(repoDir)
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = os.Create(filepath.Join(repoDir, testFile))
	if err != nil {
		t.Fatal(err.Error())
	}
}
