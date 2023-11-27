package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func EnsureDirectoryExists(directoryPath string) error {
	err := os.MkdirAll(directoryPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating directory %q:\n%w", directoryPath, err)
	}
	return nil
}

func ListTopLevelExecutablesInDirectory(directoryPath string) ([]string, error) {
	dir, err := os.Open(directoryPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}
	executables := make([]string, 0)
	for _, fileInfo := range fileInfos {
		if fileInfo.Mode().IsRegular() && fileInfo.Mode()&0111 != 0 {
			executables = append(executables, fileInfo.Name())
		}
	}
	return executables, nil
}

func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error opening source file %q:\n%w", src, err)
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("error creating destination file %q:\n%w", dst, err)
	}
	defer destinationFile.Close()
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("error copying file contents from %q to %s:\n%w", src, dst, err)
	}
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("error getting source file info %q:\n%s", src, err)
	}
	err = os.Chmod(dst, sourceInfo.Mode())
	if err != nil {
		return fmt.Errorf("error setting permissions for file %q:\n%s", dst, err)
	}
	return nil
}

func RemoveAllFilesInDirectory(directoryPath string) error {
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			err := os.Remove(path)
			if err != nil {
				return fmt.Errorf("error removing file %q:\n%w", path, err)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error walking through directory %q:\n%w", directoryPath, err)
	}
	return nil
}
