package k8s

import (
	"github.com/cheyang/fog/cloudprovider"
	"github.com/cheyang/fog/types"
	"github.com/docker/machine/libmachine/drivers"
)

type Aliyun struct {
	hosts []*types.Host
}

func New() cloudprovider.CloudInterface {
	return &Aliyun{}
}

func (this *Aliyun) SetConfigFromFlags(opts drivers.DriverOptions) error {
	return nil
}
func (this *Aliyun) SetHosts(hosts []*types.Host) {

}
func (this *Aliyun) Configure() error {
	return nil
}
