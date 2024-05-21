# simple-go-app

## The App

This is a simple golang http server, when ``/items`` endpoint is requested, a response with a message "ok" is returned after 3 seconds (which will be explained later).

## Dockerfile - Build Process

This Dockerfile uses a multi-stage build process to minimize the image size.

Initially, I used a distroless image but encountered some challenges. Therefore, I switched to using a bookworm image instead.

## Monitoring

Inside the container, you will find traces and logs (as shown in the image below) because the app is manually instrumented. 

The 3-second delay is included to allow sufficient time for proper trace collection, in case we want to visualize it in any backend system.

Adding an OpenTelemetry Collector to the application for console extraction is a straightforward process.

![alt text](<images/Screenshot 2024-05-21 191549.png>)

## Repo/Branch structure

There are three branches in this repository:

1. Main branch

1. Helm chart branch

1. Terraform branch

These branches simulate real-world production repositories. 

The main branch represents your source code production repo.

The Helm chart branch handles Kubernetes package management. 

Terraform branch is used to create, destroy, and maintain the state of the infrastructure. 

You can have them in a monorepo or separate repositories to maintain separation of concerns.



## CICD Pipeline

The entire workflow looks like the image below:

![alt text](<images/Screenshot 2023-11-11 232700.png>)

1. It starts by pushing new code or opening a PR to the main branch, which can be set up in the GitHub YAML workflow folder.

1. GitHub Actions will start running the tasks provided in the YAML files and will build your application.

1. Trivy will test your application.

1. GitHub Actions will push the image to ECR.

1. GitHub Actions will change the image tag in the Helm chart in the Helm charts branch.

1. ArgoCD listens to the Helm charts branch and, when the image tag changes, it will pull the new tag from ECR.

Finally, you have a working GitOps setup. A lot can be improved, but it's a good start!

## GHA workflows

Within the workflows directory, you'll find three YAML files.

The primary file is named main.yaml, encompassing tasks like build, test, push, and update.

The CodeQL and Trivy YAMLs execute under two conditions:

1. Scheduled via cron.

1. Triggered by Pull Requests, with CodeQL.yaml specifically executing on PRs to the main branch.

![alt text](<images/Screenshot 2024-05-22 001454.png>)

## Terraform

The Terraform branch comprises two primary directories:

1. One folder contains modules.

1. The other folder without modules.

The discrepancy lies in the fact that the folder with modules establishes a comprehensive EKS cluster, ready for production, including all IAM configurations and ECR.

Alternatively, the Terraform directory without modules configures a minimal EKS cluster suitable for development purposes.

## Helm chart

This helm chart is a default template, but there are a few modifications we need to handle post-creation. These adjustments include setting up the password and configuring the repository. These are one-time setup tasks.


## How to run in the cloud

Here's the step-by-step breakdown for replicating this setup in the cloud:

1. Start by navigating to the working directory using ``cd``, then execute ``terraform init`` followed by ``terraform apply``.

1. Update the kubeconfig with the appropriate region and name settings.

1. Create a namespace named ``argocd`` and another namespace for your application.

1. Generate an AWS secret within the application's namespace created in the previous step.

1. Port forward ArgoCD to avoid utilizing ingress or a load balancer. Note: While this isn't recommended for production applications, it's suitable for testing purposes on-the-fly.

1. Configure the Helm repository within ArgoCD.

## How to run locally 

Before delving into local Kubernetes setup, ensure you have the following software installed:

1. Docker or any of the recommended container engines.
1. Kubectl.
1. Minikube.

Once installed, initiate Minikube and ensure Ingress and Metrics Server are enabled:

    minikube addons enable ingress
    minikube addons enable metrics-server

You can verify the enabled addons using:

    minikube addons list

Next, deploy ArgoCD using the below commands, or the yaml provided in terraform branch.

    kubectl create ns argocd
    kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/v2.9.0/manifests/install.yaml

Then, port forward the ArgoCD server to access the GUI:
   
    kubectl port-forward svc/argocd-server -n argocd 8000:443

Ensure you're in a separate terminal tab for convenience.

Upon success, you should see a confirmation. Retrieve the password using:
    
    kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d; echo

To configure ArgoCD to listen on a Helm repository, you'll need an access token. After configuration, you can start deploying your application.

Upon successful setup, you should be able to access your application via the ingress endpoint.


## Additional Info 

I opted for a single pod in this setup, but in a production environment, the configuration would be vastly different.

To ensure zero downtime during updates, you might consider using Pod Disruption Budgets (PDB). However, I didn't include PDB in the Helm chart to maintain its beginner-friendly nature.

Additionally, I left the service type as LoadBalancer, which is suitable for local environments. However, in production, the setup differs significantly. On-premise environments typically utilize a combination of reverse proxies, Keepalived, and Ingress to expose applications. This approach is often referred to as the "Hard Way," or alternatively, you can leverage the new Gateway API.

If you deploy the Helm chart locally and attempt to retrieve the service's LoadBalancer IP (svc:lb), you'll notice that the external IP remains in the "pending" state. Which is the expected behaviour.

