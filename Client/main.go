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

func postFile(filename string, confYaml *utils.Config) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("code", confYaml.Code)
	bodyWriter.WriteField("computname", confYaml.ComputName)
	fileWriter, err := bodyWriter.CreateFormFile("image", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		time.Sleep(time.Duration(3) * time.Second)
		postFile(filename, confYaml)
	}

	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		time.Sleep(time.Duration(3) * time.Second)
		postFile(filename, confYaml)
	}
	defer fh.Close()

	//iocopy
	_, errs := io.Copy(fileWriter, fh)
	if errs != nil {
		// fmt.Println(err)
		time.Sleep(time.Duration(3) * time.Second)
		postFile(filename, confYaml)
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(confYaml.Host, contentType, bodyBuf)
	if err != nil {
		// fmt.Println(err)
		time.Sleep(time.Duration(10) * time.Second)
		postFile(filename, confYaml)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var m *Status
	json.Unmarshal(body, &m)
	// fmt.Println(m.Status)
	if m.Status != 1 {
		time.Sleep(time.Duration(10) * time.Second)
		postFile(filename, confYaml)
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
	confYaml, err := utils.CheckConfig(OS, CurrentPath)
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Duration(10) * time.Second)
		os.Exit(0)
	}
	postFile(fileName, confYaml)
}
