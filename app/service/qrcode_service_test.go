package service

import (
	assert2 "github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

// import assert2 "github.com/stretchr/testify/assert"
func init() {
	rootDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	rootDir = filepath.Dir(filepath.Dir(rootDir))
	err = os.Chdir(rootDir)
	if err != nil {
		panic(err)
	}
}
func TestQrcodeService_Create(t *testing.T) {
	assert := assert2.New(t)
	qrcodeService := NewQrcodeService()
	err := qrcodeService.Create("https://baidu.com", Low, 256, "./qrcode.png")
	defer os.Remove("./qrcode.png")
	assert.NoError(err)
}
func TestQrcodeService_Parser(t *testing.T) {
	TestQrcodeService_Create(t)
	qrcodePath := "./qrcode.png"
	defer os.Remove(qrcodePath)
	assert := assert2.New(t)
	qrcodeService := NewQrcodeService()
	url, err := qrcodeService.Parser(qrcodePath)
	assert.NoError(err)
	assert.Equal(*url, "https://baidu.com")

}
