package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func AddCallInfo(c *gin.Context) {
	var ctx = context.Background()
	rdb := c.MustGet("rdb").(*redis.Client)
	callID := c.Param("callID")

	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)

	rdb.Set(ctx, callID, buf.String(), time.Minute*60)
	val, err := rdb.Get(ctx, callID).Result()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not add call info.",
		})
		return
	}

	res := struct {
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
	}{
		Message: "Success",
		Data:    []byte(val),
	}

	c.JSON(http.StatusOK, res)
}

func UpdateCallInfo(c *gin.Context) {
	var ctx = context.Background()
	rdb := c.MustGet("rdb").(*redis.Client)
	callID := c.Param("callID")

	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)

	rdb.Set(ctx, callID, buf.String(), time.Minute*60)
	val, err := rdb.Get(ctx, callID).Result()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not update call info.",
		})
		return
	}

	res := struct {
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
	}{
		Message: "Success",
		Data:    []byte(val),
	}

	c.JSON(http.StatusOK, res)
}

func GetCallInfo(c *gin.Context) {
	var ctx = context.Background()
	rdb := c.MustGet("rdb").(*redis.Client)
	callID := c.Param("callID")

	val, err := rdb.Get(ctx, callID).Result()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not get call info.",
		})
		return
	}

	res := struct {
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
	}{
		Message: "Success",
		Data:    []byte(val),
	}

	c.JSON(http.StatusOK, res)
}
