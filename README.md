
# Setup

## Deploy your K8S Cluster

Crossplane and Okteto work with any Kubernetes cluster, local or remote. 

To complete this guide, you'll need cluster-admin permissions to a Kubernetes cluster.

## Install Crossplane v0.11

Follow [this instructions](https://crossplane.io/docs/v0.11/getting-started/install-configure.html#install-crossplane) to install your Crossplane instance.

> The manifests included in this guide assume that you installed crossplane in the `crossplane-system` namespace.

## Create your AWS Provider

In this guide, we'll be creating resources in AWS. An AWS user with Administrative privileges is needed to enable Crossplane to create the required resources. Once the user is provisioned, an Access Key needs to be created so the user can have API access.

```
./setup/provider.sh "AWS_KEY" "AWS_SECRET_KEY" "us-west-2"
```

## Publish your Infrastructure Offering

Instead of having every developer pick engine, parameters, DB size, version, etc... we are going to create an offering for them. This way, all the developers need to know is the name of it, and the entire infrastructure they need will be created.

First create the definitions (this is what your infrastructure consumers will use):

```
kubectl create -f setup/postgresqlinstances.yaml
```

And then the composition:
```
kubectl create -f setup/composition.yaml
```

The composition lists all the resources that will be automatically created when the user requests a PostgreSQL instance:
- VPC
- 2 Subnets
- 1 Internet Gateway
- 1 Security Group
- 1 DB Security Group
- 1 RDS instance

Give crossplane permissions over the newly created `database.example.org` api group.

```
kubectl create -f setup/rbac.yaml
```


