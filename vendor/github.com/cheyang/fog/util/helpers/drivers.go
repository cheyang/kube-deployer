package helpers

import (
	"fmt"

	"github.com/cheyang/fog/types"
	"github.com/denverdino/docker-machine-driver-aliyunecs/aliyunecs"
	"github.com/docker/machine/drivers/openstack"
	"github.com/docker/machine/drivers/softlayer"
	"github.com/docker/machine/libmachine/drivers"
)

var initFuncMaps = map[string]func(hostname, storePath string) drivers.Driver{
	"aliyun":    aliyunecs.NewDriver,
	"softlayer": softlayer.NewDriver,
	// "aws":       amazonec2.NewDriver,
	"openstack": openstack.NewDriver,
}

func InitDrivers(driverName string, hostConfig types.VMSpec, storePath string) (drivers.Driver, error) {

	driverFunc, found := initFuncMaps[driverName]
	if !found {
		return nil, fmt.Errorf("Driver %s is not found.", driverName)
	}
	d := driverFunc(hostConfig.Name, storePath)

	props := hostConfig.Properties
	opts := NewConfigFlagger(props)
	d.SetConfigFromFlags(opts)

	return d, nil
}

func InitEmptyDriver(driverName, name, storePath string) (drivers.Driver, error) {
	driverFunc, found := initFuncMaps[driverName]
	if !found {
		return nil, fmt.Errorf("Driver %s is not found.", driverName)
	}
	d := driverFunc(name, storePath)

	return d, nil
}
