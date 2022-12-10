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

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

// Node node
type Account struct {
	Account string `form:"account" json:"account" xml:"account"  binding:"required"`
}

type DelData struct {
	Accountgs string `json:"accountgs"`
	Password  string `json:"password"`
	Other     string `json:"other"`
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
	user, has := c.Get("user_id")
	if !has {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "haven't node",
		})
		return
	}
	// fmt.Println(user.(int64))
	OS := runtime.GOOS
	LinkPathStr := "/"
	if OS == "windows" {
		LinkPathStr = "\\"
	}
	CurrentPath, _ := c.Get("current_path")

	ConfigFile := strings.Join([]string{CurrentPath.(string), "config.yaml"}, LinkPathStr)
	var confYaml *Config
	yamlFile, err := os.ReadFile(ConfigFile)
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
	ID := user.(int64)
	NowTime := time.Now().Unix()
	var datalist database.ImgList
	del := postIT(form.Account, token, confYaml)
	if del == "0" {
		datalist.NewStatus = 2
		datalist.UpDateTime = NowTime
		datalist.UserID = ID
		datalist.UpdateImg(form.Account)
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

	if del == "{}" {
		datalist.NewStatus = 2
		datalist.UpDateTime = NowTime
		datalist.UserID = ID
		datalist.UpdateImg(form.Account)
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "提取失败没数据",
		})
		return
	}
	if strings.Contains(del, "参数") {
		datalist.NewStatus = 2
		datalist.UpDateTime = NowTime
		datalist.UserID = ID
		datalist.UpdateImg(form.Account)
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "提取失败没数据",
		})
		return
	}

	var datas *DelData
	json.Unmarshal([]byte(del), &datas)

	datalist.NewStatus = 1
	datalist.UpDateTime = NowTime
	datalist.UserID = ID
	datalist.Password = datas.Password

	_, err = datalist.UpdateImg(form.Account)
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
		"data":    datas,
	})
}

type List struct {
	List []string `form:"list" json:"list" xml:"list"  binding:"required"`
}

func DeleteList(c *gin.Context) {
	var form List
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}

	user, has := c.Get("user_id")
	if !has {
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
	CurrentPath, _ := c.Get("current_path")

	ConfigFile := strings.Join([]string{CurrentPath.(string), "config.yaml"}, LinkPathStr)

	var confYaml *Config
	yamlFile, err := os.ReadFile(ConfigFile)
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
	// fmt.Println(form.List)
	token = token[7:]
	var temp []*DelData
	ID := user.(int64)
	NowTime := time.Now().Unix()
	var datalist database.ImgList
	for _, item := range form.List {
		var datas *DelData
		del := postIT(item, token, confYaml)
		// fmt.Println(del)
		if del != "0" && del != "1" && del != "{}" && !strings.Contains(del, "参数") {
			json.Unmarshal([]byte(del), &datas)

			datalist.NewStatus = 1
			datalist.UpDateTime = NowTime
			datalist.UserID = ID
			datalist.Password = datas.Password

			datalist.UpdateImg(item)
			temp = append(temp, datas)
		} else {
			if del != "{}" {
				json.Unmarshal([]byte(del), &datas)

				datalist.NewStatus = 2
				datalist.UpDateTime = NowTime
				datalist.UserID = ID
				datalist.UpdateImg(item)
			}
		}
	}
	if len(temp) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "提取失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "提取成功",
		"list":    temp,
	})
}

type PostData struct {
	Account string `json:"account"`
	GameID  string `json:"gameid"`
}

func postIT(account, token string, confYaml *Config) string {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("account", account)
	bodyWriter.WriteField("gameid", confYaml.GameID)

	var user database.User
	data, _ := user.GetUser(token)
	if data.Token == token {
		info := PostData{
			Account: account,
			GameID:  confYaml.GameID,
		}
		byte, _ := json.Marshal(info)
		// fmt.Println(string(byte))
		URL := strings.Join([]string{confYaml.Host, "/api/manage/gitone"}, "")
		req, _ := http.NewRequest("POST", URL, strings.NewReader(string(byte)))
		req.Header.Set("Accept", "application/json, text/plain, */*")
		req.Header.Set("content-type", "application/json")
		req.Header.Set("sec-ch-ua-platform", "Windows")
		req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36")
		req.Header.Set("Authorization", strings.Join([]string{"Bearer", token}, " "))
		resp, err := (&http.Client{}).Do(req)
		if err != nil {
			return "0"
		}
		defer resp.Body.Close()

		respByte, _ := io.ReadAll(resp.Body)
		// fmt.Println(string(respByte))
		d := string(respByte)
		if strings.Contains(d, "Error") {
			return "0"
		}
		return d
	} else {
		return "1"
	}
	// return nil
}
