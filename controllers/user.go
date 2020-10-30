package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/drrainlab/hydropony-contest/cacher"
	"github.com/drrainlab/hydropony-contest/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetUserInfoHandler - get user info by "user-id" header
func GetUserInfoHandler(c *gin.Context) {

	var user models.User

	userID := c.GetHeader("user-id")

	// trying to get record from cache
	u, _ := cacher.Instance.Get(fmt.Sprintf("user-%s", userID))

	if u != nil {
		if err := json.Unmarshal(u.Value, &user); err != nil {
			panic(err)
		}
	} else {
		// get user info from db
		if err := models.DB.First(&user, userID).Error; errors.Is(err, gorm.ErrRecordNotFound) == true {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "User not found",
			})
			return
		}
		// caching record
		userDataEncoded, err := json.Marshal(&user)

		if err != nil {
			panic(err)
		}

		if err := cacher.Instance.Set(&memcache.Item{Key: fmt.Sprintf("user-%d", user.ID), Value: userDataEncoded}); err != nil {
			panic(err)
		}
	}

	c.JSON(200, gin.H{
		"status":   "OK",
		"fullname": user.GetFullname(),
		"address":  user.Adresses,
		"email":    user.Email,
	})

}

// create new user
func CreateUserHandler(c *gin.Context) {

	var user models.User

	if err := c.ShouldBindJSON(&user); err == nil {
		if result := models.DB.Create(&user); result.Error != nil {
			c.JSON(http.StatusConflict, gin.H{
				"error": result.Error.Error(),
			})
			return
		}
	} else {
		c.JSON(422, gin.H{
			"status": "Invalid request",
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "OK",
		"id":     user.ID,
	})

}
