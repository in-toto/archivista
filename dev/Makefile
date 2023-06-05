.PHONY: deps help

SHELL=/bin/bash
DESTDIR=/usr/local/bin
ATTESTATIONDIR=.attestations
GSBUCKET=dev-archivist-542b760
ARCHIVISTASERVER=https://archivista.testifysec.localhost
OS:=$(shell uname)
HOSTS=archivista.testifysec.localhost auth.testifysec.localhost bundle.testifysec.localhost hydra-admin.testifysec.localhost hydra.testifysec.localhost judge.testifysec.localhost phpmyadmin.testifysec.localhost registrar.testifysec.localhost kratos.testifysec.localhost kratos-admin.testifysec.localhost login.testifysec.localhost grafana.testifysec.localhost mail.testifysec.localhost

help:
	@echo "Usage: make [target]"
	@echo "Available targets:"
	@echo "  deps: install/upgrade dependencies"
	@echo "  help: show this help"
	@echo "  k3d: install/upgrade k3d"
	@echo "  skaffold: install/upgrade skaffold"
	@echo "  helm: install/upgrade helm"
	@echo "  kubectl: install/upgrade kubectl"
	@echo "  go: install/upgrade go"
	@echo "  down: delete all resources"
	@echo "  up: create all resources"
	@echo "  tpm-udev-rules: add tpm-udev-rules for current user to /etc/udev/rules.d/70-tpm.rules"

## TODO: pin specific version and check hash

deps: go gcloud colima mkcert minikube skaffold helm kubectl nsenter tmux

colima: colima-$(OS)

colima-Darwin:
	command -v docker > /dev/null || brew install colima

colima-Linux:
	echo todo

go: go-$(OS)

go-Darwin:
	command -v go > /dev/null || brew install go

go-Linux:
	echo todo

gcloud: gcloud-$(OS)

gcloud-Darwin:
	command -v gcloud > /dev/null || brew install --cask google-cloud-sdk

gcloud-Linux:
	echo todo

tmux: tmux-$(OS)

tmux-Darwin:
	command -v tmux > /dev/null || brew install tmux

tmux-Linux:
	sudo apt install tmux
	mkdir -p bin

nsenter: nsenter-$(OS)

nsenter-Linux:
	docker run --rm -v $(DESTDIR):/target jpetazzo/nsenter

nsenter-Darwin:
	echo todo

mkcert: mkcert-$(OS)

mkcert-Linux:
	sudo apt install libnss3-tools
	mkdir -p bin
	curl -JLo bin/mkcert "https://dl.filippo.io/mkcert/latest?for=linux/amd64"
	sudo install ./bin/mkcert $(DESTDIR)/mkcert

mkcert-Darwin:
	brew install mkcert
	brew install nss

skaffold: skaffold-$(OS)

skaffold-Linux:
	mkdir -p bin
	curl -Lo ./bin/skaffold https://storage.googleapis.com/skaffold/releases/v2.0.0/skaffold-linux-amd64 && \
	sudo install ./bin/skaffold $(DESTDIR)/skaffold

skaffold-Darwin:
	command -v skaffold > /dev/null || brew install skaffold

helm: helm-$(OS)

helm-Linux:
	command -v helm > /dev/null || (curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | HELM_INSTALL_DIR=$(DESTDIR) $(SHELL))

helm-Darwin:
	command -v helm > /dev/null || brew install helm

kubectl: kubectl-$(OS)

kubectl-Linux:
	command -v kubectl > /dev/null || curl -Lo ./bin/kubectl "https://dl.k8s.io/release/`curl -L -s https://dl.k8s.io/release/stable.txt`/bin/linux/amd64/kubectl"
	command -v kubectl > /dev/null || (sudo install ./bin/kubectl $(DESTDIR)/kubectl)

kubectl-Darwin:
	command -v kubectl > /dev/null || brew install kubectl

minikube: minikube-$(OS)

minikube-Linux:
	command -v minikube > /dev/null || curl -Lo ./bin/minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
	command -v minikube > /dev/null || (sudo install ./bin/minikube $(DESTDIR)/minikube)

minikube-Darwin:
	command -v minikube > /dev/null || brew install minikube

tpm-permissions:
	sudo sh -c 'echo KERNEL=="tpm[0-9]*", TAG+="systemd", MODE="0660", OWNER="tpm" > /etc/udev/rules.d/70-tpm.rules'
	sudo sh -c 'echo KERNEL=="tpmrm[0-9]*", TAG+="systemd", MODE="0660", OWNER="tpm", GROUP="tpm" >> /etc/udev/rules.d/70-tpm.rules'

	sudo usermod -a -G tpm $(USER)
	sudo chown root:root /etc/udev/rules.d/70-tpm.rules
	sudo chmod 644 /etc/udev/rules.d/70-tpm.rules
	sudo udevadm control --reload-rules
	sudo udevadm trigger

down:
	pkill skaffold || true
	pkill kubectl || true
	minikube stop

up:
	./scripts/build-web.sh
	minikube status || minikube start docker --cpus max --memory max --static-ip 192.168.50.2
	minikube addons enable ingress
	minikube addons enable ingress-dns
	kubectl wait --for=condition=Ready node --all

	if [ "$(OS)" = "Linux" ]; then \
		./scripts/docker-mount.sh; \
	fi

	skaffold dev --tail=false

port-forward:
	sudo -E kubectl -n ingress-nginx port-forward service/ingress-nginx-controller 443:443 80:80

copy-attestations:
	gcloud auth login
	mkdir -p "$(ATTESTATIONDIR)"
	gsutil -m rsync -r "gs://$(GSBUCKET)/" "$(ATTESTATIONDIR)"

load-attestations: copy-attestations build-archivistactl
	./bin/archivistactl -u "$(ARCHIVISTASERVER)" store $(ATTESTATIONDIR)/*

build-witness:
	mkdir -p ./bin
	cd ../witness && go build -o bin/witness
	mv ../witness/bin/witness ./bin/witness

build-archivistactl:
	mkdir -p ./bin
	cd ../archivista && go build -o ./bin/archivistactl ./cmd/archivistactl/main.go
	mv ../archivista/bin/archivistactl ./bin/archivistactl

build-scribe:
	mkdir -p ./bin
	cd ../scribe && go build -o .bin/scribe ./main.go
	mv ../scribe/bin/scribe ./bin/scribe

build-containers:
	skaffold build

build: build-witness build-archivistactl build-scribe build-containers

clean:
	minikube delete || true
	CAROOT=.certs mkcert -uninstall
	rm -rf bin/*
	rm -rf .certs/*
	rm -rf .minio-data/*
	rm -rf .mysql-data/*

hosts:
	grep -qxF '127.0.0.1 $(HOSTS)' /etc/hosts || echo '127.0.0.1 $(HOSTS)' | sudo tee -a /etc/hosts
