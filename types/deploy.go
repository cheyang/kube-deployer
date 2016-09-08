package types

type DeployArguments struct {
	KeyID      string
	KeySecret  string
	Region     string
	MasterSize string
	// NodeSize       string
	// ClusterName    string
	// NumNode        int
	// ImageID    string
	Retry          bool
	AnsibleVarFile string
	Arguments
}

type ScaleArguments struct {
	Arguments
}

type Arguments struct {
	NumNode     int
	ImageID     string
	NodeSize    string
	ClusterName string
}
