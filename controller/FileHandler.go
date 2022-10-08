package controller

import (
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	ocr "worldimg/OCR"
	"worldimg/database"
	"worldimg/utils"

	"github.com/gin-gonic/gin"
)

// Node node
type Node struct {
	Code       string `form:"code" json:"code" xml:"code"  binding:"required"`
	ComputName string `form:"computname" json:"computname" xml:"computname"  binding:"required"`
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

	if len(form.Code) <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "haven't node",
		})
		return
	}
	if len(form.ComputName) <= 0 {
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

	b, _ := ioutil.ReadAll(file)
	gold := GetWord(b)
	gold = strings.Replace(gold, " ", "", -1)
	gold = strings.Replace(gold, "'", "", -1)
	gold = strings.Replace(gold, ",", "", -1)
	if len(gold) == 0 {
		gold = "0"
	}
	if strings.Contains(gold, "\n") {
		gold = "0"
	}
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
	Path := strings.Join([]string{"upimg", timeStr}, SplitString)
	// fmt.Println(Path)
	imgPath := strings.Join([]string{Path, fileName}, SplitString)
	toDBPath := strings.Join([]string{"/static", timeStr}, "/")
	toDBPath = strings.Join([]string{toDBPath, fileName}, "/")
	// fmt.Println(imgPath)
	check := utils.IsExist(Path)
	if !check {
		err := os.MkdirAll(Path, 0766)
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
	newTime := nowTime.Unix()
	var imglist database.ImgList
	account, err := imglist.GetImgOne(fileNameList[0])
	imglist.Account = fileNameList[0]
	imglist.Cover = toDBPath
	imglist.Today = gold
	imglist.YesterDay = account.Today
	imglist.DateTime = newTime
	if err != nil {
		if account.ID == 0 {
			imglist.YesterDay = "0"
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

func GetWord(b []byte) string {
	// client.SetImage("/home/sun/Pictures/photo_2022-09-26_01-32-46.jpg")
	ocr.Client.SetImageFromBytes(b)
	text, err := ocr.Client.Text()
	if err != nil {
		return "0"
	}
	return text
}
