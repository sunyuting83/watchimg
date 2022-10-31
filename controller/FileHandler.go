package controller

import (
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"worldimg/database"
	"worldimg/utils"

	"github.com/gin-gonic/gin"
)

// Node node
type Node struct {
	Gold     string `form:"gold" json:"gold" xml:"gold"  binding:"required"`
	Multiple string `form:"multiple" json:"multiple" xml:"multiple"  binding:"required"`
}

func FileHandler(c *gin.Context) {
	OS := runtime.GOOS
	var SplitString string = "\\"
	if OS == "linux" {
		SplitString = "/"
	}
	var form Node
	// This will infer what binder to use depending on the content-type header.
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	if len(form.Gold) <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "haven't node",
		})
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  0,
			"message": "1上传文件失败",
		})
		return
	}

	b, _ := io.ReadAll(file)
	/*
		gold := GetWord(b)
		gold = strings.Replace(gold, " ", "", -1)
		gold = strings.Replace(gold, "'", "", -1)
		gold = strings.Replace(gold, ",", "", -1)
		if len(gold) == 0 {
			gold = "0"
		}
	*/
	// fmt.Println(form.Gold)
	var gold int64
	if strings.Contains(form.Gold, "亿") {
		gx := strings.Split(form.Gold, "亿")
		goldstr := gx[0]
		if strings.Contains(form.Gold, ".") {
			g := strings.Split(goldstr, ".")
			// fmt.Println(g)
			goldstr = strings.Join([]string{g[0], g[1]}, "")
			// fmt.Println(goldstr)
			var x int64 = 10000000
			if len(g[1]) >= 2 {
				x = 1000000
			}
			n, _ := strconv.ParseInt(goldstr, 10, 64)
			// fmt.Println(n)
			// s := n * x
			gold = n * x
			// fmt.Println("yi1")
			// fmt.Println(gold)
		} else {
			// fmt.Println(goldstr)
			n, _ := strconv.ParseInt(goldstr, 10, 64)
			gold = n * 100000000
			// fmt.Println("yi2")
			// fmt.Println(gold)
		}
	} else if strings.Contains(form.Gold, "万") {
		gx := strings.Split(form.Gold, "万")
		goldstr := gx[0]
		n, _ := strconv.ParseInt(goldstr, 10, 64)
		gold = n * 10000
		// fmt.Println("wan2")
		// fmt.Println(gold)
	} else {
		n, _ := strconv.ParseInt(form.Gold, 10, 64)
		gold = n
	}
	// fmt.Println(gold)
	fileName := header.Filename
	if strings.Contains(header.Filename, `\`) {
		nameList := strings.Split(header.Filename, `\`)
		Length := len(nameList) - 1
		fileName = nameList[Length]
	}
	fileNameList := strings.Split(fileName, ".")

	nowTime := time.Now()
	timeStr := nowTime.Format("20060102")
	// fileName = strings.Join([]string{fileNameList[0], "_", gold, ".", fileNameList[1]}, "")
	// Path := strings.Join([]string{"upimg", timeStr}, SplitString)
	// fmt.Println(Path)
	imgPath := strings.Join([]string{"upimg", fileName}, SplitString)
	// toDBPath := strings.Join([]string{"/static", timeStr}, "/")
	toDBPath := strings.Join([]string{"/static", fileName}, "/")
	// fmt.Println(imgPath)
	check := utils.IsExist("upimg")
	if !check {
		err := os.MkdirAll("upimg", 0766)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  0,
				"message": "2上传文件失败",
			})
			return
		}
	}
	err = os.WriteFile(imgPath, b, 0644)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  0,
			"message": "3上传文件失败",
		})
		return
	}
	Multiple, _ := strconv.ParseInt(form.Multiple, 10, 64)
	newTime := nowTime.Unix()
	var imglist database.ImgList
	account, err := imglist.GetImgOne(fileNameList[0])
	imglist.Account = fileNameList[0]
	imglist.Cover = toDBPath
	imglist.Today = gold
	imglist.Multiple = Multiple

	timeobj := time.Unix(int64(account.DateTime), 0)
	olDate := timeobj.Format("20060102")

	if timeStr > olDate {
		imglist.YesterDay = account.Today
	}

	imglist.DateTime = newTime
	if err != nil {
		if account.ID == 0 {
			imglist.YesterDay = 0
			imglist.Insert()
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "上传文件成功",
			})
			return
		}
	}
	imglist.UpdateStatus(fileNameList[0])
	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "上传文件成功",
	})
}

// func GetWord(b []byte) string {
// 	// client.SetImage("/home/sun/Pictures/photo_2022-09-26_01-32-46.jpg")
// 	ocr.Client.SetImageFromBytes(b)
// 	text, err := ocr.Client.Text()
// 	if err != nil {
// 		return "0"
// 	}
// 	return text
// }
