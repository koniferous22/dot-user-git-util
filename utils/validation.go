package utils

import (
	"bufio"
	"os"
)

func ValidateDirectoryExists(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}

func ValidateAllTrue(slice []bool) bool {
	for _, v := range slice {
		if !v {
			return false
		}
	}
	return true
}

func ValidateAtLeastOneTrue(slice []bool) bool {
	for _, v := range slice {
		if v {
			return true
		}
	}
	return false
}

func ValidateLineInFile(fileName string, lineToFind string) (bool, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return false, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == lineToFind {
			return true, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}
	return false, nil
}
