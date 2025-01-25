package core

import (
	"database/sql"
	"encoding/json"
	"goFlaky/core/progress"
	"os"
)

type Config struct {
	BaseDir  string          `json:"baseDir,omitempty"`
	LogDir   string          `json:"logDir,omitempty"`
	TmpDir   string          `json:"tmpDir,omitempty"`
	Worker   int8            `json:"worker,omitempty"`
	Projects []ConfigProject `json:"projects,omitempty"`
}

type ConfigProject struct {
	Identifier           string               `json:"identifier,omitempty"`
	Framework            string               `json:"framework,omitempty"`
	Language             string               `json:"language,omitempty"`
	ProjectDir           string               `json:"projectPath,omitempty"`
	TestExecutionCommand TestExecutionCommand `json:"testExecutionCommand,omitempty"`
	TestExecutionDir     string               `json:"testExecutionDir,omitempty"`
	TestDir              string               `json:"testDir,omitempty"`
	TestResultDir        string               `json:"testResultDir,omitempty"`
	Strategy             string               `json:"strategy,omitempty"`
	PreRuns              int32                `json:"preRuns,omitempty"`
	Enabled              bool                 `json:"enabled,omitempty"`
}

type TestExecutionCommand struct {
	Path    string   `json:"path,omitempty"`
	Command string   `json:"command,omitempty"`
	Args    []string `json:"args,omitempty"`
}

func LoadConfig(path string) (*Config, error) {
	file, readErr := os.ReadFile(path)
	if readErr != nil {
		return nil, readErr
	}

	var config Config
	serializeErr := json.Unmarshal(file, &config)
	if serializeErr != nil {
		return nil, serializeErr
	}

	prepareErr := os.MkdirAll(config.BaseDir, 0777)
	if prepareErr != nil {
		return nil, prepareErr
	}
	return &config, nil
}

type DependencyInjection struct {
	Config             *Config
	Progress           []progress.ProjectProgress
	ProgressChannel    chan []progress.ProjectProgress
	TerminalLogChannel chan string
	FileLogChannel     chan string
	Db                 *sql.DB
}

func CreateProgressSlice(config *Config) []progress.ProjectProgress {
	var prgs []progress.ProjectProgress
	for n := range config.Projects {
		prgs = append(prgs, progress.ProjectProgress{
			Identifier: config.Projects[n].Identifier,
			Status:     "INIT",
			Runs:       1,
			Index:      0,
		})
	}
	return prgs
}
