package util

import (
	"fmt"
	"path/filepath"
	"regexp"
)

var storeBase string = ".fog"

func GetStorePath(name string) (storePath string, err error) {
	var etc string = "/etc"
	// if etc, err = os.Getwd(); err != nil {
	// 	return storePath, err
	// } else {
	storePath = filepath.Join(etc, storeBase, name)
	// }

	return storePath, err
}

func SetStoreRoot(root string) error {
	re := regexp.MustCompile("^(\\w+)$")
	if re.MatchString(root) {
		storeBase = fmt.Sprintf(".%s", root)
	} else {
		return fmt.Errorf("store base %s is illegal.", root)
	}

	return nil
}
