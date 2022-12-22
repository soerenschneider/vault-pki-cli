FROM golang:1.19.3 AS build

WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./
ENV CGO_ENABLED=0
RUN go mod download
RUN go build -o /vault-pki-cli ./cmd

<<<<<<< HEAD
=======
RUN CGO_ENABLED=0 go build -o /vault-pki-cli ./cmd

>>>>>>> 02eaae1c2b3d5d25cbeb61cd3f6774d49c0ee207
#FROM gcr.io/distroless/static AS final
FROM debian

LABEL maintainer="soerenschneider"
#USER nonroot:nonroot

COPY --from=build --chown=nonroot:nonroot /vault-pki-cli /vault-pki-cli

ENTRYPOINT ["/vault-pki-cli"]
