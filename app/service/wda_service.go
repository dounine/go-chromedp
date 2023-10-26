package service

//go:generate mockgen -destination=./wda_service_mock.go -package=service go-chromedp/app/service IWdaService
import (
	"context"
	"errors"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/duke-git/lancet/v2/slice"
	. "github.com/electricbubble/gwda"
	extOpenCV "github.com/electricbubble/gwda-ext-opencv"
	"go-chromedp/app/middleware"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	WdaStatusNormal WdaStatus = iota + 1
	WdaStatusDownloadBegin
	WdaStatusDownloading
	WdaStatusDownloadSuccess
	WdaStatusDownloadFail
)

type (
	IWdaService interface {
		Run(ctx *context.Context, udid string, imagePath string, receive <-chan DumpJson) <-chan WdaInfo
	}
	WdaStatus int
	WdaInfo   struct {
		json   DumpJson
		status WdaStatus
		err    error
	}
	WdaService struct {
		IWdaService
	}
)

func NewWdaService() *WdaService {
	return &WdaService{}
}

func (s *WdaService) ClickAlerts(devicer *WebDriver, alertPath string) {

}

// 自动点击多种安装按钮
func (s *WdaService) autoClickInstall(devicer WebDriver, imagePath string) {
	installs, err := fileutil.ListFileNames(filepath.Join(imagePath, "install"))
	if err != nil {
		panic(err)
	}
	deviceExt, err := extOpenCV.Extend(devicer, 0.8)
	if err != nil {
		panic(err)
	}
	for _, fileName := range installs {
		_ = deviceExt.Tap(filepath.Join(imagePath, "install", fileName))
		time.Sleep(1 * time.Second)
	}
}

// 是否匹配到此应用不可下载，也许要求的系统版本高
func (s *WdaService) matchNotInstall(devicer WebDriver, imagePath string) bool {
	countrys, err := fileutil.ListFileNames(filepath.Join(imagePath, "nonsupport"))
	if err != nil {
		panic(err)
	}
	deviceExt, err := extOpenCV.Extend(devicer, 0.8)
	if err != nil {
		panic(err)
	}
	if _, ok := slice.FindBy(countrys, func(index int, fileName string) bool {
		_, err := deviceExt.FindAllImageRect(filepath.Join(imagePath, "nonsupport", fileName))
		time.Sleep(500 * time.Millisecond)
		if err == nil {
			return true
		}
		return false
	}); ok {
		return true
	}
	return false
}
func homeScreen(devicer WebDriver, count int) {
	for i := 0; i < count; i++ {
		err := devicer.PressButton(DeviceButtonHome)
		if err != nil {
			panic(err)
		}
		time.Sleep(2 * time.Second)
	}
}
func (s *WdaService) switchCountry(devicer WebDriver, imagePath string, country string) {
	if country == "" {
		panic("国家不能为空")
	}
	middleware.Logger.Infof("返回首页")
	homeScreen(devicer, 2)
	deviceExt, err := extOpenCV.Extend(devicer, 0.8)
	if err != nil {
		panic(err)
	}
	middleware.Logger.Info("点击appstore")
	err = deviceExt.Tap(filepath.Join(imagePath, "app", "appstore.jpg")) // 点击appstore
	if err != nil {
		panic(err)
	}
	countrys, err := fileutil.ListFileNames(filepath.Join(imagePath, "country"))
	if err != nil {
		panic(err)
	}
	time.Sleep(2 * time.Second)
	middleware.Logger.Infof("切换国家: %s", country)
	if countryEntry, ok := slice.FindBy(countrys, func(index int, fileName string) bool {
		return fileName == country+".jpg"
	}); ok {
		err = deviceExt.Tap(filepath.Join(imagePath, "country", countryEntry)) // 点击地区
		if err != nil {
			panic(err)
		}
	}
	time.Sleep(2 * time.Second)
	homeScreen(devicer, 2)
}
func (s *WdaService) openLink(device WebDriver, imagePath string, link string) {
	deviceExt, err := extOpenCV.Extend(device, 0.8)
	if err != nil {
		panic(err)
	}
	err = deviceExt.Tap(filepath.Join(imagePath, "app", "activelink.jpg")) // 点击appstore
	if err != nil {
		panic(err)
	}
}

func (s *WdaService) appStoreDownload(ctx *context.Context, udid string, downloadPathDirs []string, imagePath string, device WebDriver, info DumpJson, send chan<- WdaInfo) (err error) {
	deviceExt, err := extOpenCV.Extend(device, 0.8)
	if err != nil {
		panic(err)
	}
	s.switchCountry(device, imagePath, info.Country)
	for {
		appInstalled := s.appExit(ctx, info.BundleId, udid)
		if appInstalled {
			middleware.Logger.Infof("应用下载完成")
			break
		}
		matchResult := s.matchNotInstall(device, imagePath)
		if matchResult {
			err = errors.New("不可下载的应用")
			middleware.Logger.Infof("匹配到不可下载的应用")
			break
		}
		s.openLink(device, imagePath, info.Link)
		time.Sleep(2 * time.Second)
		go s.autoClickInstall(device, imagePath)
		for index, fileName := range downloadPathDirs {
			rects, err := deviceExt.FindAllImageRect(filepath.Join(imagePath, "download", fileName))
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			}
			downloadPathDirs = slice.DeleteAt(downloadPathDirs, index)
			middleware.Logger.Infof("找到了: %s %+v", fileName, rects)
			send <- WdaInfo{
				json:   info,
				status: WdaStatusDownloading,
			}
			break
		}
		time.Sleep(1 * time.Second)
	}
	return
}
func (s *WdaService) appExit(ctx *context.Context, bundleId string, udid string) bool {
	result, err := exec.CommandContext(
		*ctx,
		"tidevice",
		"-u", udid,
		"appinfo", bundleId,
	).Output()
	if err != nil {
		if err.Error() == "exit status 1" {
			return false
		}
		panic(err)
	}
	if result == nil {
		return false
	}
	return string(result) != ""
}
func (s *WdaService) Run(ctx *context.Context, udid string, imagePath string, receive <-chan DumpJson) <-chan WdaInfo {
	devices, err := DeviceList()
	if err != nil {
		panic(err)
	}
	findDevice, ok := slice.FindBy(devices, func(index int, item Device) bool {
		return item.SerialNumber() == udid
	})
	if !ok {
		panic("没有找到对应的设备")
	}
	device, err := NewUSBDriver(nil, findDevice)
	send := make(chan WdaInfo)
	if err != nil {
		panic(err)
	}
	downloadPathDirs, err := fileutil.ListFileNames(filepath.Join(imagePath, "download"))
	if err != nil {
		panic(err)
	}
	slice.SortBy(downloadPathDirs, func(i, j string) bool {
		a, _ := strconv.Atoi(strings.Split(i, ".")[0])
		b, _ := strconv.Atoi(strings.Split(j, ".")[0])
		return a < b
	})
	go s.ClickAlerts(&device, filepath.Join(imagePath, "alert"))
	go func() {
		log := middleware.Logger.WithField("module", "wda")
		for {
			select {
			case info := <-receive:
				log.Infof("收到下载请求: %+v", info)
				appInstalled := s.appExit(ctx, info.BundleId, info.BundleId)
				log.Infof("应用是否已经安装: %+v", appInstalled)
				if appInstalled {
					send <- WdaInfo{
						json:   info,
						status: WdaStatusDownloadSuccess,
						err:    nil,
					}
					continue
				}
				send <- WdaInfo{
					json:   info,
					status: WdaStatusDownloadBegin,
					err:    nil,
				}
				// 开始下载
				err = s.appStoreDownload(ctx, udid, downloadPathDirs, imagePath, device, info, send)
				if err != nil {
					send <- WdaInfo{
						json:   info,
						status: WdaStatusDownloadFail,
						err:    err,
					}
				} else {
					send <- WdaInfo{
						json:   info,
						status: WdaStatusDownloadSuccess,
						err:    nil,
					}
				}
			}
		}
	}()
	return send
}
