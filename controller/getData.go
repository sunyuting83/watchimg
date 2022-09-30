package controller

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
	"worldimg/utils"

	"github.com/gin-gonic/gin"
)

type Data struct {
	Data []List `json:"data"`
}

type List struct {
	Comput string                 `json:"comput"`
	Data   map[string]interface{} `json:"data"`
}

type Images struct {
	Title  string `json:"title"`
	ImgURL string `json:"imgurl"`
	Gold   string `json:"gold"`
}

func GetData(c *gin.Context) {
	OS := runtime.GOOS
	LinkPathStr := "/"
	if OS == "windows" {
		LinkPathStr = "\\"
	}
	var data Data
	CurrentPath, _ := utils.GetCurrentPath()
	imgPath := strings.Join([]string{CurrentPath, "upimg"}, LinkPathStr)
	getDir := GetAllFile(imgPath)
	// fmt.Println(getDir)
	for _, item := range getDir {
		img := GetAllFile(item)
		cnamePath := strings.Join([]string{item, ".cname"}, LinkPathStr)
		cname, _ := os.ReadFile(cnamePath)
		// var imgList []Images
		DateList := makeDateList(img, item, LinkPathStr)
		DateList = RemoveRepeatedElement(DateList)
		// fmt.Println(DateList)
		m := make(map[string]interface{})
		// keys := make([]string, 0, len(DateList))
		for _, v := range DateList {
			var imgList []Images
			for _, it := range img {
				titlePath := strings.Split(it, LinkPathStr)
				titleLen := len(titlePath) - 1
				titleD := titlePath[titleLen]
				titleAndGold := strings.Split(titleD, ".")[0]

				var (
					titleList []string
					title     string = titleAndGold
					gold      string = "0"
				)
				if strings.Contains(titleAndGold, "_") {
					titleList = strings.Split(titleAndGold, "_")
					titleLen := len(titleList)
					title = titleList[0]
					gold = titleList[titleLen-1]
				}
				if len(gold) >= 8 {
					l := len(gold) - 8
					gold = gold[0:l]
				}
				if len(title) > 0 {
					file := strings.Join([]string{item, titleD}, LinkPathStr)
					fileInfo, _ := os.Stat(file)
					modTime := fileInfo.ModTime()
					day := modTime.Format("2006-01-02")
					if day == v {
						ddd := Images{Title: title, ImgURL: it, Gold: gold}
						imgList = append(imgList, ddd)
					}
				}
			}
			sort.Slice(imgList, func(i, j int) bool {
				iInt, _ := strconv.ParseInt(imgList[i].Gold, 10, 64)
				jInt, _ := strconv.ParseInt(imgList[j].Gold, 10, 64)
				return iInt > jInt
			})
			// fmt.Println(imgList)
			m[v] = imgList
		}
		// fmt.Println(m)

		data.Data = append(data.Data, List{
			Comput: string(cname),
			Data:   m,
		})
	}
	c.JSON(http.StatusOK, data)
}

func makeDateList(img []string, item, LinkPathStr string) []string {
	var imgList []time.Time
	for _, it := range img {
		titlePath := strings.Split(it, LinkPathStr)
		titleLen := len(titlePath) - 1
		titleD := titlePath[titleLen]
		title := strings.Split(titleD, ".")[0]
		if len(title) > 0 {
			file := strings.Join([]string{item, titleD}, LinkPathStr)
			fileInfo, _ := os.Stat(file)
			modTime := fileInfo.ModTime()
			// day := modTime.Format("2006-01-02")
			imgList = append(imgList, modTime)
		}
	}
	sort.Slice(imgList, func(i, j int) bool { return imgList[i].Before(imgList[j]) })
	return TimeToString(imgList)
}

func TimeToString(imgList []time.Time) []string {
	var a []string
	for _, v := range imgList {
		a = append(a, v.Format("2006-01-02"))
	}
	return a
}

func RemoveRepeatedElement(personList []string) (result []string) {
	n := len(personList)
	for i := 0; i < n; i++ {
		repeat := false
		for j := i + 1; j < len(personList); j++ {
			if personList[i] == personList[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			result = append(result, personList[i])
		}
	}
	return
}

func GetAllFile(pathname string) []string {

	OS := runtime.GOOS
	LinkPathStr := "/"
	if OS == "windows" {
		LinkPathStr = "\\"
	}
	// 获取文件目录下所有文件
	a := []string{}
	rd, err := os.ReadDir(pathname)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return []string{}
	}
	PH := strings.Split(pathname, LinkPathStr)
	PHLen := len(PH)
	PathName := PH[PHLen-1]
	for _, fi := range rd {
		if fi.IsDir() {
			// fmt.Printf("[%s]\n", pathname+"/"+fi.Name())
			thisPath := strings.Join([]string{pathname, fi.Name()}, LinkPathStr)
			a = append(a, thisPath)
		} else {
			fileSuffix := path.Ext(fi.Name())
			if fileSuffix == ".jpg" || fileSuffix == ".jpeg" || fileSuffix == ".png" {
				a = append(a, strings.Join([]string{"/static", PathName, fi.Name()}, "/"))
			}
		}
	}
	return a
}
