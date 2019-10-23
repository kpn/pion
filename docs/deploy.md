# Deployment
Deployment are carried out by the `pion` Helm chart in `$PROJECTDIR/charts/pion`. 

Deployment steps are as follows:

## Install ETCD cluster
You can either choose to install an Etcd cluster manually or use [Etcd Operator](https://github.com/coreos/etcd-operator).
After install the Etcd Operator, you can bring up the Etcd cluster with `$PROJECTDIR/k8s/etcd/etcd-cluster.yaml`

## Install Minio instance/cluster
You can use Minio Helm chart provided by the TCloud services catalog. After installation, find the Access Key and Secret Key used to access
Minio in a K8S Secret.
 
## Deploy Pion
Use helm chart 

## Client Configuration
We suggest to use [AWS-CLI S3](https://docs.aws.amazon.com/cli/latest/reference/s3/index.html) at the client-side to 
upload and download files from Pion. Detail installation can be found [here](https://docs.aws.amazon.com/cli/latest/userguide/installing.html).
 
Client side configurations are as follows: 

- Login to the UI Dashboard to generate and retrieve an access/secret key-pair. You can store this key-pair in 
[config file](https://docs.aws.amazon.com/cli/latest/userguide/cli-config-files.html) 
`$HOME/.aws/credentials` or [environment variables](https://docs.aws.amazon.com/cli/latest/userguide/cli-environment.html) 
`AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY`.

- Set the `AWS_CA_BUNDLE` environment variable to the private CA certificate, e.g.: 
```bash
export AWS_CA_BUNDLE=$HOME/.aws/PrivateRootCA.crt
``` 

- Set the alias for the aws-cli so that it always uses the custom endpoint of Pion Gateway rather than default AWS S3:
```bash
alias pion='aws s3 --endpoint=https://pion-gw.example.com'
```

- Verify if the settings work well, e.g. :
```bash
17:25 $ pion ls
2018-07-25 11:53:10 mybucket
2018-07-25 13:50:33 public
2018-07-25 10:58:57 pion-devs
```

