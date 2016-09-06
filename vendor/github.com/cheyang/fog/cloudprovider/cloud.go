package cloudprovider

import "github.com/cheyang/fog/types"

type CloudInterface interface {
	// SetConfigFromFlags(opts drivers.DriverOptions) error

	SetHosts(hosts []types.Host)

	Configure() error
}
