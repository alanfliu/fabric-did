package Utils

import (
	"fmt"
	"os"
)

func ReadFile(filename string) []byte {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	info, _ := file.Stat()
	result := make([]byte, info.Size())
	file.Read(result)
	file.Close()
	return result
}

func ReadPublicPem(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("打开私钥文件 - private.pem 失败!!!")
		return "", err
	}
	fileInfo, _ := file.Stat()
	privatePemBytes := make([]byte, fileInfo.Size())
	_, err = file.Read(privatePemBytes)
	if err != nil {
		fmt.Println("读文件内容失败!!!")
		return "", err
	}

	file.Close()
	return string(privatePemBytes), nil
}
