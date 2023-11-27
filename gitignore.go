package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/koniferous22/dot-user-git-util/utils"
)

func GetGitignorePattern(targetFolder string) string {
	return fmt.Sprintf("%s/", targetFolder)
}

func GetGitignorePresence(gitRepositoryPaths []string, targetFolder string) (*[]bool, error) {
	result := make([]bool, len(gitRepositoryPaths))
	for i, gitRepositoryPath := range gitRepositoryPaths {
		gitignorePath := filepath.Join(gitRepositoryPath, ".gitignore")
		directoryGitignored, err := utils.ValidateLineInFile(gitignorePath, GetGitignorePattern(targetFolder))
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
		result[i] = directoryGitignored
	}
	return &result, nil
}
func GitignoreWritePattern(gitRepositoryPath, pattern string) error {
	gitignorePath := filepath.Join(gitRepositoryPath, ".gitignore")
	file, err := os.OpenFile(gitignorePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	// Stat for file-size
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	if fileInfo.Size() > 0 {
		lastByte := make([]byte, 1)
		_, err := file.ReadAt(lastByte, fileInfo.Size()-1)
		if err != nil {
			return err
		}
		// Fix newline before appending phrase
		if lastByte[0] != '\n' {
			if _, err := file.WriteString("\n"); err != nil {
				return err
			}
		}
	}
	if _, err := file.WriteString(pattern + "\n"); err != nil {
		return err
	}

	return nil
}
