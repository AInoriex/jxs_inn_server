package qrcode

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
)

// 解码base64
func DecodeBase64(encoded string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

// 写入图片
func WriteToImageFile(data []byte, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// DecodeBase64ToImage 将base64字符串解析成图片，并将其写入文件
// base64: base64字符串
// exportImageName: 输出图片的文件名
// 返回nil表示成功，否则返回错误
func DecodeBase64ToImage(base64String string, exportImageName string) error {
	var err error
	var filename string
	if !strings.Contains(exportImageName, ".") {
		_ext := "jpg"                                           // 默认JPEG格式
		filename = fmt.Sprintf("%s.%s", exportImageName, _ext) // 输出文件名和格式
	} else {
		filename = exportImageName
	}
	data, err := DecodeBase64(base64String)
	if err != nil {
	    return fmt.Errorf("DecodeBase64ToImage 解析base64字符串失败, %s", err.Error())
	}
	err = WriteToImageFile(data, filename)
	if err != nil {
		return fmt.Errorf("DecodeBase64ToImage 写入图片失败, %s", err.Error())
	}
	return nil
}

func PrintQrCodeImage(data string) error {
	// qr, err:= qrcode.New("https://blog.csdn.net/a6100china/article/details/137829574?spm=1001.2014.3001.5502", qrcode.Medium)
	qr, err:= qrcode.New(data, qrcode.Medium)
	if err != nil {
		return fmt.Errorf("PrintQrCodeImage 打印二维码失败, %s", err.Error())
	}
	fmt.Println(qr.ToSmallString(false))
	return nil
}
