package util

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/cheyang/fog/persist"
)

var (
	storeBase string = ".fog"
	etc       string = "/etc"
)

func GetStorePath(name string) (storePath string, err error) {
	storePath = filepath.Join(etc, storeBase, name)
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

func GetStorage(name string) (persist.Store, error) {
	storePath, err := GetStorePath(name)
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(storePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Failed to find the storage of cluster %s in %s",
			name,
			storePath)
	}
	storage := persist.NewFilestore(storePath)

	return storage, nil
}
