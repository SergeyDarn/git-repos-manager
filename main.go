package main

import (
	"github.com/SergeyDarn/git-repos-manager/manager"
	"github.com/charmbracelet/log"
)

func main() {
	manager.Init(log.DebugLevel)
}
