# Judge Platform

This is the internal monorepo project for the Judge Platform, a SaaS (Software as a Service) that provides ready-to-use open source software for supply-chain security.

## Internal Monorepo Structure

This monorepo is for internal use only and serves as the source of truth for all TestifySec work on the platform. Development should be done within this monorepo, which acts as our home. You are free to work on any project within the monorepo as you would in a normal project.

At the root of the monorepo, we can share configuration files, scripts, and more across the project. However, all subtrees must be self-contained. More details on this will be provided later.

### Private Subfolders vs Public Subtrees

The monorepo is divided into several subfolders at the root. Most subfolders contain internal private code that supports the platform (e.g., `web/`, `dev/`, `judge-api/`). Other subfolders are [Git Subtrees](https://www.atlassian.com/git/tutorials/git-subtree) linked to public repositories. These special subfolders require synchronization between the monorepo and the associated public repositories. When making changes to these subfolders, it's important to be mindful of the synchronization process.

### Prerequisites

Before you can fully contribute, make sure you have completed the following prerequisites:

1. Have a GitLab/GitHub account with access to the TestifySec Git repositories.
1. Set up SSH keys or HTTP login for your GitLab/GitHub accounts. You can try using the 1Password CLI with Google Cloud and GitHub SSH keys for convenience.
1. Have the necessary roles to access the attestations via gcloud in the `load attestations` command.
1. Have your TestifySec physical security key provisioned.
1. Strongly encouraged to have [nvm](https://github.com/nvm-sh/nvm) installed and configured on your machine to synchronize with the entire team on the Node.js version. You can try installing it with `brew install nvm` and following the setup instructions. Once you have `nvm`, run `nvm use` and `nvm install` to get in sync with the team's specified Node.js version.
1. Run `npm i` from the root of this repository to install all dependencies. This should also `go get` all go dependencies for all of our go projects.
1. Run `make hosts` from the `dev/` folder at least once to set up the hosts file for local development.

## Getting Started

Assuming you have completed all the prerequisites mentioned above, follow these steps to get started:

1. Run `npm start` to start everything in your development environment using an npm script. This includes running `make deps`, `make up`, and setting up port forwarding.

   Note: If you use this method, `minikube tunnel` will be used, which may require your password. Be sure not to miss the prompt and press Enter in the minikube tunnel pane if necessary.

   If you prefer a more manual process, you can follow these alternative steps:

1. Run `make deps` from the `dev/` subfolder to install dependencies.
1. You should now have everything you need to start contributing. Follow the instructions in the `dev/readme.md` file for the remaining steps, including running `make up`, `make load-attestations`, and `minikube tunnel`.

Once your local environment is set up, you can make changes to the repository as needed.

### Load Attestations

If you want to use a local Archivista and local Judge API with the web project, you need to load the attestations. This requires using your physical security key to log in and load all the attestation

 data to your machine.

To load the attestations, run `make load-attestations` from the `dev/` folder.

### Getting Started with web/

If you have the local Kubernetes environment running as mentioned above, the web project should already be available in a locally deployed production instance at [https://judge.testifysec.localhost/](https://judge.testifysec.localhost/).

To run a development instance locally with HMR (hot module reloading), use the following commands:

1. Run `npm run start:web` to start the web project in HMR mode.

By default, running the web project this way will use the full local Kubernetes environment, including Archivista and Judge, for comprehensive local development. If you want to connect to the remote proxies, follow the instructions below.

#### web/ with Remote Proxies

Note: This may require reintroducing Hydra with Kratos to allow for multiple domains in the login process.

To run the web project connected to production data, use the following commands:

1. Run `npm run start:web:remote-proxy` to start the web project in HMR mode, connected to the production APIs as remote proxies.

### How to generate code changes to our sub projects

Some of our projects utilize code generation to assist in abstracting away boilerplate. Namely, we have some go projects that use Ent and gqlgen. 

You can generate what you need from inside those project folders running the `go generate ./... -v` command, but we have also provided shortcuts in the root folder for you. 

- `npm run gen:archivista` will `go generate` all the archivista things
- `npm run gen:judge-api` will `go generate` all the judge-api things