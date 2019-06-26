package exporter

import (
	"github.com/miniclip/gonsul/internal/configuration"
	"github.com/miniclip/gonsul/internal/util"
)

// IExporter ...
type IExporter interface {
	Start() map[string]string
}

// exporter ...
type exporter struct {
	config configuration.IConfig
	logger util.ILogger
}

// NewExporter ...
func NewExporter(config configuration.IConfig, logger util.ILogger) IExporter {
	return &exporter{config: config, logger: logger}
}

// Start ...
func (e *exporter) Start() map[string]string {
	// Instantiate our local data map
	var localData = map[string]string{}

	// Should we clone the repo, or is it already done via 3rd party
	if e.config.IsCloning() {
		e.logger.PrintInfo("REPO: GIT cloning from: " + e.config.GetRepoURL())
		e.downloadRepo()
	} else {
		e.logger.PrintInfo("REPO: Skipping GIT clone, using local path: " + e.config.GetRepoRootDir())
	}

	// Set the path where Gonsul should start traversing files to add to Consul
	repoDir := e.config.GetRepoRootDir() + "/" + e.config.GetRepoBasePath()
	// Traverse our repo directory, filling up the data.EntryCollection structure
	e.parseDir(repoDir, localData)

	// Return our final data.EntryCollection structure
	return localData
}