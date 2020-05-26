
# Add Cloud Services to your Development Environment with Okteto and Crossplane

This is the companion code to the talk I gave as part of Crossplane's Community Day 2020.

# Setup

## Deploy your K8S Cluster

Crossplane and Okteto work with any Kubernetes cluster, local or remote. 

To complete this guide, you'll need cluster-admin permissions to a Kubernetes cluster.

## Install Crossplane v0.11

Follow [this instructions](https://crossplane.io/docs/v0.11/getting-started/install-configure.html#install-crossplane) to install your Crossplane instance.

> The manifests included in this guide assume that you installed crossplane in the `crossplane-system` namespace.

```
kubectl create namespace crossplane-system

helm repo add crossplane-master https://charts.crossplane.io/master/

helm search repo crossplane-master

helm install crossplane --namespace crossplane-system crossplane-master/crossplane --version 0.11
```

## Create your AWS Provider

In this guide, we'll be creating resources in AWS. An AWS user with Administrative privileges is needed to enable Crossplane to create the required resources. Once the user is provisioned, an Access Key needs to be created so the user can have API access.

```
./setup/provider.sh "AWS_KEY" "AWS_SECRET_KEY"
```

## Publish your Infrastructure Offering

Instead of having every developer pick engine, parameters, DB size, version, etc... we are going to create an offering for them. This way, all the developers need to know is the name of it, and the entire infrastructure they need will be created.

First create the definitions (this is the CRD that your infrastructure consumers will use):

```
kubectl apply -f setup/postgresqlinstances.yaml
```

And then the composition:
```
kubectl apply -f setup/composition.yaml
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

## Install the Okteto CLI

We're going to use the Okteto CLI to create and use our development environment while building the application. [Follow this instructions](https://okteto.com/docs/getting-started/installation/) to install it in your local computer.

# Developin' time

## Create Your Development Namespace

Create a namespace for your application and your development environment.

```
kubectl create namespace devtime
```

## Deploy the Application

Deploy the latest version of the application:
```
kubectl apply -f k8s.yaml --namespace=devtime
```

Deploying the application will create the following resources:
- AWS:
    - Postgres instance with a DB called `guestbook`
- Kubernetes:
  - Ingress
  - Service
  - Deployment

It'll take 5-6 minutes to deploy the DB. You can monitor the state with the following command:

```
kubectl describe postgresqlinstancerequirement guestbook --namespace=devtime
```

Once the DB is active, the rest of the application will finish its deployment. You can check the state by running the following command:

```
kubectl get pod -n=devtime
```

Open your browser, and access the application.

###  Via the ingress 
You'll need to create a hosts entry for guestbook.local that points to the external IP of your ingress. 

```
/usr/local/bin/kubectl get ing -n=devtime
```

###  Via port forwarding

You can also start a port-forward to expose port 8080, if running in a local cluster

```
kubectl port-forward svc/guestbook 8080:8080
```

## Deploy your Development Environment

Deploy your development environment with the `okteto up` command:

```
okteto up --namespace=devtime
```

This command will replace the deployment of the application with a fully configured development environment that includes:
* The Go 1.14 runtime
* `fresh`, a [hot-reloader](https://github.com/gravityblast/fresh) for go 
* `delve`, a [go debugger](https://github.com/go-delve/delve)
* PSQL, a CLI to access postgres
* A file synchronization service

## Access the DB

You can query the DB directly from your remote development environment. And since your development environment inherits all the settings from your application, you don't need to setup anything. Just call `psql` from your remote development environment

```
psql
```

## Integrated with Kubernetes

Go to your app, and click on the 'env' link. Notice how the value of `HOSTNAME` maps the name of your pod. This is because your development environment is running directly in Kubernetes.

## Live Coding 

The okteto CLI will automatically synchronize your repository between your local and remote development environments. If you edit a file in your local IDE, the change will be instantly reflected in your remote development environment (and viceversa). This lets us take advantage of existing tooling, like a hot reloader.

Start the application using the preinstalled hot reloader:

```
fresh
```

Go back to your browser and reload the page. You'll see the same application as before, but now connected to your remote development environment.

The application is not querying any data from the DB. Let's change chat. Open `main.go` in your favorite IDE, and uncomment the DB code (lines 66-70 and 91-95).

As soon as you save the file, `okteto` will sync the file, `fresh` will notice it, and it will recompile and reload the go app. 

Go back to your browser and use the application. Notice how the messages are now being stored and retrieved from the DB instance that we provisioned with Crossplane!

## Debugging

One of the biggest advantages of having a remote development environment is that you can use all your tools. 

Start the debugger in your remote development environment:

```
dlv debug --headless --listen=:2345 --log --api-version=2
```

And connect to it from your IDE  (this repo already includes the require configuration for VSCode).







