package agent

import (
	"context"
	"fmt"
	"log"
	"query_req/stringsman"

	apiv2 "cloud.google.com/go/dialogflow/apiv2"

	"cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
	"google.golang.org/api/option"

	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

const (
	projectID = ""
	keyJSON   = ""
)

var (
	mode = Entity{
		EntityDisplayName: "mode",
		EntityValue:       "$mode",
		EntityTypeName:    "@mode",
		EntityAttributes: map[string][]string{
			"cooling": {"cool", "cooling", "cool mode", "chill", "cold", "cooling mode"},
			"heating": {"heat", "heating", "warm", "hot mode", "warm mode", "heat up", "heating mode"},
			"auto":    {"auto", "automatic", "auto mode"},
			"sleep":   {"sleep", "sleeping", "night", "quiet mode", "rest mode"},
		},
	}

	adjust = Entity{
		EntityDisplayName: "adjust",
		EntityValue:       "$adjust",
		EntityTypeName:    "@adjust",
		EntityAttributes: map[string][]string{
			"change": {"change", "adjust", "modify", "alter", "tweak", "set", "make", "switch", "put"},
		},
	}

	temperature = Entity{
		EntityDisplayName: "temperature",
		EntityValue:       "$temperature",
		EntityTypeName:    "@sys.number",
	}

	turnOn = Entity{
		EntityDisplayName: "turnon",
		EntityValue:       "$turnon",
		EntityTypeName:    "@turnon",
		EntityAttributes: map[string][]string{
			"turn-on": {"power on", "start", "turn on", "activate", "enable", "power up", "switch on", "engage", "enable", "fire up", "going", "initiate", "kick on"},
		},
	}

	turnOff = Entity{
		EntityDisplayName: "turnoff",
		EntityValue:       "$turnoff",
		EntityTypeName:    "@turnoff",
		EntityAttributes: map[string][]string{
			"turn-off": {"power off", "shut down", "turn off", "deactivate", "disable", "switch off", "stop", "power down", " off", "cut off", "shut down", "shut"},
		},
	}

	again = Entity{
		EntityDisplayName: "again",
		EntityValue:       "$again",
		EntityTypeName:    "@again",
		EntityAttributes: map[string][]string{
			"repeat": {"again", "repeat", "do that again", "once more"},
		},
	}

	increase = Entity{
		EntityDisplayName: "increase",
		EntityValue:       "$increase",
		EntityTypeName:    "@increase",
		EntityAttributes: map[string][]string{
			"increase": {"increase", "raise", "boost", "turn up", "warmer", "hotter", "heat the", "rise", "higher", "turn up", " up"},
		},
	}

	decrease = Entity{
		EntityDisplayName: "decrease",
		EntityValue:       "$decrease",
		EntityTypeName:    "@decrease",
		EntityAttributes: map[string][]string{
			"decrease": {"decrease", "lower", "reduce", "cooler", "turn down", "drop", "cool down", "cool the", "down"},
		},
	}

	decreaseTemp = Intent{
		IntentName: "Decrease",
		TrainingPhrases: []string{
			"Cool the room by 7 degrees",
			"Reduce by 6 degrees",
			"Cool down by 5 degrees",
			"Drop the temperature by 5 degrees",
			"Drop the temperature by 4 degrees",
			"Lower the temperature by 3 degrees",
			"Decrease the AC to 20 degrees",
			"Turn down the thermostat",
			"Make it cooler",
			"Reduce the cooling level",
			"Set the temperature lower",
			"Bring down the temperature by 2 degrees",
			"Cool the room down by 3 degrees",
			"Cool the air by 4 degrees",
			"Decrease the temperature",
			"Lower the temperature by 3 degrees",
			"Turn down the heat",
			"Set the temperature lower",
			"Make it cooler",
			"Reduce the temperature",
			"Lower the AC setting",
			"Adjust the temperature down",
			"Make the room cooler",
			"Turn the temperature down by 2 degrees",
			"Set it to a lower temperature",
			"Decrease the heat by 4 degrees",
			"Lower the air conditioning temperature",
			"Adjust the AC to a cooler setting",
			"Decrease the temperature on the AC",
			"Make it 22 degrees",
			"Set the air conditioner to a cooler temperature",
			"Reduce it to 20 degrees",
			"Turn down the temperature on the AC",
			"Decrease the cooling level",
			"Lower the thermostat setting",
		},
		Entities:         []Entity{decrease, temperature},
		RequiredEntities: []Entity{temperature},
	}

	increaseTemp = Intent{
		IntentName: "Increase",
		TrainingPhrases: []string{
			"Make it 4 degrees hotter",
			"Heat the room by 5 degrees",
			"Set the temperature to rise by 3 degrees",
			"Heat the room by 3 degrees",
			"Raise the temperature",
			"Increase the AC temperature",
			"Turn up the heat by 5 degrees",
			"Make it warmer",
			"Boost the temperature",
			"Set it to a higher temperature",
			"Make it warmer by 2 degrees",
			"Turn the temperature up by 4 degrees",
			"Heat up by 4 degrees",
			"Warm up the AC by 5 degrees",
			"Increase the temperature",
			"Raise the temperature by 5 degrees",
			"Turn up the heat",
			"Set the temperature higher",
			"Make it warmer",
			"Boost the temperature",
			"Raise the AC temperature",
			"Increase the thermostat setting",
			"Make the room warmer",
			"Turn the temperature up by 3 degrees",
			"Set it to a higher temperature",
			"Increase the heat by 2 degrees",
			"Raise the air conditioning temperature",
			"Adjust the AC to a warmer setting",
			"Increase the temperature on the AC",
			"Make it 24 degrees",
			"Set the air conditioner to a warmer temperature",
			"Raise it to 26 degrees",
			"Turn up the temperature on the AC",
			"Increase the cooling setting",
		},
		Entities:         []Entity{increase, temperature},
		RequiredEntities: []Entity{temperature},
	}

	turnOnIntent = Intent{
		IntentName: "TurnOn",
		TrainingPhrases: []string{
			"Turn on the AC",
			"Activate the air conditioner",
			"Power up the AC",
			"Switch on the air conditioning",
			"Start the AC at 24 degrees",
			"Can you turn on the AC",
			"Fire up the air conditioner",
			"Get the cooling system going",
			"Engage the climate control unit",
			"Fire up the climate control",
			"Turn the air conditioning unit to on",
			"Switch the cooling unit on",
			"Turn the climate control system on",
			"Power the AC on",
			"Turn the cooling system on",
			"Get the air conditioner working",
			"Turn the air conditioner on now",
			"Engage the climate control system",
			"Initiate the cooling system",
			"Power on the climate control system",
			"Begin the air conditioning",
			"Get the cooling fan started",
			"Turn on the AC",
			"Switch on the air conditioner",
			"Power up the AC",
			"Activate the air conditioning",
			"Start the air conditioner",
			"Can you turn on the AC?",
			"Please turn on the air conditioning",
			"Switch on the cooling system",
			"Enable the air conditioner",
			"Turn on the cooling unit",
			"Power on the air conditioning",
			"Start the AC unit",
			"Make sure the air conditioner is on",
			"Turn on the air conditioning unit",
			"Activate the cooling system",
			"Turn on the climate control",
			"Start the air conditioner now",
			"Switch the AC on",
			"Get the air conditioner running",
			"Turn the cooling system on",
		},
		Entities:         []Entity{turnOn, mode, temperature},
		RequiredEntities: []Entity{turnOn},
	}

	turnOffIntent = Intent{
		IntentName: "TurnOff",
		TrainingPhrases: []string{
			"Stop the air control system",
			"Stop the climate control system now",
			"Stop the cooling unit now",
			"Stop the HVAC unit immediately",
			"Turn off the AC",
			"Switch off the air conditioner",
			"Power down the air conditioning",
			"Deactivate the AC",
			"Stop the air conditioner",
			"Can you turn off the AC?",
			"Please turn off the air conditioning",
			"Shut down the air conditioner",
			"Disable the AC",
			"Turn off the air conditioning unit",
			"Stop the AC immediately",
			"Cut off the air conditioning",
			"Turn off the cooling system",
			"Turn off the AC unit",
			"Make sure the air conditioner is off",
			"Stop the AC from running",
			"Deactivate the cooling",
			"Switch off the AC",
			"Turn off the climate control",
			"Power off the air conditioner",
			"Turn off the AC system",
			"Cut off the air conditioning",
			"Turn off the climate control",
			"Make sure the AC is off",
			"Deactivate the cooling system",
			"Switch off the air conditioning",
			"Stop the air conditioning unit",
			"Shut off the AC",
			"Turn the AC off",
			"Turn off the AC system",
			"End the air conditioning",
			"Stop the AC immediately",
			"Turn off the air conditioner completely",
			"Make the AC go off",
			"Switch off the cooling system",
			"Power down the air conditioning unit",
		},
		Entities:         []Entity{turnOff},
		RequiredEntities: []Entity{turnOff},
	}

	changeModeIntent = Intent{
		IntentName: "ChangeMode",
		TrainingPhrases: []string{
			"Switch the AC to cooling mode",
			"Set the air conditioner to cooling",
			"Switch to cooling",
			"Put the AC in cooling mode",
			"Switch AC to cooling mode",
			"Set the air conditioner to cooling",
			"Set the mode to cooling",
			"Set the mode to heat",
			"Switch ac to cooling mode",
			"Switch to cooling mode",
			"Change the AC to heating mode",
			"Set to sleep mode",
			"Adjust the air conditioner to auto mode",
			"Change to cooling mode",
			"Switch the air conditioner to heating mode",
			"Set the AC to auto mode",
			"Change the mode to sleep",
			"Switch to cooling",
			"Set the air conditioning to auto",
			"Change to heating mode",
			"Switch to sleep mode",
			"Change the AC mode to cooling",
			"Set the air conditioner to heating",
			"Switch the mode to auto",
			"Change the mode to cooling",
			"Set the air conditioning to sleep mode",
			"Change to auto mode",
			"Switch the air conditioner to cooling mode",
			"Change the AC to sleep mode",
			"Set the mode to heating",
			"Switch to auto mode",
			"Change the air conditioner to auto mode",
			"Set it to cooling mode",
			"Change the air conditioning to heating mode",
		},
		Entities:         []Entity{mode, adjust},
		RequiredEntities: []Entity{mode},
	}

	changeTemperatureIntent = Intent{
		IntentName: "ChangeDegree",
		TrainingPhrases: []string{
			"Set the AC to 32 degrees",
			"Change the AC to 22 degrees",
			"Set the temperature to 24 degrees",
			"Adjust the AC to 22 degrees",
			"Make it 25 degrees",
			"Set the air conditioning to 26 degrees",
			"Change the temperature to 21 degrees",
			"Set the thermostat to 23 degrees",
			"Make the room 27 degrees",
			"Adjust the temperature to 20 degrees",
			"Set it to 28 degrees",
			"Change the air conditioner to 22 degrees",
			"Set the cooling system to 24 degrees",
			"Adjust the temperature to 19 degrees",
			"Make it 26 degrees",
			"Set the AC to 23 degrees",
			"Change the temperature on the air conditioner to 21 degrees",
			"Set the temperature to 25 degrees",
			"Adjust it to 24 degrees",
			"Change the setting to 22 degrees",
			"Make the AC 25 degrees",
			"Set the air conditioning to 27 degrees",
			"Adjust the AC to 20 degrees",
			"Change the room temperature to 26 degrees",
			"Set the air conditioner to 24 degrees",
			"Make it 23 degrees on the AC",
			"Adjust the thermostat to 28 degrees",
		},
		Entities:         []Entity{temperature, adjust},
		RequiredEntities: []Entity{temperature},
	}

	RepeatIntent = Intent{
		IntentName: "Again",
		TrainingPhrases: []string{
			"Do that again",
			"Repeat the last action",
			"Execute the previous command",
			"Run the last command again",
		},
		Entities:         []Entity{again},
		RequiredEntities: []Entity{again},
	}

	Intents  = []Intent{decreaseTemp, increaseTemp, turnOnIntent, turnOffIntent, changeModeIntent, changeTemperatureIntent, RepeatIntent}
	Entities = []Entity{mode, temperature, turnOn, turnOff, adjust, again, increase, decrease}
)

type Intent struct {
	IntentName       string
	TrainingPhrases  []string
	Entities         []Entity
	RequiredEntities []Entity
}

type Entity struct {
	EntityTypeName    string
	EntityDisplayName string
	EntityValue       string
	EntityAttributes  map[string][]string
}

func CreateIntent(intentName string) {
	ctx := context.Background()

	client, err := apiv2.NewIntentsClient(ctx, option.WithCredentialsFile(keyJSON))
	if err != nil {
		log.Fatalf("Failed to create Dialogflow client: %v", err)
	}
	defer client.Close()

	req := &dialogflowpb.CreateIntentRequest{
		Parent: "projects/" + projectID + "/agent",
		Intent: &dialogflowpb.Intent{
			DisplayName: intentName,
		},
	}

	resp, err := client.CreateIntent(ctx, req)
	if err != nil {
		fmt.Println("Intent already exists.")
		return
	}
	log.Printf("Created intent: %v", resp)

}

func getWordsFromEntities(intent Intent) []string {
	var word []string
	for _, entity := range intent.Entities {
		for _, attributes := range entity.EntityAttributes {
			word = append(word, attributes...)
		}

	}
	return word
}

func TrainAgent() {

	ctx := context.Background()

	client, err := apiv2.NewAgentsClient(ctx, option.WithCredentialsFile(keyJSON))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	defer client.Close()

	client.TrainAgent(ctx, &dialogflowpb.TrainAgentRequest{
		Parent: "projects/" + projectID + "/agent",
	})

}

func combineTrainingPhraseParts(results []string, intent Intent) []*dialogflowpb.Intent_TrainingPhrase {
	var parts []*dialogflowpb.Intent_TrainingPhrase_Part
	for _, str := range results {
		flag := false
		var Name string
		var typeName string
		for _, entity := range intent.Entities {
			for _, attributes := range entity.EntityAttributes {

				for _, attribute := range attributes {
					if str == attribute {
						flag = true
						Name = entity.EntityDisplayName
						typeName = entity.EntityTypeName
						break
					}
				}

			}

		}

		if flag {
			parts = append(parts, &dialogflowpb.Intent_TrainingPhrase_Part{

				Text:       str,
				EntityType: typeName,
				Alias:      Name,
			})

		} else {
			before, number, after := stringsman.FindNumberInString(str)
			if number != "" {
				parts = append(parts, &dialogflowpb.Intent_TrainingPhrase_Part{
					Text: before,
				})
				parts = append(parts, &dialogflowpb.Intent_TrainingPhrase_Part{

					Text:       number,
					EntityType: "@sys.number",
					Alias:      "temperature",
				})
				parts = append(parts, &dialogflowpb.Intent_TrainingPhrase_Part{
					Text: after,
				})
			} else {
				parts = append(parts, &dialogflowpb.Intent_TrainingPhrase_Part{
					Text: str,
				})
			}
		}

	}

	trainingPhrase := &dialogflowpb.Intent_TrainingPhrase{
		Parts: parts,
	}

	trainingPhrases := []*dialogflowpb.Intent_TrainingPhrase{
		trainingPhrase,
	}

	return trainingPhrases

}

func CreateEntityType(entityTypeName string, entities map[string][]string) {
	ctx := context.Background()
	client, err := apiv2.NewEntityTypesClient(ctx, option.WithCredentialsFile(keyJSON))
	if err != nil {
		log.Fatalf("Failed to create intents client: %v", err)
	}
	defer client.Close()

	reqGet := &dialogflowpb.ListEntityTypesRequest{
		Parent: "projects/" + projectID + "/agent",
	}
	it := client.ListEntityTypes(ctx, reqGet)

	var existingEntityType *dialogflowpb.EntityType
	for {
		entityType, err := it.Next()
		if err != nil {
			break
		}
		if entityType.DisplayName == entityTypeName {
			existingEntityType = entityType
			break
		}
	}

	if existingEntityType != nil {

		for value, synonyms := range entities {
			entityExists := false
			for _, existingEntity := range existingEntityType.Entities {
				if existingEntity.Value == value {
					existingEntity.Synonyms = synonyms
					entityExists = true
					break
				}
			}
			if !entityExists {
				existingEntityType.Entities = append(existingEntityType.Entities, &dialogflowpb.EntityType_Entity{
					Value:    value,
					Synonyms: synonyms,
				})
			}
		}

		reqUpdate := &dialogflowpb.UpdateEntityTypeRequest{
			EntityType: existingEntityType,
		}
		updatedResp, err := client.UpdateEntityType(ctx, reqUpdate)
		if err != nil {
			log.Fatalf("Failed to update entity type: %v", err)
		}
		log.Printf("Updated entity type: %v", updatedResp)

	} else {

		var entityList []*dialogflowpb.EntityType_Entity
		for value, synonyms := range entities {
			entityList = append(entityList, &dialogflowpb.EntityType_Entity{
				Value:    value,
				Synonyms: synonyms,
			})
		}

		reqCreate := &dialogflowpb.CreateEntityTypeRequest{
			Parent: "projects/" + projectID + "/agent",
			EntityType: &dialogflowpb.EntityType{
				Kind:        dialogflowpb.EntityType_KIND_MAP,
				DisplayName: entityTypeName,
				Entities:    entityList,
			},
		}

		resp, err := client.CreateEntityType(ctx, reqCreate)
		if err != nil {
			log.Fatalf("Failed to create entity type: %v", err)
		}
		log.Printf("Created entity type: %v", resp)
	}
}

func AddPhrases(intent Intent, trainingPhrases []string) {

	var trainingPhraseList []*dialogflowpb.Intent_TrainingPhrase

	ctx := context.Background()

	client, err := apiv2.NewIntentsClient(ctx, option.WithCredentialsFile(keyJSON))
	if err != nil {
		log.Fatalf("Failed to create Dialogflow client: %v", err)
	}
	defer client.Close()

	intentID, err := getIntentIDByDisplayName(projectID, intent.IntentName)
	if err != nil {
		log.Fatalf("Failed to get intent ID: %v", err)
	}

	req := &dialogflowpb.GetIntentRequest{
		Name: fmt.Sprint(intentID),
	}

	resp, err := client.GetIntent(ctx, req)
	if err != nil {
		log.Fatalf("Failed to get intent: %v", err)
	}

	for _, phrase := range trainingPhrases {
		words := stringsman.SplitByWords(phrase, getWordsFromEntities(intent))
		trainingPhraseList = append(trainingPhraseList, combineTrainingPhraseParts(words, intent)...)
	}

	resp.TrainingPhrases = append(resp.TrainingPhrases, trainingPhraseList...)

	updateReq := &dialogflowpb.UpdateIntentRequest{
		Intent:     resp,
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"training_phrases"}},
	}

	updatedIntent, err := client.UpdateIntent(ctx, updateReq)
	if err != nil {
		log.Fatalf("Failed to update intent: %v", err)
	}

	fmt.Printf("Updated intent with new training phrases: %v\n", updatedIntent)

}

func includes(slice []Entity, item string) bool {
	for _, entity := range slice {
		if entity.EntityDisplayName == item {
			return true
		}
	}
	return false
}

func ClearParams(intent Intent) {

	ctx := context.Background()

	client, err := apiv2.NewIntentsClient(ctx, option.WithCredentialsFile(keyJSON))
	if err != nil {
		log.Fatalf("Failed to create Dialogflow client: %v", err)
	}
	defer client.Close()
	intentID, err := getIntentIDByDisplayName(projectID, intent.IntentName)
	if err != nil {
		log.Fatalf("Failed to get intent ID: %v", err)
	}

	req := &dialogflowpb.GetIntentRequest{
		Name: fmt.Sprint(intentID),
	}

	resp, err := client.GetIntent(ctx, req)
	if err != nil {
		log.Fatalf("Failed to get intent: %v", err)
	}
	var empty []*dialogflowpb.Intent_Parameter
	resp.Parameters = empty

	updateReq := &dialogflowpb.UpdateIntentRequest{
		Intent:     resp,
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"parameters"}},
	}

	updatedIntent, err := client.UpdateIntent(ctx, updateReq)
	if err != nil {
		fmt.Println("Entity already exists in the intent.")
	}

	fmt.Printf("Cleared parameter: %v\n", updatedIntent)
}

func AddEntityToIntent(intent Intent, displayName string, entityType string, value string) {

	ctx := context.Background()

	client, err := apiv2.NewIntentsClient(ctx, option.WithCredentialsFile(keyJSON))
	if err != nil {
		log.Fatalf("Failed to create Dialogflow client: %v", err)
	}
	defer client.Close()
	intentID, err := getIntentIDByDisplayName(projectID, intent.IntentName)
	if err != nil {
		log.Fatalf("Failed to get intent ID: %v", err)
	}

	req := &dialogflowpb.GetIntentRequest{
		Name: fmt.Sprint(intentID),
	}

	resp, err := client.GetIntent(ctx, req)
	if err != nil {
		log.Fatalf("Failed to get intent: %v", err)
	}

	isRequired := false

	if includes(intent.RequiredEntities, displayName) {
		isRequired = true
	}

	newParameter := &dialogflowpb.Intent_Parameter{
		DisplayName:           displayName,
		EntityTypeDisplayName: entityType,
		Value:                 value,
		Mandatory:             isRequired,
	}

	resp.Parameters = append(resp.Parameters, newParameter)

	updateReq := &dialogflowpb.UpdateIntentRequest{
		Intent:     resp,
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"parameters"}},
	}

	updatedIntent, err := client.UpdateIntent(ctx, updateReq)
	if err != nil {
		fmt.Println("Entity already exists in the intent.")
	}

	fmt.Printf("Updated intent: %v\n", updatedIntent)
}

func getIntentIDByDisplayName(projectID, displayName string) (string, error) {

	ctx := context.Background()

	client, err := apiv2.NewIntentsClient(ctx, option.WithCredentialsFile(keyJSON))
	if err != nil {
		log.Fatalf("Failed to create Dialogflow client: %v", err)
	}
	defer client.Close()
	parent := fmt.Sprintf("projects/%s/agent", projectID)

	req := &dialogflowpb.ListIntentsRequest{
		Parent: parent,
	}

	it := client.ListIntents(ctx, req)

	for {
		intent, err := it.Next()
		if err != nil {
			break
		}

		if intent.GetDisplayName() == displayName {
			return intent.GetName(), nil
		}
	}

	return "", fmt.Errorf("intent with display name %s not found", displayName)
}
