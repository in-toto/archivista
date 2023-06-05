# Judge Development Automation

This repo contains scripts and automation to help with the development and testing of the Judge platform.
It will install all deps, clone all sub-repos, and help you start up and tear down the Judge Platform locally in docker.
It will also help you build and deploy changes to gitlab.

## Getting Started

### Installing Deps

1. `make deps` will Install Dependencies

## Frequent Rituals

### Starting Judge Platform up

**NOTE**: You will want to leave all of these commands in the terminal running in the background.

1. `make up` will start the Judge platform using minikube and skaffold.

1. `make port-forward` will Port Forward the Judge platform to your local network.

1. `make load-attestations` will Load the attestation data into the Judge platform.

Now you have the the Platform running and you can get access Judge from your localhost.

### Build and Deploy Containers to Gitlab

1. `make build` will build client applications and push containers to gitlab.

### Tearing Judge Platform down and cleaning it up

1. `make down` will tear down the Judge platform but leave the persistent volumes intact.

1. `make clean` will delete all transient data including the minikube cluster, the persistent volumes, and attestation data.

### Dev the Web-UI with Hot Module Reloading (HMR)

**NOTE**: You will want to have the Judge Platform up before moving on.

By default, with the Judge Platform up on your localhost, it will not have HMR turned on. To enable HMR, you will need to run the web application locally. 

1. open the `web/` folder in your favorite editor + terminal
1. `npm run start' to start up the app in development mode

You should notice that the web application at <https://judge.testifysec.localhost> starts to update with your local changes to the web project. You may notice a delay between saving your changes and seeing your changes reflected in your browser, and you may need to refresh and/or clear cache for all changes to reflect.

## Persistance

Persistent data is stored in the `.mysql-data` and `.minio-data` directories. These directories are mounted into the mysql and minio containers. If you want to start with a clean slate, delete these directories with `make clean`.

## TLS Certificates

The creation and configuration of TLS certs is handled by `mkcert`. mkcert generates a root CA that is injected into the minikube cluster and local trust stores. The root CA is then used to generate certificates for the Judge platform using `cert-manager`. The certificates are stored in the `certs` directory. If you want to start with a clean slate, delete the `certs` directory with `make clean`.

## Development

### Web Dev

The web application is avalaible at <https://judge.testifysec.localhost>. This instance will not have HMR enabled unless you follow the steps to enable it above.

### Backend Microservice Dev

Microservices are specified in both the `skaffold.yaml` file and the kubernetes manifests in the K8s directory. Microservices are build and redeployed when changes are detected.

### Tools

Tools are utilities that are installed on your local machine to help devops.

Tools will be installed automatically with the `make deps` command.

- [nsenter](https://github.com/jpetazzo/nsenter) is used to enter the minikube container to mount the persistent volumes.

- [Skaffold](https://skaffold.dev/) is a command line tool that facilitates continuous development for Kubernetes applications. Skaffold handles the workflow for building, pushing and deploying your application, allowing you to focus on what matters most: writing code. It can be used with any container registry, continuous integration or deployment system.

- [Kustomize](https://kustomize.io/) is a tool to customize Kubernetes objects through a kustomization file. Kustomize has several new features that make it easier to manage Kubernetes applications in production.

- [Helm](https://helm.sh/) is a tool for managing Kubernetes charts. Charts are packages of pre-configured Kubernetes resources.

- [MiniKube](https://minikube.sigs.k8s.io/docs/) is a tool that makes it easy to run Kubernetes locally. MiniKube runs a single-node Kubernetes cluster inside a Virtual Machine (VM) on your laptop for users looking to try out Kubernetes or develop with it day-to-day.

- [Colima](https://github.com/abiosoft/colima) is a Docker daemon for MacOS that doesn't require Docker Desktop. You can start it up after installing it for this proj easily with  `colima start --cpu 4 --memory 8`

### Components

Components are dependencies to the Judge Platform that will be running in containers.

Components will be installed automatically with the `make deps` command.

- [Hydra](https://www.ory.sh/hydra/) - OpenID Certified OAuth 2.0 Server and OpenID Connect Provider optimized for low-latency, high throughput, and low resource consumption.

- [Spire](https://gitlab.com/testifysec/judge-platform/spire-internal) - SPIFFE is a framework for identifying and securing communications between services. SPIRE is a reference implementation of SPIFFE. We use an internal fork of SPIRE to support bindings to our accounting systems.

- [MySQL](https://www.mysql.com/) - MySQL is an open-source relational database management system (RDBMS). We use MySQL as the database for the Judge platform.

- [Registrar](https://gitlab.com/testifysec/judge-platform/registrar) - Rigistrar handles making and verifying node registrations with the Judge platform.

- [auth-login-provider](https://gitlab.com/testifysec/judge-platform/auth-login-provider) - auth-login-provider is a simple web application that provides a login page for the Judge platform.
