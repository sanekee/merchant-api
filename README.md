# merchant-api

A Merchant API Demo Using [Limi HTTP Router](https://github.com/sanekee/limi)

## Using this demo

### This demo consists of 2 containers

1. Postgres image with default DB Init (db/db_init.sql) script, listen to tcp port 5432
1. API Backend code, listen to tcp port 8123
1. Environment variables

   | ENV | Type | Description |
   | --- | ---- | ----------- |
   | `APP_PORT` | Integer | App HTTP listening port |
   | `APP_SPEC` | String | Path contains openapi.yaml |
   | `POSTGRES_DB` | String | postgresql database |
   | `POSTGRES_USER` | String | postgresql database user |
   | `POSTGRES_PASSWORD` | String | postgresql database password |
   | `POSTGRES_HOST` | String | postgresql hostname |
   | `POSTGRES_PORT` | Integre | postgresql port |

### Quickstart

1. Run docker compose build

   ```shell
      docker compose build
   ```

1. Run docker compose

   ```shell
      docker compose up
   ```

### Accessing API Doc

1. View online [API Doc](https://redocly.github.io/redoc/?url=https://raw.githubusercontent.com/sanekee/merchant-api/main/spec/openapi.yaml)

   Or

1. Open [http://127.0.0.1:8123/docs](http://127.0.0.1:8123/docs) after starting the backend.

### Test Token

```code
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMTExMTExMSIsIm5hbWUiOiJVc2VyIEEiLCJpYXQiOjE2NDA0MTY3MTN9.SbdB7XjwUDk2iNKegVPG7OEvodf5btXP1UjVCXXWHo0
```

---

## Development

### Required softwares

1. [Go 1.20](https://go.dev/)
1. [Docker / Docker Desktop](https://www.docker.com/)
1. [OpenAPI Codegen](https://github.com/deepmap/oapi-codegen) - optional for generating schemas

### Generate API Schemas

To generate models schema after updated the API specification.

1. With installed `oapi-codegen`

    1. Install oapi-codegen

    ```shell
      go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
    ```

    2. Generate Schema

    ```shell
      oapi-codegen -generate types -o backend/internal/model/schema.go -package model ./doc/spec/openapi.yaml
    ```

1. With `oapi-codegen` docker image

```shell
    docker run -v $(pwd):/src --rm --name oapi-codegen aaguerocashin/oapi-codegen oapi-codegen \
        -generate types \
        -o /src/backend/internal/model/schema.go \
        -package model \
        /src/doc/spec/openapi.yaml
```
