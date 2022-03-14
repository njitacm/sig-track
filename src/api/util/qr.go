package util

import (
	_ "github.com/njitacm/sig-track/src/api/util"
	qrcode "github.com/skip2/go-qrcode"
)

func QRCodeGen(content, filename string) {
	png, err := qrcode.Encode(content, qrcode.Medium, 256)
	Check(err)
	FileWrite(filename, png)
}
