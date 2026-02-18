# DevOps App

Simple Express.js API containerized and deployed via Kubernetes (Minikube + Helm) and Ansible (Docker + PM2).

## Prerequisites

- Node.js >= 18
- Docker
- Minikube (for Kubernetes deployment)
- Helm (for Kubernetes deployment)
- Ansible (for Ansible deployment)

## Project Structure

```
devops/
  app.js                  # Express API (port 3000)
  Dockerfile              # Multi-stage Node.js 18 Alpine build
  ecosystem.config.js     # PM2 cluster config (2 instances)
  k8s/                    # Raw Kubernetes manifests
    deployment.yaml
    service.yaml
    ingress.yaml
    hpa.yaml
  helm/devops-app/        # Helm chart
  ansible/                # Ansible playbooks + Docker target
    setup-runtime.yml     # Install Node.js runtime
    deploy-app-local.yml  # Deploy app with PM2
    docker-compose.yml    # Target server container
```

## API Endpoints

| Endpoint  | Description                              |
|-----------|------------------------------------------|
| `GET /`   | App status, timestamp, hostname          |
| `GET /health` | Health check with uptime and memory info |

## Quick Start (Local)

```bash
pnpm install
pnpm start
```

## Deployment Options

### Kubernetes (Minikube + Helm)

```bash
make minikube-start     # Start Minikube cluster
make docker-build       # Build Docker image
make deploy-helm        # Deploy with Helm
make status             # Check pods/services
make url                # Get service URL
make validate-k8s       # Validate full deployment
```

Full deployment with Ingress + HPA:

```bash
make deploy-helm-full
```

Teardown:

```bash
make undeploy-helm      # Remove Helm release
make stop-k8s           # Remove all K8s resources
make minikube-delete    # Delete Minikube cluster
```

### Ansible (Docker + PM2)

```bash
make ansible-up         # Start target server container
make ansible-setup      # Install Node.js runtime
make ansible-deploy     # Deploy app with PM2
make ansible-status     # Check PM2 status
make ansible-test       # Test endpoints
make validate-ansible   # Validate full deployment
```

Teardown:

```bash
make ansible-down       # Stop target container
```

## Testing

### Test Endpoints Manually

```bash
curl http://localhost:3000/ | python3 -m json.tool
curl http://localhost:3000/health | python3 -m json.tool
```

### Validate Deployments

```bash
make validate-all       # Validate both Ansible and Kubernetes
make validate-ansible   # Validate Ansible only
make validate-k8s       # Validate Kubernetes only
```

Validation checks: container/pod status, process health, Node.js version, health endpoint response, and error logs.

### Full Workflow Test

```bash
make all                # minikube-start → docker-build → deploy-helm → status
```

## All Make Commands

Run `make help` for the full list, or:

| Command | Description |
|---------|-------------|
| `make minikube-start` | Start Minikube with 2 CPUs, 4GB RAM |
| `make minikube-stop` | Stop Minikube |
| `make minikube-delete` | Delete Minikube cluster |
| `make docker-build` | Build image inside Minikube |
| `make deploy` | Deploy with raw K8s manifests |
| `make deploy-helm` | Deploy with Helm |
| `make deploy-helm-full` | Deploy with Ingress + HPA |
| `make undeploy` | Remove raw K8s resources |
| `make undeploy-helm` | Remove Helm release |
| `make logs` | Tail K8s pod logs |
| `make status` | Show pods, services, deployment |
| `make url` | Get Minikube service URL |
| `make ansible-up` | Start Ansible target container |
| `make ansible-down` | Stop Ansible target container |
| `make ansible-setup` | Run setup-runtime playbook |
| `make ansible-deploy` | Run deploy-app playbook |
| `make ansible-status` | Check PM2 status |
| `make ansible-logs` | View PM2 logs |
| `make ansible-test` | Test app endpoints |
| `make validate-k8s` | Validate Kubernetes deployment |
| `make validate-ansible` | Validate Ansible deployment |
| `make validate-all` | Validate both deployments |
| `make stop-all` | Stop everything |
| `make clean` | Undeploy + delete Minikube |

## Tech Stack

- **Runtime:** Node.js 18, Express
- **Containerization:** Docker (multi-stage Alpine)
- **Process Manager:** PM2 (cluster mode, 2 instances)
- **Orchestration:** Kubernetes (Minikube), Helm
- **Provisioning:** Ansible
