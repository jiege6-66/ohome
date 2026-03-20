package utils

import (
	"errors"
	"fmt"
)

func AppendError(err1, err2 error) error {
	if err1 == nil {
		return err2
	}
	if err2 == nil {
		return err1
	}
	return errors.New(fmt.Sprintf("%s; %s", err1.Error(), err2.Error()))
}
