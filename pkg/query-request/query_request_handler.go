package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	agent "query_req/agent"
	ai "query_req/ai"
	dbmanager "query_req/dbmanager"
	speech "query_req/speech"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/data", dataHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/agent", InitializeAgent)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Failed to start server:", err)
		os.Exit(1)
	}
}

func InitializeAgent(w http.ResponseWriter, r *http.Request) {

	for _, entity := range agent.Entities {
		agent.CreateEntityType(entity.EntityDisplayName, entity.EntityAttributes)
	}

	for _, intent := range agent.Intents {

		agent.CreateIntent(intent.IntentName)
		agent.AddPhrases(intent, intent.TrainingPhrases)
		agent.ClearParams(intent)

		for _, entity := range intent.Entities {
			agent.AddEntityToIntent(intent, entity.EntityDisplayName, entity.EntityTypeName, entity.EntityValue)
		}
	}
	agent.TrainAgent()

}

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	byteData := r.FormValue("byte_data")
	if byteData == "" {
		http.Error(w, "No byte data found", http.StatusBadRequest)
		return
	}

	dbmanager.SessionID = r.FormValue("sessionID")
	dbmanager.UserID = r.FormValue("userID")
	responseHandler(w, "success", ai.QueryDialogFlow(speech.TranscribeAudio([]byte(byteData)), true))

}

func dataHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	result := handleRequest(data)

	responseHandler(w, "success", result)
}

func responseHandler(w http.ResponseWriter, status string, message string) {
	w.Header().Set("Content-Type", "application/json")
	data := Response{
		Status:  status,
		Message: message,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	w.Write(jsonData)
}

func handleRequest(data map[string]string) string {
	dbmanager.SessionID = data["sessionID"]
	dbmanager.UserID = data["userID"]
	request := data["request"]

	switch request {
	case "record":

		return ai.QueryDialogFlow(speech.TranscribeAudio([]byte(data["buffer"])), true)
	case "history":
		return dbmanager.GetGeneralHistory()
	case "sessionhistory":
		return dbmanager.GetSessionHistory()
	case "write":
		return ai.QueryDialogFlow(data["query"], true)
	default:
		return "Unknown request"
	}
}
