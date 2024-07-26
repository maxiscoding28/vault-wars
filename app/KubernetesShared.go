package app

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
	"vault-wars/util"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

func createKubernetesClient() (*kubernetes.Clientset, *rest.Config, error) {
	kubeconfig := os.Getenv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create Kubernetes client config: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create Kubernetes client: %v", err)
	}

	return clientset, config, nil
}

func getStatefulSet(clientset *kubernetes.Clientset, name string, namespace string) (*appsv1.StatefulSet, error) {
	ctx := context.Background()
	statefulSet, err := clientset.AppsV1().StatefulSets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get StatefulSet %s: %v", name, err)
	}
	return statefulSet, nil
}

func getPods(ctx context.Context, clientset *kubernetes.Clientset, namespace string, selector *metav1.LabelSelector) (*corev1.PodList, error) {
	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(selector),
	})

	return pods, err
}

func execOnPod(clientset *kubernetes.Clientset, config *rest.Config, podName string, namespace string, containerName string, cmdString string, ignoreNonZeroCode bool) (string, error) {
	command := []string{"/bin/sh", "-c", cmdString}
	req := clientset.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		Param("stdin", "false").
		Param("stdout", "true").
		Param("stderr", "true").
		Param("tty", "false").
		Param("container", containerName)

	for _, cmd := range command {
		req.Param("command", cmd)
	}
	executor, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return "", err
	}

	stdout, stderr := &strings.Builder{}, &strings.Builder{}
	err = executor.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: stdout,
		Stderr: stderr,
		Tty:    false,
	})

	if err != nil && ignoreNonZeroCode {
		return stdout.String(), nil
	}
	if err != nil {
		return stderr.String(), fmt.Errorf("error running command, %s", stderr.String())
	}
	return stdout.String(), nil
}

func AllPodsRunning(clientset *kubernetes.Clientset, releaseName string, namespace string) error {
	const retryDelay = 2 * time.Second

	statefulSet, err := getStatefulSet(clientset, fmt.Sprintf("%s-vault", releaseName), namespace)
	if err != nil {
		return err
	}
	for {
		ctx := context.Background()

		pods, err := getPods(ctx, clientset, namespace, statefulSet.Spec.Selector)
		if err != nil {
			return fmt.Errorf("failed to list pods: %v", err)
		}

		allRunning := true
		for _, pod := range pods.Items {
			if pod.Status.Phase != corev1.PodRunning {
				allRunning = false
				break
			}
		}

		if allRunning {
			util.LogInfo("All pods are running. Initializing cluster...")
			return nil
		}

		util.LogInfo("Not all pods are running. Retrying...")
		time.Sleep(retryDelay)
	}
}

func IsPodInitialized(clientset *kubernetes.Clientset, config *rest.Config, podName string, namespace string) (bool, error) {
	command := "vault status -format=json"
	output, err := execOnPod(clientset, config, podName, namespace, "vault", command, true)
	if err != nil {
		return false, fmt.Errorf("failed to execute command in pod %s: %v", podName, err)
	}

	var status VaultStatus
	if err := json.Unmarshal([]byte(output), &status); err != nil {
		return false, fmt.Errorf("failed to unmarshal Vault status for pod %s: %v", podName, err)
	}

	return status.Initialized, nil
}

func EnsureNoPodsInitialized(clientset *kubernetes.Clientset, config *rest.Config, releaseName string, namespace string) error {
	const maxRetries = 3
	const retryDelay = 3 * time.Second

	statefulSet, err := getStatefulSet(clientset, fmt.Sprintf("%s-vault", releaseName), namespace)
	if err != nil {
		return err
	}

	ctx := context.Background()
	pods, err := getPods(ctx, clientset, namespace, statefulSet.Spec.Selector)
	if err != nil {
		return fmt.Errorf("failed to list pods: %v", err)
	}

	for i := 0; i < maxRetries; i++ {
		for _, pod := range pods.Items {
			command := "vault status -format=json"
			output, err := execOnPod(clientset, config, pod.Name, namespace, "vault", command, true)
			if err != nil {
				return fmt.Errorf("failed to execute command in pod %s: %v", pod.Name, err)
			}

			var status VaultStatus
			if err := json.Unmarshal([]byte(output), &status); err != nil {
				return fmt.Errorf("failed to unmarshal Vault status for pod %s: %v", pod.Name, err)
			}

			if status.Initialized {
				return fmt.Errorf("pod %s is already initialized", pod.Name)
			}
		}

		util.LogWarn("Checking if any pods are initialized...")
		time.Sleep(retryDelay)
	}

	util.LogInfo(fmt.Sprintf("No pods already initialized after %d checks, attempting to initialize...", maxRetries))
	return nil
}

func CountInitializedPods(clientset *kubernetes.Clientset, config *rest.Config, releaseName string, namespace string) (int, int, error) {
	const maxRetries = 5
	const retryDelay = 3 * time.Second

	statefulSet, err := getStatefulSet(clientset, fmt.Sprintf("%s-vault", releaseName), namespace)
	if err != nil {
		return 0, 0, err
	}

	ctx := context.Background()
	pods, err := getPods(ctx, clientset, namespace, statefulSet.Spec.Selector)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to list pods: %v", err)
	}

	podStatusMap := make(map[string]bool)
	for _, pod := range pods.Items {
		podStatusMap[pod.Name] = false
	}

	for i := 0; i < maxRetries; i++ {
		for podName, initialized := range podStatusMap {
			if initialized {
				continue
			}

			command := "vault status -format=json"
			output, err := execOnPod(clientset, config, podName, namespace, "vault", command, true)
			if err != nil {
				return 0, 0, fmt.Errorf("failed to execute command in pod %s: %v", podName, err)
			}

			var status VaultStatus
			if err := json.Unmarshal([]byte(output), &status); err != nil {
				return 0, 0, fmt.Errorf("failed to unmarshal Vault status for pod %s: %v", podName, err)
			}

			if status.Initialized {
				podStatusMap[podName] = true
			}
		}

		allInitialized := true
		for _, initialized := range podStatusMap {
			if !initialized {
				allInitialized = false
				break
			}
		}

		if allInitialized {
			break
		}

		util.LogWarn("Not all pods are initialized. Retrying...")
		time.Sleep(retryDelay)
	}

	initializedCount := 0
	for _, initialized := range podStatusMap {
		if initialized {
			initializedCount++
		}
	}

	return len(podStatusMap), initializedCount, nil
}

func UnsealPods(clientset *kubernetes.Clientset, config *rest.Config, releaseName string, namespace string, excludePodName string, unsealKey string) error {
	const delay = 2 * time.Second
	statefulSet, err := getStatefulSet(clientset, fmt.Sprintf("%s-vault", releaseName), namespace)
	if err != nil {
		return fmt.Errorf("failed to get StatefulSet: %v", err)
	}

	ctx := context.Background()
	pods, err := getPods(ctx, clientset, namespace, statefulSet.Spec.Selector)
	if err != nil {
		return fmt.Errorf("failed to list pods: %v", err)
	}

	for _, pod := range pods.Items {
		if pod.Name == excludePodName {
			continue
		}

		command := fmt.Sprintf("vault operator unseal %s", unsealKey)
		out, err := execOnPod(clientset, config, pod.Name, namespace, "vault", command, true)
		if err != nil {
			return fmt.Errorf("failed to execute unseal command on pod %s: %v", pod.Name, err)
		}
		time.Sleep(delay)
		fmt.Printf("Unseal command executed on pod %s\n%v", pod.Name, out)
	}

	return nil
}
