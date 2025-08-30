# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

# RSA Keys
# 	To generate a private/public key PEM file.
# 	$ openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
# 	$ openssl rsa -pubout -in private.pem -out public.pem

run:
	go run apis/services/sales/main.go | go run apis/tooling/logfmt/main.go

help:
	go run apis/services/sales/main.go --help

version:
	go run apis/services/sales/main.go --version

curl-liveness:
	curl -il -X GET http://localhost:3000/liveness

curl-readiness:
	curl -il -X GET http://localhost:3000/readiness

curl-test-error:
	curl -il -X GET http://localhost:3000/testerror

curl-test-panic:
	curl -il -X GET http://localhost:3000/testpanic


# Admin Token
# export TOKEN=eyJhbGciOiJSUzI1NiIsImtpZCI6ImRjNzVhMzE2LWU4NjItNDVjYS1hNDhiLTBkNjdmMjI5ZDYyYiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzZXJ2aWNlIHByb2plY3QiLCJzdWIiOiI0ODAxYjg1MC1lNzBmLTRiMWYtOGZhNy1kOThhYTJkYWM2ZDEiLCJleHAiOjE3ODgwMjMxOTMsImlhdCI6MTc1NjQ4NzE5MywiUm9sZXMiOlsiQURNSU4iXX0.jMA07zrG8QbJ10lNX-xG9BRKeBFnxiJYWnu2Kk6qamiUmoLyZXSVWtmJwXHbIkXFFIvhvSPtxLFKshbaeIGIWAaQmZ0zaGgAoeM6nS8V0iDN6lBcW53Ij2tEt0mGuGy_Ds4hh34skp_rAp4gB-NK42QIdqA-oflbiGlSnTjpaQbcHBtduOzce3JyLbw9Lo2Om00kOypWjsKnM7Fm1Jylo8WIIXSxfN5JFvGzep7Ss_9qewDrNcC5tga_jGS37sVigI_nRBf0tcuDLBTMjpQ5KD8ACIb7l41SJgC0CsDEiLcL2N4MJ7pOp6saDDjZOnJuCU0zQIm3ruAs3DpW4CuPVA
# make curl-test-auth
# to generate the token, use make admin-genjwt {USER OR ADMIN}


# USER Token
# export TOKEN=eyJhbGciOiJSUzI1NiIsImtpZCI6ImRjNzVhMzE2LWU4NjItNDVjYS1hNDhiLTBkNjdmMjI5ZDYyYiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzZXJ2aWNlIHByb2plY3QiLCJzdWIiOiI0ODAxYjg1MC1lNzBmLTRiMWYtOGZhNy1kOThhYTJkYWM2ZDEiLCJleHAiOjE3ODgwMzg4MTgsImlhdCI6MTc1NjUwMjgxOCwiUm9sZXMiOlsiVVNFUiJdfQ.AEZTKMLNYOSfLPXjyL8gk4sfTuxeD44dZJ5xlTJREsp2cSkpj0WaMZo-W-WIarwsU6ZCFjm1xhDKnPaXddUlrn53Exk-5EUMMWfvS73P9Pz4yEppHr02-jsG0Qj_S1Gw8B085BjUfxMsaH5imTcAEV8YPnrEsFwFOM7Z_bk0Lv_uRPWImL0x2EG5iw8uEqqVnP1_cDzg9u6emMDOWG4zOco3fVDeUmYFgmTbOizRM-fyPhQr_Mbt43TO_65rHH8szYIAR8LwtpFiC7frTVW9tByevPwYf_7rDYlRNoSJfckHuypt2jk3--zZyCF6B5veY_pARK7ArIIpJR5tsy5MNA

curl-test-auth:
	curl -il \
	-H "Authorization: Bearer ${TOKEN}" "http://localhost:3000/testauth"

admin-genkey:
	go run apis/tooling/admin/main.go genkey

admin-genjwt: admin-genjwt-admin-role

admin-genjwt-user-role:
	go run apis/tooling/admin/main.go genjwt USER

admin-genjwt-admin-role:
	go run apis/tooling/admin/main.go genjwt ADMIN

admin-tools: admin-genkey admin-genjwt

# ==============================================================================
# Define dependencies

GOLANG          := golang:1.24.3
ALPINE          := alpine:3.21.3
KIND            := kindest/node:v1.33.1
POSTGRES        := postgres:17.5
GRAFANA         := grafana/grafana:12.0.1
PROMETHEUS      := prom/prometheus:v3.4.0
TEMPO           := grafana/tempo:2.7.2
LOKI            := grafana/loki:3.5.1
PROMTAIL        := grafana/promtail:3.5.1

KIND_CLUSTER    := sales-starter-cluster
NAMESPACE       := sales-system
SALES_APP       := sales
AUTH_APP        := auth
BASE_IMAGE_NAME := localhost/sales
VERSION         := 0.0.1
# Some containers systems needs a url-based image name
SALES_IMAGE     := $(BASE_IMAGE_NAME)/$(SALES_APP):$(VERSION)
METRICS_IMAGE   := $(BASE_IMAGE_NAME)/metrics:$(VERSION)
AUTH_IMAGE      := $(BASE_IMAGE_NAME)/$(AUTH_APP):$(VERSION)

# ==============================================================================
# Running from within k8s/kind

build: sales

sales:
	docker build \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ") \
		-t $(SALES_IMAGE) \
		-f zarf/docker/dockerfile.sales \
		.

dev-up:


	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/dev/kind-config.yaml

	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner

dev-down:
	kind delete cluster --name $(KIND_CLUSTER)

dev-status-all:
	kubectl get nodes -o wide
	kubectl get svc -o wide --all-namespaces
	kubectl get pods -o wide --watch --all-namespaces

dev-status:
	kubectl get pods -o wide --all-namespaces --watch

# ==============================================================================

dev-update-kustomization:
	# Update the dev kustomization file with the current git reference
	# This ensures the deployment uses the correct image tag that matches the built image
	sed -i '' 's|newTag:.*|newTag: $(VERSION)|' zarf/k8s/dev/sales/kustomization.yaml

dev-load:
	# Load the sales service image into the Kind cluster
	# - kind: The Kubernetes IN Docker tool
	# - load: Command to import an image from local storage (works offline)
	# - docker-image: Specifies we're loading a Docker image
	# - $(SALES_IMAGE): The sales service image reference
	# - --name: Specify which cluster to load into
	# - $(KIND_CLUSTER): The name of our Kind cluster
	kind load docker-image $(SALES_IMAGE) --name $(KIND_CLUSTER)

dev-apply:
	kustomize build zarf/k8s/dev/sales | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --for=condition=Ready --timeout=120s --selector app=$(SALES_APP)

dev-restart:
	kubectl rollout restart deployment $(SALES_APP) --namespace=$(NAMESPACE)

dev-update: build dev-load dev-apply dev-restart

dev-logs:
	kubectl logs --namespace=$(NAMESPACE) --selector app=$(SALES_APP) --all-containers=true -f --tail=100 --max-log-requests=6 | go run apis/tooling/logfmt/main.go

dev-describe-deployment:
	kubectl describe deployment $(SALES_APP) --namespace=$(NAMESPACE)

dev-describe-sales:
	kubectl describe pod --selector app=$(SALES_APP) --namespace=$(NAMESPACE)


# ==============================================================================
# Metrics and Tracing

metrics:
	go tool expvarmon -ports="localhost:3010" -vars="build,requests,goroutines,errors,panics,mem:memstats.HeapAlloc,mem:memstats.HeapSys,mem:memstats.Sys"

statsviz:
	open -a "Google Chrome" http://localhost:3010/debug/statsviz

# ==============================================================================
# Modules support

tidy:
	go mod tidy
	go mod vendor

# ==============================================================================
# Install tools
# Uses Go 1.24's new tool dependency management

tools-add-expvarmon:
	go get -tool github.com/divan/expvarmon@latest

tools-add-staticcheck:
	go get -tool honnef.co/go/tools/cmd/staticcheck@latest

tools-add-govulncheck:
	go get -tool golang.org/x/vuln/cmd/govulncheck@latest

tools-list:
	go list tool

tools-upgrade:
	go get tool

tools-add: tools-add-expvarmon tools-add-staticcheck tools-add-govulncheck
	@echo "All tools added to go.mod!"

# ==============================================================================
# Running tests within the local computer

test-r:
	CGO_ENABLED=1 go test -race -count=1 ./...

test-only:
	CGO_ENABLED=0 go test -count=1 ./...

lint:
	CGO_ENABLED=0 go vet ./...
	go tool staticcheck -checks=all ./...

vuln-check:
	go tool govulncheck ./...

test: test-only lint vuln-check

test-race: test-r lint vuln-check
