package slack

import (
	"fmt"

	"k8bot/internal/k8s"

	"github.com/slack-go/slack"
	"k8s.io/client-go/kubernetes"
)

// sendInteractiveMenu builds the Double-Dropdown UI
// NOTICE: We added clientset to the parameters!
func sendInteractiveMenu(api *slack.Client, clientset *kubernetes.Clientset, channelID string) {
	// 1. Header Text
	headerText := slack.NewTextBlockObject("mrkdwn", "*🤖 Kubebot Command Center*\nSelect your parameters and click Fetch:", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// 2. Dropdown 1: Namespaces (DYNAMIC!)
	var namespaceOptions []*slack.OptionBlockObject
	
	// Fetch live namespaces from the cluster
	liveNamespaces, err := k8s.GetNamespaceNames(clientset)
	if err != nil || len(liveNamespaces) == 0 {
		// Fallback if the RBAC fails or cluster is empty
		liveNamespaces = []string{"default"} 
	}

	// Slack limits dropdowns to 100 items. Let's cap it just in case.
	limit := 99
	if len(liveNamespaces) < limit {
		limit = len(liveNamespaces)
	}

	// Loop through the live cluster namespaces to build the UI array
	for i := 0; i < limit; i++ {
		ns := liveNamespaces[i]
		namespaceOptions = append(namespaceOptions, slack.NewOptionBlockObject(ns, slack.NewTextBlockObject("plain_text", ns, false, false), nil))
	}

	namespaceDropdown := slack.NewOptionsSelectBlockElement("static_select", slack.NewTextBlockObject("plain_text", "Select Namespace...", false, false), "select_namespace", namespaceOptions...)
	dropdown1Text := slack.NewTextBlockObject("mrkdwn", "*Namespace:*", false, false)
	dropdown1Section := slack.NewSectionBlock(dropdown1Text, nil, slack.NewAccessory(namespaceDropdown))
	dropdown1Section.BlockID = "block_namespace"

	// 3. Dropdown 2: Resource Type
	resourceOptions := []*slack.OptionBlockObject{
		slack.NewOptionBlockObject("namespaces", slack.NewTextBlockObject("plain_text", "📂 Namespaces", false, false), nil),
		slack.NewOptionBlockObject("pods", slack.NewTextBlockObject("plain_text", "📦 Pods", false, false), nil),
		slack.NewOptionBlockObject("deployments", slack.NewTextBlockObject("plain_text", "🚀 Deployments", false, false), nil),
		slack.NewOptionBlockObject("services", slack.NewTextBlockObject("plain_text", "🌐 Services", false, false), nil),
		slack.NewOptionBlockObject("nodes", slack.NewTextBlockObject("plain_text", "🖥️ Nodes", false, false), nil),
		slack.NewOptionBlockObject("events", slack.NewTextBlockObject("plain_text", "⚠️ Warning Events", false, false), nil),
		slack.NewOptionBlockObject("node-metrics", slack.NewTextBlockObject("plain_text", "📊 Node Metrics (top)", false, false), nil),
		slack.NewOptionBlockObject("pod-metrics", slack.NewTextBlockObject("plain_text", "📈 Pod Metrics (top)", false, false), nil),
	}
	
	resourceDropdown := slack.NewOptionsSelectBlockElement("static_select", slack.NewTextBlockObject("plain_text", "Select Resource...", false, false), "select_resource", resourceOptions...)
	dropdown2Text := slack.NewTextBlockObject("mrkdwn", "*Resource:*", false, false)
	dropdown2Section := slack.NewSectionBlock(dropdown2Text, nil, slack.NewAccessory(resourceDropdown))
	dropdown2Section.BlockID = "block_resource"

	// 4. The Single "Fetch" Button
	fetchBtn := slack.NewButtonBlockElement("action_fetch", "fetch", slack.NewTextBlockObject("plain_text", "🔍 Fetch Data", false, false))
	fetchBtn.Style = slack.StylePrimary
	actionBlock := slack.NewActionBlock("block_actions", fetchBtn)

	// 5. Send the payload to Slack
	_, _, err = api.PostMessage(
		channelID,
		slack.MsgOptionBlocks(
			headerSection,
			slack.NewDividerBlock(),
			dropdown1Section,
			dropdown2Section,
			slack.NewDividerBlock(),
			actionBlock,
		),
	)

	if err != nil {
		fmt.Printf("Failed to send menu: %v\n", err)
	}
}