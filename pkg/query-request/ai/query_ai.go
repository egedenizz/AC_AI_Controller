package ai

import (
	dbmanager "query_req/dbmanager"
	"strconv"

	"context"
	"fmt"
	"log"

	apiv2 "cloud.google.com/go/dialogflow/apiv2"
	"cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
	"google.golang.org/api/option"
)

const (
	projectID  = ""
	tempEntity = "temperature"
	modeEntity = "mode"
	keyJSON    = ""
	minVal     = 18
	maxVal     = 32
)

func QueryDialogFlow(query string, isAgain bool) string {

	ctx := context.Background()

	client, err := apiv2.NewSessionsClient(ctx, option.WithCredentialsFile(keyJSON))
	if err != nil {
		log.Fatalf("Failed to create Dialogflow client: %v", err)
	}
	defer client.Close()

	sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, dbmanager.SessionID)

	textInput := &dialogflowpb.TextInput{
		Text:         query,
		LanguageCode: "en-US",
	}
	queryInput := &dialogflowpb.QueryInput{
		Input: &dialogflowpb.QueryInput_Text{
			Text: textInput,
		},
	}

	request := &dialogflowpb.DetectIntentRequest{
		Session:    sessionPath,
		QueryInput: queryInput,
	}

	response, err := client.DetectIntent(ctx, request)
	if err != nil {
		dbmanager.AddToHistory("", "Default Fallback Intent", "", "Sorry, can you say that again?")
		return "Sorry, can you say that again?"
	}

	queryResult := response.GetQueryResult()

	return reactToQuery(queryResult, isAgain)
}

func reactToQuery(queryResult *dialogflowpb.QueryResult, isRecursive bool) string {
	var responseText string

	intentName := queryResult.GetIntent().GetDisplayName()
	parameters := queryResult.GetParameters().AsMap()
	mode := fmt.Sprint(parameters[modeEntity])
	temp := fmt.Sprint(parameters[tempEntity])
	dbmanager.CreateUserDoc()

	switch intentName {
	case "TurnOn":
		if dbmanager.GetDocumentAttribute("status") == "On" {
			responseText = "The AC is already turned on."
		} else {
			responseText = "Turning the air conditioner on."
			dbmanager.EditData("status", "On")
		}

	case "TurnOff":

		if dbmanager.GetDocumentAttribute("status") == "Off" {
			responseText = "The AC is already turned off"
		} else {
			responseText = "Turning the air conditioner off."
			dbmanager.EditData("status", "Off")
		}

	case "ChangeDegree":

		if temp == "" {
			responseText = queryResult.GetFulfillmentText()
			goto end
		} else if temp == dbmanager.GetDocumentAttribute("temp") {
			responseText = "The temperature is already " + temp + " degrees."
		} else {
			currDegInt, err := strconv.Atoi(temp)
			if err != nil {
				fmt.Printf("Error converting string to int: %v\n", err)

			}

			if currDegInt >= minVal && currDegInt <= maxVal {
				dbmanager.EditData("temp", temp)
				responseText = fmt.Sprintf("The temperature is set to %s", temp)
			} else {
				responseText = "The temperature is out of range."
			}
		}

	case "ChangeMode":
		if mode == "" {
			responseText = queryResult.GetFulfillmentText()
			goto end
		} else if mode == dbmanager.GetDocumentAttribute("mode") {
			responseText = "The mode is already " + mode + " mode."
		} else {

			dbmanager.EditData("mode", mode)
			responseText = "Changed the mode to " + mode + "."
		}

	case "Again":

		newQuery, _ := dbmanager.GetLastHistoryInputAndIntent()

		if newQuery == "" {
			responseText = "There is no action that I can repeat"
		} else {
			responseText = QueryDialogFlow(newQuery, false)

		}

	case "Decrease":

		degree, err := strconv.Atoi(dbmanager.GetDocumentAttribute("temp"))
		if err != nil {
			fmt.Printf("Error converting string to int: %v\n", err)

		}

		if temp == "" {
			if degree-2 >= minVal {
				newVal := degree - 2
				dbmanager.EditData("temp", fmt.Sprint(newVal))
				responseText = "Decreased the temperature by 2 degrees"

			} else {
				responseText = "Cannot Decrease the temperature outside the range of the air conditioner."
			}

		} else {
			adjustment := int(parameters[tempEntity].(float64))
			if degree-adjustment >= minVal {
				newVal := degree - adjustment
				dbmanager.EditData("temp", fmt.Sprint(newVal))
				responseText = fmt.Sprintf("Decreased the temperature by %d degrees", adjustment)
			} else {
				responseText = "Cannot Decrease the temperature outside the range of the air conditioner."
			}

		}
	case "Increase":

		degree, err := strconv.Atoi(dbmanager.GetDocumentAttribute("temp"))
		if err != nil {
			fmt.Printf("Error converting string to int: %v\n", err)

		}

		if temp == "" {
			if degree+2 <= maxVal {
				newVal := degree + 2
				dbmanager.EditData("temp", fmt.Sprint(newVal))
				responseText = "Increased the temperature by 2 degrees"

			} else {
				responseText = "Cannot Increase the temperature outside the range of the air conditioner."
			}

		} else {
			adjustment := int(parameters[tempEntity].(float64))
			if degree+adjustment <= maxVal {
				newVal := degree + adjustment
				dbmanager.EditData("temp", fmt.Sprint(newVal))
				responseText = fmt.Sprintf("Increased the temperature by %d degrees", adjustment)
			} else {
				responseText = "Cannot Increase the temperature outside the range of the air conditioner."
			}

		}
	default:
		responseText = queryResult.GetFulfillmentText()
	}

end:
	if isRecursive {
		dbmanager.AddToHistory(queryResult.QueryText, intentName, fmt.Sprint(queryResult.GetParameters().AsMap()), responseText)
	}

	return responseText

}
