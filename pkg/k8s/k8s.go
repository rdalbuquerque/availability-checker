package k8s

import (
	"context"
	"fmt"
	"log"
	"time"

	"errors"
	"math"

	"github.com/mitchellh/go-homedir"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sClient struct {
	clientset *kubernetes.Clientset
}

func NewK8sClient() (*K8sClient, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		// Try to use kubeconfig at user's .kube folder
		kubeconfig, err := homedir.Expand("~/.kube/config")
		if err != nil {
			return nil, err
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	// Create a Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &K8sClient{clientset: clientset}, nil
}

func (kc *K8sClient) GetDeployment(namespace, deploymentName string) (*appsv1.Deployment, error) {
	// Get the deployment
	deployment, err := kc.clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return deployment, nil
}

func (kc *K8sClient) ScaleDeployment(namespace, deploymentName string, replicas int32) error {
	// Get the deployment
	deployment, err := kc.GetDeployment(namespace, deploymentName)
	if err != nil {
		return err
	}

	// Set the replicas to the desired number
	deployment.Spec.Replicas = &replicas

	// Update the deployment
	_, err = kc.clientset.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("Deployment %s scaled to %d replicas\n", deploymentName, replicas)

	return nil
}

func (kc *K8sClient) ScaleDeploymentToZero(namespace, deploymentName string) error {
	// Get the deployment
	deployment, err := kc.GetDeployment(namespace, deploymentName)
	if err != nil {
		return err
	}

	// Set the replicas to 0
	deployment.Spec.Replicas = new(int32)
	*deployment.Spec.Replicas = 0

	// Update the deployment
	_, err = kc.clientset.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	// Wait for the deployment to scale down
	err = retryWithBackoffAndTimeout(func() error {
		pods, err := kc.clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: fmt.Sprintf("app=%s", deploymentName),
			FieldSelector: "status.phase=Running",
		})
		if err != nil {
			return err
		}

		if len(pods.Items) > 0 {
			return errors.New("deployment still has pods")
		}

		return nil
	}, 5*time.Minute, 1*time.Second)

	if err != nil {
		return err
	}

	fmt.Printf("Deployment %s scaled to 0 replicas\n", deploymentName)

	return nil
}

func retryWithBackoffAndTimeout(f func() error, timeout time.Duration, initialBackoff time.Duration) error {
	backoff := initialBackoff
	deadline := time.Now().Add(timeout)

	for {
		err := f()
		if err == nil {
			return nil
		}

		if time.Now().After(deadline) {
			return fmt.Errorf("timeout exceeded: %v", err)
		}

		log.Printf("Retrying after %v with backoff %v\n", err, backoff)
		time.Sleep(backoff)
		backoff = time.Duration(math.Min(float64(backoff*2), float64(time.Second*30)))
	}
}

func (kc *K8sClient) ScaleDeploymentToDesiredReplicas(namespace, deploymentName string, desiredReplicas int32) error {
	// Scale deployment to zero
	err := kc.ScaleDeploymentToZero(namespace, deploymentName)
	if err != nil {
		return err
	}

	// Scale deployment to desired replicas
	err = kc.ScaleDeployment(namespace, deploymentName, desiredReplicas)
	if err != nil {
		return err
	}

	return nil
}
