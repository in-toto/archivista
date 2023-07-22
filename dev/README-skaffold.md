# Skaffold Configuration Guide

This document aims to help you understand our Skaffold configuration setup, add new modules, and start up different combinations of modules based on your needs.

## Understanding the Configuration

Our Skaffold configuration is modularized, divided into different `yaml` files. Each module encapsulates a specific set of services or functionalities, such as:

- `init.yaml`: Sets up namespaces.
- `local.yaml`: Deploys services necessary for local development.
- `auth.yaml`: Handles authentication.
- `core.yaml`: Deploys the core application services.
- `observability.yaml`: Provides observability tools (like Grafana and Loki) for monitoring.

The base module `skaffold.yaml` requires all the basic configurations mentioned above. We also have a module `with-observability.yaml`, which includes observability services on top of the base configuration.

## How to Add New Modules

If you need to add a new set of services or functionalities, follow these steps:

1. Create a new `yaml` file under the `skaffold` directory. The file name should describe the purpose of the module, e.g., `new-service.yaml`.

2. Inside the new yaml file, define your Skaffold configuration. Ensure you have `apiVersion: skaffold/v3alpha1`, `kind: Config`, and under `metadata`, set the `name` to your new module's name.

3. In the `requires` section of the module where you want to include this new module (e.g., `base` or `with-observability`), add a new path to your module. For example: `- path: ./skaffold/new-service.yaml`.

Now your new module is integrated into the Skaffold configuration and will be used whenever the parent module is invoked.

## How to Start Different Modules

To start a specific module, use the `skaffold dev --module <module-name>` command. For instance, to start the base module, use:

```
skaffold dev --module base --tail=false
```

You can also start multiple modules at once by specifying each module's name separated by a comma:

```
skaffold dev --module base,new-service --tail=false
```

In this example, both `base` and `new-service` modules will be started.

Note: The `--tail=false` flag is used to prevent Skaffold from streaming logs.

## How to Mix and Match Modules

You can selectively start the services you need by choosing the right modules. If you want to start the core services with observability, run:

```
skaffold dev --module with-observability --tail=false
```

This command will start all services included in the `base` module, plus the observability services.

Remember, the power of this modular configuration lies in the ability to start just what you need without having to manage a complex monolithic Skaffold configuration file. Happy developing!