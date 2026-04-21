package main

import (
	"fmt"
	"log"
	"os"

	"k8bot/internal/k8s"
	"k8bot/internal/slack"
)

func main() {
	// 1. Grab Config (Tokens)
	appToken := os.Getenv("SLACK_APP_TOKEN")
	botToken := os.Getenv("SLACK_BOT_TOKEN")

	if appToken == "" || botToken == "" {
		log.Fatal("Must set SLACK_APP_TOKEN and SLACK_BOT_TOKEN")
	}

	// 2. Init Standard Kubernetes Client (for Pods, Nodes, Deployments)
	k8sClient, err := k8s.NewClient()
	if err != nil {
		log.Fatalf("Failed to create K8s client: %v", err)
	}
	fmt.Println("✅ Connected to Standard Kubernetes API")

	// 3. Init Metrics Kubernetes Client (for CPU/RAM top commands)
	metricsClient, err := k8s.NewMetricsClient()
	if err != nil {
		log.Fatalf("Failed to create K8s Metrics client: %v", err)
	}
	fmt.Println("✅ Connected to Kubernetes Metrics API")

	// 4. Start Slack (passing BOTH clients now)
	slack.Start(appToken, botToken, k8sClient, metricsClient)
}