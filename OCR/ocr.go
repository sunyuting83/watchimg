package ocr

import (
	"github.com/otiai10/gosseract/v2"
)

var (
	Client *gosseract.Client
)

func init() {
	Client = gosseract.NewClient()
	Client.SetWhitelist("0123456789")
}
