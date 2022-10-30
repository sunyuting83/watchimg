package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
	"worldimg/database"
	"worldimg/utils"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

// Node node
type Login struct {
	UserName string `form:"username" json:"username" xml:"username"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password"  binding:"required"`
}

type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

type Config struct {
	Host   string `yaml:"Host"`
	GameID string `yaml:"GameID"`
}

func Sgin(c *gin.Context) {
	var form Login
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}

	if len(form.UserName) <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "haven't username",
		})
		return
	}
	if len(form.Password) <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "haven't password",
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
	yamlFile, err := os.ReadFile(ConfigFile)
	if err != nil {
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "登陆失败",
			})
			return
		}
	}
	err = yaml.Unmarshal(yamlFile, &confYaml)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "登陆失败",
		})
		return
	}
	if len(confYaml.Host) <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "登陆失败",
		})
		return
	}
	login := postIt(form.UserName, form.Password, CurrentPath, confYaml)
	if login == "0" {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "登陆失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "登陆成功",
		"token":   login,
	})
}

func postIt(username, password, CurrentPath string, confYaml *Config) string {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("username", username)
	bodyWriter.WriteField("password", password)

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	URL := strings.Join([]string{confYaml.Host, "/api/manage/login"}, "")
	resp, err := http.Post(URL, contentType, bodyBuf)
	if err != nil {
		return "0"
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "0"
	}
	var m *Status
	json.Unmarshal(body, &m)
	// fmt.Println(m.Status)
	if m.Code != 1 {
		return "0"
	} else {
		nowTime := time.Now().Unix()
		var user database.User
		user.Username = username
		user.Token = m.Token
		user.DateTime = nowTime
		user.Insert()
		return m.Token
	}
}
