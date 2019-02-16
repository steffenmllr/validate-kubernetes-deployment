# validate-kubernetes-deployment [![](https://images.microbadger.com/badges/image/steffenmllr/validate-kubernetes-deployment.svg)](https://hub.docker.com/r/steffenmllr/validate-kubernetes-deployment/ "This image on Docker Hub")

> Validates kubernetes deployments and (optional) sends message to slack, exits with a non-zero status code if the deployment is not valid

### Github Action
You can either pass in a path to your `KUBECONFIG` (Tipp: Save it in your repo and encrypt it with [git-crypt](https://github.com/AGWA/git-crypt))

```hcl
action "Validate K8s Deployment" {
  uses = "docker://steffenmllr/validate-kubernetes-deployment:latest"
  secrets = ["SLACK_HOOK_URL"]
  env = {
    KUBECONFIG  = "/path/to/KUBECONFIG/within/workspace",
    NAMESPACE   = "staging",
    DEPLOYMENTS = "backend,frontend"
  }
}

```

Or you can pass in `KUBE_CONFIG_DATA` via SECRETS as base64 encoded [like in the aws action example](https://github.com/actions/aws/tree/master/kubectl#secrets)

```hcl
action "Validate K8s Deployment" {
  uses = "docker://steffenmllr/validate-kubernetes-deployment:latest"
  secrets = ["SLACK_HOOK_URL", "KUBE_CONFIG_DATA]
  env = {
    NAMESPACE   = "staging",
    DEPLOYMENTS = "backend,frontend"
  }
}

```


### Docker Usage
```bash
docker run -it \
-v path/to/kubeconfig.conf:/kubeconfig.conf \
-e KUBECONFIG=/kubeconfig.conf \
-e NAME="steffenmllr/validate-kubernetes-deployment v1.1.10" \
-e SLACK_HOOK_URL="https://hooks.slack.com/services/XXXX" \
-e NAMESPACE=staging \
-e DEPLOYMENTS=backend,frontend \
steffenmllr/validate-kubernetes-deployment:latest
```

### Local Usage

```bash
make setup
make

NAME="steffenmllr/validate-kubernetes-deployment v1.1.10" \
SLACK_HOOK_URL="https://hooks.slack.com/services/XXXX" \
NAMESPACE=staging \
DEPLOYMENTS=backend,frontend \
./validate

```

### Preview
![Preview](./.github/sample.png)
