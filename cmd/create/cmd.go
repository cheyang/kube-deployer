package create

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cheyang/fog/cluster"
	"github.com/cheyang/fog/types"
	"github.com/cheyang/fog/util/yaml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type DeploymentArguments struct {
	KeyID       string
	KeySecret   string
	ImageID     string
	Region      string
	MasterSize  string
	NodeSize    string
	ClusterName string
	NumNode     int
	Retry       bool
}

var (
	Cmd = &cobra.Command{
		Use:   "create",
		Short: "Create a k8s cluster in Aliyun",
		RunE: func(cmd *cobra.Command, args []string) error {

			deployArgs := parseDeployArgs(cmd, args)
			fmt.Printf("args: %+v", deployArgs)
		},
	}
)

func init() {
	flags := Cmd.Flags()
	flags.StringP("key-id", "", "", "ECS Access Key id")
	flags.StringP("key-secret", "", "", "ECS Access Key secret")
	flags.StringP("image-id", "", "entos7u2_64_40G_cloudinit_20160520.raw", "ECS image id to create k8s cluster")
	flags.StringP("region", "", "cn-hongkong", "The region to create cluster")
	flags.StringP("master-size", "", "ecs.n1.small", "The size of master virtual machine")
	flags.StringP("node-size", "", "ecs.n1.small", "The size of node virtual machine")
	flags.StringP("cluster-name", "", "mycluster", "The k8s cluster name")
	flags.IntP("num-nodes", "", 2, "The number of k8s node")
	flags.BoolP("retry", "r", false, "retry to create k8s cluster.")
}

func parseDeployArgs(cmd *cobra.Command, args []string) (*DeploymentArguments, error) {
	if !cmd.Flags().Changed("key-id") {
		return nil, errors.New("--key-id are mandatory")
	}
	keyId, err := cmd.Flags().GetString("key-id")
	if err != nil {
		return nil, err
	}

	if !cmd.Flags().Changed("key-secret") {
		return nil, errors.New("--key-secret are mandatory")
	}
	keySecret, err := cmd.Flags().GetString("key-secret")
	if err != nil {
		return nil, err
	}

	viper.BindEnv("key-id", "ALIYUNECS_KEY_ID")
	viper.BindEnv("key-secret", "ALIYUNECS_KEY_SECRET")
	viper.BindEnv("image-id", "ALIYUNECS_IMAGE_ID")
	viper.BindEnv("region", "ALIYUNECS_REGION")
	viper.BindEnv("master-size", "ALIYUNECS_MASTER_SIZE")
	viper.BindEnv("node-size", "ALIYUNECS_NODE_SIZE")
	viper.BindEnv("cluster-name", "ALIYUNECS_CLUSTER_NAME")

	flags := cmd.Flags()
	viper.BindPFlag("key-id", flags.Lookup("key-id"))
	viper.BindPFlag("key-secret", flags.Lookup("key-secret"))
	viper.BindPFlag("image-id", flags.Lookup("image-id"))
	viper.BindPFlag("region", flags.Lookup("region"))
	viper.BindPFlag("master-size", flags.Lookup("master-size"))
	viper.BindPFlag("node-size", flags.Lookup("node-size"))
	viper.BindPFlag("cluster-name", flags.Lookup("cluster-name"))
	viper.BindPFlag("num-nodes", flags.Lookup("num-nodes"))

	return &DeploymentArguments{
		KeyID:       viper.GetString("key-id"),
		KeySecret:   viper.GetString("key-secret"),
		ImageID:     viper.GetString("image-id"),
		Region:      viper.GetString("region"),
		MasterSize:  viper.GetString("master-size"),
		NodeSize:    viper.GetString("node-size"),
		ClusterName: viper.GetString("cluster-name"),
		NumNode:     viper.GetInt("num-nodes"),
	}, nil
}

func runDeploy(configFile string) error {
	//read and parse the config file
	spec := types.Spec{}
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return err
	}
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	decoder := yaml.NewYAMLToJSONDecoder(bytes.NewReader(data))
	err = decoder.Decode(&spec)
	if err != nil {
		return err
	}

	retry, err := cmd.Flags().GetBool("retry")
	if err != nil {
		return err
	}
	spec.Update = retry

	return cluster.Bootstrap(spec)
}
