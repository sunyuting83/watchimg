package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
	"worldimg/Client/utils"
)

type Status struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func postFile(filename, gold, multiple string, confYaml *utils.Config) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("gold", gold)
	bodyWriter.WriteField("multiple", multiple)
	fileWriter, err := bodyWriter.CreateFormFile("image", filename)
	if err != nil {
		fmt.Println("0")
	}

	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("0")
	}
	defer fh.Close()

	//iocopy
	_, errs := io.Copy(fileWriter, fh)
	if errs != nil {
		fmt.Println("0")
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(confYaml.Host, contentType, bodyBuf)
	if err != nil {
		fmt.Println("0")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var m *Status
	json.Unmarshal(body, &m)
	// fmt.Println(m.Status)
	if m.Status != 1 {
		fmt.Println("0")
	} else {
		fmt.Println("1")
	}
	// return nil
}

// sample usage
func main() {
	var (
		file     string
		gold     string
		multiple string
	)
	flag.StringVar(&file, "file", "3000", "image path")
	flag.StringVar(&gold, "gold", "0", "gold")
	flag.StringVar(&multiple, "multiple", "0", "gold")
	flag.Parse()

	OS := runtime.GOOS
	LinkPathStr := "/"
	if OS == "windows" {
		LinkPathStr = "\\"
	}
	CurrentPath, _ := utils.GetCurrentPath()
	fileName := strings.Join([]string{CurrentPath, file}, LinkPathStr)
	confYaml, err := utils.CheckConfig(OS, CurrentPath)
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Duration(10) * time.Second)
		os.Exit(0)
	}
	postFile(fileName, gold, multiple, confYaml)
}
