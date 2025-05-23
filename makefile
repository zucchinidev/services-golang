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
SALES_IMAGE     := $(BASE_IMAGE_NAME)/$(SALES_APP):$(VERSION)
METRICS_IMAGE   := $(BASE_IMAGE_NAME)/metrics:$(VERSION)
AUTH_IMAGE      := $(BASE_IMAGE_NAME)/$(AUTH_APP):$(VERSION)

# ==============================================================================
# Running from within k8s/kind

dev-up:
	# System configuration for Kubernetes cluster performance and stability
	# Increase the maximum number of inotify instances per user
	# This helps with file system monitoring and prevents "too many open files" errors
	sudo sysctl -w fs.inotify.max_user_instances=524288

	# Increase the maximum number of inotify watches per user
	# Required for proper file system event monitoring in Kubernetes
	sudo sysctl -w fs.inotify.max_user_watches=524288

	# Set the maximum number of file handles that can be opened system-wide
	# Critical for handling many concurrent connections in Kubernetes
	sudo sysctl -w fs.file-max=2097152

	# Increase the maximum number of memory map areas a process can use
	# Required for proper memory management in containerized applications
	sudo sysctl -w vm.max_map_count=262144

	# Increase the maximum number of connection requests that can be queued
	# Helps prevent connection drops under high load
	sudo sysctl -w net.core.somaxconn=32768

	# Increase the maximum number of remembered connection requests
	# Improves handling of TCP connection requests
	sudo sysctl -w net.ipv4.tcp_max_syn_backlog=8192

	# Expand the range of local ports available for outgoing connections
	# Prevents port exhaustion in high-traffic scenarios
	sudo sysctl -w net.ipv4.ip_local_port_range="1024 65535"

	# Set TCP connection timeout to 30 seconds
	# Helps clean up stale connections more quickly
	sudo sysctl -w net.ipv4.tcp_fin_timeout=30

	# Set TCP keepalive time to 30 minutes
	# Helps maintain connection health
	sudo sysctl -w net.ipv4.tcp_keepalive_time=1800

	# Set number of TCP keepalive probes to 5
	# Determines how many times to retry before dropping a connection
	sudo sysctl -w net.ipv4.tcp_keepalive_probes=5

	# Set TCP keepalive interval to 15 seconds
	# Time between keepalive probes
	sudo sysctl -w net.ipv4.tcp_keepalive_intvl=15

	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/dev/kind-config.yaml

	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner


dev-down:
	kind delete cluster --name $(KIND_CLUSTER)

dev-status-all:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

dev-status:
	watch -n 2 kubectl get pods -o wide --all-namespaces

# ==============================================================================
# Modules support

tidy:
	go mod tidy
	go mod vendor

# ==============================================================================
# Running tests within the local computer

test-r:
	CGO_ENABLED=1 go test -race -count=1 ./...

test-only:
	CGO_ENABLED=0 go test -count=1 ./...

lint:
	CGO_ENABLED=0 go vet ./...
	staticcheck -checks=all ./...

vuln-check:
	govulncheck ./...

test: test-only lint vuln-check

test-race: test-r lint vuln-check
