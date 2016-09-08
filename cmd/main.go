package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/cheyang/fog/util"
	"github.com/cheyang/kube-deployer/cmd/create"
	"github.com/cheyang/kube-deployer/cmd/remove"
	"github.com/cheyang/kube-deployer/cmd/scale"
	"github.com/cheyang/kube-deployer/helper"
	"github.com/docker/machine/libmachine/log"
	"github.com/spf13/cobra"
)

/**

 */

const (
	defaultDockerClientVersion = "1.18"
)

func main() {
	if err := mainCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

var mainCmd = &cobra.Command{
	Use:          os.Args[0],
	Short:        "control a kubernetes cluster in aliyun!",
	SilenceUsage: true,
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		flags := cmd.Flags()

		flag, err := flags.GetString("log-level")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		level, err := logrus.ParseLevel(flag)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		logrus.SetLevel(level)

		debugFlag, err := flags.GetBool("debug-docker-machine")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		if debugFlag {
			log.SetDebug(true)
			fmt.Printf("Enable docker machine debug %b\n", debugFlag)
		}

		versionSet := os.Getenv("DOCKER_API_VERSION")
		if versionSet == "" {
			if err := os.Setenv("DOCKER_API_VERSION", defaultDockerClientVersion); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		}

		if flags.Changed("docker-version") {
			dockerVersion, err := flags.GetString("docker-version")
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			} else {
				if err := os.Setenv("DOCKER_API_VERSION", dockerVersion); err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}
			}
		}

		util.SetStoreRoot(helper.Root)
	},
}

func init() {
	mainCmd.PersistentFlags().StringP("log-level", "l", "info", "Log level (options \"debug\", \"info\", \"warn\", \"error\", \"fatal\", \"panic\")")
	mainCmd.PersistentFlags().BoolP("debug-docker-machine", "D", false, "Debug the docker machine library")
	mainCmd.PersistentFlags().StringP("docker-version", "d", "1.23", "Set the docker client version")
	mainCmd.AddCommand(
		create.Cmd,
		scale.Cmd,
		remove.Cmd,
	)
}
