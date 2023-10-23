#!/bin/bash

# Define usage function to display how to use the script
usage() {
  echo "Usage: $0 -a <architecture> [-l]"
  echo "Options:"
  echo "  -a <architecture>   Choose the an architecture (number) (Required)"
  echo "  -l                  List of available architectures"
  exit 1
}

list_archs() {
    echo "1. server"
    echo "2. relay-server"
    echo "3. relay-server-relay"
    echo "4. relay-server-relay-relay"
    echo "5. relay-server-relay (Multi-cluster)"
}

apply_modified_relay() {
    local service_name="$1"
    local deployment_name="$2"
    local port="$3"
    local multi="${4:-cluster}"

    # Replace placeholders with actual values
    sed -e "s/SERVICE_NAME/$service_name/g" \
        -e "s/DEPLOYMENT_NAME/$deployment_name/g" \
        -e "s/PORT_NUMBER/$port/g" \
        -e "s/cluster/$multi/g" \
        relay.yaml > modified-relay.yaml

    # Apply the modified YAML using kubectl (assuming you have kubectl configured)
    kubectl apply -f modified-relay.yaml
    wait 5
    kubectl wait --for=condition=ready pod -l app="$deployment_name" --timeout=2m

    # Clean up the temporary modified YAML file
    rm modified-relay.yaml
}

basic() {
    kubectl delete all --all
    # This is the cert and the key certificate
    kubectl apply -f secret.yaml
    kubectl apply -f server.yaml
    server_ip=$(wait_for_lb "quicrq-server-lb")
    echo "$server_ip"
}

wait_for_lb() {
    local lb="$1"
    local external_ip=""
    while [ -z "$external_ip" ]; do
        external_ip=$(kubectl get svc "$lb" -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
        # external_ip=$(kubectl get svc $lb --template="{{range .status.loadBalancer.ingress}}{{.ip}}{{end}}")
        if [ -z "$external_ip" ]; then
            sleep 2
        fi
    done
    echo "$external_ip"
}

arch=""
list_flag=0

# Use getopts to parse command-line options and arguments
while getopts "a:l:" opt; do
  case "$opt" in
    a) arch="$OPTARG" ;;
    l) list_flag=1 ;;  # Set the flag to true when -l is provided
    \?) usage ;;
  esac
done

if [ "$list_flag" = 1 ]; then
    list_archs
fi

# Check if required options are provided
if [ -z "$arch" ]; then
  usage
fi

server_ip=""

if [ "$arch" = 1 ]; then
    echo "1. server"
    server_ip=$(basic)
elif [ "$arch" = 2 ]; then
    echo "2. relay-server"
    server_ip=$(basic)
    apply_modified_relay "quicrq-relay-lb" "quicrq-relay" 30900
    relay_ip=$(wait_for_lb "quicrq-relay-lb")
    echo "Relay IP $relay_ip"
elif [ "$arch" = 3 ]; then
    echo "3. relay-server-relay"
    server_ip=&(basic)
    apply_modified_relay "quicrq-relay-lb-1" "quicrq-relay-1" 30901
    apply_modified_relay "quicrq-relay-lb-2" "quicrq-relay-2" 30902
    relay1_ip=$(wait_for_lb "quicrq-relay-lb-1")
    echo "Relay IP $relay1_ip"
    relay2_ip=$(wait_for_lb "quicrq-relay-lb-2")
    echo "Relay IP $relay2_ip"
elif [ "$arch" = 4 ]; then
    echo "4. relay-server-relay-relay"
    server_ip=$(basic)
    apply_modified_relay "quicrq-relay-lb" "quicrq-relay" 30900
    apply_modified_relay "quicrq-relay-lb-1" "quicrq-relay-1" 30901
    apply_modified_relay "quicrq-relay-lb-2" "quicrq-relay-2" 30902
    relay_ip=$(wait_for_lb "quicrq-relay-lb")
    echo "Relay IP $relay_ip"
    relay1_ip=$(wait_for_lb "quicrq-relay-lb-1")
    echo "Relay IP $relay1_ip"
    relay2_ip=$(wait_for_lb "quicrq-relay-lb-2")
    echo "Relay IP $relay2_ip"
elif [ "$arch" = 5 ]; then
    echo "5. relay-server-realy (Multi-cluster)"
    kubectl config use-context quicrq-eu
    server_ip=$(basic)
    kubectl apply -f server-export.yaml
    apply_modified_relay "quicrq-relay-lb-eu" "quicrq-relay-eu" 30900
    kubectl config use-context quicrq-us
    kubectl apply -f secret.yaml
    apply_modified_relay "quicrq-relay-lb-us" "quicrq-relay-us" 30900 "clusterset"
else
    echo "Unknown architecture"
    list_archs
fi
