## Build image

# 1. use golang image with 1.17-alpine tag as base builder for `deploy image`
FROM golang:1.17-alpine as builder
# 2. define exposed environment variable
ENV APP_HOME $GOPATH/src/github.com/go-seidon/hippo
# 3. update os index packages
RUN apk update
# 4. install package (git, ca-certificates, update-ca-certificates)
RUN apk add --no-cache git
RUN apk add --no-cache ca-certificates && update-ca-certificates
# 5. add sepcialized user for running the container
RUN adduser -D -g '' app-user
# 6. set working directory
WORKDIR "$APP_HOME"
# 7. copy all files and folder to workdir
COPY . .
# 8. validate, download, verify and make vendor folder for go dependencies
RUN go mod tidy && go mod verify && go mod vendor
# 9. build golang code into binary app
RUN go build -o /main ./cmd/hybrid-app/main.go
# 10. copying app required files to root folder
COPY ./config /config
COPY ./migration /migration

## Deploy image

# 1. use apline image with 3.14 tag
FROM alpine:3.14
# 2. install package (bash)
RUN apk add --no-cache bash
# 3. copying system required files
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /main /main
COPY --from=builder /config /config
COPY --from=builder /migration /migration
# 5. set user for running container
USER app-user
# 6. expose ports used by the rest and grpc app
EXPOSE 3000 5000
# 7. set entry point to the binary file generated from previous step
ENTRYPOINT ["/main"]
