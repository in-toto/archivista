# Contributing to Archivista

We welcome contributions from the community and first want to thank you for
taking the time to contribute!

Please familiarize yourself with the [Code of Conduct](CODE_OF_CONDUCT.md)


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

### Detailed installation instructions

#### Getting the Archivista source code

[Fork](https://docs.github.com/en/get-started/quickstart/fork-a-repo>) the repository on GitHub and
[clone](https://docs.github.com/en/repositories/creating-and-managing-repositories/cloning-a-repository) it to
your local machine:

```console
    git clone git@github.com:YOUR-USERNAME/archivista.git
```

Add a [remote](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/configuring-a-remote-for-a-fork) and
regularly [sync](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/syncing-a-fork) to make sure
you stay up-to-date with our repository:

```console
    git remote add upstream https://github.com/in-toto/archivista.git
    git checkout main
    git fetch upstream
    git merge upstream/main
```

#### Setting up your environment

Archivista is written in [Go](https://golang.org/). You will need to have Go.

Archivista runs using Container and we provide some shortucts to help you
during the development process.

Some tools are required to be installed on your system to help you with
the easy development process. You can install them:

* pre-commit: It is a framework for managing and maintaining multi-language
  pre-commit hooks. It is used to run some checks on the code before
  committing it. You can install it by following the instructions
  [here](https://pre-commit.com/#install).

* jq: It is a lightweight and flexible command-line JSON processor. It is
  used to parse JSON files. You can install it by following the instructions
  [here](https://stedolan.github.io/jq/download/).

* Docker Engine: It is used to build and run the Archivista container. You
  can install it by following the instructions
  [here](https://docs.docker.com/engine/install/).

### Running Archivista Developement Environment

Archivista uses Docker to run the development environment. You can run the command:

```console
    make run-dev
```

This will run the Archivista container and its required dependencies services.
Archivista will be available at http://localhost:8082

Any change in the code will be automatically reflected in the container.

To stop the running containers, run the command:

```console
    make stop
```

To clean all Archivista containers in your environment run the command:

```console
    make clean
```


### Running Tests

You can run all the tests by running the command:

```console
    make test
```
