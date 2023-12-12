# Contributing to Archivista

We welcome contributions from the community and first want to thank you for
taking the time to contribute!

Before starting, please take some time to familiarize yourself with the [Code of Conduct](CODE_OF_CONDUCT.md).


## Getting started

We welcome many different types of contributions and not all of them need a
Pull Request. Contributions may include:

* New features and proposals
* Documentation
* Bug fixes
* Issue Triage
* Answering questions and giving feedback
* Helping to onboard new contributors
* Other related activities

### Setting up your environment

### Required Tooling
Some tools are required on your system in order to help you with
the development process:

* Git: Archivista is hosted on GitHub, so you will need to have Git installed. For
 more information, please follow [this guide](https://github.com/git-guides/install-git).

* GNU Make: The root of the directory contains a `Makefile` for automating development
 processes. The `make` CLI tool is usually installed by default on most systems
 (excluding Windows), but you can check if it is installed by running `make --version`
 on your terminal. If this command is unsuccessful, you will need to find the standard
 method for installing it for your system. For installing `make` on Windows, please see
 [here](https://gnuwin32.sourceforge.net/packages/make.html).
 
* Go v1.19: Archivista is written in [Go](https://golang.org/), so you 
 will need this installed in order to compile and run the source code.

* pre-commit: It is a framework for managing and maintaining multi-language
 pre-commit hooks. It is used to run some checks on the code before it is 
 committed. You can install it by following the instructions
 [here](https://pre-commit.com/#install).

* jq: It is a lightweight and flexible command-line JSON processor. It is
 used to parse JSON files. You can install it by following the instructions
 [here](https://stedolan.github.io/jq/download/).

* A Container Runtime: Archivista is a container-based tool, and therefore you might
 need a container runtime tool installed for development. The most common solution
 for this is [Docker](https://docs.docker.com/engine/install/), however this is not
 free on all platforms. If you are using MacOS, [Colima](https://github.com/abiosoft/colima)
 is a free good alternative, as well as [Rancher Desktop](https://github.com/rancher-sandbox/rancher-desktop).
 Please note that Archivista is not tested on these alternative tools, so compatibility
 is not guaranteed.

#### Getting the Archivista source code

[Fork](https://docs.github.com/en/get-started/quickstart/fork-a-repo>) the repository on GitHub and
[clone](https://docs.github.com/en/repositories/creating-and-managing-repositories/cloning-a-repository) it to
your local machine: 
```console
    git clone git@github.com:YOUR-USERNAME/archivista.git
```
*The command above uses SSH to clone the repository, which we recommend. You can find out more
about how to set SSH up with Github [here](https://docs.github.com/en/authentication/connecting-to-github-with-ssh).*


Add a [remote](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/configuring-a-remote-for-a-fork) and
regularly [sync](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/syncing-a-fork) to make sure
you stay up-to-date with our repository:

```console
    git remote add upstream https://github.com/in-toto/archivista.git
    git checkout main
    git fetch upstream
    git merge upstream/main
```

### Running Archivista Development Environment

*Please note that the following `make` commands make use of both the `docker` and
`docker-compose` commands, so you may need to modify this locally if using tools
such as [nerdctl](https://github.com/containerd/nerdctl) or [podman](https://github.com/containers/podman).*

To start the Archivista development environment, simply execute the command:
```console
    make run-dev
```
This will run the Archivista container and its required dependent services.
Archivista will be available at http://localhost:8082

Any changes made to the source code will be reflected in the development environment
while it is running, so there is no need to restart it.

To stop the development environment, run:

```console
    make stop
```

To clean all Archivista containers in your environment execute the command:

```console
    make clean
```


### Running Tests

You can run all the tests by executing the command:

```console
    make test
```
