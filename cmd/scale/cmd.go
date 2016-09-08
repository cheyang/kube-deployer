package scale

import (
	"errors"

	"github.com/cheyang/fog/persist"
	"github.com/cheyang/fog/types"
	"github.com/cheyang/fog/util"
	"github.com/cheyang/kube-deployer/helper"
	deployer_type "github.com/cheyang/kube-deployer/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	name = ""
	Cmd  = &cobra.Command{
		Use:   "scale",
		Short: "scale out/in a k8s cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				storage   persist.Store
				err       error
				spec      types.Spec
				slaveSpec types.VMSpec
			)

			if len(args) != 1 {
				return errors.New("scale out/in command takes only 1 argument")
			}
			name = args[len(args)-1]

			storage, err = util.GetStorage(name)
			if err != nil {
				return err
			}
			scaleArgs, err := parseScaleArgs(cmd, args)
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

			for _, vmSpec := range spec.VMSpecs {

			}
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

func parseScaleArgs(cmd *cobra.Command, args []string) (*deployer_type.ScaleArguments, error) {
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

	return &deployer_type.ScaleArguments{
		Arguments: deployer_type.Arguments{
			NumNode:     scaleNumNode,
			ImageID:     viper.GetString("image-id"),
			NodeSize:    viper.GetString("node-size"),
			ClusterName: name,
		},
	}, nil

}
