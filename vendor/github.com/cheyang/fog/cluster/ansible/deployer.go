package ansible

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/Sirupsen/logrus"
	"github.com/cheyang/fog/cluster/deploy"
	"github.com/cheyang/fog/persist"
	"github.com/cheyang/fog/types"
	"github.com/cheyang/fog/util"
	docker_client "github.com/docker/engine-api/client"
	docker "github.com/docker/engine-api/types"
)

const ansibleEtc = "/etc/ansible"

var dockerClient *docker_client.Client

var (
	ansibleHostFile   string
	ansibleSSHkeysDir string
)

func init() {
	ansibleHostFile = filepath.Join(ansibleEtc, "hosts")
	ansibleSSHkeysDir = filepath.Join(ansibleEtc, "machines")
}

type ansibleManager struct {
	name                  string
	hosts                 []types.Host
	containerCreateConfig *docker.ContainerCreateConfig
	run                   []string
	roleMap               map[string][]types.Host
	store                 persist.Store
}

func NewDeployer(name string) (deploy.Deployer, error) {
	storePath, err := util.GetStorePath(name)
	if err != nil {
		return nil, err
	}

	return &ansibleManager{
		name:  name,
		store: persist.NewFilestore(storePath),
	}, nil
}

func (this ansibleManager) Run() error {

	inventoryFile, err := this.createInventoryFile()
	if err != nil {
		return err
	}
	logrus.Infof("inventory file: %s\n", inventoryFile)

	if this.containerCreateConfig != nil {
		return this.dockerRun()
	} else {
		// run command
	}

	return nil
}

func (this *ansibleManager) dockerRun() error {
	err := this.initDockerClient()
	if err != nil {
		return err
	}
	err = this.pullImage()
	if err != nil {
		// return err
		if !this.imageExist() {
			return err
		} else {
			logrus.WithError(err).Warnf("the image %s is not downloaded from remote repo",
				this.containerCreateConfig.Config.Image)
		}
	}

	id, err := this.startContainer()
	if err != nil {
		return err
	}
	err = this.printContainerLogs(id)
	if err != nil {
		return err
	}
	c, err := this.inspectContainer(id)
	if err != nil {
		return err
	}
	if c.State.ExitCode != 0 {
		// logrus.Errorf("Exit failed %v, rc is %d", c.State.Error, c.State.ExitCode)
		return fmt.Errorf("Exit failed %v, rc is %d", c.State.Error, c.State.ExitCode)
	} else {
		logrus.Infoln("Exit successfully.")
	}
	return nil
}

func (this *ansibleManager) SetCommander(cmd interface{}) error {
	switch cmd.(type) {
	case []string:
		this.run = cmd.([]string)
	case *docker.ContainerCreateConfig:
		this.containerCreateConfig = cmd.(*docker.ContainerCreateConfig)
	default:
		return fmt.Errorf("Unrecongized type %v", cmd)
	}
	return nil
}

func (this *ansibleManager) SetHosts(hosts []types.Host) {

	this.hosts = hosts
	this.roleMap = make(map[string][]types.Host)

	for _, host := range hosts {

		for _, role := range host.Roles {

			if _, found := this.roleMap[role]; !found {
				this.roleMap[role] = make([]types.Host, 0)
			}

			this.roleMap[role] = append(this.roleMap[role], host)
		}

	}

	for _, value := range this.roleMap {
		sort.Sort(byHostName(value))
	}
}

// create the inventory file which is used by ansible
func (this *ansibleManager) createInventoryFile() (path string, err error) {
	storePath, err := util.GetStorePath(this.name)
	if err != nil {
		return path, err
	}
	storage := persist.NewFilestore(storePath)
	err = storage.CreateDeploymentDir()
	if err != nil {
		return path, err
	}

	deploymentDir := storage.GetDeploymentDir()
	path = filepath.Join(deploymentDir, "inventory")
	f, err := os.Create(path)
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()
	for k, hosts := range this.roleMap {
		_, err = w.WriteString(fmt.Sprintf("[%s]\n", k))
		if err != nil {
			return path, err
		}

		for _, h := range hosts {
			_, err = w.WriteString(fmt.Sprintf("%s ansible_host=%s ansible_user=%s ansible_ssh_private_key_file=%s\n",
				// h.Name,
				h.SSHHostname,
				h.SSHHostname,
				h.SSHUserName,
				this.mappingKeyPath(h.SSHKeyPath)))
			if err != nil {
				return path, err
			}
		}

		_, err = w.WriteString("\n")
		if err != nil {
			return path, err
		}
	}

	if this.name != "" {
		_, err = w.WriteString(fmt.Sprintf("[%s:children]\n", this.name))
		if err != nil {
			return path, err
		}

		for k, _ := range this.roleMap {
			_, err := w.WriteString(fmt.Sprintf("%s\n", k))
			if err != nil {
				return path, err
			}
		}
	}

	return path, err
}
