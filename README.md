# Application

This is an example on how to setup an application for handling some basic media and tags operations.

For a full view of the provided API endpoints take a look at the opanapi specification in the [api/](api/openapi.yaml) folder
(for simplicity you can copy such configuration into the [swagger-editor](https://editor.swagger.io/) for a more graphical visualization).

# Architecture

The architecture of this application is kept pretty simple, consisting of a golang backend exposing 4 rest API endpoints,
a Postgres DB for storing the media and tag information and a MinIO container mimicking a Cloud Object Storage for the actual media storage.

The backend exposes 2 endpoints (Get/Post) for the media service and 2 endpoints (Get/Post) for the tags service.
For a complete api specification take a look [here](api/openapi.yaml).

The Postgres database schema can be found in the [migrations](migrations/00001_initial.sql) folder and consist of 3 tables:
* media: for the media entities information
* tag: for storing tag entities
* media_tag: a table representing a many to many relationship between media and corresponding tags linked to it.

This migration is automatically applied to the locally running postgres container on startup by the db helper functions [here](internal/db/db.go).

For the storage a MinIO based implementation is provided [here](internal/storage/minio_storage.go). This is mostly done as POC for being able to run the application locally, in place of an actual cloud based object storage (AWS S3 or GCP bucket for example would be the production solution).

The application configuration is done via [env variables](internal/config/config.go). For local spin-up the configured defaults are fine enough.

All the functionality except for the storage is always split into two entities, one tied to the tags and the other to the media.

Structured logging is used to provide meaningful information in the logs.

# Setup

To build/run and test the application you will need the following tools

- [Make >= 4.3](https://www.gnu.org/software/make/)
- [Go >= 1.23](https://go.dev/doc/install)
- [Docker >= 27.3.1 and Docker compose >= v2.29.7](https://docs.docker.com/compose/install/)
- [Curl >= 8.5.0 ](https://curl.se/)

Then simply run:

```
make run
```
This will spin up the Postgres and MinIO containers, build the go application and run the built binary.

Once you see the server log `INF Starting server at :8080` the application is ready to go.

The easiest way to test the application form a user perspective is to use the provided API specs in the swagger editor and generate some curl
snippets to send requests to the server.

To create a new tag you can run the following curl snippet: 
```
curl -X 'POST' \
  'http://localhost:8080/tags' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "Wembley Stadium"
}'
```

This will return you a JSON response similar to this: `{"id":"7c44027e-fc3f-46df-9d72-13ec022e668a","name":"Wembley Stadium"}`

To list the created tags you can use: 
```
curl -X 'GET' \
  'http://localhost:8080/tags' \
  -H 'accept: application/json'
```
this will return you a list of the created tags.

Now on to the more complex ones, to create a new media you need a picture first (either png or jpeg, enforced on the backend).
Some samples can be found in the [samples folder](samples/) lets assume the tennis image should be uploaded:
```
cd samples && \
curl -X 'POST' \
  'http://localhost:8080/media' \
  -H 'accept: application/json' \
  -H 'Content-Type: multipart/form-data' \
  -F 'name=Tennis player' \
  -F 'tags=7c44027e-fc3f-46df-9d72-13ec022e668a' \
  -F 'file=@tennis.png;type=image/png'
```
Keep in mind that the provided tag ids will be validated, so they need to be created first, otherwise the request will fail. 
The response will look something like this:

```
{
  "id": "be6b13d2-bef2-4b26-8200-bf379d80473f",
  "name": "Tennis player",
  "tags": [
    "Wembley Stadium"
  ],
  "fileUrl": "http://localhost:9000/example/66824288-6867-4898-a655-845a91b8a341.png"
}
```
Since the MinIO [bucket](http://localhost:9001/) is configured to be public the fileUrl can be directly used to verify that the correct picture was uploaded to the bucket.

Last but not least, filtering media based on tag id:
```
curl -X 'GET' \
  'http://localhost:8080/media?tag=7c44027e-fc3f-46df-9d72-13ec022e668a' \
  -H 'accept: application/json'
```
This will return a list of all media that is in association with the specified tag.

To run the test two make targets are provide:
```
make test
```
and 

```
make test-coverage 
```
for testing and generating a coverage report. 

`make test` is also used in the GitHub [action](.github/workflows/run_tests.yml) for testing the current state on each commit.

To shut down the running containers the following make target can be used:
```
make stop
```

### Possible improvements / Assumptions / Shortcuts taken
* MinIO should be swapped with actual S3 bucket.
* Metrics should be exposed in order for something like a Prometheus scraper to collect them.
* Pagination on both tags listing and media filtering could be implemented, since it was not required in this task it was considered out of scope.
* For tags the id was implemented using a uuid, seemed unnecessary to maintain 2 different types of id's in such a small scope.
* File size upload was limited to 10MB per request. 
* File MIME type was limited to jpeg and png for simplicity, much more different types could be used but would require more validation.
* When uploading to the storage, file names are dropped in favor of uuid+fileType, so that it does not lead to any name clashes during upload on the bucket.
* File compression was skipped but would definitely be an improvement with more time(could be important for cost saving on large volumes of data)

#### IMPORTANT
* Due to time constraints skipped writing test for the tags_handler.go in favor of the more complex media_handler.go, if tests would have been added they would have been done in a similar fashion to what was done in the media_handler_test.go

* Also similarly due to time constraints skipping service package testing since its basically reduced to calling the repository methods, and for that an integration test was already provided in the repository package.