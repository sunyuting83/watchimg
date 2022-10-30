package controller

import (
	"net/http"
	"strconv"
	"worldimg/database"

	"github.com/gin-gonic/gin"
)

func GetData(c *gin.Context) {
	var page string = c.DefaultQuery("page", "0")
	pageInt, _ := strconv.ParseInt(page, 10, 64)

	var datalist database.ImgList
	count, err := datalist.GetCount()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  0,
			"message": "失败",
		})
		return
	}
	dataList, err := datalist.GetImgList(pageInt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  0,
			"message": "失败",
		})
		return
	}
	Data := gin.H{
		"status": 1,
		"data":   dataList,
		"total":  count,
	}
	c.JSON(http.StatusOK, Data)
}

func GetOldData(c *gin.Context) {
	var datalist database.ImgList
	dataList, err := datalist.GetOldImgList()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  0,
			"message": "失败",
		})
		return
	}
	// fmt.Println(dataList)
	Data := gin.H{
		"status": 1,
		"data":   dataList,
	}
	c.JSON(http.StatusOK, Data)
}
