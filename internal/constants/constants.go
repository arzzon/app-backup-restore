package constants

const (
	// Kubernetes
	KUBECONFIG_PATH = "~/.kube/config"

	// store
	//STORE_DIR   = "store"
	APPS_DIR    = "store/apps"
	BACKUPS_DIR = "store/backups"

	// tasks
	//TASK_POLL_INTERVAL = 2
	//TASK_TIMEOUT       = 60

	// worker pool
	NUM_BACKUP_WORKERS   = 10
	BACKUP_JOB_POOL_SIZE = 100
)
