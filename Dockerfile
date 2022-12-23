FROM golang:1.19.3 AS build

WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./
ENV CGO_ENABLED=0
RUN go mod download
RUN go build -o /vault-pki-cli ./cmd
RUN CGO_ENABLED=0 go build -o /vault-pki-cli ./cmd


FROM gcr.io/distroless/static AS final

LABEL maintainer="soerenschneider"
USER nonroot:nonroot
COPY --from=build --chown=nonroot:nonroot /vault-pki-cli /vault-pki-cli

ENTRYPOINT ["/vault-pki-cli"]
