package main

import (
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/chai2010/webp"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

type FileType int

const (
	PNG FileType = iota
	JPG
	GIF
	WEBP
	BMP
	TIFF
	ERR
)

func getFileType(input string) FileType {
	switch input {
	case "jpg":
		fallthrough
	case "jpeg":
		return JPG
	case "png":
		return PNG
	case "gif":
		return GIF
	case "bmp":
		return BMP
	case "webp":
		return WEBP
	case "tiff":
		return TIFF
	default:
		return ERR
	}
}

func getFileExtension(input FileType) string {
	switch input {
	case JPG:
		return "jpg"
	case PNG:
		return "png"
	case GIF:
		return "gif"
	case BMP:
		return "bmp"
	case WEBP:
		return "webp"
	case TIFF:
		return "tiff"
	default:
		return ""
	}
}

func convert(files []string, outputDir string, fileType FileType) {
	var wg sync.WaitGroup
	for _, currPath := range files {
		wg.Add(1)
		go convertFile(&wg, currPath, outputDir, fileType)
	}

	wg.Wait()
}

func convertFile(wg *sync.WaitGroup, currPath string, outputDir string, fileType FileType) {
	// call done when finished
	defer wg.Done()

	ext := strings.ToLower(filepath.Ext(currPath))
	newExt := getFileExtension(fileType)

	_, filename := filepath.Split(currPath)
	filenameNoExt := filename[0 : len(filename)-len(ext)]
	newFileName := filenameNoExt + "." + newExt
	newFilePath := outputDir + "/" + newFileName

	// validate file type
	startType := getFileType(ext[1:])
	if startType == ERR {
		logAndExit(errors.New("input file type not valid"))
	}

	// open files
	file, err := os.Open(currPath)
	if err != nil {
		logAndExit(err)
	}
	defer file.Close()

	outFile := openOrCreate(newFilePath)
	defer outFile.Close()

	// decode
	imageData, _, err := image.Decode(file)
	if err != nil {
		logAndExit(errors.New("error decoding image"))
	}

	// encode in new type
	switch fileType {
	case JPG:
		err := jpeg.Encode(outFile, imageData, nil)
		if err != nil {
			logAndExit(errors.New("error converting to jpeg"))
		}
	case PNG:
		err := png.Encode(outFile, imageData)
		if err != nil {
			logAndExit(err)
		}
	case WEBP:
		if err := webp.Encode(outFile, imageData, nil); err != nil {
			logAndExit(errors.New("error converting to webp"))
		}
	case GIF:
		err := gif.Encode(outFile, imageData, nil)
		if err != nil {
			logAndExit(errors.New("error converting to gif"))
		}
	case BMP:
		err := bmp.Encode(outFile, imageData)
		if err != nil {
			logAndExit(errors.New("error converting to bmp"))
		}
	case TIFF:
		err := tiff.Encode(outFile, imageData, nil)
		if err != nil {
			fmt.Println(err)
			logAndExit(errors.New("error converting to tiff"))
		}
	}
}

func openOrCreate(filename string) *os.File {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			logAndExit(errors.New("error creating output file"))
		}
		return file
	} else {
		file, err := os.Open(filename)
		if err != nil {
			logAndExit(errors.New("error opening output file"))
		}
		return file
	}
}

func logAndExit(a error) {
	fmt.Println(a)
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
			if fileSuffix == ".jpg" || fileSuffix == ".bmp" || fileSuffix == ".png" {
				a = append(a, strings.Join([]string{pathname, PathName, fi.Name()}, "/"))
			}
		}
	}
	return a
}

func main() {
	path := "/home/sun/Pictures/dzz/20221009/"
	AllFile := GetAllFile(path)
	outDir := "/home/sun/Pictures/dzz/outfile"
	convert(AllFile, outDir, TIFF)
}
