DOCKER_USERNAME ?= erizzardi.mine.io
APPLICATION_NAME ?= vault-secrets-operator
VERSION ?= latest

DOCKERFILE_PATH ?= docker/Dockerfile

GOOS ?= linux
GOARCH ?= amd64

VAULT_URL ?= http://localhost:8200
VAULT_TOKEN ?= myroot
KUBECONFIG ?= .kube/config
NAMESPACE ?= vault-secrets
 
image:
	docker build -t ${DOCKER_USERNAME}/${APPLICATION_NAME}:${VERSION} -f ${DOCKERFILE_PATH} .

test: clean
	VAULT_URL=${VAULT_URL} VAULT_TOKEN=${VAULT_TOKEN} KUBECONFIG=${KUBECONFIG} NAMESPACE=${NAMESPACE} go test ./...

clean:
	go clean -testcache

all: test image
