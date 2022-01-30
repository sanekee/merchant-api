# merchant-api
A Merchant API Demo

# Using this demo

## Required softwares
1. [Go 1.17 ](https://go.dev/)
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
1. API Backend code, listen to tcp port 8123 

## Building
1. Edit backend.env to update the environment variables

   | ENV | Description |
   | --- | ----------- |
   | `POSTGRES_DB` | postgresql database |
   | `POSTGRES_USER` | postgresql database user |
   | `POSTGRES_PASSWORD` | postgresql database password |

1. Run docker compose build
   ```
      docker compose build
   ```

1. Run docker compose
   ```
      docker compose up
   ```

## Accessing API Doc
1. After docker containers are running.
   1. SwaggerUI - Browse to [http://127.0.0.1:8123/docs](http://127.0.0.1:8123/docs) to access SwaggerUI documentations.
   1. RapiDoc - Browse to [http://127.0.0.1:8123/docs/rapi](http://127.0.0.1:8123/docs/rapi) to access RapiDoc documentations.
   1. ReDoc - Browse to [http://127.0.0.1:8123/docs/redoc](http://127.0.0.1:8123/docs/redoc) to access Redoc documentations.
   1. To test the API, authorize with the following token
      ```
      eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMTExMTExMSIsIm5hbWUiOiJVc2VyIEEiLCJpYXQiOjE2NDA0MTY3MTN9.SbdB7XjwUDk2iNKegVPG7OEvodf5btXP1UjVCXXWHo0
      ```

## Testing 
1. Sample test cases for handlers:
   1. backend/internal/handler/merchant_test.go
   1. backend/internal/handler/team_member_test.go
2. Start app with environment variable APP_MOCKDB=true to start API server with a simple map backend without Postgresql DB