package create

var AliyunTemplate string = `
--- 
ClusterName: [[[ .ClusterName ]]]
ClusterType: k8s
DockerRun: 
  Config: 
    Cmd: 
      - ./deploy-cluster.sh
    Env: 
      - ANSIBLE_HOST_KEY_CHECKING=False
    Image: "cheyang/k8s-ansible:on-build"
    Tty: true
  HostConfig: 
    Binds: 
      - "[[[ .AnsibleVarFile ]]]:/etc/ansible/group_vars/[[[ .ClusterName ]]].yml"
Driver: aliyun
Properties: 
  aliyunecs-access-key-id: [[[ .KeyID ]]]
  aliyunecs-access-key-secret: [[[ .KeySecret ]]]
  aliyunecs-image-id: [[[ .ImageID ]]]
  aliyunecs-internet-max-bandwidth: 100
  aliyunecs-private-address-only: false
  aliyunecs-region: [[[ .Region ]]]
  aliyunecs-security-group: k8s
  aliyunecs-system-disk-category: cloud_efficiency
  aliyunecs-io-optimized: optimized
Vmspecs: 
  - 
    Instances: 1
    Name: kube-master
    Properties: 
      aliyunecs-description: "kube-master,etcd"
      aliyunecs-instance-type: [[[ .MasterSize ]]]
      aliyunecs-tag: 
        - kube-master
    Roles: 
      - masters
      - etcd
  - 
    Instances: 2
    Name: kube-slave
    Properties: 
      aliyunecs-description: kube-slave
      aliyunecs-instance-type: [[[ .NodeSize ]]]
      aliyunecs-tag: 
        - kube-slave
    roles: 
      - nodes

`
