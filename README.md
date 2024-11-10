# Application

This application is an example on how to setup an application for handling some basic media and tags operations.

For a full view of the provided API endpoints take a look at the opanapi specification in the `api/` folder
(for simplicity you can copy such configuration into the [Swagger-editor](https://editor.swagger.io/) )

# Architecture

TODO: Architecture description

# Setup

To run and test the application you will need the following tools

- [Make >= 4.3](https://www.gnu.org/software/make/)
- [Go >= 1.23](https://go.dev/doc/install)
- [Docker >= 27.3.1 and Docker compose >= v2.29.7](https://docs.docker.com/compose/install/)

Then simply run:

```
make run
```



### Possible improvements / Shortcuts taken
Minio swapped with actual S3 bucket and reading/writing to S3 in production

Pagination on both tags listing and media filtering. Bot since it was not required in this task it was considered out of scope

Tag id == uuid, its easier to pick a form of id over all entities and stick with it (if we do not need some additional temporal connotation to it, like with increasing sequences)