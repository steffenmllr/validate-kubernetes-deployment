# Stage 1: Build executable
FROM golang:1.14 as buildImage


WORKDIR $GOPATH/src/github.com/steffenmllr/validate-kubernetes-deployment
COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o validate

# Stage 2: Create release image
FROM alpine:3
RUN apk --no-cache add ca-certificates

COPY --from=buildImage /go/src/github.com/steffenmllr/validate-kubernetes-deployment/validate ./validate
RUN chmod +x ./validate
COPY entrypoint.sh /entrypoint.sh

CMD ["/entrypoint.sh"]
