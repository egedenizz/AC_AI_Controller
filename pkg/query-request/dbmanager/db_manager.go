package db

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	projectID    = ""
	statusEntity = "status"
	tempEntity   = "temp"
	modeEntity   = "mode"
	dbJSON       = ""
)

var SessionID = ""
var UserID = ""

type AirConditioner struct {
	Status string `firestore:"status"`
	Mode   string `firestore:"mode"`
	Temp   string `firestore:"temp"`
}

type Date struct {
	Date string `firestore:"time-code"`
}

type queryRecord struct {
	Date       time.Time `firestore:"date"`
	Query      string    `firestore:"query"`
	Intent     string    `firestore:"intent"`
	Parameters string    `firestore:"parameters"`
	Response   string    `firestore:"response"`
	Status     string    `firestore:"status"`
	Mode       string    `firestore:"mode"`
	Temp       string    `firestore:"temp"`
}

func GetGeneralHistory() string {
	result := ""
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(dbJSON))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	docIter := client.Collection("AC").Doc(UserID).Collection("sessions").Documents(ctx)

	for {
		doc, err := docIter.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}

		if doc == nil {
			break
		}
		histIter := doc.Ref.Collection("history").Documents(ctx)
		for {
			history, err := histIter.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				break
			}

			if history == nil {
				break
			}

			result += fmt.Sprintf("Document ID: %s\nDate: %v\nIntent: %s\nQuery: %s\nMode: %s\nStatus: %s\nTemp: %s\nParameters: %s\n------------------------------\n", history.Ref.ID, history.Data()["date"], history.Data()["intent"], history.Data()["query"], history.Data()["mode"], history.Data()["status"], history.Data()["temp"], history.Data()["parameters"])
		}

	}

	return result
}

func GetSessionHistory() string {
	result := ""
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(dbJSON))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()
	docRef := client.Collection("AC").Doc(UserID).Collection("sessions").Doc(SessionID).Collection("history").Documents(ctx)
	for {
		doc, err := docRef.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}

		if doc == nil {
			break
		}

		result += fmt.Sprintf("Document ID: %s\nDate: %v\nIntent: %s\nQuery: %s\nMode: %s\nStatus: %s\nTemp: %s\nParameters: %s\n------------------------------\n", doc.Ref.ID, doc.Data()["date"], doc.Data()["intent"], doc.Data()["query"], doc.Data()["mode"], doc.Data()["status"], doc.Data()["temp"], doc.Data()["parameters"])
	}

	return result
}

func GetLastHistoryInputAndIntent() (string, error) {

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(dbJSON))
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	historiesRef := client.Collection("AC").Doc(UserID).Collection("sessions").Doc(SessionID).Collection("history")

	query := historiesRef.OrderBy("date", firestore.Desc)
	iter := query.Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err != nil {
			if err == iterator.Done {

				return "", nil
			}
			return "", nil
		}

		intentName, ok := doc.Data()["intent"].(string)
		if !ok {
			return "", errors.New("field 'intentName' not found or is not a string")
		}

		if intentName != "Again" {
			input, ok := doc.Data()["query"].(string)
			if !ok {
				return "", errors.New("field 'input' not found or is not a string")
			}
			return input, nil
		}
	}

}
func GetDocumentAttribute(attributeName string) string {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(dbJSON))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()
	docRef := client.Collection("AC").Doc(UserID)
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		log.Fatalf("Failed to get document: %v", err)
	}
	if !docSnap.Exists() {
		log.Println("No such document!")

	}

	attributeValue, ok := docSnap.Data()[attributeName]
	if !ok {
		log.Printf("Attribute %s not found", attributeName)

	}
	return fmt.Sprint(attributeValue)

}

func AddToHistory(query string, intent string, parameters string, response string) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(dbJSON))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()
	sessionDoc := client.Collection("AC").Doc(UserID).Collection("sessions").Doc(SessionID)

	record := queryRecord{
		Date:       time.Now(),
		Query:      query,
		Intent:     intent,
		Parameters: parameters,
		Response:   response,
		Status:     GetDocumentAttribute(statusEntity),
		Temp:       GetDocumentAttribute(tempEntity),
		Mode:       GetDocumentAttribute(modeEntity),
	}
	subColRef := sessionDoc.Collection("history")
	subColDocRef := subColRef.Doc(time.Now().Format("20060102_150405"))

	_, err = subColDocRef.Set(ctx, record)
	if err != nil {
		log.Fatalf("Failed to add record to subcollection: %v", err)
	}

	log.Println("Subcollection document created successfully")
	if err != nil {
		log.Fatalf("Failed to create placeholder document: %v", err)
	}

}

func doesUserExist() bool {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(dbJSON))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()
	docRef := client.Collection("AC").Doc(UserID)

	doc, err := docRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return false
		} else {

			log.Fatalf("Failed to get document: %v", err)
		}
	} else {
		fmt.Println("Document exists:", doc.Ref.ID)
	}
	return true
}

func CreateUserDoc() {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(dbJSON))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	if !doesUserExist() {
		docRef := client.Collection("AC").Doc(UserID)
		ac := AirConditioner{
			Status: "Off",
			Mode:   "Auto",
			Temp:   "24",
		}

		_, err = docRef.Set(ctx, ac)
		if err != nil {
			log.Fatalf("Failed to create document: %v", err)
		}
	}

	colRef := client.Collection("AC").Doc(UserID).Collection("sessions").Doc(SessionID)
	date := Date{
		Date: fmt.Sprint(time.Now().Format("20060102_150405")),
	}

	_, err = colRef.Set(ctx, date)
	if err != nil {
		log.Fatalf("Failed to create document: %v", err)
	}

	log.Println("Placeholder document created in 'sessions' collection.")

	log.Println("Subcollection created with a placeholder document")

	log.Println("Document created successfully")
}

func EditData(key string, value string) {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(dbJSON))
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	docRef := client.Collection("AC").Doc(UserID)

	doc, err := docRef.Get(ctx)
	if err != nil {
		log.Fatalf("Failed to get document: %v", err)
	}

	if !doc.Exists() {
		log.Fatalf("No such document!")
	}

	updates := []firestore.Update{
		{Path: key, Value: value},
	}

	_, err = docRef.Update(ctx, updates)
	if err != nil {
		log.Fatalf("Failed to update document: %v", err)
	}

	fmt.Println("Document updated successfully.")

}
