#!/bin/bash



kubectl delete -f rabbitmq-cluster.yaml -n rabbitmq
kubectl delete -f cluster-operator.yml

kubectl delete ns rabbitmq


