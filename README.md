# ChefBook Backend Service Template
Default project structure for ChefBook microservice

### Structure
* `api` - go module for inter service communication
* `assets` - non-code resources required for service
* `cmd` - application entry point
* `deployments` - Kubernetes setup
* `internal` - microservice code
* `migrations` - directory for data storage (usually database) migrations
* `pkg` - common code that can be placed in other projects; reusable packages usually implemented by dependency on [common repository](https://github.com/mephistolie/chefbook-backend-common)
* `scripts` - useful `.sh` scripts
* `sercrets` - directory for sensitive data
