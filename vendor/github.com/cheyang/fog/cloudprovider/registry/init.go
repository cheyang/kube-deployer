package registry

import (
	"github.com/Sirupsen/logrus"
	"github.com/cheyang/fog/cloudprovider"
	aliyun_k8s "github.com/cheyang/fog/cloudprovider/providers/aliyun/k8s"
	"github.com/cheyang/fog/persist"
)

func GetProvider(provider, clusterType string, storage persist.Store) cloudprovider.CloudInterface {

	providerFunc := providerFuncMap[provider][clusterType]

	if providerFunc == nil {
		logrus.Infof("Not able to find provider %s for %s, ignore it...", provider, clusterType)
		return nil
	}

	return providerFunc(storage)
}

func RegisterProvider(cloudDriverName string, clusterType string, method func(s persist.Store) cloudprovider.CloudInterface) error {

	providerFuncMap[cloudDriverName] = map[string]func(s persist.Store) cloudprovider.CloudInterface{
		clusterType: method,
	}

	return nil
}

var providerFuncMap = map[string](map[string]func(s persist.Store) cloudprovider.CloudInterface){
	"aliyun": map[string]func(s persist.Store) cloudprovider.CloudInterface{
		"k8s": aliyun_k8s.New,
	},
}
