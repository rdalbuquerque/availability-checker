# Availability Checker
- [Availability Checker](#availability-checker)
  - [Description](#description)
  - [Pre-requisites](#pre-requisites)
  - [Credential Providers](#credential-providers)
  - [Project Structure](#project-structure)
  - [Overview](#overview)
  - [Example usage](#example-usage)
    - [Adding new checks](#adding-new-checks)
    - [Web interface](#web-interface)
  - [Core Concepts](#core-concepts)
  - [Future Considerations](#future-considerations)

## Description
This project aims to explore Go extensability through the usage of interfaces. `Checker` interface establishes methods for new checkers to be added and `CredentialProvider` interface provides a way to add new credential providers with ease, here we implemented both Hachicorp's Vault and Azure Key Vault providers. Interfaces are also used here to enable dependency injection, allowing us to properly mock and test database connections.

Another Go feature this project explore is concurrency. Each check is executed in it's own goroutine reporting back to a channel when it's done, allowing us to add more checks with minimal impact to verification time.

This project also explores Go templating capabilities to render and serve a simple HTML page summarizing the status of each checker.

## Pre-requisites
It depends on you credential provider of choice.
For HCP Vault you need to set the following environment variables:
```bash
HCPVAULT_ADDR
VAULT_TOKEN
```
For Azure Key Vault you need to set the following environment variables:
```bash
AZURE_TENANT_ID
AZURE_CLIENT_ID
AZURE_CLIENT_SECRET
AZURE_KEYVAULT
```

## Credential Providers
In the current version of this project there are two credential providers implemented, Hachicorp's Vault and Azure Key Vault. Both of them implement the `CredentialProvider` interface, which is responsible for retrieving credentials for a given service/resource.

To chosse a credential provider, simply set one of the pre-requisites environment variables. If `HCPVAULT_ADDR` is set, the Hachicorp's Vault provider will be used, otherwise the Azure Key Vault provider will be used.

The Azure Key Vault provider expects the environment variable `AZURE_KEYVAULT` to be set with the name of the vault that will hold all credentials for the checkers. The credentials are expected to be stored as secrets in the vault in the following format: `{checkerType}-user` and `{checkerType}-pwd`.

The Hachicorp's Vault provider expects the environment variable `HCPVAULT_ADDR` to be set with the address of the vault. The credentials for each checker type are expected to be stored on `admin` namespace and under it's own folder under `secret`. This structure should contain the `user` and `pwd` secrets.

Example:
```
.
└── admin (namespace)
    └── secret
        ├── mysql
        │   ├── pwd
        │   └── user
        └── postgres
            ├── pwd
            └── user
```
## Project Structure

```
.
├── LICENSE
├── Set-EnvVars.ps1
├── config.yaml
├── coverage.out
├── go.mod
├── go.sum
├── main.go
├── pkg
│   ├── checker
│   │   ├── checker.go
│   │   ├── httpchecker.go
│   │   ├── httpchecker_test.go
│   │   ├── mysqlchecker.go
│   │   ├── mysqlchecker_test.go
│   │   ├── pgchecker.go
│   │   └── pgchecker_test.go
│   ├── credentialprovider
│   │   ├── azurekeyvault.go
│   │   ├── credentialprovider.go
│   │   ├── hcpvault.go
│   │   └── mock.go
│   ├── database
│   │   ├── connection.go
│   │   └── sql.go
│   ├── k8s
│   │   └── k8s.go
│   └── server
│       └── server.go
├── template.gotmpl
└── test-deployments
    ├── mysql-deployment.yaml
    └── postgres-deployment.yaml
```

## Overview

- **Checkers**: The core functionality is provided by different "checkers". Each checker is responsible for verifying the availability of a particular service/resource (e.g., HTTP, MySQL, PostgreSQL).
  
- **Server**: Hosts an interface to view the status of all checkers, providing real-time feedback on each service's availability and the ability to trigger corrective actions for specific services.

- **Credential Providers**: To securely connect and verify services, credential providers such as Azure Key Vault and HashiCorp Vault are utilized.

- **Database**: Contains files related to managing database connections and executing necessary SQL statements.

- **Kubernetes**: Contains the `k8s.go` file, which manages interactions with Kubernetes deployments and resources. Used to to validate the `Fix` functionality for MySQL and PostgreSQL checkers.

- **Test Deployments**: Provides YAML files for deploying services like MySQL and PostgreSQL in Kubernetes environments. These are used to deploy and validate SQL Checkers.

## Example usage
### Adding new checks
After implementing the checker type, you can add it to the config.yaml file and it will be automatically added to the list of checks.\
Also, if it's a checker that needs credentials, make sure to add the credentials to the corresponding credential provider.

example:
```yaml
checkers:
  - type: http
    url: https://google.com
  - type: http
    url: https://microsoft.com
  - type: postgres
    server: mypostgres.net
    port: 5432
  - type: mysql
    server: mysql.net
    port: 3306
```

### Web interface
A web-based interface provides users with a clear overview of the status of each service/resource. Each entry in the table corresponds to a checker, and its current status is color-coded for clarity (green for available, red for unavailable). If a service/resource is unavailable and fixable, a "Fix" button is available to attempt corrective action.
![checks](https://github.com/rdalbuquerque/availability-checker/blob/master/.attachments/image.png)

## Core Concepts

1. **Continuous Checking**: The system periodically verifies the availability of each service/resource.
  
2. **Concurrency**: Utilizes Go's concurrency features (goroutines) to check multiple services/resources simultaneously.
   
3. **Self-healing**: The ability to fix known issues with services automatically or with minimal user intervention.

## Future Considerations

Given the extensible nature of the project, more checkers can be added for different services and resources as needed. Integration with other monitoring tools, notification systems, or even automated scaling solutions could be potential next steps.
