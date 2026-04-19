# k8bot

A Slack bot for querying Kubernetes clusters in real time. Built in Go, it connects via Slack Socket Mode and provides both text-based ChatOps commands and an interactive Block Kit UI for browsing cluster resources and metrics.

## Features

- **ChatOps commands** — mention the bot in any channel to query pods, deployments, services, nodes, namespaces, events, and resource metrics.
- **Interactive UI** — a dropdown-based command center posted via Slack Block Kit, with dynamic namespace discovery.
- **Metrics support** — CPU and memory usage for nodes and pods via the Kubernetes Metrics API.
- **In-cluster & local** — auto-detects in-cluster config or falls back to `~/.kube/config` for local development.

## Requirements

- Go 1.25+
- A Kubernetes cluster with [metrics-server](https://github.com/kubernetes-sigs/metrics-server) installed (for `top` commands)
- A Slack app configured with **Socket Mode** enabled, **Event Subscriptions** (`app_mention`), and **Interactivity**

## Environment Variables

| Variable | Description |
|---|---|
| `SLACK_APP_TOKEN` | Slack app-level token (starts with `xapp-`) |
| `SLACK_BOT_TOKEN` | Slack bot user OAuth token (starts with `xoxb-`) |

## Getting Started

### Run Locally

```bash
export SLACK_APP_TOKEN="xapp-..."
export SLACK_BOT_TOKEN="xoxb-..."
go run ./cmd/k8bot
```

### Build

```bash
go build -o k8bot ./cmd/k8bot
```

## Usage

### Text Commands

Mention the bot followed by a command and an optional namespace (defaults to `default`):

```
@kubebot ping
@kubebot pods kube-system
@kubebot deployments argocd
@kubebot services default
@kubebot nodes
@kubebot namespaces
@kubebot events kube-system
@kubebot top nodes
@kubebot top pods kube-system
@kubebot menu
```

| Command | Alias | Description |
|---|---|---|
| `ping` | | Health check — replies `pong` |
| `namespaces` | `ns` | List all namespaces |
| `pods <ns>` | | List pods and status in a namespace |
| `deployments <ns>` | `deps` | List deployments and ready replicas |
| `services <ns>` | `svc` | List services and ClusterIPs |
| `nodes` | | List cluster nodes |
| `events <ns>` | | Show recent warning events |
| `top nodes` | | Node CPU and memory usage |
| `top pods <ns>` | | Pod CPU and memory usage |
| `menu` / `help` | | Show interactive Block Kit UI |

### Interactive UI

Type `@kubebot menu` to receive a dropdown-based control panel where you can select a namespace and resource type, then click **Fetch Data**.

## Project Structure

```
cmd/k8bot/main.go          Entry point — loads config, inits clients, starts Slack listener
internal/
  k8s/
    client.go               Kubernetes & Metrics client initialization
    helpers.go              Helper functions (e.g. GetNamespaceNames for UI)
    metrics.go              Node and pod metrics via metrics-server
    resources.go            Core resource queries (pods, nodes, deployments, services, events, namespaces)
  slack/
    client.go               Slack Socket Mode connection and event loop
    router.go               Command parser and interaction handler
    ui.go                   Block Kit interactive menu builder
deploy/
  Dockerfile                Container image build (placeholder)
  helm/                     Helm chart (placeholder)
```

## Deployment

The `deploy/` directory contains a Dockerfile and Helm chart scaffolding for deploying k8bot into a Kubernetes cluster. When running in-cluster, the bot uses the pod's service account to authenticate with the Kubernetes API — ensure the service account has appropriate RBAC permissions for the resources you want to query.

## License

This project is licensed under the [MIT License](LICENSE).
