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
	"time"
	"worldimg/Client/utils"
)

type Status struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func postFile(filename, gold, multiple, exptime string, confYaml *utils.Config) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("0")
		return
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", filename)
	if err != nil {
		fmt.Println("0")
		return
	}
	_, err = io.Copy(part, file)
	_ = writer.WriteField("gold", gold)
	_ = writer.WriteField("multiple", multiple)
	_ = writer.WriteField("exptime", exptime)
	err = writer.Close()
	if err != nil {
		fmt.Println("0")
		return
	}

	client := &http.Client{
		Timeout: time.Duration(10 * time.Second),
	}
	req, _ := http.NewRequest("POST", confYaml.Host, body)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("sec-ch-ua-platform", "Windows")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("0")
		return
	}
	defer resp.Body.Close()
	respByte, _ := io.ReadAll(resp.Body)

	var m *Status
	json.Unmarshal(respByte, &m)
	if m.Status != 1 {
		fmt.Println("0")
		return
	} else {
		fmt.Println("1")
		return
	}
	// return nil
}

// sample usage
func main() {
	var (
		file     string
		gold     string
		multiple string
		exptime  string
	)
	flag.StringVar(&file, "file", "3000", "image path")
	flag.StringVar(&gold, "gold", "0", "gold")
	flag.StringVar(&multiple, "multiple", "0", "multiple")
	flag.StringVar(&exptime, "exptime", "0", "exptime")
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
	postFile(fileName, gold, multiple, exptime, confYaml)
}
