# Appvia Terranetes for Azure Demo

This repository contains a simple sample project and yaml to help you understand and get started using Appvia's Terranetes controller to provision Azure Cloud Resources.

## Prerequisites for the Terraform Controller
* You must have a Kubernetes cluster.
* Helm installed.
* If you don't have an Azure subscription, create a [free account](https://portal.azure.com/#create/Microsoft.StorageAccount-ARM) before you begin.

Start by adding the Appvia terraform controller repo to helm:
```bash
helm repo add appvia https://terraform-controller.appvia.io
helm repo update
```

Next Deploy the controller:
```bash
helm install -n terraform-system terraform-controller appvia/terraform-controller --create-namespace
```

Check that the services and pods were created:
```bash
kubectl -n terraform-system get pods
```

Configure and apply the secret that the controller will use to communicate with your Azure account:
```bash
kubectl -n terraform-system create secret generic azure \
  --from-literal=ARM_CLIENT_ID=$ARM_CLIENT_ID \
  --from-literal=ARM_CLIENT_SECRET=$ARM_CLIENT_SECRET \
  --from-literal=ARM_SUBSCRIPTION_ID=$ARM_SUBSCRIPTION_ID \
  --from-literal=ARM_TENANT_ID=$ARM_TENANT_ID
```
You'll replace the $ items with your configuration.

## Deploying the Azure provider
You'll deploy the Azure provider into the same namespace that the controller is in. In this case, terraform-system:
```bash
kubectl apply -f ./yaml/azure-provider.yaml -n terraform-system
```

## Creating Cloud Resources
In this demo, we're creating a storage account using the Terraform controller.
* You will need to edit the 'variables:' section of ./yaml/azure-deployment.yaml. Ensure the following sections match your configuration:

```yaml
  variables:
    resource_group_name: legendary-robot 
    name: appviastor3
    location: eastus
    replication_type: LRS
    shared_access_key_enabled: true
    tags: 
      name: test
```
Once configured, apply the deployment to a namespace of your choosing (in this case, I'm using a demo namespace.)

```bash
kubectl apply -f ./yaml/azure-deployment.yaml -n demo
```

Now you can watch the progress of the apply, checking for the completion of the pods in the namespace you just deployed to:
```bash
watch kubectl get pods -n demo
```

Once you see the 'apply' pod completed, you can check it's logs to see the outcome:
```bash
kubectl logs $pod_name_here -n demo
```

## Using the cloud resource
The Go application in this repo was written to check for enviornment variables that are exposed to the pod as secrets from the output of the Terranetes deployment. In this example, you can simply export the two keys locally to your shell and go run the application. To use in a kubernetes cluster, ensure that these secrets are created and available to the pod on deployment. An example deployment of this application is included at ./yaml/container.yaml

## Prerequisites for the Go Application (Local Development)

To work with this application locally:

* Install [Go](https://golang.org/dl/) 1.8 or later

## Download and Install the Azure Storage Blob SDK for Go
This is only required if you wish to edit and run the Go application locally.

From your GOPATH, execute the following command:

```bash
go get github.com/Azure/azure-sdk-for-go/sdk/azidentity

go get github.com/Azure/azure-sdk-for-go/sdk/storage/azblob
```

At this point, you can run this application. It creates an Azure storage container and blob object.

## Run the application

Open the `storage-quickstart.go` file.

Export your storage account name and storage account key as enviornment variables.

Run the application with the `go run` command:

```bash
go run storage-quickstart.go
```
## More information
For more information on Appvia's Terranetes project, please visit [https://terranetes.appvia.io/](https://terranetes.appvia.io/)