package create

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cheyang/fog/cluster"
	"github.com/cheyang/fog/types"
	"github.com/cheyang/fog/util"
	"github.com/cheyang/fog/util/yaml"
	"github.com/cheyang/kube-deployer/helper"
	"github.com/cheyang/kube-deployer/templates/create"
	deployer_type "github.com/cheyang/kube-deployer/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Cmd = &cobra.Command{
		Use:   "create",
		Short: "create a k8s cluster in Aliyun",
		RunE: func(cmd *cobra.Command, args []string) error {

			deployArgs, err := parseDeployArgs(cmd, args)
			if err != nil {
				return err
			}
			fmt.Printf("args: %+v\n", deployArgs)

			deployFile, paramFile, err := generateConfigFiles(deployArgs)
			if err != nil {
				return err
			}
			fmt.Printf("deployFile %s\n", deployFile)
			fmt.Printf("paramFile %s\n", paramFile)
			// return runDeploy(deployFile)
			return nil
		},
	}

	retry = false
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

func parseDeployArgs(cmd *cobra.Command, args []string) (*deployer_type.DeploymentArguments, error) {
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
	viper.BindPFlag("retry", flags.Lookup("retry"))

	if viper.GetString("key-id") == "" {
		return nil, errors.New("--key-id is mandatory")
	}
	if viper.GetString("key-secret") == "" {
		return nil, errors.New("--key-secret is mandatory")
	}

	name := viper.GetString("cluster-name")
	if name == "" {
		return nil, errors.New("--cluster-name is mandatory")
	}
	var err error
	retry, err = flags.GetBool("retry")
	if err != nil {
		return nil, err
	}

	return &deployer_type.DeploymentArguments{
		KeyID:       viper.GetString("key-id"),
		KeySecret:   viper.GetString("key-secret"),
		ImageID:     viper.GetString("image-id"),
		Region:      viper.GetString("region"),
		MasterSize:  viper.GetString("master-size"),
		NodeSize:    viper.GetString("node-size"),
		ClusterName: name,
		NumNode:     viper.GetInt("num-nodes"),
		// AnsibleVarFile: config.GetClusterInputPath(name, "create"),
	}, nil
}

func generateConfigFiles(args *deployer_type.DeploymentArguments) (deployFileName, paramFileName string, err error) {

	workingDir := filepath.Join(helper.GetRootDir(), args.ClusterName)
	_, err = os.Stat(workingDir)
	if !(os.IsNotExist(err) || retry) {
		return deployFileName, paramFileName, fmt.Errorf("working dir %s is not clean, can't work in create mode", workingDir)
	}

	inputDir := filepath.Join(workingDir, "input", "create")
	err = os.MkdirAll(inputDir, 0700)
	if err != nil {
		return deployFileName, paramFileName, err
	}

	deployFileName = filepath.Join(inputDir, "aliyun-create.yaml")
	paramFileName = filepath.Join(inputDir, "ansible-create.yaml")
	args.AnsibleVarFile = paramFileName

	deployFile, err := os.OpenFile(deployFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return deployFileName, paramFileName, err
	}
	err = helper.RenderTemplateToFile(create.AliyunTemplate, deployFile, args)
	if err != nil {
		return deployFileName, paramFileName, err
	}

	paramFile, err := os.OpenFile(paramFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return deployFileName, paramFileName, err
	}
	err = helper.RenderTemplateToFile(create.AnsibleTemplate, paramFile, args)
	if err != nil {
		return deployFileName, paramFileName, err
	}

	return deployFileName, paramFileName, nil
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

	retry := viper.GetBool("retry")
	spec.Update = retry

	util.SetStoreRoot(helper.Root)

	return cluster.Bootstrap(spec)
}
