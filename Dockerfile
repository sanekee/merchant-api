FROM golang:1.20 as build

WORKDIR /build
COPY backend .
RUN go build -a -ldflags "-linkmode external -extldflags -static" -o /build/merchant-api cmd/merchant-api/main.go

FROM alpine
WORKDIR /app
COPY ./doc /app/doc
COPY --from=build /build/merchant-api /app

ENTRYPOINT [ "/app/merchant-api" ]
