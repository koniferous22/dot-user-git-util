package main

import (
	"fmt"
	"path/filepath"

	"github.com/caarlos0/env"
	"github.com/koniferous22/dot-user-git-util/utils"
	"github.com/ogier/pflag"
)

const DotGitDirectory = ".git"

type Config struct {
	TemplateDirectory         string `env:"DOT_USER_GIT_UTIL_TEMPLATE_DIRECTORY,required"`
	TargetFolder              string `env:"DOT_USER_GIT_UTIL_TARGET_FOLDER,required"`
	FlagPerRepoMode           bool   `env:"DOT_USER_GIT_UTIL_PER_REPO_MODE"`
	FlagYesInitialPrompt      bool   `env:"DOT_USER_GIT_UTIL_YES_INITIAL_PROMPT"`
	FlagGitignoreInclude      bool   `env:"DOT_USER_GIT_UTIL_GITIGNORE_INCLUDE"`
	FlagGitignoreOmit         bool   `env:"DOT_USER_GIT_UTIL_GITIGNORE_OMIT"`
	FlagSkipWhereTargetExists bool   `env:"DOT_USER_GIT_UTIL_SKIP_WHERE_TARGET_EXISTS"`
	FlagSkipWhereGitignored   bool   `env:"DOT_USER_GIT_UTIL_SKIP_WHERE_GITIGNORED"`
	FlagForceReinitialize     bool   `env:"DOT_USER_GIT_UTIL_FORCE_REINITIALIZE"`
	FlagUnionPreselections    bool   `env:"DOT_USER_GIT_UTIL_UNION_PRESELECTIONS"`
}

type Input struct {
	GitRepositories []string
}

type AppConfig struct {
	Config Config
	Input  Input
}

func validateConfig(appConfig AppConfig) error {
	if _, err := utils.ValidateDirectoryExists(appConfig.Config.TemplateDirectory); err != nil {
		return fmt.Errorf("template directory %q doesn't exists", appConfig.Config.TemplateDirectory)
	}
	for _, gitRepository := range appConfig.Input.GitRepositories {
		gitRepositoryAbsPath, err := filepath.Abs(gitRepository)
		if err != nil {
			return fmt.Errorf("error resolving absolute path of input argument %q", gitRepository)
		}
		gitRepositoryDotGitDirectory := filepath.Join(gitRepositoryAbsPath, DotGitDirectory)
		if result, validationErr := utils.ValidateDirectoryExists(gitRepositoryDotGitDirectory); !result {
			err := fmt.Errorf("%q is not a git repository", gitRepositoryAbsPath)
			if validationErr != nil {
				err = fmt.Errorf("unable to stat %q\n%w", gitRepositoryDotGitDirectory, validationErr)
			}
			return err
		}
		if err != nil {
			return fmt.Errorf("error validating %q", gitRepository)
		}
	}
	return nil
}

func InitializeConfig() (*AppConfig, error) {
	config := Config{}
	if err := env.Parse(&config); err != nil {
		return nil, fmt.Errorf("error parsing config from env variables\n%+v", err)
	}
	defaultCliArgs := []string{"."}
	pflag.StringVarP(&config.TargetFolder, "target-folder", "t", config.TargetFolder, "Target folder in .git repositories")
	pflag.BoolVarP(&config.FlagPerRepoMode, "per-repo-mode", "p", config.FlagPerRepoMode, "Run prompts for each repository")
	pflag.BoolVarP(&config.FlagYesInitialPrompt, "yes", "y", config.FlagYesInitialPrompt, "Yes for initial prompt")
	pflag.BoolVarP(&config.FlagForceReinitialize, "force-reinit", "f", config.FlagForceReinitialize, "Force removal of all previous contents on visit + disables preselection")
	pflag.BoolVarP(&config.FlagSkipWhereTargetExists, "skip-where-target-exists", "e", config.FlagSkipWhereTargetExists, "Skip for arguments where target already exists - otherwise trigger update")
	pflag.BoolVarP(&config.FlagSkipWhereGitignored, "skip-where-gitignored", "g", config.FlagSkipWhereGitignored, "Skip for arguments where target directory is .gitignored - otherwise trigger update")
	pflag.BoolVarP(&config.FlagUnionPreselections, "union-preselections", "u", config.FlagUnionPreselections, "Pre-select if script occurs in at least one arg (repository), doesn't work with \"per-repo-mode\"")
	pflag.StringVar(&config.TemplateDirectory, "template-dir", config.TemplateDirectory, "Template directory")
	pflag.BoolVar(&config.FlagGitignoreInclude, "gitignore-yes", config.FlagGitignoreInclude, "Yes for gitignore y/n prompt")
	pflag.BoolVar(&config.FlagGitignoreOmit, "gitignore-no", config.FlagGitignoreOmit, "No for gitignore y/n prompt")
	pflag.Parse()
	gitRepositories := pflag.Args()
	if len(gitRepositories) == 0 {
		gitRepositories = defaultCliArgs
	}
	var gitRepositoryAbsolutePaths []string
	for _, gitRepository := range gitRepositories {
		gitRepositoryAbsPath, err := filepath.Abs(gitRepository)
		if err != nil {
			return nil, fmt.Errorf("error resolving absolute path of input argument %q", gitRepository)
		}
		gitRepositoryAbsolutePaths = append(gitRepositoryAbsolutePaths, gitRepositoryAbsPath)
	}
	app := &AppConfig{Config: config, Input: Input{GitRepositories: gitRepositoryAbsolutePaths}}
	return app, nil
}
