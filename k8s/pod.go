package k8s

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var clientset *kubernetes.Clientset

type Pod struct {
	Name string
	Port string
}

func GetPods() []Pod {
	pods := []Pod{}

	podList, _ := clientset.CoreV1().Pods("default").List(context.TODO(), v1.ListOptions{})

	for _, pod := range podList.Items {
		pods = append(pods, Pod{
			Name: pod.Name,
			Port: fmt.Sprint(pod.Spec.Containers[0].Ports[0].ContainerPort),
		})
	}

	return pods
}

func init() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err = kubernetes.NewForConfig(config)
}
