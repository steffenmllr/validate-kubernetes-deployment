# validate-kubernetes-deployment ![https://cloud.docker.com/u/steffenmllr/repository/docker/steffenmllr/steffenmllr/validate-kubernetes-deployment](https://img.shields.io/docker/automated/steffenmllr/validate-kubernetes-deployment.svg?style=flat-square)

> Validates kubernetes deployments and (optional) sends message to slack

### Github Action
```
action "Validate K8s Deployment" {
  uses = "steffenmllr/validate-kubernetes-deployment/action@master"
  secrets = ["SLACK_HOOK_URL]
  env = {
    KUBECONFIG  = "/path/to/KUBECONFIG/within/workspace",
    NAMESPACE   = "staging",
    DEPLOYMENTS = "backend,frontend"
  }
}

```

### Docker Usage
```bash
docker run -it
-v path/to/kubeconfig.conf:/kubeconfig.conf \
-e KUBECONFIG=/kubeconfig.conf
-e LINK="https://github.com/steffenmllr/validate-kubernetes-deployment/commit/ed96dd695757e272cb5fdddd4262b77f012e6486" \
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

LINK="https://github.com/steffenmllr/validate-kubernetes-deployment/commit/ed96dd695757e272cb5fdddd4262b77f012e6486" \
NAME="steffenmllr/validate-kubernetes-deployment v1.1.10" \
SLACK_HOOK_URL="https://hooks.slack.com/services/XXXX" \
NAMESPACE=staging \
DEPLOYMENTS=backend,frontend \
./validate

```
