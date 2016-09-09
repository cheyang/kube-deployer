# Kubernetes on aliyun

## Create Cluster

### Step 1: Install kube-aliyun

```
go get github.com/cheyang/kube-deployer/cmd
cp $GOPATH/bin/cmd /usr/local/bin/kube-aliyun
```

### Step 2: configure Environment variables

```
export ALIYUNECS_KEY_ID=<your key>
export ALIYUNECS_KEY_SECRET=<your secret>
export ALIYUNECS_CLUSTER_NAME=mycluster
export ALIYUNECS_IMAGE_ID=entos7u2_64_40G_cloudinit_20160520.raw
export ALIYUNECS_MASTER_SIZE=ecs.n1.small
export ALIYUNECS_NODE_SIZE=ecs.n1.small
```

### Step 3: Run deploy

```
kube-aliyun create --num-nodes 2
```