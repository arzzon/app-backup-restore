package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/arzzon/app-backup-restore/internal/utils/fileUtils"
	"net/http"

	"github.com/arzzon/app-backup-restore/internal/types"
)

func ApplicationDataHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		StoreAppData(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// StoreAppData stores the application data
func StoreAppData(w http.ResponseWriter, r *http.Request) {
	var app types.Application
	err := json.NewDecoder(r.Body).Decode(&app)
	if err != nil {
		fmt.Printf("[Application Handler] Error unmarshalling request body: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate app ID
	appID := fmt.Sprintf("%s-%s", app.Name, app.Namespace)

	// Store application metadata
	if err := storeApplicationMetaData(appID, app); err != nil {
		fmt.Printf("[Application Handler] Error storing application metadata: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"appId": appID}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("[Application Handler] Error encoding response: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	fmt.Printf("[Application Handler] Application data stored successfully")
}

// storeApplicationMetaData saves the application data
func storeApplicationMetaData(appID string, app types.Application) error {
	fileName := "store/apps/" + appID
	if fileUtils.CheckFile(fileName) {
		return nil
	}
	// Encode the object to JSON format
	appData, err := json.Marshal(app)
	if err != nil {
		return fmt.Errorf("Error encoding object to JSON: %v", err)
	}
	if err := fileUtils.WriteFile(fileName, appData); err != nil {
		return fmt.Errorf("Error storing application details: %v", err)
	}
	return nil
}
