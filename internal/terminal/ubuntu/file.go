package ubuntu

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func AddTestGoFIle() error {
	file, err := os.Create("test.go")
	if err != nil {
		return fmt.Errorf("can't create test.go file: %v", err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()
	return nil
}
