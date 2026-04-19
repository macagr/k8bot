# Architecture

## Overview

k8bot is a Go application that bridges Slack and Kubernetes. It listens for Slack events over Socket Mode and translates them into Kubernetes API calls, returning formatted results back to the Slack channel.

```
┌──────────┐   Socket Mode   ┌──────────┐   K8s API   ┌────────────────┐
│  Slack   │ ◄─────────────► │  k8bot   │ ──────────► │  K8s Cluster   │
│  Users   │                 │          │ ──────────► │  metrics-server │
└──────────┘                 └──────────┘             └────────────────┘
```

## Components

### `cmd/k8bot/main.go`

Application entry point. Reads `SLACK_APP_TOKEN` and `SLACK_BOT_TOKEN` from the environment, initializes both the standard Kubernetes clientset and the Metrics clientset, then hands them off to the Slack listener.

### `internal/slack/`

- **client.go** — Creates the Slack API client and Socket Mode client. Spawns a goroutine that routes incoming events to the appropriate handler.
- **router.go** — Two handlers:
  - `handleMention` parses `@kubebot <command> [namespace]` text commands and dispatches to the corresponding `k8s` package function.
  - `handleInteraction` processes Block Kit button clicks from the interactive UI, extracting the selected namespace and resource type from Slack's block state.
- **ui.go** — `sendInteractiveMenu` builds a Block Kit message with two dropdowns (namespace, resource type) and a Fetch button. Namespaces are populated dynamically from the cluster.

### `internal/k8s/`

- **client.go** — `NewClient` and `NewMetricsClient` establish connections to the Kubernetes API. They attempt in-cluster config first and fall back to `~/.kube/config`.
- **helpers.go** — Utility functions like `GetNamespaceNames` that return raw data for UI components rather than formatted strings.
- **resources.go** — Query functions for core resources: `GetPods`, `GetNodes`, `GetDeployments`, `GetServices`, `GetWarningEvents`, `GetNamespaces`. Each returns a Slack-formatted string.
- **metrics.go** — `GetNodeMetrics` and `GetPodMetrics` query the Metrics API (v1beta1) and format CPU (millicores) and memory (MB) usage.

## Event Flow

1. A Slack user mentions the bot: `@kubebot pods kube-system`
2. Socket Mode delivers an `AppMentionEvent` to `client.go`
3. `client.go` routes it to `handleMention` in `router.go`
4. `handleMention` parses the command and namespace, calls `k8s.GetPods`
5. `GetPods` calls the Kubernetes API, formats the result
6. The reply is posted back to the originating Slack channel

For interactive UI:

1. User types `@kubebot menu` — `sendInteractiveMenu` posts a Block Kit message
2. User selects options and clicks **Fetch Data**
3. Socket Mode delivers an `InteractionCallback` to `client.go`
4. `handleInteraction` extracts the selections and calls the matching `k8s` function
5. Result is posted back to the channel

## Authentication

- **Slack**: App-level token (`xapp-`) for Socket Mode, bot token (`xoxb-`) for posting messages.
- **Kubernetes**: In-cluster service account when deployed, or local kubeconfig for development.
