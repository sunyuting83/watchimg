package controller

import (
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

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
	// fmt.Println(file)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  0,
			"message": "上传文件失败",
		})
		return
	}
	fileName := header.Filename
	if strings.Contains(header.Filename, `\`) {
		nameList := strings.Split(header.Filename, `\`)
		Length := len(nameList) - 1
		fileName = nameList[Length]
	}
	Path := strings.Join([]string{"upimg", form.Code}, SplitString)
	imgPath := strings.Join([]string{Path, fileName}, SplitString)
	check := utils.IsExist(Path)
	if !check {
		err := os.MkdirAll(Path, 0766)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  0,
				"message": "上传文件失败",
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
	out, err := os.Create(imgPath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  0,
			"message": "上传文件失败",
		})
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  0,
			"message": "上传文件失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "上传文件成功",
	})
}
