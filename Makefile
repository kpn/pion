REGISTRY = quay.io/canhnt


# ONLY use 'latest' tag in development
#TAG = ${shell git describe --always}
TAG = latest
PKG = github.com/kpn/pion

# -----------------------------------
# Security token service
STS_IMAGE = $(REGISTRY)/pion-sts
STS_DOCKERFILE = build/pion-sts/Dockerfile
STS_BUILD_OUTPUT := build/pion-sts/pion-sts

.PHONY: sts-build
sts-build: ${STS_BUILD_OUTPUT}
${STS_BUILD_OUTPUT}:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ${STS_BUILD_OUTPUT} ${PKG}/cmd/pion-sts

.PHONY: sts-image
sts-image: sts-build
		docker build --tag $(STS_IMAGE):$(TAG) -f $(STS_DOCKERFILE) build
		docker push $(STS_IMAGE):$(TAG)

.PHONY: sts-clean
sts-clean:
		rm -f ${STS_BUILD_OUTPUT}
# -----------------------------------
# TCloud-OSS UI
UI_IMAGE = $(REGISTRY)/pion-ui
UI_DOCKERFILE = build/pion-ui/Dockerfile
UI_BUILD_OUTPUT := build/pion-ui/pion-ui

.PHONY: ui-build
ui-build: ui-build-be ui-build-fe

.PHONY: ui-build-be
ui-build-be: ${UI_BUILD_OUTPUT}

.PHONY: ui-build-fe-dev
ui-build-fe-dev:
		make -C ui build-remote
		mv ui/dist build/pion-ui/dist

.PHONY: ui-build-fe-prod
ui-build-fe-prod:
		make -C ui build-prod
		mv ui/dist build/pion-ui/dist

ifeq ($(BUILD_ENV),dev)
ui-build-fe: ui-build-fe-dev
else
ui-build-fe: ui-build-fe-prod
endif

${UI_BUILD_OUTPUT}:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ${UI_BUILD_OUTPUT} ${PKG}/cmd/pion-ui

.PHONY: ui-image
ui-image: ui-build
		docker build --tag $(UI_IMAGE):$(TAG) -f $(UI_DOCKERFILE) build
		docker push $(UI_IMAGE):$(TAG)

.PHONY: ui-clean
ui-clean:
		rm -f ${UI_BUILD_OUTPUT}
		rm -rf build/pion-ui/dist
		make -C ui clean
# -----------------------------------
# TCloud-OSS Proxy
PROXY_IMAGE = $(REGISTRY)/pion-proxy
PROXY_DOCKERFILE = build/pion-proxy/Dockerfile
PROXY_BUILD_OUTPUT := build/pion-proxy/pion-proxy

.PHONY: proxy-build
proxy-build: ${PROXY_BUILD_OUTPUT}
${PROXY_BUILD_OUTPUT}:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ${PROXY_BUILD_OUTPUT} ${PKG}/cmd/pion-proxy

.PHONY: proxy-image
proxy-image: ${PROXY_BUILD_OUTPUT}
		docker build --tag $(PROXY_IMAGE):$(TAG) -f $(PROXY_DOCKERFILE) build
		docker push $(PROXY_IMAGE):$(TAG)

.PHONY: proxy-clean
proxy-clean:
		rm -f ${PROXY_BUILD_OUTPUT}
# -----------------------------------
# TCloud-OSS Authz
AUTHZ_IMAGE = $(REGISTRY)/pion-authz
AUTHZ_DOCKERFILE = build/pion-authz/Dockerfile
AUTHZ_BUILD_OUTPUT := build/pion-authz/pion-authz

.PHONY: authz-build
authz-build: ${AUTHZ_BUILD_OUTPUT}
${AUTHZ_BUILD_OUTPUT}:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ${AUTHZ_BUILD_OUTPUT} ${PKG}/cmd/pion-authz

.PHONY: authz-image
authz-image: ${AUTHZ_BUILD_OUTPUT}
		docker build --tag $(AUTHZ_IMAGE):$(TAG) -f $(AUTHZ_DOCKERFILE) build
		docker push $(AUTHZ_IMAGE):$(TAG)

.PHONY: authz-clean
authz-clean:
		rm -f ${AUTHZ_BUILD_OUTPUT}
# -----------------------------------
# TCloud-OSS Manager
MGR_IMAGE = $(REGISTRY)/pion-manager
MGR_DOCKERFILE = build/pion-manager/Dockerfile
MGR_BUILD_OUTPUT := build/pion-manager/pion-manager
MGR_BOOTSTRAP := build/pion-manager/bootstrap

.PHONY: mgr-build
mgr-build: ${MGR_BUILD_OUTPUT}
${MGR_BUILD_OUTPUT}:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ${MGR_BUILD_OUTPUT} ${PKG}/cmd/pion-manager
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ${MGR_BOOTSTRAP} ${PKG}/cmd/bootstrap

.PHONY: mgr-image
mgr-image: ${MGR_BUILD_OUTPUT}
		docker build --tag $(MGR_IMAGE):$(TAG) -f $(MGR_DOCKERFILE) build
		docker push $(MGR_IMAGE):$(TAG)

.PHONY: mgr-clean
mgr-clean:
		rm -f ${MGR_BUILD_OUTPUT}
		rm -f ${MGR_BOOTSTRAP}
# -----------------------------------

.PHONY: test
test:
		go test -v $$(go list ./... | grep -ve 'test/e2e')

.PHONY: image
image: sts-image ui-image authz-image proxy-image mgr-image

.PHONY: build
build: sts-build authz-build proxy-build mgr-build ui-build


.PHONY: clean
clean: sts-clean ui-clean proxy-clean authz-clean mgr-clean



