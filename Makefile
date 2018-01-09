ifeq ($(origin VERSION), undefined)
  VERSION=$(git rev-parse --short HEAD)
endif

GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
REPOPATH = kubernetes-ldap

build: vendor
	GOOS=linux GOARCH=amd64 go build -o bin/linux/kubectllogin cmd/kubelogin.go
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin/kubectllogin cmd/kubelogin.go
	GOOS=windows GOARCH=amd64 go build -o bin/windows/kubectllogin.exe cmd/kubelogin.go

run:
	./bin/${GOOS}/kubectllogin

dep:
	curl -o dep -L https://github.com/golang/dep/releases/download/v0.3.2/dep-${GOOS}-${GOARCH}
	chmod +x dep

vendor: dep
	./dep ensure
	./dep status
