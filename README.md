# Go Application Template

This is a template for a Go application using layered architecture and dependency injection.

## Structure
The goal of this project structure is to maintain separation between the business logic,
application input/output mechanisms, and communication with external systems. 

### Stores
Stores are adapters to outside systems. They are the base layer of the project architecture.
This is where logic for interacting with databases, caches, and logging systems happens. These
store implementations should not depend on any of the higher layers in the project. Entities
that require store implementations should define the required store interface for that entity.

### Services
Services are where ALL of the business logic lives. Services can depend on stores to interact with
outside systems. Services should be entirely independent from any details about the binding
of the application.

### Bindings
Bindings define how the project can be interacted with by users. The binding bridges the gap between
users and the business logic of the application. This could be a REST API, command line application,
or anything else.