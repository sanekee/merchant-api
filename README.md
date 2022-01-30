# merchant-api
A Merchant API Demo

# Using this demo

## Required softwares
1. [Docker / Docker Desktop](https://www.docker.com/)
1. [OpenAPI Codegen](https://github.com/deepmap/oapi-codegen)

## Generate API Schema
Install oapi-codegen
```
  go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
```
Generate Schema
```
  oapi-codegen -generate types -o backend/internal/model/schema.go -package model spec/openapi.yaml
```

## This demo contains 2 containers
1. Postgres image with default DB Init (db/db_init.sql) script, listen to tcp port 5432
2. Golang Backend code, listen to tcp port 8123 

## Building
1. Edit backend.env to update the environment variables

   | ENV | Description |
   | --- | ----------- |
   | `POSTGRES_DB` | postgresql database |
   | `POSTGRES_USER` | postgresql database user |
   | `POSTGRES_PASSWORD` | postgresql database password |

2. Run docker compose build
   ```
      docker compose build
   ```

3. Run docker compose
   ```
      docker compose up
   ```

## Accessing API Doc
1. After docker containers are running.
2. SwaggerUI - Browse to [http://127.0.0.1:8123/docs](http://127.0.0.1:8123/docs) to access SwaggerUI documentations.
2. RapiDoc - Browse to [http://127.0.0.1:8123/docs/rapi](http://127.0.0.1:8123/docs/rapi) to access RapiDoc documentations.
3. ReDoc - Browse to [http://127.0.0.1:8123/docs/redoc](http://127.0.0.1:8123/docs/redoc) to access Redoc documentations.
4. To test the API, authorize with the following token
   ```
   eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMTExMTExMSIsIm5hbWUiOiJVc2VyIEEiLCJpYXQiOjE2NDA0MTY3MTN9.SbdB7XjwUDk2iNKegVPG7OEvodf5btXP1UjVCXXWHo0
   ```