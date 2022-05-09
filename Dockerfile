FROM golang:1.18.1 AS build

WORKDIR /src
COPY ./go.mod ./go.sum ./
COPY ./ ./

RUN CGO_ENABLED=0 go build -o /vault-pki-cli ./cmd

FROM gcr.io/distroless/static AS final

LABEL maintainer="soerenschneider"
USER nonroot:nonroot

COPY --from=build --chown=nonroot:nonroot /vault-pki-cli /vault-pki-cli

ENTRYPOINT ["/vault-pki-cli"]
