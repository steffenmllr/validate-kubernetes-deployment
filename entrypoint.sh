#!/bin/sh

# Extract the base64 encoded config data and write this to the KUBECONFIG if it is not defined
# Used with github actions: https://github.com/actions/aws/tree/master/kubectl#secrets
if [ -n "$KUBECONFIG" ]; then
    echo "Using KUBECONFIG"
else
    echo "$KUBE_CONFIG_DATA" | base64 -d > /tmp/config
    export KUBECONFIG=/tmp/config
    echo "Using KUBE_CONFIG_DATA"
fi

/validate
