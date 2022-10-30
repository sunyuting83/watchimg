package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"runtime"
	"strings"
	"worldimg/Client/utils"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type Status struct {
	Code int      `json:"code"`
	Data []string `json:"data"`
}

func postFile(filename string) {
	var (
		Host     string = "http://www.iinside.cn:7001/api_req"
		Password string = "8907"
		Reqmode  string = "ocr_pp"
	)
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("password", Password)
	bodyWriter.WriteField("reqmode", Reqmode)
	fileWriter, err := bodyWriter.CreateFormFile("image_ocr_pp", filename)
	if err != nil {
		fmt.Println("")
	}

	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("")
	}
	defer fh.Close()

	//iocopy
	_, errs := io.Copy(fileWriter, fh)
	if errs != nil {
		fmt.Println("")
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(Host, contentType, bodyBuf)
	if err != nil {
		fmt.Println("")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var m *Status
	json.Unmarshal(body, &m)
	// fmt.Println(m.Status)
	if m.Code == -1 {
		fmt.Println("")
	} else {
		if len(m.Data) > 0 {
			GBKStr, _ := Utf8ToGbk([]byte(m.Data[0]))
			fmt.Println(string(GBKStr))
		} else {
			fmt.Println("")
		}
	}
	// return nil
}

// sample usage
func main() {
	var (
		file string
	)
	flag.StringVar(&file, "file", "3000", "image path")
	flag.Parse()

	OS := runtime.GOOS
	LinkPathStr := "/"
	if OS == "windows" {
		LinkPathStr = "\\"
	}
	CurrentPath, _ := utils.GetCurrentPath()
	fileName := strings.Join([]string{CurrentPath, file}, LinkPathStr)
	postFile(fileName)
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
