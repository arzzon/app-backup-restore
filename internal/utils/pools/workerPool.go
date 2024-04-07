package pools

import (
	"github.com/arzzon/app-backup-restore/internal/constants"
	"sync"
)

var BackUpWorkerPool chan BackupJob

type BackupJob struct {
	Kind                        string
	BackupID                    string
	Namespace                   string
	completionStatusUpdateMutex *sync.Mutex
}

func (*BackupJob) FetchAndStore() {
	// Fetch the backup data and store it in the backup directory
}

func (backupJob *BackupJob) UpdateStatus() {
	backupJob.completionStatusUpdateMutex.Lock()
	defer backupJob.completionStatusUpdateMutex.Unlock()
	// Update the status of the backup job
}

func init() {
	BackUpWorkerPool = make(chan BackupJob, constants.BACKUP_JOB_POOL_SIZE)
	for i := 0; i < constants.NUM_BACKUP_WORKERS; i++ {
		go func(backupChanPool chan BackupJob) {
			for job := range backupChanPool {
				job.FetchAndStore()
				job.UpdateStatus()
			}
		}(BackUpWorkerPool)
	}
}
