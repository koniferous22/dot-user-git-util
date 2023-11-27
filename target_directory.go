package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/koniferous22/dot-user-git-util/utils"
)

func checkTargetDirectoryPresent(gitRepositoryPath string, targetFolder string) (bool, error) {
	gitRepositoryDotGitDirectory := filepath.Join(gitRepositoryPath, targetFolder)
	result, err := utils.ValidateDirectoryExists(gitRepositoryDotGitDirectory)
	if os.IsNotExist(err) {
		return false, nil
	}
	return result, err
}

func GetTargetDirectoryPresence(gitRepositoryPaths []string, targetFolder string) (*[]bool, error) {
	result := make([]bool, len(gitRepositoryPaths))
	for i, gitRepositoryPath := range gitRepositoryPaths {
		directoryPresent, err := checkTargetDirectoryPresent(gitRepositoryPath, targetFolder)
		if err != nil {
			return nil, err
		}
		result[i] = directoryPresent
	}
	return &result, nil
}

func CheckExecutableInTargetDirectories(gitRepositoryPaths []string, targetFolder string, targetName string) (*[]bool, error) {
	result := make([]bool, len(gitRepositoryPaths))
	var errors []error
	for i, gitRepositoryPath := range gitRepositoryPaths {
		targetPath := filepath.Join(gitRepositoryPath, targetFolder, targetName)
		targetFileInfo, err := os.Stat(targetPath)
		if os.IsNotExist(err) {
			result[i] = false
			continue
		}
		if targetFileInfo.Mode().IsRegular() {
			if targetFileInfo.Mode()&0111 != 0 {
				result[i] = true
			} else {
				errors = append(errors, fmt.Errorf("no exec permissions on %q", targetPath))
			}
		}
	}
	var err error
	if len(errors) > 0 {
		err = utils.AggregateErrors(errors)
	}
	return &result, err
}
