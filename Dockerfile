FROM golang:1.17 as build

WORKDIR /build
COPY backend .
RUN go build -a -ldflags "-linkmode external -extldflags -static" -o /build/merchant-api cmd/merchant-api/main.go

FROM alpine
WORKDIR /app
COPY spec /app/spec
COPY --from=build /build/merchant-api /app

ENTRYPOINT [ "/app/merchant-api" ]
