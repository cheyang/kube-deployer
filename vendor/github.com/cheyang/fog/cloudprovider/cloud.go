package cloudprovider

import (
	"github.com/cheyang/fog/types"
	"github.com/docker/machine/libmachine/drivers"
)

type CloudInterface interface {
	SetConfigFromFlags(opts drivers.DriverOptions) error

	SetHosts(hosts []types.Host)

	Configure() error
}
