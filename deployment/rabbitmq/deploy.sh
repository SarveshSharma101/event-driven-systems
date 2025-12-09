#!/bin/bash

kubectl create ns rabbitmq

kubectl apply -f cluster-operator.yml

kubectl apply -f rabbitmq-cluster.yaml -n rabbitmq

