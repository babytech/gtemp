package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func ReadFile(filename string) (string, error) {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("ReadFile Fail", err)
		return "", err
	}
	return string(f), nil
}

func WriteFile(filename string, data []byte) {
	fmt.Println("Write File: ", filename)
	fmt.Println("Write data[]: ", data)
	err := ioutil.WriteFile(filename, data, 0666)
	if err != nil {
		fmt.Println("WriteFile Fail", err)
	} else {
		fmt.Println("WriteFile Success!")
	}
}

func CheckFile(fileName string) string {
	dir := filepath.Dir(fileName)
	base := filepath.Base(fileName)
	if _, err := os.Stat(dir); err == nil {
		//fmt.Println("Directory path exists", dir)
	} else {
		fmt.Println("Directory path not exists ", dir)
		err := os.MkdirAll(dir, 0711)
		if err != nil {
			log.Println("Error creating directory")
			log.Println(err)
			return ""
		}
	}
	return filepath.Join(dir, base)
}

func IntToBytes(value uint32) []byte {
	bytebuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytebuffer, binary.BigEndian, value)
	return bytebuffer.Bytes()
}

func BytesToInt(bs []byte) uint32 {
	bytebuffer := bytes.NewBuffer(bs)
	var data uint32
	_ = binary.Read(bytebuffer, binary.BigEndian, &data)
	return data
}

func UploadFile(fileName string, ip string, port string) (io.ReadCloser, error) {
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)
	fileWriter, _ := bodyWriter.CreateFormFile("files", fileName)
	file, _ := os.Open(fileName)
	defer file.Close()
	if _, err := io.Copy(fileWriter, file); err != nil {
		log.Fatal(err)
	}
	contentType := bodyWriter.FormDataContentType()
	err := bodyWriter.Close()
	if err != nil {
		return nil, err
	}
	url := "http://" + ip + ":" + port + "/upload"
	fmt.Println("URL: ", url)
	resp, err := http.Post(url, contentType, bodyBuffer)
	if resp == nil || err != nil {
		return nil, err
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	log.Println("response from server status: ", resp.Status)
	log.Println("response string: ", string(respBody))
	return resp.Body, nil
}
