package service

import (
	"context"
	"github.com/golang/mock/gomock"
	assert2 "github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	assert := assert2.New(t)
	//ctx := context.Background()
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	//mockDumpService := NewMockIDumpService(mockCtl)
	dumpService := NewDumpService()
	config := dumpService.Config()
	assert.Equal(config.Section("qiniu").Key("bucket").String(), "ipadump-ipas")
}
func TestWriteFile(t *testing.T) {
	assert := assert2.New(t)
	filePath := "../../dump/123.ipa"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		assert.Error(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
		err = os.Remove(filePath)
		if err != nil {
			panic(err)
		}
	}(file)
	size, err := file.WriteString("test")
	if err != nil {
		assert.Error(err)
	}
	assert2.Equal(t, size, 4)
}
func TestUploadIpa(t *testing.T) {
	assert := assert2.New(t)
	ctx := context.Background()
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	//mockDumpService := NewMockIDumpService(mockCtl)
	dumpService := NewDumpService()
	//create file
	filePath := "../../dump/123.ipa"
	err := os.Remove(filePath)
	if err != nil {
		assert.Error(err)
	}
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		assert.Error(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
		err = os.Remove(filePath)
		if err != nil {
			panic(err)
		}
	}(file)
	_, err = file.WriteString("test")
	if err != nil {
		assert.Error(err)
	}
	err = dumpService.UploadIpa(&ctx, "../../dump", DumpJson{
		BundleId: "com.a.b",
		Appid:    "123",
		Country:  "test",
		Version:  "1.0",
		Name:     "测试上传",
		Icon:     "https://abc.png",
	})
	if err != nil {
		assert.Error(err)
	}
}
