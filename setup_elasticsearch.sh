#!/bin/bash

# Set up the CRDs required by ECK
echo "Creating CRDs for ECK..."
kubectl create -f https://download.elastic.co/downloads/eck/2.12.1/crds.yaml

# Deploy the ECK operator
echo "Deploying the ECK operator..."
kubectl apply -f https://download.elastic.co/downloads/eck/2.12.1/operator.yaml

# Wait for the ECK operator to be ready
echo "Waiting for the ECK operator to be ready..."
kubectl wait --for=condition=available --timeout=600s deployment/elastic-operator -n elastic-system

# Deploy Elasticsearch cluster
echo "Deploying Elasticsearch cluster..."
kubectl apply -f deployment.yaml

echo "Elasticsearch setup completed."
