FROM golang:1.19.3 AS builder

WORKDIR /vault-secrets-operator
COPY . .
RUN go mod download
RUN GOOS=linux GOARCH=amd64 GCO_ENABLED=0 go build -o operator.out main/main.go
################
FROM scratch

WORKDIR /vault-secrets-operator
COPY --from=builder /vault-secrets-operator/operator.out .
CMD ["./operator"]