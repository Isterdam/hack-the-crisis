package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"os"

	"github.com/Isterdam/hack-the-crisis-backend/src/utils/random"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type GTTSAudioConfig struct {
	AudioEncoding string  `json:"audioEncoding"`
	Pitch         float32 `json:"pitch"`
	SpeakingRate  float32 `json:"speakingRate"`
}

type GTTSInput struct {
	Text string `json:"text"`
}

type GTTSVoice struct {
	LanguageCode string `json:"languageCode"`
	Name         string `json:"name"`
}

type GTTSConfig struct {
	AudioConfig GTTSAudioConfig `json:"audioConfig"`
	Input       GTTSInput       `json:"input"`
	Voice       GTTSVoice       `json:"voice"`
}

type GTTSResponse struct {
	AudioContent string `json:"audioContent"`
}

func TextToSpeech(c *gin.Context) {
	var ctx = context.Background()

	rdb := c.MustGet("rdb").(*redis.Client)
	var data GTTSConfig

	err := c.ShouldBindJSON(&data)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid JSON.",
		})
		return
	}

	GTTSToken := os.Getenv("GTTS_KEY")

	url := "https://texttospeech.googleapis.com/v1beta1/text:synthesize?key=" + GTTSToken 

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(data)

	req, err := http.NewRequest("POST", url, buf)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error.",
		})
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(buf.Len()))

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error.",
		})
		return
	}

	defer res.Body.Close()

	audioResponse := GTTSResponse{}

	fileID := utils.RandStringBytes(10)

	err = json.NewDecoder(res.Body).Decode(&audioResponse)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	rdb.Set(ctx, fileID, audioResponse.AudioContent, time.Minute*5)

	var response struct{
		GTTSResponse
		URL string `json:"url"`
	}

	response.AudioContent = audioResponse.AudioContent
	response.URL = "https://api.shopalone/private/mp3/" + fileID + ".mp3"

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    response,
	})

}

func GetMP3(c *gin.Context) {
	fileID := c.Param("fileID")

	var ctx = context.Background()
	rdb := c.MustGet("rdb").(*redis.Client)

	split := strings.Split(fileID, ".")

	val, err := rdb.Get(ctx, split[0]).Result()

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "File not found.",
		})
		return
	}

	buf, err := base64.StdEncoding.DecodeString(val)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "File not found.",
		})
		return
	}

	c.Header("Content-Type", "audio/mp3")
	c.Header("Content-Length", strconv.Itoa(len(buf)))
	c.Status(200)
	c.Stream(func(w io.Writer) bool {
		w.Write(buf)
		return false
	})

}
