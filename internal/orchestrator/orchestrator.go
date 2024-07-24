package orchestrator

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/0jk6/freight/internal/db"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// spin up a job in k8s
func getK8sClient() *kubernetes.Clientset {

	config, err := rest.InClusterConfig()

	if err != nil {
		//dont panic, check if the env var is dev
		if os.Getenv("ENV") == "dev" {
			kubeconfigPath := "/Users/user1/.kube/config"
			config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return clientset
}

func RunJob(namespace, jobName, image, code string, command []string) {
	clientset := getK8sClient()

	jobsClient := clientset.BatchV1().Jobs(namespace)

	backoffLimit := int32(0)           // number of retries
	activeDeadlineSeconds := int64(10) // time limit for the job

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: namespace,
		},
		Spec: batchv1.JobSpec{
			BackoffLimit:          &backoffLimit,
			ActiveDeadlineSeconds: &activeDeadlineSeconds,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    jobName,
							Image:   image,
							Command: command,
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									// corev1.ResourceCPU:    resource.MustParse("250m"),
									// corev1.ResourceMemory: resource.MustParse("256Mi"),
								},
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
	}

	_, err := jobsClient.Create(context.Background(), job, metav1.CreateOptions{})

	if err != nil {
		panic(err)
	}
}

func ListJobs(namespace string) {
	// log.Println("processing jobs")
	clientset := getK8sClient()
	jobsClient := clientset.BatchV1().Jobs(namespace)

	jobs, err := jobsClient.List(context.Background(), metav1.ListOptions{})

	if err != nil {
		panic(err)
	}

	for _, job := range jobs.Items {
		go func() {
			log.Printf("job name: %s\n", job.Name)
			log.Printf("completions: %d\n", *job.Spec.Completions)
			log.Printf("succeeded: %d\n", job.Status.Succeeded)
			log.Printf("failed: %d\n", job.Status.Failed)
			log.Println("------------")

			if *job.Spec.Completions == 1 || job.Status.Succeeded == 1 || job.Status.Failed == 1 {
				logs, err := getLogs(job.Name, namespace)
				if err != nil {
					log.Println(err)
				} else {
					if logs == "" {
						logs = "No logs found, process might've ran for more than 10 seconds."
					}
					// log.Printf("logs: %s\n", logs)
					log.Println("job completed, deleting job")
					pool := db.GetConnectionPool()
					pool.Exec(context.Background(), "UPDATE submissions SET output = $1 WHERE job_id = $2", logs, job.Name)
					deleteJob(namespace, job)
				}
			}
		}()
	}
}

func getLogs(jobName, namespace string) (string, error) {
	clientset := getK8sClient()

	podList, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("job-name=%s", jobName),
	})

	if err != nil {
		panic(err)
	}

	if len(podList.Items) == 0 {
		log.Println("No pods found for the job")
		return "", nil
	}

	podName := podList.Items[0].Name
	log.Println("pod name:", podName)

	// Get the logs of the pod
	logs, err := clientset.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{}).Stream(context.Background())
	if err != nil {
		return "", err
	}
	defer logs.Close()

	// Print the logs
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(logs)
	if err != nil {
		return "", nil
	}

	return buf.String(), nil
}

func deleteJob(namespace string, job batchv1.Job) {
	clientset := getK8sClient()
	jobsClient := clientset.BatchV1().Jobs(namespace)

	backgroundDeletion := metav1.DeletePropagationBackground

	err := jobsClient.Delete(context.Background(), job.Name, metav1.DeleteOptions{
		PropagationPolicy: &backgroundDeletion,
	})

	if err != nil {
		panic(err)
	}
}
