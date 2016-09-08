package types

import (
	fog "github.com/cheyang/fog/types"
)

type DeployArguments struct {
	KeyID      string
	KeySecret  string
	Region     string
	MasterSize string
	// NodeSize       string
	// ClusterName    string
	// NumNode        int
	// ImageID    string
	Retry bool
	Arguments
}

type ScaleArguments struct {
	Arguments
}

type Arguments struct {
	NumNode        uint
	ImageID        string
	NodeSize       string
	ClusterName    string
	AnsibleVarFile string
}

func (this *ScaleArguments) UpdateVMSpec(vmSpec *fog.VMSpec) {
	vmSpec.Instances = this.NumNode
	if this.ImageID != "" {
		vmSpec.Properties["aliyunecs-image-id"] = this.ImageID
	}
	if this.NodeSize != "" {
		vmSpec.Properties["aliyunecs-aliyunecs-instance-type"] = this.NodeSize
	}
}
