REGISTRY_NAME ?= erizzardi.mine.io
APPLICATION_NAME ?= vault-secrets-operator
VERSION ?= latest

DOCKERFILE_PATH ?= docker/Dockerfile

GOOS ?= linux
GOARCH ?= amd64

VAULT_URL ?= http://localhost:8200
VAULT_TOKEN ?= root
KUBECONFIG ?= .kube/config
NAMESPACE ?= vault-secrets
 
push: 
	docker push ${REGISTRY_NAME}/${APPLICATION_NAME}:${VERSION}

image: test
	@echo
	@echo
	@echo Building image
	@echo ===============================================
	docker build -t ${REGISTRY_NAME}/${APPLICATION_NAME}:${VERSION} -f ${DOCKERFILE_PATH} .


test: clean-test
	@echo
	@echo 
	@echo Executing unit tests
	@echo ===============================================
	VAULT_URL=${VAULT_URL} VAULT_TOKEN=${VAULT_TOKEN} KUBECONFIG=${KUBECONFIG} NAMESPACE=${NAMESPACE} go test ./...
	

clean-test:
	@echo
	@echo 
	@echo Cleaning test cache
	@echo ===============================================
	go clean -testcache