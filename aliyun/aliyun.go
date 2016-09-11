package aliyun

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/cheyang/fog/cloudprovider"
	"github.com/cheyang/fog/persist"
	fog "github.com/cheyang/fog/types"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/docker-machine-driver-aliyunecs/aliyunecs"
)

const (
	ipRange = "0.0.0.0/0"
)

var (
	kubeletPort = 10250
	flannelPort = 8472
	client      *ecs.Client
	once        sync.Once
)

type aliyunProvider struct {
	hosts   []*fog.Host
	storage persist.Store
}

func New(s persist.Store) cloudprovider.CloudInterface {
	return &aliyunProvider{storage: s}
}

func (this *aliyunProvider) SetHosts(hosts []*fog.Host) {
	this.hosts = hosts

}

func (this *aliyunProvider) Configure() error {
	for _, host := range this.hosts {
		err := this.configureHost(host)
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *aliyunProvider) configureHost(host *fog.Host) (err error) {
	if driver, ok := host.Driver.(*aliyunecs.Driver); ok {
		host.PrivateIPAddress = driver.PrivateIPAddress
		if !driver.PrivateIPOnly {
			host.PublicIPAddress = driver.IPAddress
		}
		// set private ip for SSHHostname
		host.SSHHostname = host.PrivateIPAddress
		err = this.configureSecurityGroup(driver)
		if err != nil {
			return err
		}

		err = this.storage.Save(host)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("failed to parse %v", host.Driver)
	}

	return nil
}

func (this *aliyunProvider) configureSecurityGroup(d *aliyunecs.Driver) error {
	var securityGroup *ecs.DescribeSecurityGroupAttributeResponse

	args := ecs.DescribeSecurityGroupAttributeArgs{
		SecurityGroupId: d.SecurityGroupId,
		RegionId:        d.Region,
	}
	securityGroup, err := getClient(d).DescribeSecurityGroupAttribute(&args)

	if err != nil {
		return err
	}
	if securityGroup == nil {
		return fmt.Errorf("Failed to configure the security group %s", d.SecurityGroupId)
	}

	perms := this.configureSecurityGroupPermissions(securityGroup)
	for _, permission := range perms {
		args := permission.createAuthorizeSecurityGroupArgs(d.Region, d.SecurityGroupId)
		if err := getClient(d).AuthorizeSecurityGroup(args); err != nil {
			return err
		}
	}

	return nil
}

type IpPermission struct {
	IpProtocol ecs.IpProtocol
	FromPort   int
	ToPort     int
	IpRange    string
}

func (p *IpPermission) createAuthorizeSecurityGroupArgs(regionId common.Region, securityGroupId string) *ecs.AuthorizeSecurityGroupArgs {
	args := ecs.AuthorizeSecurityGroupArgs{
		RegionId:        regionId,
		SecurityGroupId: securityGroupId,
		IpProtocol:      p.IpProtocol,
		SourceCidrIp:    p.IpRange,
		PortRange:       fmt.Sprintf("%d/%d", p.FromPort, p.ToPort),
	}
	return &args
}

func (d *aliyunProvider) configureSecurityGroupPermissions(group *ecs.DescribeSecurityGroupAttributeResponse) []IpPermission {
	hasSSHPort := false
	hasKuberletPort := false
	hasFlannelPort := false
	hasAllIncomingPort := false
	for _, p := range group.Permissions.Permission {
		portRange := strings.Split(p.PortRange, "/")

		fromPort, _ := strconv.Atoi(portRange[0])
		switch fromPort {
		case -1:
			if portRange[1] == "-1" && p.IpProtocol == "ALL" && p.Policy == "Accept" {
				hasAllIncomingPort = true
			}
		case 22:
			hasSSHPort = true
		case kubeletPort:
			hasKuberletPort = true
		case flannelPort:
			hasFlannelPort = true
		}
	}

	perms := []IpPermission{}

	if !hasSSHPort {
		perms = append(perms, IpPermission{
			IpProtocol: ecs.IpProtocolTCP,
			FromPort:   22,
			ToPort:     22,
			IpRange:    ipRange,
		})
	}

	if !hasKuberletPort {
		perms = append(perms, IpPermission{
			IpProtocol: ecs.IpProtocolTCP,
			FromPort:   kubeletPort,
			ToPort:     kubeletPort,
			IpRange:    ipRange,
		})
	}

	if !hasFlannelPort {
		perms = append(perms, IpPermission{
			IpProtocol: ecs.IpProtocolUDP,
			FromPort:   flannelPort,
			ToPort:     flannelPort,
			IpRange:    ipRange,
		})
	}

	if !hasAllIncomingPort {
		perms = append(perms, IpPermission{
			IpProtocol: ecs.IpProtocolAll,
			FromPort:   -1,
			ToPort:     -1,
			IpRange:    ipRange,
		})
	}

	return perms
}

func getClient(d *aliyunecs.Driver) *ecs.Client {
	once.Do(func() {
		client := ecs.NewClient(d.AccessKey, d.SecretKey)
		if d.APIEndpoint != "" {
			client.SetEndpoint(d.APIEndpoint)
		}
	})

	return client
}
