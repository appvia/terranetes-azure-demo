#!bin/bash
helm repo add appvia https://terraform-controller.appvia.io
helm repo update

helm install -n terraform-system terraform-controller appvia/terraform-controller --create-namespace
kubectl -n terraform-system get pods

kubectl -n terraform-system create secret generic azure \
  --from-literal=ARM_CLIENT_ID=$ARM_CLIENT_ID \
  --from-literal=ARM_CLIENT_SECRET=$ARM_CLIENT_SECRET \
  --from-literal=ARM_SUBSCRIPTION_ID=$ARM_SUBSCRIPTION_ID \
  --from-literal=ARM_TENANT_ID=$ARM_TENANT_ID