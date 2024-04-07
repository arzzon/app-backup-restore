package handlers

//
//import (
//	"github.com/arzzon/app-backup-restore/internal/types"
//	"net/http"
//	"time"
//)
//
//func TaskHandler(w http.ResponseWriter, r *http.Request) {
//	switch r.Method {
//	case http.MethodPost:
//		StoreAppData(w, r)
//	default:
//		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//	}
//}
//
//// pollTaskStatus polls the status of the backup task and updates it accordingly
//func pollTaskStatus(task *types.Task) {
//	// Poll every second until the task is completed or failed
//	for {
//		if task.Status == "completed" || task.Status == "failed" {
//			return
//		}
//
//		// Simulate polling (sleep for 1 second)
//		time.Sleep(1 * time.Second)
//
//		// Check if the backup has taken more than 10 minutes
//		if time.Since(time.Now()) > 10*time.Minute {
//			task.Status = "failed"
//			return
//		}
//	}
//}
