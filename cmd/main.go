package main

import (
	"fmt"
	"github.com/arzzon/app-backup-restore/internal/constants"
	"github.com/arzzon/app-backup-restore/internal/handlers"
	"github.com/arzzon/app-backup-restore/internal/utils/fileUtils"
	"log"
	"net/http"
)

func init() {
	fmt.Println("Initializing application...")
	// Creating store directory if it does not exist
	if fileUtils.CreateDir(constants.APPS_DIR) != nil {
		log.Fatal("Error creating store")
	}
	if fileUtils.CreateDir(constants.BACKUPS_DIR) != nil {
		log.Fatal("Error creating store")
	}
}

func main() {
	http.HandleFunc("/application/", handlers.ApplicationDataHandler)
	http.HandleFunc("/backup/", handlers.BackupHandler)
	http.HandleFunc("/restore/", handlers.RestoreBackupHandler)

	fmt.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
