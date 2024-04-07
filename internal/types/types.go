package types

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"time"
)

type Application struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type BackupRequest struct {
	AppID string `json:"app"`
}

type BackupResponse struct {
	AppID    string `json:"app"`
	BackupID string `json:"backupId"`
	Message  string `json:"message"`
}

type RestoreRequest struct {
	Namespace string `json:"namespace"`
	BackupID  string `json:"backupId"`
}

type RestoreResponse struct {
	Namespace string `json:"namespace"`
	BackupID  string `json:"backupId"`
	Message   string `json:"message"`
}

// enums for task status
type TaskStatus string

const (
	Completed  TaskStatus = "completed"
	Failed     TaskStatus = "failed"
	InProgress TaskStatus = "in-progress"
)

type Task struct {
	ID              string
	Status          TaskStatus
	StartTime       timestamp.Timestamp
	TimeoutDuration time.Duration
}

// Resource kind

type ResourceKind string

const (
	Pod            ResourceKind = "Pod"
	Delpoyment     ResourceKind = "Deployment"
	ReplicaSet     ResourceKind = "ReplicaSet"
	StatefulSet    ResourceKind = "StatefulSet"
	Service        ResourceKind = "Service"
	Secret         ResourceKind = "Secret"
	ConfigMap      ResourceKind = "ConfigMap"
	PV             ResourceKind = "PV"
	PVC            ResourceKind = "PVC"
	ServiceAccount ResourceKind = "ServiceAccount"
)

type BackupChan struct {
	BackupID  string
	Namespace string
	Kind      ResourceKind
}

var AllResources = []ResourceKind{Pod, Delpoyment, StatefulSet, Service, Secret, ConfigMap, ReplicaSet, PV, PVC, ServiceAccount}
