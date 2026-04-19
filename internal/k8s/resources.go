package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// 1. GetPods lists pods and their status in a specific namespace
func GetPods(clientset *kubernetes.Clientset, namespace string) (string, error) {
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", err
	}

	if len(pods.Items) == 0 {
		return fmt.Sprintf("✅ No pods found in `%s`.", namespace), nil
	}

	reply := fmt.Sprintf("📦 *Pods in `%s` (%d total):*\n", namespace, len(pods.Items))
	
	// Limit output so we don't hit Slack's message size limits on massive namespaces
	limit := 15
	if len(pods.Items) < 15 {
		limit = len(pods.Items)
	}

	for i := 0; i < limit; i++ {
		p := pods.Items[i]
		// This will output: • my-app-7d48f9-abcde - Running
		reply += fmt.Sprintf("• `%s` - %s\n", p.Name, p.Status.Phase)
	}
	
	if len(pods.Items) > limit {
		reply += fmt.Sprintf("... and %d more.", len(pods.Items)-limit)
	}
	
	return reply, nil
}

// 2. GetNodes fetches a quick summary of the cluster nodes
func GetNodes(clientset *kubernetes.Clientset) (string, error) {
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", err
	}

	reply := fmt.Sprintf("🖥️ *Cluster Nodes (%d total):*\n", len(nodes.Items))
	for _, node := range nodes.Items {
		// Just grab the OS image and name as a quick health check
		osImage := node.Status.NodeInfo.OSImage
		reply += fmt.Sprintf("• `%s` (%s)\n", node.Name, osImage)
	}
	return reply, nil
}

// 3. GetWarningEvents fetches recent warnings in a specific namespace
func GetWarningEvents(clientset *kubernetes.Clientset, namespace string) (string, error) {
	// FieldSelector is K8s magic to only bring back warnings
	events, err := clientset.CoreV1().Events(namespace).List(context.TODO(), metav1.ListOptions{
		FieldSelector: "type=Warning",
	})
	if err != nil {
		return "", err
	}

	if len(events.Items) == 0 {
		return fmt.Sprintf("✅ No warning events found in namespace `%s` right now.", namespace), nil
	}

	reply := fmt.Sprintf("⚠️ *Recent Warnings in `%s`:*\n", namespace)
	// Only show the last 5 to not spam Slack
	limit := 5
	if len(events.Items) < 5 {
		limit = len(events.Items)
	}

	for i := 0; i < limit; i++ {
		evt := events.Items[i]
		reply += fmt.Sprintf("• *%s*: %s\n", evt.InvolvedObject.Kind, evt.Message)
	}
	return reply, nil
}

// 4. GetDeployments lists deployments and their ready status
func GetDeployments(clientset *kubernetes.Clientset, namespace string) (string, error) {
	// Notice we use AppsV1() for Deployments, not CoreV1()
	deps, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", err
	}

	if len(deps.Items) == 0 {
		return fmt.Sprintf("No deployments found in `%s`.", namespace), nil
	}

	reply := fmt.Sprintf("🚀 *Deployments in `%s`:*\n", namespace)
	for _, d := range deps.Items {
		reply += fmt.Sprintf("• `%s`: %d/%d ready\n", d.Name, d.Status.ReadyReplicas, d.Status.Replicas)
	}
	return reply, nil
}

// 5. GetServices lists services and their ClusterIPs
func GetServices(clientset *kubernetes.Clientset, namespace string) (string, error) {
	svcs, err := clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", err
	}

	if len(svcs.Items) == 0 {
		return fmt.Sprintf("No services found in `%s`.", namespace), nil
	}

	reply := fmt.Sprintf("🌐 *Services in `%s`:*\n", namespace)
	for _, s := range svcs.Items {
		reply += fmt.Sprintf("• `%s` (%s): %s\n", s.Name, s.Spec.Type, s.Spec.ClusterIP)
	}
	return reply, nil
}

// 6. GetNamespaces lists all namespaces in the cluster (used to populate our dropdown menu)
func GetNamespaces(clientset *kubernetes.Clientset) (string, error) {
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", err
	}

	if len(namespaces.Items) == 0 {
		return "✅ No namespaces found.", nil
	}

	reply := fmt.Sprintf("📂 *Namespaces (%d total):*\n", len(namespaces.Items))
	for _, ns := range namespaces.Items {
		reply += fmt.Sprintf("• `%s`\n", ns.Name)
	}
	return reply, nil
}