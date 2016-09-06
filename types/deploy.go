package types

type DeploymentArguments struct {
	KeyID          string
	KeySecret      string
	ImageID        string
	Region         string
	MasterSize     string
	NodeSize       string
	ClusterName    string
	NumNode        int
	Retry          bool
	AnsibleVarFile string
}
