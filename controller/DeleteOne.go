package controller

import (
	"net/http"
	"os"
	"runtime"
	"strings"
	"worldimg/utils"

	"github.com/gin-gonic/gin"
)

// Node node
type Image struct {
	Path string `form:"path" json:"path" xml:"path"  binding:"required"`
}

func DeleteOne(c *gin.Context) {
	var form Image
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error() + "sssssss",
		})
		return
	}

	if len(form.Path) <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "haven't node",
		})
		return
	}
	OS := runtime.GOOS
	LinkPathStr := "/"
	if OS == "windows" {
		LinkPathStr = "\\"
	}
	CurrentPath, _ := utils.GetCurrentPath()
	// fmt.Println(form.Path)
	newPath := strings.Replace(form.Path, "static", "upimg", -1)
	newPath = newPath[1:]
	imgPath := strings.Join([]string{CurrentPath, newPath}, LinkPathStr)
	// fmt.Println(imgPath)
	err := os.Remove(imgPath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "删除文件失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "删除文件成功",
	})
}
