FROM golang:1.20 AS build

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o /webhook

## Deploy
FROM gcr.io/distroless/base-debian10
WORKDIR /
COPY --from=build /webhook /webhook
EXPOSE 8080
USER nonroot:nonroot

ENTRYPOINT ["/webhook"]
