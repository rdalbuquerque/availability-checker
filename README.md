# Availability Checker
- [Availability Checker](#availability-checker)
  - [Description](#description)
  - [Project Structure](#project-structure)
  - [Overview](#overview)
  - [Web Interface](#web-interface)
  - [Core Concepts](#core-concepts)
  - [Future Considerations](#future-considerations)

## Description
This project aims to explore Go extensability through the usage of interfaces. `Checker` interface establishes methods for new checkers to be added and `CredentialProvider` interface provides a way to add new credential providers with ease, here we implemented both Hachicorp's Vault and Azure Key Vault providers. Interfaces are also used here to enable dependency injection, allowing us to properly mock and test database connections.

Another Go feature this project explore is concurrency. Each check is executed in it's own goroutine reporting back to a channel when it's done, allowing us to add more checks with minimal impact to verification time.

This project also explores Go templating capabilities to render and serve a simple HTML page summarizing the status of each checker.

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

## Web Interface

A web-based interface provides users with a clear overview of the status of each service/resource. Each entry in the table corresponds to a checker, and its current status is color-coded for clarity (green for available, red for unavailable). If a service/resource is unavailable and fixable, a "Fix" button is available to attempt corrective action.

## Core Concepts

1. **Continuous Checking**: The system periodically verifies the availability of each service/resource.
  
2. **Concurrency**: Utilizes Go's concurrency features (goroutines) to check multiple services/resources simultaneously.
   
3. **Self-healing**: The ability to fix known issues with services automatically or with minimal user intervention.

## Future Considerations

Given the extensible nature of the project, more checkers can be added for different services and resources as needed. Integration with other monitoring tools, notification systems, or even automated scaling solutions could be potential next steps.
