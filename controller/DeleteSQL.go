package controller

import (
	"net/http"
	"worldimg/database"

	"github.com/gin-gonic/gin"
)

func DeleteOneSQL(c *gin.Context) {
	var form Account
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}

	if len(form.Account) <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "haven't node",
		})
		return
	}

	var delit database.ImgList

	delit.DeleteOne(form.Account)
	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "删除成功",
	})
}

func DeleteListSQL(c *gin.Context) {
	var form List
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}

	var delit database.ImgList
	for _, item := range form.List {
		delit.DeleteOne(item)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "删除成功",
	})
}
