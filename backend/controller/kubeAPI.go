package controller

import (
	"fmt"
	"net/http"
	"time"
	"os"

	"context"
	"github.com/google/uuid"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"encoding/json"
	"strings"
	//ingcorev1 "k8s.io/client-go/kubernetes"
	//"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var (
	PREFIX    = ""
	NAMESPACE = os.Getenv("NAMESPACE")
)

func CreateNewWetty(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
	id := uuid.New()
	PREFIX = "x"+strings.Split(id.String(), "-")[0]

	kc := GetKubeClient()

	svc := getServiceObject()
	_, err := kc.CoreV1().Services(NAMESPACE).Create(context.TODO(), svc, metav1.CreateOptions{})
	if err != nil {
		returnError(w)
		panic(err)
	}
	fmt.Println("Service Created successfully!")

	ing := getIngressObject()
	_, err = kc.NetworkingV1().Ingresses(NAMESPACE).Create(context.TODO(), ing, metav1.CreateOptions{})
	if err != nil {
		returnError(w)
		panic(err)
	}
	fmt.Println("Ingress Created successfully!")

	deploy := getDeployObject()
	_, err = kc.AppsV1().Deployments(NAMESPACE).Create(context.TODO(), deploy, metav1.CreateOptions{})
	if err != nil {
		returnError(w)
		panic(err)
	}

	success := podRunning()
	if !success {
		returnError(w)
		panic(err)
	}
	fmt.Println("Deployment Created successfully!")

	returnOKUID(w)
}

func podRunning() bool{
	cnt := 0
	kc := GetKubeClient()
	labelSelector := "component=" + PREFIX + "-wetty"
	watch, err := kc.CoreV1().Pods(NAMESPACE).Watch(context.TODO(), metav1.ListOptions{
	    LabelSelector: labelSelector,
	})
	if err != nil {
	    fmt.Println(err.Error())
	}

	for event := range watch.ResultChan() {
	    //fmt.Printf("Type: %v\n", event.Type)
	    p, ok := event.Object.(*corev1.Pod)
	    if !ok || cnt > 12 {
	        fmt.Println("unexpected type")
			return false
	    }

	    if p.Status.Phase == "Running"{
			return true
		}

		time.Sleep(5 * time.Second)
		cnt = cnt + 1
	}

	return false
}

func returnError(w http.ResponseWriter){
	resp := make(map[string]string)
	w.WriteHeader(http.StatusBadRequest)
	resp["result"] = "server error"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		panic("Error happened in JSON marshal")
	}
	w.Write(jsonResp)
}

func returnOKUID(w http.ResponseWriter){
	resp := make(map[string]string)
	w.WriteHeader(http.StatusOK)
	resp["uid"] = PREFIX
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		panic("Error happened in JSON marshal. Err")
	}
	w.Write(jsonResp)
}

func getDeployObject() *appsv1.Deployment {
	deploy := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      PREFIX + "-wetty-deploy",
			Namespace: NAMESPACE,
			Labels: map[string]string{
				"component": PREFIX + "-wetty",
			},
			Annotations: map[string]string{
				"deployment.kubernetes.io/revision": "2",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: ptrint32(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"component": PREFIX + "-wetty",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"component": PREFIX + "-wetty",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						corev1.Container{
							Name:  "master",
							Image: os.Getenv("IMAGE"),
							Ports: []corev1.ContainerPort{
								corev1.ContainerPort{
									Name:          "http",
									HostPort:      0,
									ContainerPort: 3000,
									Protocol:      corev1.Protocol("TCP"),
								},
							},
							Env: []corev1.EnvVar{
								corev1.EnvVar{
									Name: "BASE",
									Value: "/" + PREFIX,
								},
								corev1.EnvVar{
									Name: "MY_POD_IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											APIVersion: "v1",
											FieldPath:  "status.podIP",
										},
									},
								},
							},
							/*Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									"cpu":    *resource.NewQuantity(600, resource.DecimalSI)+"m",
									"memory": *resource.NewQuantity(524288000, resource.BinarySI),
								},
								Requests: corev1.ResourceList{
									"cpu":    *resource.NewQuantity(50, resource.DecimalSI)+"m",
									"memory": *resource.NewQuantity(157286400, resource.BinarySI),
								},
							},*/
							TerminationMessagePath:   "/dev/termination-log",
							TerminationMessagePolicy: corev1.TerminationMessagePolicy("File"),
							ImagePullPolicy:          corev1.PullPolicy("IfNotPresent"),
						},
					},
					RestartPolicy:                 corev1.RestartPolicy("Always"),
					TerminationGracePeriodSeconds: ptrint64(30),
					DNSPolicy:                     corev1.DNSPolicy("ClusterFirst"),
					SchedulerName:                 "default-scheduler",
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.DeploymentStrategyType("RollingUpdate"),
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Type(1),
						IntVal: 0,
						StrVal: "25%",
					},
					MaxSurge: &intstr.IntOrString{
						Type:   intstr.Type(1),
						IntVal: 0,
						StrVal: "25%",
					},
				},
			},
			MinReadySeconds:         0,
			RevisionHistoryLimit:    ptrint32(1),
			ProgressDeadlineSeconds: ptrint32(600),
		},
	}
	return deploy
}

func getServiceObject() *corev1.Service {
	svc := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      PREFIX + "-wetty-svc",
			Namespace: NAMESPACE,
			Labels: map[string]string{
				"component": PREFIX + "-wetty",
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				corev1.ServicePort{
					Name: "http",
					Port: 3000,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Type(0),
						IntVal: 0,
					},
					NodePort: 0,
				},
			},
			Selector: map[string]string{
				"component":                  PREFIX + "-wetty",
			},
		},
	}
	return svc
}

func getIngressObject() *networkingv1.Ingress {
	ing := &networkingv1.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "networking.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      PREFIX+"-wetty-ing",
			Namespace: NAMESPACE,
			Labels: map[string]string{
				"component": PREFIX + "-wetty",
			},
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: ptrstring("nginx-internal"),
			Rules: []networkingv1.IngressRule{
				networkingv1.IngressRule{
					Host: "k8tty.yad2.io",
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								networkingv1.HTTPIngressPath{
									Path:     "/"+PREFIX,
									PathType: ptrPathType("ImplementationSpecific"),
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: PREFIX+"-wetty-svc",
											Port: networkingv1.ServiceBackendPort{
												Name:   "http",
												Number: 0,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return ing
}

func ptrstring(p string) *string {
	return &p
}

func ptrPathType(p networkingv1.PathType) *networkingv1.PathType {
	return &p
}
func ptrint64(p int64) *int64 {
	return &p
}

func ptrint32(p int32) *int32 {
	return &p
}
