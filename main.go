package main

import (
	"fmt"
	"os"

	"github.com/koniferous22/dot-user-git-util/utils"
)

func main() {
	handleError := func(err error) {
		fmt.Printf("%s%s%s\n", utils.ColorRed, err.Error(), utils.Reset)
	}
	appConfig, err := InitializeConfig()
	if err != nil {
		handleError(fmt.Errorf("error initializing app config\n%w", err))
		os.Exit(1)
	}
	if err := validateConfig(*appConfig); err != nil {
		handleError(fmt.Errorf("error validating app config\n%w", err))
		os.Exit(1)
	}
	processingContext, err := InitializeProcessingContext(appConfig.Config)
	if err != nil {
		handleError(fmt.Errorf("error initializing processing context\n%w", err))
		os.Exit(1)
	}
	if appConfig.Config.FlagPerRepoMode {
		for _, repository := range appConfig.Input.GitRepositories {
			result, err := RunInitializationOnRepositories(appConfig.Config, *processingContext, []string{repository})
			if err != nil {
				handleError(err)
				os.Exit(1)
			}
			if result != nil && result.ShouldExit {
				os.Exit(1)
			}
		}
	} else {
		_, err := RunInitializationOnRepositories(appConfig.Config, *processingContext, appConfig.Input.GitRepositories)
		if err != nil {
			handleError(err)
			os.Exit(1)
		}
	}
}
