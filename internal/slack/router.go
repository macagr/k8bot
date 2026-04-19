package slack

import (
	"fmt"
	"strings"

	"k8bot/internal/k8s"

	slackgo "github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"k8s.io/client-go/kubernetes"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

// handleMention processes typed text commands (ChatOps)
func handleMention(api *slackgo.Client, clientset *kubernetes.Clientset, metricsClient *metricsv.Clientset, ev *slackevents.AppMentionEvent) {
	// The text usually comes in as "<@U12345BOT> command namespace"
	text := strings.TrimSpace(ev.Text)
	parts := strings.Fields(text) // fields cleanly splits by any whitespace

	// Set defaults
	command := "menu"
	namespace := "default"

	// Parse command (e.g. "pods")
	if len(parts) > 1 {
		command = strings.ToLower(parts[1])
	}
	// Parse namespace (e.g. "kube-system")
	if len(parts) > 2 {
		namespace = strings.ToLower(parts[2])
	}

	var reply string
	var err error

	// Route the typed command
	switch command {
	case "ping":
		reply = "pong 🏓"
	case "namespaces", "ns":
		reply, err = k8s.GetNamespaces(clientset)	
	case "pods":
		reply, err = k8s.GetPods(clientset, namespace)
	case "deployments", "deps": // Aliases are great for ChatOps!
		reply, err = k8s.GetDeployments(clientset, namespace)
	case "services", "svc":
		reply, err = k8s.GetServices(clientset, namespace)
	case "nodes":
		reply, err = k8s.GetNodes(clientset)
	case "events":
		reply, err = k8s.GetWarningEvents(clientset, namespace)
	case "top":
			// If they type "@kubebot top nodes" or "@kubebot top pods argocd"
			target := "pods" // default
			if len(parts) > 2 {
				target = strings.ToLower(parts[2])
			}
			
			// If they provided a namespace, it will be the 4th word: "@kubebot top pods kube-system"
			if len(parts) > 3 {
				namespace = strings.ToLower(parts[3])
			}

			if target == "nodes" {
				reply, err = k8s.GetNodeMetrics(metricsClient)
			} else {
				reply, err = k8s.GetPodMetrics(metricsClient, namespace)
			}	
case "menu", "help":
		// This triggers the UI block. Pass the clientset now!
		sendInteractiveMenu(api, clientset, ev.Channel)
		return
	default:
		reply = "❌ Unknown command. Type `@kubebot menu` to see the UI or use commands like `@kubebot pods kube-system`."
	}

	// Catch any Kubernetes API errors
	if err != nil {
		reply = fmt.Sprintf("❌ API Error: %v", err)
	}

	// Post the final text response back to Slack
	api.PostMessage(ev.Channel, slackgo.MsgOptionText(reply, false))
}

// handleInteraction processes Block Kit button clicks
func handleInteraction(api *slackgo.Client, clientset *kubernetes.Clientset, metricsClient *metricsv.Clientset, callback slackgo.InteractionCallback) {
    // ... rest of your code stays the same
	if callback.Type != slackgo.InteractionTypeBlockActions {
		return
	}

	for _, action := range callback.ActionCallback.BlockActions {
		// We only care when they click the main "Fetch Data" button
		if action.ActionID == "action_fetch" {
			
			// 1. Set our defaults in case the user didn't pick anything
			namespace := "default"
			resourceType := "pods"

			// 2. Extract State from Slack's weird nested map: state.Values[BlockID][ActionID]
			state := callback.BlockActionState.Values
			
			// Grab Namespace
			if nsBlock, ok := state["block_namespace"]; ok {
				if nsAction, ok := nsBlock["select_namespace"]; ok && nsAction.SelectedOption.Value != "" {
					namespace = nsAction.SelectedOption.Value
				}
			}

			// Grab Resource Type
			if resBlock, ok := state["block_resource"]; ok {
				if resAction, ok := resBlock["select_resource"]; ok && resAction.SelectedOption.Value != "" {
					resourceType = resAction.SelectedOption.Value
				}
			}

			// 3. Route to the correct K8s function based on the dropdown choice
			var reply string
			var err error

			switch resourceType {
			case "namespaces":
				reply, err = k8s.GetNamespaces(clientset)
			case "pods":
				count, err := k8s.GetPods(clientset, namespace)
				if err == nil {
					reply = fmt.Sprintf("📦 There are *%d* pods in `%s`.", count, namespace)
				}
			case "deployments":
				reply, err = k8s.GetDeployments(clientset, namespace)
			case "services":
				reply, err = k8s.GetServices(clientset, namespace)
			case "nodes":
				reply, err = k8s.GetNodes(clientset) // Nodes ignore namespace
			case "events":
				reply, err = k8s.GetWarningEvents(clientset, namespace)
			case "node-metrics":
				reply, err = k8s.GetNodeMetrics(metricsClient) // ignores namespace dropdown
			case "pod-metrics":
				reply, err = k8s.GetPodMetrics(metricsClient, namespace)
			default:
				reply = "❌ Unknown resource type selected."
			}

			// 4. Handle errors gracefully and send the response back
			if err != nil {
				reply = fmt.Sprintf("❌ API Error: %v", err)
			}
			
			api.PostMessage(callback.Channel.ID, slackgo.MsgOptionText(reply, false))
		}
	}
}