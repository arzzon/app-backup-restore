package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/arzzon/app-backup-restore/internal/constants"
	. "github.com/arzzon/app-backup-restore/internal/types"
	"github.com/arzzon/app-backup-restore/internal/utils/fileUtils"
	"github.com/arzzon/app-backup-restore/internal/utils/orchestratorClient"
	v12 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"net/http"
	"os"
	"path/filepath"
	"sigs.k8s.io/yaml"
)

// Order of restoring backups
var restoreOrder = []ResourceKind{ServiceAccount, Secret, ConfigMap, PV, PVC, StatefulSet, Delpoyment, ReplicaSet, Pod, Service}

// RestoreBackupHandler handles the restore backup request
func RestoreBackupHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		RestoreBackup(w, r)
	case http.MethodGet:

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// RestoreBackup restores the backup into the given namespace
func RestoreBackup(w http.ResponseWriter, r *http.Request) {
	var restoreReq RestoreRequest
	err := json.NewDecoder(r.Body).Decode(&restoreReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the backup exists
	if !checkIfBackupStored(restoreReq.BackupID) {
		restoreResponse := getRestoreResponse(restoreReq.Namespace, restoreReq.BackupID, "Backup not found")
		jsonResponse, err := json.Marshal(restoreResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonResponse)
		return
	}

	// Restore resources in the order specified
	for _, resourceKind := range restoreOrder {
		fmt.Println("[Restore] Restoring resource: ", resourceKind)
		err := parseAndRestore(restoreReq.BackupID, restoreReq.Namespace, resourceKind)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Return response
	restoreResponse := getRestoreResponse(restoreReq.Namespace, restoreReq.BackupID, "Backup restored successfully")
	jsonResponse, err := json.Marshal(restoreResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusOK)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// checkIfBackupStored checks if the backup is stored in the backups directory
func checkIfBackupStored(backupID string) bool {
	dirPath := fmt.Sprintf("%s/%s", constants.BACKUPS_DIR, backupID)
	// Check if the backup exists
	return fileUtils.CheckDirectory(dirPath)
}

// parseAndRestore parses the YAML files and restores the resources
func parseAndRestore(backupID, namespace string, resourceKind ResourceKind) error {
	yamlDir := fmt.Sprintf("%s/%s/%s", constants.BACKUPS_DIR, backupID, resourceKind)
	if !fileUtils.CheckDirectory(yamlDir) {
		return nil
	}

	// Get list of YAML files in directory
	files, err := os.ReadDir(yamlDir)
	if err != nil {
		return err
	}

	var clientset *kubernetes.Clientset
	clientset, err = orchestratorClient.GetClientFromKubeconfig("")
	if err != nil {
		panic(err.Error())
	}

	// Iterate over YAML files
	for _, file := range files {
		filePath := filepath.Join(yamlDir, file.Name())

		// Read YAML file
		yamlDataBytes, err := fileUtils.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("error reading YAML file %s: %v", filePath, err)
		}

		// Parse YAML
		switch resourceKind {
		case Pod:
			pod := &v1.Pod{}
			err := yaml.Unmarshal(yamlDataBytes, pod)
			if err != nil {
				return fmt.Errorf("error unmarshalling Pod YAML: %v", err)
			}
			pod.ResourceVersion = ""
			_, err = clientset.CoreV1().Pods(namespace).Create(context.Background(), pod, metav1.CreateOptions{})
			if err != nil && !errors.IsAlreadyExists(err) {
				return fmt.Errorf("error creating Pod: %v", err)
			}
		case StatefulSet:
			ss := &v12.StatefulSet{}
			err := yaml.Unmarshal(yamlDataBytes, ss)
			if err != nil {
				return fmt.Errorf("error unmarshalling Pod StatefulSet: %v", err)
			}
			ss.ResourceVersion = ""
			_, err = clientset.AppsV1().StatefulSets(namespace).Create(context.Background(), ss, metav1.CreateOptions{})
			if err != nil && !errors.IsAlreadyExists(err) {
				return fmt.Errorf("error creating StatefulSet: %v", err)
			}
		case Delpoyment:
			deployment := &v12.Deployment{}
			err = yaml.Unmarshal(yamlDataBytes, deployment)
			if err != nil {
				return fmt.Errorf("error unmarshalling Pod YAML: %v", err)
			}
			deployment.ResourceVersion = ""
			_, err = clientset.AppsV1().Deployments(namespace).Create(context.Background(), deployment, metav1.CreateOptions{})
			if err != nil && !errors.IsAlreadyExists(err) {
				return fmt.Errorf("error restoring Deployment: %v", err)
			}
		case Service:
			svc := &v1.Service{}
			err = yaml.Unmarshal(yamlDataBytes, svc)
			if err != nil {
				return fmt.Errorf("error unmarshalling Service YAML: %v", err)
			}
			svc.ResourceVersion = ""
			_, err = clientset.CoreV1().Services(namespace).Create(context.Background(), svc, metav1.CreateOptions{})
			if err != nil && !errors.IsAlreadyExists(err) {
				return fmt.Errorf("error creating Service: %v", err)
			}
		case ConfigMap:
			cm := &v1.ConfigMap{}
			err = yaml.Unmarshal(yamlDataBytes, cm)
			if err != nil {
				return fmt.Errorf("error unmarshalling ConfigMap YAML: %v", err)
			}
			cm.ResourceVersion = ""
			_, err = clientset.CoreV1().ConfigMaps(namespace).Create(context.Background(), cm, metav1.CreateOptions{})
			if err != nil && !errors.IsAlreadyExists(err) {
				return fmt.Errorf("error creating ConfigMap: %v", err)
			}
		case ReplicaSet:
			rs := &v12.ReplicaSet{}
			err = yaml.Unmarshal(yamlDataBytes, rs)
			if err != nil {
				return fmt.Errorf("error unmarshalling ReplicaSet YAML: %v", err)
			}
			rs.ResourceVersion = ""
			_, err = clientset.AppsV1().ReplicaSets(namespace).Create(context.Background(), rs, metav1.CreateOptions{})
			if err != nil && !errors.IsAlreadyExists(err) {
				return fmt.Errorf("error restoring ReplicaSet: %v", err)
			}
		case PV:
			pv := &v1.PersistentVolume{}
			err := yaml.Unmarshal(yamlDataBytes, pv)
			if err != nil {
				return fmt.Errorf("error unmarshalling PV YAML: %v", err)
			}
			pv.ResourceVersion = ""
			_, err = clientset.CoreV1().PersistentVolumes().Create(context.Background(), pv, metav1.CreateOptions{})
			if err != nil && !errors.IsAlreadyExists(err) {
				return fmt.Errorf("error creating PersistentVolume: %v", err)
			}
		case PVC:
			pvc := &v1.PersistentVolumeClaim{}
			err := yaml.Unmarshal(yamlDataBytes, pvc)
			if err != nil {
				return fmt.Errorf("error unmarshalling PVC YAML: %v", err)
			}
			pvc.ResourceVersion = ""
			_, err = clientset.CoreV1().PersistentVolumeClaims(namespace).Create(context.Background(), pvc, metav1.CreateOptions{})
			if err != nil && !errors.IsAlreadyExists(err) {
				return fmt.Errorf("error creating PersistentVolumeClaims: %v", err)
			}
		case ServiceAccount:
			sa := &v1.ServiceAccount{}
			err := yaml.Unmarshal(yamlDataBytes, sa)
			if err != nil {
				return fmt.Errorf("error unmarshalling ServiceAccount YAML: %v", err)
			}
			sa.ResourceVersion = ""
			_, err = clientset.CoreV1().ServiceAccounts(namespace).Create(context.Background(), sa, metav1.CreateOptions{})
			if err != nil && !errors.IsAlreadyExists(err) {
				return fmt.Errorf("error creating ServiceAccount: %v", err)
			}
		case Secret:
			secret := &v1.Secret{}
			err := yaml.Unmarshal(yamlDataBytes, secret)
			if err != nil {
				return fmt.Errorf("error unmarshalling Secret YAML: %v", err)
			}
			secret.ResourceVersion = ""
			_, err = clientset.CoreV1().Secrets(namespace).Create(context.Background(), secret, metav1.CreateOptions{})
			if err != nil && !errors.IsAlreadyExists(err) {
				return fmt.Errorf("error creating Secret: %v", err)
			}
		default:
			fmt.Printf("[Restore] Invalid resource type: %s", resourceKind)
		}

	}
	return nil
}

// getRestoreResponse returns the restore response
func getRestoreResponse(namespace, backUpID, msg string) RestoreResponse {
	// For simplicity, just return a success message
	restoreResponse := RestoreResponse{
		Namespace: namespace,
		BackupID:  backUpID,
		Message:   msg,
	}
	return restoreResponse
}
