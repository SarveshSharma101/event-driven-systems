#!/bin/bash
set -euo pipefail

# === Config ===
NAMESPACE="kafka"
YAML_DIR="."

# === Colors ===
RED="\033[0;31m"
GREEN="\033[0;32m"
YELLOW="\033[1;33m"
BLUE="\033[0;34m"
CYAN="\033[0;36m"
BOLD="\033[1m"
RESET="\033[0m"

# === ASCII Logging Function ===
banner() {
  echo -e "${CYAN}"
  echo "======================================"
  echo -e "🚀 ${BOLD}$1${RESET}${CYAN}"
  echo "======================================"
  echo -e "${RESET}"
}

step() {
  echo -e "${YELLOW}👉 $1${RESET}"
}

success() {
  echo -e "${GREEN}✅ $1${RESET}"
}

warn() {
  echo -e "${RED}⚠️  $1${RESET}"
}

sleep_progress() {
  secs=$1
  msg=$2
  echo -ne "${BLUE}$msg${RESET} "
  for ((i=secs; i>0; i--)); do
    echo -ne "•"
    sleep 1
  done
  echo -e " ${GREEN}done!${RESET}"
}

# === Ensure Namespace ===
banner "Ensuring namespace: $NAMESPACE"
kubectl create ns "$NAMESPACE" 2>/dev/null && success "Namespace $NAMESPACE created" || success "Namespace $NAMESPACE already exists"

# === Apply Cluster Operator ===
banner "Deploying strimzi Operator"
step "Applying $YAML_DIR/strimzi.yaml"
kubectl apply -f "$YAML_DIR/strimzi.yaml"

sleep_progress 5 "Sleeping before checking pods..."
step "Waiting for operator pods in $NAMESPACE"
kubectl wait --for=condition=ready pod --all -n $NAMESPACE --timeout=120s
success "Cluster Operator is ready"

# === Apply RabbitMQ Cluster ===
banner "Deploying Kafka Cluster"
step "Applying $YAML_DIR/kafka-single-node.yaml"
kubectl apply -f "$YAML_DIR/kafka-single-node.yaml" -n $NAMESPACE

sleep_progress 10 "Sleeping before checking pods..."
step "Waiting for kafka pods in $NAMESPACE"
kubectl wait --for=condition=ready pod --all -n "$NAMESPACE" --timeout=120s
success "Kafka Cluster is ready"
sleep_progress 10 "Sleeping before checking pods..."

kubectl apply -f "$YAML_DIR/kafka-topic.yaml" -n $NAMESPACE

echo
banner "port forwarding rabbitmq svc, access rabbitmq broker at localhost:9092"
kubectl port-forward -n $NAMESPACE svc/kafka-cluster-kafka-bootstrap 9092:9092