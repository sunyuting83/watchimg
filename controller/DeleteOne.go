package controller

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"runtime"
	"strings"
	"worldimg/database"
	"worldimg/utils"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

// Node node
type Account struct {
	Account string `form:"account" json:"account" xml:"account"  binding:"required"`
}

func DeleteOne(c *gin.Context) {
	var form Account
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error() + "sssssss",
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
	OS := runtime.GOOS
	LinkPathStr := "/"
	if OS == "windows" {
		LinkPathStr = "\\"
	}
	CurrentPath, _ := utils.GetCurrentPath()

	ConfigFile := strings.Join([]string{CurrentPath, "config.yaml"}, LinkPathStr)

	var confYaml *Config
	yamlFile, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "提取失败1",
			})
			return
		}
	}
	err = yaml.Unmarshal(yamlFile, &confYaml)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "提取失败2",
		})
		return
	}
	if len(confYaml.Host) <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "提取失败3",
		})
		return
	}
	token := c.GetHeader("Authorization")
	if len(token) < 10 {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "提取失败6",
		})
		return
	}
	token = token[7:]
	del := postIT(form.Account, CurrentPath, token, confYaml)
	if del == "0" {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "提取失败4",
		})
		return
	}
	if del == "1" {
		c.JSON(http.StatusOK, gin.H{
			"status":  -1,
			"message": "登陆失效5",
		})
		return
	}
	var datalist database.ImgList
	datalist.DeleteOne(form.Account)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "提取失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "提取成功6",
		"data":    del,
	})
}

func postIT(account, CurrentPath, token string, confYaml *Config) string {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("account", account)
	bodyWriter.WriteField("gameid", confYaml.GameID)

	OS := runtime.GOOS
	LinkPathStr := "/"
	if OS == "windows" {
		LinkPathStr = "\\"
	}

	ConfigFile := strings.Join([]string{CurrentPath, ".token"}, LinkPathStr)

	TokenFile, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		return "0"
	}
	if string(TokenFile) == token {
		contentType := bodyWriter.FormDataContentType()
		bodyWriter.Close()
		URL := strings.Join([]string{confYaml.Host, "/api/gitone"}, "")
		resp, err := http.Post(URL, contentType, bodyBuf)
		resp.Header.Add("Bearer ", "zhaofan")
		if err != nil {
			return "0"
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "0"
		}
		// fmt.Println(string(body))
		return string(body)
	} else {
		return "1"
	}
	// return nil
}
