package ocr

import (
	"github.com/otiai10/gosseract/v2"
)

func main() {
	// gin.SetMode(gin.ReleaseMode)

	// app := gin.New()
	// app.GET("/", func(c *gin.Context) {
	// 	c.String(403, "403")
	// })
	// app.POST("/ocr", func(c *gin.Context) {
	// 	file, _, err := c.Request.FormFile("image")
	// 	defer file.Close()
	// 	if err != nil {
	// 		c.String(200, "1data")
	// 		return
	// 	}

	// 	buf := bytes.NewBuffer(nil)
	// 	if _, err := io.Copy(buf, file); err != nil {
	// 		c.String(200, "2data")
	// 		return
	// 	}
	// 	data := GetWord(buf.Bytes())
	// 	c.String(200, data)
	// })
	// app.Run(strings.Join([]string{":", "65530"}, ""))
}

func GetWord(b []byte) string {
	client := gosseract.NewClient()
	defer client.Close()
	// client.SetImage("/home/sun/Pictures/photo_2022-09-26_01-32-46.jpg")
	client.SetImageFromBytes(b)
	text, err := client.Text()
	if err != nil {
		return err.Error()
	}
	return text
}
