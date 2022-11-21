package controller

import (
	"net/http"
	"worldimg/database"

	"github.com/gin-gonic/gin"
)

func GetSearch(c *gin.Context) {
	var (
		datalist  database.ImgList
		st        string = c.DefaultQuery("status", "0")
		key       string = c.DefaultQuery("key", "0")
		datayList []*database.ImgLists
		dataList  []*database.ImgList
		err       error
	)
	if st == "1" {
		datayList, err = datalist.SearchKeyY(key)
		// fmt.Println(datayList)
	} else {
		dataList, err = datalist.SearchKey(key)
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  0,
			"message": err.Error(),
		})
		return
	}

	// fmt.Println(dataList)
	Data := gin.H{
		"status": 1,
		"data":   dataList,
	}

	if st == "1" {
		Data = gin.H{
			"status": 1,
			"data":   datayList,
		}
	}
	c.JSON(http.StatusOK, Data)
}
