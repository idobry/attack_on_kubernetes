package controller

import (
    "fmt"
	"net/http"

    "context"
    core "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreatePod(w http.ResponseWriter, r *http.Request) {
    // build the pod definition
	fmt.Println("11111")
    desiredPod := getPodObjet()

    pod, err := GetKubeClient().CoreV1().Pods(desiredPod.Namespace).Create( context.TODO(), desiredPod , metav1.CreateOptions{})
    if err != nil {
		fmt.Println("Failed to create the static pod")
    }
    fmt.Println("Created Pod: ", pod.Name)
}

func getPodObjet() *core.Pod {
    pod := &core.Pod{
        ObjectMeta: metav1.ObjectMeta{
            Name:      "test-pod",
            Namespace: "default",
            Labels: map[string]string{
                "app": "test-pod",
            },
        },
        Spec: core.PodSpec{
            Containers: []core.Container{
                {
                    Name:            "busybox",
                    Image:           "busybox",
                    ImagePullPolicy: core.PullIfNotPresent,
                    Command: []string{
                        "sleep",
                        "3600",
                    },
                },
            },
        },
    }
    return pod
}