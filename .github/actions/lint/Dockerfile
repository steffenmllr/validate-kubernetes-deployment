FROM golangci/golangci-lint:v1.14

LABEL version="1.0.0"
LABEL com.github.actions.name="Linter"
LABEL com.github.actions.description="Lints the golang api"
LABEL com.github.actions.icon="hash"
LABEL com.github.actions.color="green"

RUN go get -u github.com/golang/dep/cmd/dep

ADD entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
