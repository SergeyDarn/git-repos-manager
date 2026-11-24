// Package manager provides the core application logic
package manager

import (
	"flag"
	"maps"

	"github.com/SergeyDarn/git-repos-manager/utils"
	"github.com/charmbracelet/log"
)

type ReposMap map[string]string

var reposMap = ReposMap{}

func (m ReposMap) ToSlice() []string {
	var repos []string

	for name := range m {
		repos = append(repos, name)
	}

	return repos
}

func GetReposMap() ReposMap {
	return reposMap
}

func Init(logLevel log.Level) {
	log.SetLevel(logLevel)

	utils.InitEnv()

	showReposStatus := flag.Bool("status", false, "show repos with unsaved git changes")
	pullRepos := flag.Bool("pull", false, "perform git pull for all repos")
	parallelInstall := flag.Bool("parallel-install", true, "clone repos in parallel (true/false)")
	safeInstall := flag.Bool("safe-install", true, "if this is set to false, repos folder will be deleted on start. Default version - each repo is checked separately, and only if repos has no unsaved changes, it's deleted and then recloned")
	flag.Parse()

	if *showReposStatus {
		ShowReposStatus()
		return
	}

	if *pullRepos {
		PullRepos()
		return
	}

	configJSON := ReadConfigJSON()
	dependencies := configJSON.Dependencies
	maps.Copy(dependencies, configJSON.DevDependencies)

	for key, value := range dependencies {
		if utils.IsGitURL(value) {
			reposMap[key] = value
		}
	}

	if len(reposMap) == 0 {
		utils.LogWarning("no git repos to install.")
		return
	}

	if !*safeInstall {
		RemoveReposDir()
	}

	CreateReposDir()
	CloneRepos(reposMap, *parallelInstall)
}
