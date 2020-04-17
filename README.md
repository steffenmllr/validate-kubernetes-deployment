# validate-kubernetes-deployment [![](https://images.microbadger.com/badges/image/steffenmllr/validate-kubernetes-deployment.svg)](https://hub.docker.com/r/steffenmllr/validate-kubernetes-deployment/ "This image on Docker Hub")

> Validates kubernetes deployments and (optional) sends message to slack, exits with a non-zero status code if the deployment is not valid

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
