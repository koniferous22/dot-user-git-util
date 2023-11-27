package utils

import (
	"errors"
	"fmt"
)

func AggregateErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	var combinedMessages string
	for _, err := range errs {
		if err != nil {
			combinedMessages += fmt.Sprintf("* %s\n", err.Error())
		}
	}

	return errors.New(combinedMessages)
}
