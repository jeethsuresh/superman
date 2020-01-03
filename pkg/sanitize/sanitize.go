package sanitize

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// InputContextKey is used in other structs to find the input data from the context
const InputContextKey = "INPUT"

// InputData is used to sanitize the input into usable information and carry it along the program flow
type InputData struct {
	Username  string `json:"username" binding:"required"`
	Timestamp int    `json:"unix_timestamp" binding:"required"`
	UUID      string `json:"event_uuid" binding:"required"`
	IP        string `json:"ip_address" binding:"required"`
}

// Middleware is the handler func that slots into gin's middleware chain
func Middleware(c *gin.Context) {
	fmt.Println("Sanitization Middleware")

	var data InputData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}

	//TODO: check UUID formatting

	//TODO: should we check the timestamp for reasonably high values?

	c.Set(InputContextKey, data)

	fmt.Printf("Sanitized input: %+v\n", data)
	c.Next()
	fmt.Println("==========================")
}
