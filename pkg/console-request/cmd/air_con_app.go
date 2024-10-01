package main

import (
	"bufio"
	"bytes"
	"console_request/speech"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	url      = ""
	filePath = "output.wav"
)

var (
	sessionID   = ""
	borderColor = "\033[48;2;29;54;88m"

	optionBgColor  = "\033[48;2;68;122;156m"
	textColor      = "\033[38;2;167;218;220m"
	highlightColor = "\033[48;2;242;250;239m"
	reset          = "\033[0m"
)

type UploadResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Upload struct {
	UserID    string `json:"userID"`
	SessionID string `json:"sessionID"`
	Request   string `json:"request"`
	Query     string `json:"query"`
	Buffer    []byte `json:"buffer"`
}

func main() {
	sessionID = time.Now().Format("20060102_150405")
	fmt.Println(borderColor + textColor + "====================================" + reset)
	fmt.Println(highlightColor + textColor + "Enter user ID:                      " + reset)
	fmt.Println(borderColor + textColor + "====================================" + reset)

	reader := bufio.NewReader(os.Stdin)

	userID, _ := reader.ReadString('\n')
	userID = strings.TrimSpace(userID)

	for {

		fmt.Println(borderColor + "====================================" + reset)
		fmt.Println(highlightColor + textColor + "           Choose an option         " + reset)
		fmt.Println(borderColor + "====================================" + reset)
		fmt.Println(optionBgColor + textColor + "|  record        -> Record Voice   |" + reset)
		fmt.Println(optionBgColor + textColor + "|  write         -> Text Input     |" + reset)
		fmt.Println(optionBgColor + textColor + "|  history       -> Show History   |" + reset)
		fmt.Println(optionBgColor + textColor + "|  sessionhistory-> Show Session   |" + reset)
		fmt.Println(optionBgColor + textColor + "|  train         -> Train the Agent|" + reset)
		fmt.Println(optionBgColor + textColor + "|  exit          -> Exit           |" + reset)
		fmt.Println(borderColor + "====================================" + reset)

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch strings.ToLower(input) {
		case "record":
			speech.RecordVoice()

			data := speech.ConvertToWav()

			stringData := map[string]string{
				"userID":    userID,
				"sessionID": sessionID,
			}

			sendByteArray(data, stringData)
		case "train":
			resp, err := http.Get(url + "agent")
			if err != nil {
				panic(err)
			}

			err = resp.Body.Close()
			if err != nil {
				fmt.Println("Failed to close resource:", err)
				return
			}

		case "history", "sessionhistory":
			data := Upload{
				UserID:    userID,
				Request:   input,
				SessionID: sessionID,
			}

			httpHandler(data)

		case "exit":
			os.Exit(0)
		case "write":
			reader := bufio.NewReader(os.Stdin)
			fmt.Println(borderColor + textColor + "====================================" + reset)
			fmt.Println(optionBgColor + textColor + "Enter the query:                    " + reset)
			fmt.Println(borderColor + textColor + "====================================" + reset)
			query, _ := reader.ReadString('\n')

			data := Upload{
				UserID:    userID,
				Request:   input,
				Query:     query,
				SessionID: sessionID,
			}
			httpHandler(data)

		default:
			fmt.Println("Entered wrong value, try again.")
		}

	}

}

func sendByteArray(data []byte, stringData map[string]string) error {

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	part, err := writer.CreateFormField("byte_data")
	if err != nil {
		return err
	}
	part.Write(data)

	for key, value := range stringData {
		err := writer.WriteField(key, value)
		if err != nil {
			return err
		}
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url+"upload", &b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send data, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var uploadResponse UploadResponse
	err = json.Unmarshal(body, &uploadResponse)
	if err != nil {
		fmt.Println("Error unmarshaling response:", err)
		return err
	}
	fmt.Println(borderColor + textColor + "====================================" + reset)
	fmt.Println(optionBgColor+textColor+"Upload Status:", uploadResponse.Status+reset)
	fmt.Println(optionBgColor+textColor+"Message:", uploadResponse.Message+reset)
	fmt.Println(borderColor + textColor + "====================================" + reset)
	return nil
}

func httpHandler(data Upload) {

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	resp, err := http.Post(url+"data", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var uploadResponse UploadResponse
	err = json.Unmarshal(responseBody, &uploadResponse)
	if err != nil {
		fmt.Println("Error unmarshaling response:", err)
		return
	}
	fmt.Println(borderColor + textColor + "====================================" + reset)
	fmt.Println(optionBgColor+textColor+"Upload Status:", uploadResponse.Status+reset)
	fmt.Println(optionBgColor+textColor+"Message:", uploadResponse.Message+reset)
	fmt.Println(borderColor + textColor + "====================================" + reset)

}
