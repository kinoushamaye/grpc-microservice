#!/bin/bash

# Arrêter et supprimer les conteneurs existants (si déjà déployés)
docker stop inventory-service order-service payment-service || true
docker rm inventory-service order-service payment-service || true

# Exécuter les nouveaux conteneurs
docker run -d --name inventory-service -p 8081:8081 $DOCKER_USERNAME/inventory-service:latest
docker run -d --name order-service -p 8082:8082 $DOCKER_USERNAME/order-service:latest
docker run -d --name payment-service -p 8083:8083 $DOCKER_USERNAME/payment-service:latest

echo "Déploiement effectué avec succès."
