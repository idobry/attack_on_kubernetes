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
	"k8s.io/client-go/kubernetes"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"encoding/json"
	"strings"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var (
	PREFIX    = ""
	NAMESPACE = os.Getenv("NAMESPACE") //man
	HOST = os.Getenv("HOST") //man
	IMAGE = os.Getenv("IMAGE") //man
	SSHPASS = os.Getenv("SSHPASS") //man
	SSHUSER = os.Getenv("SSHUSER") //man
	SAN = os.Getenv("SERVICE_ACCOUNT_NAME") //optional 
)

func DeleteWetty(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")

	decoder := json.NewDecoder(r.Body)
	deleteRequest := struct {
		UID string `json:"uid"` // pointer so we can test for field absence
	}{}
    err := decoder.Decode(&deleteRequest)
    if err != nil {
        panic(err)
    }
    fmt.Println(deleteRequest.UID)

	kc := GetKubeClient()

	err = kc.CoreV1().Services(NAMESPACE).Delete(context.TODO(), deleteRequest.UID + "-wetty-svc", metav1.DeleteOptions{})
	if err != nil {
		returnError(w)
		panic(err)
	}
	fmt.Println("Service " + deleteRequest.UID + " Deleted successfully!")

	err = kc.NetworkingV1().Ingresses(NAMESPACE).Delete(context.TODO(), deleteRequest.UID + "-wetty-ing", metav1.DeleteOptions{})
	if err != nil {
		returnError(w)
		panic(err)
	}
	fmt.Println("Ingress " + deleteRequest.UID + " Deleted successfully!")

	err = kc.AppsV1().Deployments(NAMESPACE).Delete(context.TODO(), deleteRequest.UID + "-wetty-deploy", metav1.DeleteOptions{})
	if err != nil {
		returnError(w)
		panic(err)
	}
	fmt.Println("Deployment " + deleteRequest.UID + " Deleted successfully!")

	returnOKUID(w)
}

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
	fmt.Println("Service " + PREFIX + " Created successfully!")

	ing := getIngressObject()
	_, err = kc.NetworkingV1().Ingresses(NAMESPACE).Create(context.TODO(), ing, metav1.CreateOptions{})
	if err != nil {
		returnError(w)
		panic(err)
	}
	fmt.Println("Ingress " + PREFIX + " Created successfully!")

	deploy := getDeployObject()
	_, err = kc.AppsV1().Deployments(NAMESPACE).Create(context.TODO(), deploy, metav1.CreateOptions{})
	if err != nil {
		returnError(w)
		panic(err)
	}

	success := podRunning(kc)
	if !success {
		returnError(w)
		panic(err)
	}
	fmt.Println("Deployment " + PREFIX + " Created successfully!")

	returnOKUID(w)
}

func podRunning(kc *kubernetes.Clientset) bool{
	count := 0
	labelSelector := "component=" + PREFIX + "-wetty"
	fmt.Println(labelSelector)
	fmt.Println(NAMESPACE)
	watch, err := kc.CoreV1().Pods(NAMESPACE).Watch(context.TODO(), metav1.ListOptions{
	    LabelSelector: labelSelector,
	})
	if err != nil {
	    fmt.Println(err.Error())
	}
	for event := range watch.ResultChan() {
	    //fmt.Printf("Type: %v\n", event.Type)
	    p, ok := event.Object.(*corev1.Pod)
	    if !ok || count > 12 {
	        fmt.Println("unexpected type")
			return false
	    }

	    if p.Status.Phase == "Running"{
			return true
		}

		time.Sleep(5 * time.Second)
		count = count + 1
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
		panic("Error happened in JSON marshal")
	}
	w.Write(jsonResp)
}

func getDeployObject() *appsv1.Deployment {

	if SAN == ""{
		SAN = "defailt"
	}

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
							Image: IMAGE,
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
									Name: "SSHUSER",
									Value: SSHUSER,
								},
								corev1.EnvVar{
									Name: "SSHPASS",
									Value: SSHPASS,
								},
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
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									"cpu":    *resource.NewQuantity(int64(1), resource.DecimalSI),
									"memory": *resource.NewQuantity(524288000, resource.BinarySI),
								},
								/*Requests: corev1.ResourceList{
									"cpu":    *resource.NewQuantity(.5, resource.DecimalSI),
									"memory": *resource.NewQuantity(157286400, resource.BinarySI),
								},*/
							},
							TerminationMessagePath:   "/dev/termination-log",
							TerminationMessagePolicy: corev1.TerminationMessagePolicy("File"),
							ImagePullPolicy:          corev1.PullPolicy("IfNotPresent"),
						},
					},
					ServiceAccountName: SAN,
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
				"component": PREFIX + "-wetty",
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
					Host: HOST,
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
