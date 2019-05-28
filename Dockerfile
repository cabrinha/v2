FROM golang:1.12 as build

WORKDIR /app
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

FROM gcr.io/distroless/base
COPY --from=build /app /
COPY --from=build /app/config.yaml /config.yaml
ENTRYPOINT [ "/app" ]
