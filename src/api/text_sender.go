package api

import (
	"github.com/gin-gonic/gin"

	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

// PhoneUser and PhonePass have been initialized

// to_phone structure: "+46....", confirmation is the string that is sent
func Send_text(c *gin.Context, to_phone string, confirmation string) {
	data := url.Values{
		"from":    {"ShopAlone"},
		"to":      {to_phone},
		"message": {confirmation},
	}

	c.Request, _ = http.NewRequest("POST", "https://api.46elks.com/a1/SMS", bytes.NewBufferString(data.Encode()))
	c.Request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	c.Request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	c.Request.SetBasicAuth(os.Getenv("PHONEUSER"), os.Getenv("PHONEPASS"))

	client := &http.Client{}
	resp, err := client.Do(c.Request)
	if err != nil {
		fmt.Println("Could not do request!")
	}
	defer resp.Body.Close()
}
