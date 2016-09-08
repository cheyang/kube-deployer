package scale

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cheyang/fog/cluster"
	"github.com/cheyang/fog/persist"
	fog "github.com/cheyang/fog/types"
	"github.com/cheyang/fog/util"
	"github.com/cheyang/kube-deployer/helper"
	"github.com/cheyang/kube-deployer/templates/scale"
	"github.com/cheyang/kube-deployer/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const slaveName = "kube-slave"

var (
	name string
	Cmd  = &cobra.Command{
		Use:   "scale",
		Short: "scale out/in a k8s cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				storage   persist.Store
				err       error
				spec      fog.Spec
				slaveSpec *fog.VMSpec
				scaleArgs *types.ScaleArguments
			)

			if len(args) != 1 {
				return errors.New("scale out/in command takes only 1 argument")
			}
			name = args[len(args)-1]

			scaleArgs, err = parseScaleArgs(cmd, args)
			if err != nil {
				return err
			}
			storage, err = util.GetStorage(name)
			if err != nil {
				return err
			}
			hostList, err := helper.GetCurrentHosts(storage)
			if err != nil {
				return err
			}
			spec, err = storage.LoadSpec()
			if err != nil {
				return err
			}

			// build vmspec for scaling out
			for _, vmSpec := range spec.VMSpecs {
				if vmSpec.Name == "" {
					slaveSpec = &vmSpec
					for k, v := range spec.Properties {
						if _, found := slaveSpec.Properties[k]; !found {
							slaveSpec.Properties[k] = v
						}
					}
					break
				}
			}

			// scale out
			if scaleArgs.NumNode > 0 {
				deployFile, paramFile, err := generateConfigFiles(scaleArgs)
				if err != nil {
					return err
				}
				fmt.Printf("deployFile %s\n", deployFile)
				fmt.Printf("paramFile %s\n", paramFile)

				newSpec, err := fog.LoadSpec(deployFile)
				if err != nil {
					return err
				}
				newSpec.VMSpecs[0] = *slaveSpec
				roleMap := map[string]bool{
					"masters": true,
					"etcd":    true,
				}

				return cluster.ExpandCluster(storage, newSpec, roleMap)
				// scale in
			} else if scaleArgs.NumNode < 0 {
				return errors.New("Not implmented yet.")
			}
			return nil
		},
	}
)

func init() {
	flags := Cmd.Flags()
	flags.StringP("scale-num-nodes", "", "", "The number of k8s node to scale out or in, can be +1 or -1")
	flags.StringP("key-id", "", "", "ECS Access Key id")
	flags.StringP("key-secret", "", "", "ECS Access Key secret")
	flags.StringP("image-id", "", "entos7u2_64_40G_cloudinit_20160520.raw", "ECS image id to create k8s cluster")
	flags.StringP("node-size", "", "ecs.n1.small", "The size of node virtual machine")
}

func parseScaleArgs(cmd *cobra.Command, args []string) (*types.ScaleArguments, error) {
	viper.BindEnv("key-id", "ALIYUNECS_KEY_ID")
	viper.BindEnv("key-secret", "ALIYUNECS_KEY_SECRET")
	viper.BindEnv("image-id", "ALIYUNECS_IMAGE_ID")
	viper.BindEnv("node-size", "ALIYUNECS_NODE_SIZE")

	flags := cmd.Flags()
	viper.BindPFlag("key-id", flags.Lookup("key-id"))
	viper.BindPFlag("key-secret", flags.Lookup("key-secret"))
	viper.BindPFlag("image-id", flags.Lookup("image-id"))
	viper.BindPFlag("node-size", flags.Lookup("node-size"))
	viper.BindPFlag("scale-num-nodes", flags.Lookup("scale-num-nodes"))

	if viper.GetString("key-id") == "" {
		return nil, errors.New("--key-id is mandatory")
	}
	if viper.GetString("key-secret") == "" {
		return nil, errors.New("--key-secret is mandatory")
	}

	scaleNumNode, err := helper.ParseScaleFlag(viper.GetString("scale-num-nodes"))
	if err != nil {
		return nil, err
	}

	return &types.ScaleArguments{
		Arguments: types.Arguments{
			NumNode:     scaleNumNode,
			ImageID:     viper.GetString("image-id"),
			NodeSize:    viper.GetString("node-size"),
			ClusterName: name,
		},
	}, nil

}

func generateConfigFiles(args *types.ScaleArguments) (deployFileName, paramFileName string, err error) {
	//check if working dir as expected
	workingDir := filepath.Join(helper.GetRootDir(), args.ClusterName)
	_, err = os.Stat(workingDir)
	if os.IsNotExist(err) {
		return deployFileName, paramFileName, fmt.Errorf("working dir %s doesn't exist, can't scale out or in", workingDir)
	}

	// create input dir for input file generation
	t := time.Now()
	timestamp := fmt.Sprint(t.Format("20060102150405"))
	inputDir := filepath.Join(workingDir,
		"input",
		"scale_"+timestamp)
	err = os.MkdirAll(inputDir, 0700)
	if err != nil {
		return deployFileName, paramFileName, err
	}

	deployFileName = filepath.Join(inputDir, "aliyun-scale.yaml")
	paramFileName = filepath.Join(inputDir, "ansible-scale.yaml")
	args.AnsibleVarFile = paramFileName

	deployFile, err := os.OpenFile(deployFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return deployFileName, paramFileName, err
	}
	err = helper.RenderTemplateToFile(scale.AliyunTemplate, deployFile, args)
	if err != nil {
		return deployFileName, paramFileName, err
	}

	paramFile, err := os.OpenFile(paramFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return deployFileName, paramFileName, err
	}
	err = helper.RenderTemplateToFile(scale.AnsibleTemplate, paramFile, args)
	if err != nil {
		return deployFileName, paramFileName, err
	}

	return deployFileName, paramFileName, nil
}
