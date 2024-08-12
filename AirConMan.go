package airconman

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	status        = "Turned Off"
	degree uint16 = 24
)

type DialogflowRequest struct {
	QueryResult struct {
		Parameters map[string]interface{} `json:"parameters"`
		Intent     struct {
			DisplayName string `json:"displayName"`
		} `json:"intent"`
	} `json:"queryResult"`
}

type DialogflowResponse struct {
	FulfillmentText string `json:"fulfillmentText"`
}

func main() {
	r := gin.Default()

	r.POST("/webhook", func(c *gin.Context) {
		var req DialogflowRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var responseText string

		switch req.QueryResult.Intent.DisplayName {
		case "TurnOn":

			status = "Turned On"
			responseText = "Turning on the air conditioner."

		case "TurnOff":

			status = "Turned Off"
			responseText = "Turning off the air conditioner."

		case "ChangeDegree":

			if 18 <= req.QueryResult.Parameters["value"].(uint16) && req.QueryResult.Parameters["value"].(uint16) <= 32 {
				degree = req.QueryResult.Parameters["value"].(uint16)
				responseText = "Changing the degree to " + fmt.Sprintf("%v", degree)
			} else {
				responseText = "The degree you entered is not in the range of the air conditioner."
			}
		default:
			responseText = "Sorry, I didn't understand that."
		}

		c.JSON(http.StatusOK, DialogflowResponse{
			FulfillmentText: responseText,
		})

	})

	r.Run(":8080")

}
