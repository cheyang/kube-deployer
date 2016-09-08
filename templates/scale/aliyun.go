package scale

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
`
