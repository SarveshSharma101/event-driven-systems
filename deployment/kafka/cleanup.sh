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

# === ASCII Logging Functions ===
banner() {
  echo -e "${RED}"
  echo "======================================"
  echo -e "🔥 ${BOLD}$1${RESET}${RED}"
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

# === Delete Kafka Cluster ===
banner "Deleting Kafka Cluster"
if kubectl get ns "$NAMESPACE" >/dev/null 2>&1; then
  if kubectl get -f "$YAML_DIR/kafka-single-node.yaml" -n "$NAMESPACE" >/dev/null 2>&1; then
    step "Deleting RabbitMQ cluster resources"
    kubectl delete -f "$YAML_DIR/kafka-single-node.yaml" -n "$NAMESPACE" || warn "Cluster resources already gone"
    sleep_progress 5 "Waiting after deleting cluster..."
  else
    warn "Cluster manifest not found/applied, skipping"
  fi
else
  warn "Namespace $NAMESPACE not found, skipping cluster deletion"
fi



# === Delete Cluster Operator ===
banner "Deleting Strimzi Operator"
if kubectl get -f "$YAML_DIR/strimzi.yaml" -n $NAMESPACE>/dev/null 2>&1; then
  kubectl delete -f "$YAML_DIR/strimzi.yaml" -n $NAMESPACE|| warn "Cluster Operator already removed"
  success "Cluster Operator deleted"
else
  warn "Cluster Operator manifest not found/applied, skipping"
fi


# === Delete Namespace ===
banner "Deleting Namespace: $NAMESPACE"
if kubectl get ns "$NAMESPACE" >/dev/null 2>&1; then
  kubectl delete ns "$NAMESPACE"
  success "Namespace $NAMESPACE deleted"
else
  warn "Namespace $NAMESPACE already deleted"
fi
