package host

import (
	"fmt"
	"os"

	"github.com/cheyang/fog/persist"
	"github.com/cheyang/fog/types"
	"github.com/cheyang/fog/util"
	"github.com/cheyang/fog/util/helpers"
)

var storePath string

// Create hosts in batch
func CreateInBatch(vmSpecs []types.VMSpec, hostBus chan<- types.Host) (err error) {

	// make working directory

	for _, vm := range vmSpecs {
		// go create(vm, driverName, hostBus)

		driverName := vm.CloudDriverName
		if driverName == "" {
			return fmt.Errorf("driver name is not specified.")
		}
		driver, err := helpers.InitDrivers(driverName, vm, storePath)
		if err != nil {
			return err
		}

		h := &HostHandler{
			Name:      vm.Name,
			Driver:    driver,
			VMSpec:    vm,
			createBus: hostBus,
		}

		go h.createOrGet()
	}

	return nil

}

func createStorePath(specs types.Spec) error {
	var err error
	storePath, err = util.GetStorePath(specs.Name)
	if err != nil {
		return err
	}

	// if the dir exists and not update mode
	if _, err := os.Stat(storePath); !os.IsNotExist(err) {
		if !specs.Update {
			return fmt.Errorf("working dir %s is not clean, can't work in create mode", storePath)
		}
	}

	if err := os.MkdirAll(storePath, 0700); err != nil {
		return err
	}

	return nil
}

// step 1
func BuildHostConfigs(specs types.Spec, save bool) (vmSpecs []types.VMSpec, err error) {
	if err := createStorePath(specs); err != nil {
		return vmSpecs, err
	}

	storage := persist.NewFilestore(storePath)
	if save {
		if err := storage.SaveSpec(&specs); err != nil {
			return vmSpecs, err
		}
	}

	dup := make(map[string]bool)
	for _, vmSpec := range specs.VMSpecs {

		if vmSpec.Name == "" {
			return vmSpecs, fmt.Errorf("Please specify the name")
		}

		if _, found := dup[vmSpec.Name]; found {
			return nil, fmt.Errorf("duplicate name %s in configuration file.", vmSpec.Name)
		} else {
			dup[vmSpec.Name] = true
		}

		// if the attribute 'instances' is not specified, set it as 1
		if vmSpec.Instances == 0 {
			vmSpec.Instances = 1
		}

		for i := 0; i < vmSpec.Instances; i++ {
			vm := vmSpec
			id := i + vmSpec.Start
			vm.Name = fmt.Sprintf("%s-%d", vm.Name, id)
			vm.Properties = mergeProperties(specs.Properties, vm.Properties)
			if len(vm.Roles) == 0 {
				return vmSpecs, fmt.Errorf("please specify the role of %s", vmSpec.Name)
			}
			// Set common cloud driver name if not specified
			if vm.CloudDriverName == "" {
				vm.CloudDriverName = specs.CloudDriverName
			}
			vmSpecs = append(vmSpecs, vm)
		}

	}

	return vmSpecs, nil
}

func mergeProperties(global, current map[string]interface{}) map[string]interface{} {

	merged := make(map[string]interface{})
	// logrus.Infof("global: %+v", global)
	for k, v := range global {
		merged[k] = v
	}

	// logrus.Infof("current: %+v", current)
	for k, v := range current {
		merged[k] = v
	}

	// logrus.Infof("merged: %+v", merged)

	return merged
}
