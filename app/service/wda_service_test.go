package service

import (
	"context"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/electricbubble/gwda"
	extOpenCV "github.com/electricbubble/gwda-ext-opencv"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
	"time"
)

func TestWda(t *testing.T) {
	ctx := context.Background()
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	//mockDumpService := NewMockIDumpService(mockCtl)
	wdaService := NewWdaService()
	send := make(chan DumpJson)
	wdaService.Run(&ctx, "", "../../images", send)
}

func TestDevicer(t *testing.T) {
	ass := assert.New(t)
	devices, err := gwda.DeviceList()
	ass.NoError(err)
	findDevice, ok := slice.FindBy(devices, func(index int, item gwda.Device) bool {
		return item.SerialNumber() == "00008030-001D24223E20802E"
	})
	ass.Equal(ok, true)
	device, err := gwda.NewUSBDriver(nil, findDevice)
	ass.NoError(err)
	err = device.PressButton(gwda.DeviceButtonHome)
	ass.NoError(err)
}

func TestClickCountry(t *testing.T) {
	ass := assert.New(t)
	imagePath, err := filepath.Abs("../../images")
	ass.NoError(err)
	devices, err := gwda.DeviceList()
	ass.NoError(err)
	findDevice, ok := slice.FindBy(devices, func(index int, item gwda.Device) bool {
		return item.SerialNumber() == "00008030-001D24223E20802E"
	})
	ass.Equal(ok, true)
	device, err := gwda.NewUSBDriver(nil, findDevice)
	ass.NoError(err)
	err = device.PressButton(gwda.DeviceButtonHome)
	time.Sleep(1 * time.Second)
	ass.NoError(err)
	deviceExt, err := extOpenCV.Extend(device, 0.8)
	ass.NoError(err)
	err = deviceExt.Tap(filepath.Join(imagePath, "app", "appstore.jpg")) // 点击appstore
	ass.NoError(err)
	time.Sleep(1 * time.Second)
	err = deviceExt.Tap(filepath.Join(imagePath, "country", "us.jpg"))
	ass.NoError(err)
}
