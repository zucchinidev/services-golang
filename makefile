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

dev-update: build dev-load dev-apply

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
