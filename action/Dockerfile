FROM steffenmllr/validate-kubernetes-deployment:latest

LABEL version="1.0.0"
LABEL com.github.actions.name="Validate Deployment"
LABEL com.github.actions.description="Validates a kubernetes deployment"
LABEL com.github.actions.icon="check-circle"
LABEL com.github.actions.color="blue"

RUN apk --no-cache add jq
ADD entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

