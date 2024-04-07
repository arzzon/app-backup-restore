package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/arzzon/app-backup-restore/internal/constants"
	. "github.com/arzzon/app-backup-restore/internal/types"
	"github.com/arzzon/app-backup-restore/internal/utils/fileUtils"
	"github.com/arzzon/app-backup-restore/internal/utils/orchestratorClient"
	"github.com/google/uuid"
	"k8s.io/client-go/kubernetes"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"sigs.k8s.io/yaml"
)

// BackupHandler handles the backup request
func BackupHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		BackUpApplication(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// BackUpApplication backs up the application
func BackUpApplication(w http.ResponseWriter, r *http.Request) {
	var backupReq BackupRequest
	var backupResponse BackupResponse
	err := json.NewDecoder(r.Body).Decode(&backupReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate a unique backup ID
	backUpID, err := uuid.NewRandom()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the app data is saved
	appNamespace := getAppNamespace(backupReq.AppID)
	if appNamespace == "" {
		backupResponse = getBackUpResponse(backupReq.AppID, backUpID.String(), "Application data not found")
		jsonResponse, err := json.Marshal(backupResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResponse)
		return
	}

	// List of all resources to backup
	//allResources := []types.ResourceKind{types.Pod}
	//backupCompletionUpdateMutex := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	for _, resource := range AllResources {
		wg.Add(1)
		backupChan := BackupJob{
			Kind:      resource,
			BackupID:  backUpID.String(),
			Namespace: appNamespace,
			Wg:        wg,
			//completionStatusUpdateMutex: backupCompletionUpdateMutex,
		}
		BackUpWorkerPool <- backupChan
	}
	wg.Wait() // TODO: Implement timeout/asynchronous status update

	backupResponse = getBackUpResponse(backupReq.AppID, backUpID.String(), "Backup created successfully")
	jsonResponse, err := json.Marshal(backupResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// BackUpWorkerPool is the pool of workers that back up resources
var BackUpWorkerPool chan BackupJob

// BackupJob represents a backup job
type BackupJob struct {
	Kind      ResourceKind
	BackupID  string
	Namespace string
	Wg        *sync.WaitGroup
	//completionStatusUpdateMutex *sync.Mutex
}

// FetchAndStore fetches the resources and stores them in the backup directory
func (backupJob *BackupJob) FetchAndStore() []error {
	// Fetch the backup data and store it in the backup directory
	var err error
	var errorList []error
	var clientset *kubernetes.Clientset
	clientset, err = orchestratorClient.GetClientFromKubeconfig("")
	if err != nil {
		panic(err.Error())
	}
	switch backupJob.Kind {
	case Pod:
		list, err := clientset.CoreV1().Pods(backupJob.Namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		// Iterate over fetched resources
		for _, item := range list.Items {
			// Convert unstructured object to YAML
			item.APIVersion = "v1"
			item.Kind = "Pod"
			err = ParseAndStoreResource(item, item.GetName(), backupJob)
			if err != nil {
				errorList = append(errorList, err)
			}
		}
	case StatefulSet:
		list, err := clientset.AppsV1().StatefulSets(backupJob.Namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errorList = append(errorList, err)
		}
		// Iterate over fetched resources
		for _, item := range list.Items {
			// Convert unstructured object to YAML
			item.APIVersion = "apps/v1"
			item.Kind = "StatefulSet"
			err = ParseAndStoreResource(item, item.GetName(), backupJob)
			if err != nil {
				errorList = append(errorList, err)
			}
		}
	case Delpoyment:
		list, err := clientset.AppsV1().Deployments(backupJob.Namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errorList = append(errorList, err)
		}
		// Iterate over fetched resources
		for _, item := range list.Items {
			// Convert unstructured object to YAML
			item.APIVersion = "apps/v1"
			item.Kind = "Deployment"
			err = ParseAndStoreResource(item, item.GetName(), backupJob)
			if err != nil {
				errorList = append(errorList, err)
			}
		}
	case Service:
		list, err := clientset.CoreV1().Services(backupJob.Namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errorList = append(errorList, err)
		}
		// Iterate over fetched resources
		for _, item := range list.Items {
			// Convert unstructured object to YAML
			item.APIVersion = "v1"
			item.Kind = "Service"
			err = ParseAndStoreResource(item, item.GetName(), backupJob)
			if err != nil {
				errorList = append(errorList, err)
			}
		}
	case ConfigMap:
		list, err := clientset.CoreV1().ConfigMaps(backupJob.Namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errorList = append(errorList, err)
		}
		// Iterate over fetched resources
		for _, item := range list.Items {
			// Convert unstructured object to YAML
			item.APIVersion = "v1"
			item.Kind = "ConfigMap"
			err = ParseAndStoreResource(item, item.GetName(), backupJob)
			if err != nil {
				errorList = append(errorList, err)
			}
		}
	case ReplicaSet:
		list, err := clientset.AppsV1().ReplicaSets(backupJob.Namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errorList = append(errorList, err)
		}
		// Iterate over fetched resources
		for _, item := range list.Items {
			// Convert unstructured object to YAML
			item.APIVersion = "apps/v1"
			item.Kind = "ReplicaSet"
			err = ParseAndStoreResource(item, item.GetName(), backupJob)
			if err != nil {
				errorList = append(errorList, err)
			}
		}
	case PVC:
		list, err := clientset.CoreV1().PersistentVolumeClaims(backupJob.Namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errorList = append(errorList, err)
		}
		// Iterate over fetched resources
		for _, item := range list.Items {
			// Convert unstructured object to YAML
			item.APIVersion = "v1"
			item.Kind = "PersistentVolumeClaim"
			err = ParseAndStoreResource(item, item.GetName(), backupJob)
			if err != nil {
				errorList = append(errorList, err)
			}
		}
	case PV:
		list, err := clientset.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errorList = append(errorList, err)
		}
		// Iterate over fetched resources
		for _, item := range list.Items {
			// Convert unstructured object to YAML
			item.APIVersion = "v1"
			item.Kind = "PersistentVolumeClaim"
			err = ParseAndStoreResource(item, item.GetName(), backupJob)
			if err != nil {
				errorList = append(errorList, err)
			}
		}
	case ServiceAccount:
		list, err := clientset.CoreV1().ServiceAccounts(backupJob.Namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errorList = append(errorList, err)
		}
		// Iterate over fetched resources
		for _, item := range list.Items {
			// Convert unstructured object to YAML
			item.APIVersion = "v1"
			item.Kind = "ServiceAccount"
			err = ParseAndStoreResource(item, item.GetName(), backupJob)
			if err != nil {
				errorList = append(errorList, err)
			}
		}
	case Secret:
		list, err := clientset.CoreV1().Secrets(backupJob.Namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errorList = append(errorList, err)
		}
		// Iterate over fetched resources
		for _, item := range list.Items {
			// Convert unstructured object to YAML
			item.APIVersion = "v1"
			item.Kind = "Secret"
			err = ParseAndStoreResource(item, item.GetName(), backupJob)
			if err != nil {
				errorList = append(errorList, err)
			}
		}
	default:
		err := fmt.Errorf("Invalid resource type: %s", backupJob.Kind)
		errorList = append(errorList, err)
	}
	return nil

}

//func (backupJob *BackupJob) UpdateStatus() error {
//	backupJob.completionStatusUpdateMutex.Lock()
//	defer backupJob.completionStatusUpdateMutex.Unlock()
//	// Update the status of the backup job
//	return nil
//}

// Initialize the backup worker pool
func init() {
	BackUpWorkerPool = make(chan BackupJob, constants.BACKUP_JOB_POOL_SIZE)
	for i := 0; i < constants.NUM_BACKUP_WORKERS; i++ {
		go func(backupChanPool chan BackupJob) {
			for job := range backupChanPool {
				errs := job.FetchAndStore()
				if errs != nil && len(errs) > 0 {
					fmt.Printf("[Backup] Error fetching and storing %s: %v\n", job.Kind, errs)
				}
				job.Wg.Done()
				// TODO: Asynchronously update the status of the backup job
				//err := job.UpdateStatus()
				//if err != nil {
				//	fmt.Printf("Error updating status of %s: %v\n", job.Kind, err)
				//}
			}
		}(BackUpWorkerPool)
	}
}

// getAppNamespace returns the namespace of the application
// TODO: Handle errors
func getAppNamespace(appID string) string {
	// check if appID exists in the database
	filePath := fmt.Sprintf("%s/%s", constants.APPS_DIR, appID)
	if !fileUtils.CheckDirectory(filePath) {
		return ""
	}
	// read the namespace from the file
	appData, err := fileUtils.ReadFile(filePath)
	if err != nil {
		return ""
	}
	var app Application
	err = json.Unmarshal(appData, &app)
	if err != nil {
		return ""
	}
	return app.Namespace
}

func ParseAndStoreResource(item interface{}, resourceName string, backupJob *BackupJob) error {
	// Parse the resource and store it in the backup directory
	itemYAML, err := yaml.Marshal(item)
	if err != nil {
		return fmt.Errorf("Error converting %s to YAML: %v\n", backupJob.Kind, err)
	}
	// Write YAML to file
	dirPath := fmt.Sprintf("%s/%s/%s/", constants.BACKUPS_DIR, backupJob.BackupID, backupJob.Kind)
	err = fileUtils.CreateDir(dirPath)
	if err != nil {
		return fmt.Errorf("Error creating directory: %v\n", err)
	}
	fileName := fmt.Sprintf("%s.yaml", resourceName)
	err = fileUtils.WriteFile(dirPath+fileName, itemYAML)
	if err != nil {
		return fmt.Errorf("Error writing %s to file: %v\n", backupJob.Kind, err)
	}

	fmt.Printf("[Backup] Resource %s written to %s\n", resourceName, dirPath+fileName)
	return nil
}

func getBackUpResponse(appID, backUpID, msg string) BackupResponse {
	// For simplicity, just return a success message
	backupResponse := BackupResponse{
		AppID:    appID,
		BackupID: backUpID,
		Message:  msg,
	}
	return backupResponse
}
