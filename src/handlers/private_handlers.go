package handlers

import (
	api "github.com/Isterdam/hack-the-crisis-backend/src/api/private"
	"github.com/Isterdam/hack-the-crisis-backend/src/middleware"
	"github.com/gin-gonic/gin"
)

func PrivateRoutes(r *gin.Engine) {
	private := r.Group("/private", middleware.PrivateAuthenticationHandler)
	{
		private.POST("/gtts", api.TextToSpeech)
		private.POST("/booking", api.CreateBooking)
	}
	r.GET("/private/mp3/:fileID", api.GetMP3)
}
