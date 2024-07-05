package survey

import (
	"errors"
	"fmt"
	"os"
)

//ArgumentAsFile - returm error if file is not valid as a source argument
func Files(paths ...string) error {
	if len(paths) == 0 {
		return fmt.Errorf("nothing to survey")
	}
	for _, path := range paths {
		fi, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("os.Stat (%v): %v", path, err)
		}
		if fi.IsDir() {
			return errors.New("input arg is dir")
		}
		f, err := os.OpenFile(path, os.O_RDONLY, 0777)
		if err != nil {
			return fmt.Errorf("open for read-write (%v): %v", path, err)
		}
		defer f.Close()
	}
	return nil
}
