# Stage 1: Build executable
FROM golang:1.11 as buildImage

# We start with migrate so this gets cached most of the time
RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR $GOPATH/src/github.com/steffenmllr/validate-kubernetes-deployment
COPY . .

RUN dep ensure
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o validate

# Stage 2: Create release image
FROM alpine:3.6
RUN apk --no-cache add ca-certificates

COPY --from=buildImage /go/src/github.com/steffenmllr/validate-kubernetes-deployment/validate ./validate
RUN chmod +x ./validate

CMD ["/validate"]
