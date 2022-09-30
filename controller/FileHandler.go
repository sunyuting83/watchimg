package controller

import (
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"

	ocr "worldimg/OCR"
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
	fileName := header.Filename
	if strings.Contains(header.Filename, `\`) {
		nameList := strings.Split(header.Filename, `\`)
		Length := len(nameList) - 1
		fileName = nameList[Length]
	}
	fileNameList := strings.Split(fileName, ".")
	fileName = strings.Join([]string{fileNameList[0], "_", gold, ".", fileNameList[1]}, "")
	Path := strings.Join([]string{"upimg", form.Code}, SplitString)
	// fmt.Println(Path)
	imgPath := strings.Join([]string{Path, fileName}, SplitString)
	// fmt.Println(imgPath)
	// fmt.Println(gold)
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
	cName := strings.Join([]string{Path, ".cname"}, SplitString)
	checkName := utils.IsExist(cName)
	if checkName {
		old, _ := os.ReadFile(cName)
		if string(old) != form.ComputName {
			os.WriteFile(cName, []byte(form.ComputName), 0644)
		}
	} else {
		os.Create(cName)
		os.WriteFile(cName, []byte(form.ComputName), 0644)
	}
	err = os.WriteFile(imgPath, b, 0644)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  0,
			"message": "3上传文件失败",
		})
		return
	}
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
