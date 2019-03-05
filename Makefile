NAME   := hotelsdotcom/kube-graffiti
TAG    := $(shell git describe --tags)
IMG    := ${NAME}:${TAG}
LATEST := ${NAME}:latest

export AZURE_STORAGE_ACCOUNT=msahelmrepos
# For old auth method
# export AZURE_STORAGE_KEY=$(az storage account keys list --resource-group ${AZURE_STORAGE_ACCOUNT} --account-name ${AZURE_STORAGE_ACCOUNT} --query "[0].value")
export AZURE_STORAGE_CONTAINER=incubator
export HELM_REPO_NAME=msa-incubator
export HELM_PACKAGE="kube-graffiti-0.8.3.tgz"


build:test
	@docker build -t "${IMG}" .
	@docker tag ${IMG} ${LATEST}
 
push:
	@docker push ${NAME}
 
login:
	@docker log -u ${DOCKER_USER} -p ${DOCKER_PASS}

chart:
	@helm package --app-version ${TAG} ./helm/kube-graffiti

upload_chart: chart
	@az storage blob upload --file ${HELM_PACKAGE} --container-name ${AZURE_STORAGE_CONTAINER} --name ${HELM_REPO_NAME}/${HELM_PACKAGE} && \
    az storage blob download-batch --source ${AZURE_STORAGE_CONTAINER} --pattern '${HELM_REPO_NAME}/*' --destination . && \
    helm repo index ${HELM_REPO_NAME}/ && \
    az storage blob upload --file ${HELM_REPO_NAME}/index.yaml --container-name ${AZURE_STORAGE_CONTAINER} --name ${HELM_REPO_NAME}/index.yaml
	
test:
	@go test ./...
