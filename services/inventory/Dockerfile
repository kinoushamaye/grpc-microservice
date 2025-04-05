# Étape 1 : Utiliser une image Go officielle pour la compilation
FROM golang:1.23-alpine AS builder

# Définir le répertoire de travail dans le conteneur
WORKDIR /usr/src/app

# Copier les fichiers source dans le conteneur
COPY . .

# Installer les dépendances et compiler le binaire pour un environnement Linux sans CGO
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o inventory ./cmd/main.go

# Étape 2 : Utiliser une image "scratch" (image vide) pour la phase de production
FROM scratch

# Copier le binaire compilé à partir de l'étape builder
COPY --from=builder /usr/src/app/inventory ./inventory

# Exposer le port sur lequel l'application Go écoute (par exemple, 8080)
EXPOSE 8080

# Démarrer l'application Go en utilisant le binaire copié
CMD ["./inventory"]
