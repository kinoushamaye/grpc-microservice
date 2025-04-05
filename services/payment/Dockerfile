# Étape 1 : Builder l'application avec l'image Go
FROM golang:1.23-alpine AS builder
WORKDIR /usr/src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o payment ./cmd/main.go

# Étape 2 : Créer une image de production minimaliste
FROM scratch
COPY --from=builder /usr/src/app/payment ./payment
CMD ["./payment"]
