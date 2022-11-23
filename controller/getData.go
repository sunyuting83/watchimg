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

func GetDateTimeData(c *gin.Context) {
	// var datalist database.ImgList
	var st string = c.DefaultQuery("status", "0")
	var date string = c.DefaultQuery("date", "0")
	dateList, err := database.GetDateTimeData(st)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  0,
			"message": "失败",
		})
		return
	}
	var datayList []*database.ImgLists
	var dataList []*database.ImgList
	if len(dateList) > 0 {
		if len(date) <= 0 || date == "0" {
			date = dateList[0]
		}

		if st == "0" {
			dataList, err = database.GetDateTimeDataNList(date)
			// fmt.Println(datayList)
		} else {
			datayList, err = database.GetDateTimeDataYList(date, st)
		}
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  0,
				"message": "失败",
			})
			return
		}
	}

	// fmt.Println(dataList)
	Data := gin.H{
		"status":   1,
		"datelist": dateList,
		"data":     datayList,
	}

	if st == "0" {
		Data = gin.H{
			"status":   1,
			"datelist": dateList,
			"data":     dataList,
		}
	}
	c.JSON(http.StatusOK, Data)
}
