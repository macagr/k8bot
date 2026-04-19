package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

// GetNodeMetrics retrieves metrics for all nodes and formats them cleanly
func GetNodeMetrics(metricsClient *metricsv.Clientset) (string, error) {
	nodeMetrics, err := metricsClient.MetricsV1beta1().NodeMetricses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		// Helpful error message in case metrics-server dies
		return "", fmt.Errorf("failed to fetch metrics (is metrics-server running?): %v", err)
	}

	if len(nodeMetrics.Items) == 0 {
		return "✅ No node metrics found.", nil
	}

	reply := fmt.Sprintf("📊 *Node Resource Usage (%d total):*\n", len(nodeMetrics.Items))
	for _, nm := range nodeMetrics.Items {
		// 1. CPU: Get value in millicores (e.g., 250m)
		cpuMilli := nm.Usage.Cpu().MilliValue()

		// 2. Memory: Get raw bytes and convert to Megabytes (MB)
		memBytes := nm.Usage.Memory().Value()
		memMB := memBytes / (1024 * 1024)

		reply += fmt.Sprintf("• `%s` ➔ CPU: *%dm* | RAM: *%d MB*\n", nm.Name, cpuMilli, memMB)
	}
	return reply, nil
}

// 2. GetPodMetrics retrieves aggregated metrics for pods in a specific namespace
func GetPodMetrics(metricsClient *metricsv.Clientset, namespace string) (string, error) {
	podMetrics, err := metricsClient.MetricsV1beta1().PodMetricses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to fetch pod metrics: %v", err)
	}

	if len(podMetrics.Items) == 0 {
		return fmt.Sprintf("✅ No pod metrics found in `%s`.", namespace), nil
	}

	reply := fmt.Sprintf("📈 *Pod Resource Usage in `%s`:*\n", namespace)
	
	// Limit output to prevent Slack spam on massive namespaces
	limit := 15
	if len(podMetrics.Items) < 15 {
		limit = len(podMetrics.Items)
	}

	for i := 0; i < limit; i++ {
		pm := podMetrics.Items[i]
		var totalCPU int64 = 0
		var totalMem int64 = 0
		
		// K8s pods can have multiple containers (e.g., an app + a sidecar). 
		// We must sum them up to get the total Pod usage.
		for _, container := range pm.Containers {
			totalCPU += container.Usage.Cpu().MilliValue()
			totalMem += container.Usage.Memory().Value()
		}
		
		memMB := totalMem / (1024 * 1024)
		reply += fmt.Sprintf("• `%s` ➔ CPU: *%dm* | RAM: *%d MB*\n", pm.Name, totalCPU, memMB)
	}
	
	if len(podMetrics.Items) > limit {
		reply += fmt.Sprintf("... and %d more.", len(podMetrics.Items)-limit)
	}

	return reply, nil
}