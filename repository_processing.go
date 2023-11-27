package main

import (
	"fmt"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/koniferous22/dot-user-git-util/prompts"
	"github.com/koniferous22/dot-user-git-util/utils"
)

type ProcessingContext struct {
	TemplateDirectoryContents []string
}

type RepositoryFragmentContext struct {
	InputGitRepositories           []string
	TargetDirectoryPresence        []bool
	GitignorePresence              []bool
	TemplateDirectoryPreselections []bool
}

type InitializationResult struct {
	ShouldExit bool
}

func runInitialPrompt(config Config, processingContext ProcessingContext, repositoryFragmentContext RepositoryFragmentContext) (*prompts.YesNoModel, error) {
	promptMessage := "-------------------------------------------------------\n" +
		"Do you want to initialize/update following repositories\n"
	for i, gitRepository := range repositoryFragmentContext.InputGitRepositories {
		var targetDirectoryOperation string
		if repositoryFragmentContext.TargetDirectoryPresence[i] {
			targetDirectoryOperation = fmt.Sprintf("[%sUPDATE%s]", utils.ColorYellow, utils.Reset)
		} else {
			targetDirectoryOperation = fmt.Sprintf("[%sCREATE%s]", utils.ColorBlue, utils.Reset)
		}
		promptMessage += fmt.Sprintf("* %s%q%s %s\n", utils.FontBold, gitRepository, utils.Reset, targetDirectoryOperation)
	}
	initialPromptModel := prompts.CreateYesNoModel(promptMessage, !config.FlagPerRepoMode)
	program := tea.NewProgram(initialPromptModel)
	result, err := program.Run()
	if err != nil {
		return nil, fmt.Errorf("error during initial prompt:\n%w", err)
	}
	if result, ok := result.(prompts.YesNoModel); ok {
		return &result, nil
	}
	return nil, fmt.Errorf("error retrieving prompt results:\n%w", err)
}

func runSelectionPrompt(config Config, processingContext ProcessingContext, repositoryFragmentContext RepositoryFragmentContext) (*prompts.MultiSelectModel, error) {
	promptMessage := "---------------------------------------\n" +
		"Pick entries for following repositories\n"
	for _, gitRepository := range repositoryFragmentContext.InputGitRepositories {
		promptMessage += fmt.Sprintf("* %s%q%s\n", utils.FontBold, gitRepository, utils.Reset)
	}
	selectionPromptModel := prompts.CreateMultiSelectModel(promptMessage, processingContext.TemplateDirectoryContents, repositoryFragmentContext.TemplateDirectoryPreselections)
	program := tea.NewProgram(selectionPromptModel)
	result, err := program.Run()
	if err != nil {
		return nil, fmt.Errorf("error during selection prompt:\n%w", err)
	}
	if result, ok := result.(prompts.MultiSelectModel); ok {
		return &result, nil
	}
	return nil, fmt.Errorf("error retrieving prompt results:\n%w", err)
}

func runGitignorePrompt(config Config, processingContext ProcessingContext, repositoryFragmentContext RepositoryFragmentContext) (*prompts.YesNoModel, error) {
	gitignorePromptModel := prompts.CreateYesNoModel(fmt.Sprintf("Do you want to add %q to .gitignore", config.TargetFolder), false)
	program := tea.NewProgram(gitignorePromptModel)
	result, err := program.Run()
	if err != nil {
		return nil, fmt.Errorf("error during gitignore prompt:\n%w", err)
	}
	if result, ok := result.(prompts.YesNoModel); ok {
		return &result, nil
	}
	return nil, fmt.Errorf("error retrieving prompt results:\n%w", err)
}

func processInitialization(templateDirectory string, templateDirectoryContents []string, gitRepositories []string, targetDirectory string, templateSelections []bool) error {
	for _, gitRepository := range gitRepositories {
		targetPath := filepath.Join(gitRepository, targetDirectory)

		err := utils.EnsureDirectoryExists(targetPath)
		if err != nil {
			return err
		}

		for j, isSelected := range templateSelections {
			if isSelected {
				templateFile := templateDirectoryContents[j]
				sourceFile := filepath.Join(templateDirectory, templateFile)
				destinationFile := filepath.Join(targetPath, filepath.Base(templateFile))
				err := utils.CopyFile(sourceFile, destinationFile)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func processGitignore(gitRepositories []string, targetDirectory string, gitignoreReferencesFound []bool) error {
	gitignorePattern := GetGitignorePattern(targetDirectory)
	for i, gitRepository := range gitRepositories {
		if gitignoreReferencesFound[i] {
			continue
		}

		if err := GitignoreWritePattern(gitRepository, gitignorePattern); err != nil {
			return fmt.Errorf("error appending or creating .gitignore:\n%w", err)
		}
	}
	return nil
}

// Pre-selection algorithm
// 1. If "Force reinitialize" Flag is set, all contents will be purged an reinitialized, therefore preselection is empty
// 2. Iterate template directory executables
// 3. Determine pre-selection set (preselect on found targets)
// 4. Select template directory executables that are present in ALL directories from pre-selection set
func resolveTemplatePreselections(config Config, processingContext ProcessingContext, gitRepositories []string) (*[]bool, error) {
	result := make([]bool, len(processingContext.TemplateDirectoryContents))
	if config.FlagForceReinitialize {
		return &result, nil
	}
	if config.FlagUnionPreselections {
		for i, templateDirectoryEntry := range processingContext.TemplateDirectoryContents {
			entryOccurenceInGitRepositories, err := CheckExecutableInTargetDirectories(gitRepositories, config.TargetFolder, templateDirectoryEntry)
			if err != nil {
				return nil, err
			}
			result[i] = utils.ValidateAtLeastOneTrue(*entryOccurenceInGitRepositories)
		}
	} else {
		for i, templateDirectoryEntry := range processingContext.TemplateDirectoryContents {
			entryOccurenceInGitRepositories, err := CheckExecutableInTargetDirectories(gitRepositories, config.TargetFolder, templateDirectoryEntry)
			if err != nil {
				return nil, err
			}
			result[i] = utils.ValidateAllTrue(*entryOccurenceInGitRepositories)
		}
	}
	return &result, nil
}

func initializeRepositorySequenceContext(config Config, processingContext ProcessingContext, gitRepositories []string) (*RepositoryFragmentContext, error) {
	targetDirectoryPresence, err := GetTargetDirectoryPresence(gitRepositories, config.TargetFolder)
	if err != nil {
		return nil, fmt.Errorf("error resolving target directory presence\n%w", err)
	}
	gitignorePresence, err := GetGitignorePresence(gitRepositories, config.TargetFolder)
	if err != nil {
		return nil, fmt.Errorf("error resolving target directory presence in .gitignore\n%w", err)
	}
	templateDirectoryPreselections, err := resolveTemplatePreselections(config, processingContext, gitRepositories)
	if err != nil {
		return nil, fmt.Errorf("error resolving template preselections\n%w", err)
	}
	return &RepositoryFragmentContext{
		InputGitRepositories:           gitRepositories,
		TargetDirectoryPresence:        *targetDirectoryPresence,
		GitignorePresence:              *gitignorePresence,
		TemplateDirectoryPreselections: *templateDirectoryPreselections,
	}, nil
}

func InitializeProcessingContext(config Config) (*ProcessingContext, error) {

	templateDirectoryContents, err := utils.ListTopLevelExecutablesInDirectory(config.TemplateDirectory)
	if err != nil {
		return nil, fmt.Errorf("error listing executables in template directory %q - %w", config.TemplateDirectory, err)
	}
	return &ProcessingContext{
		TemplateDirectoryContents: templateDirectoryContents,
	}, nil
}

func RunInitializationOnRepositories(config Config, processingContext ProcessingContext, gitRepositories []string) (*InitializationResult, error) {

	handlePromptError := func(err error) error {
		return fmt.Errorf("encountered prompt error:\n%w", err)
	}
	repositoryFragmentContext, err := initializeRepositorySequenceContext(config, processingContext, gitRepositories)
	if err != nil {
		return nil, fmt.Errorf("error initializing repository fragment context\n%w", err)
	}
	if config.FlagSkipWhereTargetExists && utils.ValidateAllTrue(repositoryFragmentContext.TargetDirectoryPresence) {
		if config.FlagPerRepoMode {
			fmt.Printf(
				"%sSkipping repository %q - target directory found%s\n",
				utils.ColorGreen,
				// NOTE - size of "repositories" expected to be 1 in per-repo mode
				repositoryFragmentContext.InputGitRepositories[0],
				utils.Reset,
			)
		} else {
			fmt.Printf("%sSKIPPING ALL - TARGET DIRECTORIES PRESENT%s\n", utils.ColorGreen, utils.Reset)
		}
		return nil, nil
	}
	targetDirectoryPresentInAllGitignores := utils.ValidateAllTrue(repositoryFragmentContext.GitignorePresence)
	if config.FlagSkipWhereGitignored && targetDirectoryPresentInAllGitignores {
		if config.FlagPerRepoMode {
			fmt.Printf(
				"%sSkipping repository %q - already found in .gitignore%s\n",
				utils.ColorGreen,
				// NOTE - size of "repositories" expected to be 1 in per-repo mode
				repositoryFragmentContext.InputGitRepositories[0],
				utils.Reset,
			)
		} else {
			fmt.Printf("%sSKIPPING ALL - TARGET DIRECTORIES .GITIGNORED%s\n", utils.ColorGreen, utils.Reset)
		}
		return nil, nil
	}
	// 1. Initial Prompt
	if !config.FlagYesInitialPrompt {
		initialPromptOutput, err := runInitialPrompt(config, processingContext, *repositoryFragmentContext)
		if err != nil {
			return nil, handlePromptError(err)
		}
		if initialPromptOutput.ShouldExit {
			return &InitializationResult{ShouldExit: true}, nil
		}
		if !initialPromptOutput.Result {
			return nil, nil
		}
	}

	// 2. Selection Prompt
	selectionPromptOutput, err := runSelectionPrompt(config, processingContext, *repositoryFragmentContext)
	if err != nil {
		return nil, handlePromptError(err)
	}
	if selectionPromptOutput.ShouldExit {
		return &InitializationResult{ShouldExit: true}, nil
	}

	// 3. .gitignore Prompt
	shouldInitializeGitignore := false
	if config.FlagGitignoreInclude {
		shouldInitializeGitignore = true
	} else if !config.FlagGitignoreOmit && !targetDirectoryPresentInAllGitignores {
		gitignorePromptOutput, err := runGitignorePrompt(config, processingContext, *repositoryFragmentContext)
		if err != nil {
			return nil, handlePromptError(err)
		}
		if gitignorePromptOutput.ShouldExit {
			return &InitializationResult{ShouldExit: true}, nil
		}
		shouldInitializeGitignore = gitignorePromptOutput.Result
	}

	// 4. Process
	err = processInitialization(
		config.TemplateDirectory,
		processingContext.TemplateDirectoryContents,
		repositoryFragmentContext.InputGitRepositories,
		config.TargetFolder,
		selectionPromptOutput.Selected,
	)
	if err != nil {
		return nil, fmt.Errorf("repository initialization error:\n%w", err)
	}
	if shouldInitializeGitignore {
		processGitignore(repositoryFragmentContext.InputGitRepositories, config.TargetFolder, repositoryFragmentContext.GitignorePresence)
		if err != nil {
			return nil, fmt.Errorf("gitignore initialization error:\n%w", err)
		}
	}
	return nil, nil
}
