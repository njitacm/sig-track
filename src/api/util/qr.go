package util

import (
	qrcode "github.com/skip2/go-qrcode"
)

func QRCodeGenToFile(content, filename string) {
	png, err := qrcode.Encode(content, qrcode.Medium, 256)
	Check(err)
	FileWrite(filename, png)
}

func QRCodeGen(content string) []byte {
	png, err := qrcode.Encode(content, qrcode.Medium, 256)
	Check(err)
	return png
}
