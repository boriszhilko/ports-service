## Ports Service

The Ports Service is a microservice that follows the DDD architecture. <br>Upon launch, it populates a Redis database with information about seaports, using a JSON file as the default input source. 
<p>
The service runs within a Docker container.
<p/>
Default database is a Redis database.
<p/>
Extension to other input sources and databases is described in the [Extending the Ports Service](#extending-the-ports-service) section.

## Installation

To run this microservice, you need to have Docker and Docker Compose installed on your machine.

1. Clone this repository and navigate to the root directory.

2. Place a `ports.json` file into the `assets` folder. This file should contain a list of seaports in JSON format.

3. Rename the `.env.example` file to `.env` and set the environment variables.

4. Run the following command to start the service:

```bash
make run
```

This will start the service in a Docker container and populate the repository with ports.

## Usage

You can use the following commands to interact with the Ports Service:

### Running Tests

Use the following command to run tests:

```bash
make test
```

Tests are run against a test database. The database is created and dropped automatically.

The test suite includes:

- Unit tests for the service layer.
- Integration tests for the input adapter.
- Integration tests for the output repository adapter.

### Linting

Use the following command to lint the code:

```bash
make lint
```

## Extending the Ports Service

The application is designed to be easily extensible by adding different input sources and repositories.

### Adding a new input source

To add a new input source, create a new input source in the `adapter/in` folder. 
Input source should implement `PortInput` interface.

### Adding a new repository

To add a new repository, create a new repository in the `adapter/out` folder.
Repository should implement `PortRepository` interface.