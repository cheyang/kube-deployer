package ansible

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/stdcopy"
	docker_client "github.com/docker/engine-api/client"
	docker "github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"golang.org/x/net/context"
)

func (this *ansibleManager) startContainer() (id string, err error) {
	ctx := context.Background()
	dockerClient, err = docker_client.NewEnvClient()
	if err != nil {
		return id, err
	}

	config := this.containerCreateConfig.Config
	hostConfig := this.containerCreateConfig.HostConfig
	if hostConfig == nil {
		hostConfig = &container.HostConfig{}
	}
	newtworkConfig := this.containerCreateConfig.NetworkingConfig
	hostConfig.Binds = append(hostConfig.Binds, this.genBindsForAnsible()...)
	config.Env = append(config.Env, this.genEnvsForAnsible()...)
	resp, err := dockerClient.ContainerCreate(ctx, config, hostConfig, newtworkConfig, "")
	if err != nil {
		return id, err
	}
	for _, w := range resp.Warnings {
		logrus.Warnf("Docker create: %v\n", w)
	}

	id = resp.ID
	logrus.Infof("Container ID is %s\n", id)
	startOpt := docker.ContainerStartOptions{}
	err = dockerClient.ContainerStart(ctx, id, startOpt)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (this *ansibleManager) printContainerLogs(id string) error {
	ctx := context.Background()
	c, err := this.inspectContainer(id)
	if err != nil {
		return err
	}

	options := docker.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Follow:     true,
	}
	response, err := dockerClient.ContainerLogs(ctx, id, options)
	defer response.Close()
	if err != nil {
		return err
	}

	if c.Config.Tty {
		_, err = io.Copy(logrus.StandardLogger().Out, response)
	} else {
		_, err = stdcopy.StdCopy(logrus.StandardLogger().Out,
			logrus.StandardLogger().Out,
			response)
	}
	return err

}

func (this *ansibleManager) inspectContainer(id string) (docker.ContainerJSON, error) {
	ctx := context.Background()

	return dockerClient.ContainerInspect(ctx, id)
}

func (this *ansibleManager) genBindsForAnsible() (binds []string) {

	binds = append(binds,
		fmt.Sprintf("%s:%s:ro", filepath.Join(this.store.GetDeploymentDir(), "inventory"), ansibleHostFile),
		fmt.Sprintf("%s:%s:ro", this.store.GetMachinesDir(), ansibleSSHkeysDir),
	)

	return binds
}

func (this *ansibleManager) genEnvsForAnsible() []string {
	return []string{
		"ANSIBLE_HOST_KEY_CHECKING=False",
	}
}

func (this *ansibleManager) mappingKeyPath(keyPath string) string {
	if this.containerCreateConfig != nil {
		return strings.Replace(keyPath, this.store.GetMachinesDir(), ansibleSSHkeysDir, 1)
	}
	return keyPath
}
