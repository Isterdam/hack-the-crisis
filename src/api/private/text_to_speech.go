package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	data := GTTSConfig{
		AudioConfig: GTTSAudioConfig{
			AudioEncoding: "MP3",
			Pitch:         0,
			SpeakingRate:  0.66,
		},
		Input: GTTSInput{
			Text: "bananer",
		},
		Voice: GTTSVoice{
			LanguageCode: "sv-SE",
			Name:         "sv-SE-Wavenet-A",
		},
	}

	url := "https://texttospeech.googleapis.com/v1beta1/text:synthesize?key=AIzaSyCggB1IlCl46gw-8Cjb6o6CMltqcb3Tnqk"

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(data)

	req, _ := http.NewRequest("POST", url, buf)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(buf.Len()))

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	defer res.Body.Close()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	audioResponse := GTTSResponse{}

	err = json.NewDecoder(res.Body).Decode(&audioResponse)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    audioResponse,
	})

}
