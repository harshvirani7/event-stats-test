#!/bin/bash

echo "Deploying Elasticsearch..."
kubectl apply -f elasticsearch-deployment.yaml

echo "Deploying Kibana..."
kubectl apply -f kibana-deployment.yaml

echo "Deploying EventStats and other services..."
kubectl apply -f deployment.yaml

kubectl apply -f service.yaml

echo "Setup complete"

# kubectl port-forward service/elasticsearch 9200:9200

# kubectl port-forward service/event-stats-test 8080:8080 
