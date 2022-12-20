FROM golang:1.19.0 AS build

WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 go build -o /vault-pki-cli ./cmd

#FROM gcr.io/distroless/static AS final
FROM debian

LABEL maintainer="soerenschneider"
#USER nonroot:nonroot

COPY --from=build --chown=nonroot:nonroot /vault-pki-cli /vault-pki-cli

ENTRYPOINT ["/vault-pki-cli"]
