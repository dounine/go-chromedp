package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/duke-git/lancet/v2/slice"
	"go-chromedp/app/middleware"
	"go-chromedp/app/models"
	"time"
)

//go:generate mockgen -destination=./dump_engine_service_mock.go -package=service go-chromedp/app/service IDumpEngineService

const (
	EngineStatusNormal EngineStatus = iota + 1
	EngineStatusUploading
	EngineStatusDumping
)

type (
	IDumpEngineService interface {
		Run(ctx *context.Context, udid string, imagePath string) error
		Dumps(ctx *context.Context) ([]DumpJson, error)
	}
	DumpAppInfo struct {
		info       DumpJson
		wdaStatus  WdaStatus
		dumpStatus DumpStatus
	}
	EngineStatus      int
	DumpEngineService struct {
		IDumpEngineService
		dumpService IDumpService
		wdaService  IWdaService
		http        *models.P8Http
	}
)

func NewDumpEngineService() *DumpEngineService {
	return &DumpEngineService{
		dumpService: NewDumpService(),
		wdaService:  NewWdaService(),
		http:        models.NewP8Http(),
	}
}
func (s *DumpEngineService) Run(ctx *context.Context, udid string, imagePath string) error {
	dumpSend := make(chan DumpJson)
	dumpReceive := s.dumpService.Run(ctx, udid, imagePath, dumpSend)
	wdaSend := make(chan DumpJson)
	wdaReceive := s.wdaService.Run(ctx, udid, imagePath, wdaSend)
	go func() {
		log := middleware.Logger
		var activeDumpInfo *DumpAppInfo
		list := make([]DumpAppInfo, 0)
		status := EngineStatusNormal
		for {
			listDumps, err := s.Dumps(ctx)
			if err != nil {
				log.Errorf("获取 dump 数据失败: %s", err)
				continue
			}
			for _, info := range listDumps {
				exitInfo, ok := slice.FindBy(list, func(index int, item DumpAppInfo) bool {
					return item.info.Appid == info.Appid && item.info.BundleId == info.BundleId && item.info.Country == info.Country
				})
				if !ok {
					exitInfo = DumpAppInfo{
						info:       info,
						wdaStatus:  WdaStatusNormal,
						dumpStatus: DumpStatusNormal,
					}
					list = append(list, exitInfo)
				}
				if exitInfo.wdaStatus != WdaStatusDownloadSuccess && s.dumpService.AppExit(ctx, udid, info.BundleId) {
					exitInfo.wdaStatus = WdaStatusDownloadSuccess
				}
			}
			select {
			case dumpInfo := <-dumpReceive:
				log.Infof("receive dumpInfo: %+v", dumpInfo)
				switch dumpInfo.status {
				case DumpStatusBegin:
					list = slice.Map(list, func(index int, item DumpAppInfo) DumpAppInfo {
						if item.info.Appid == dumpInfo.json.Appid && item.info.BundleId == dumpInfo.json.BundleId && item.info.Country == dumpInfo.json.Country {
							item.dumpStatus = DumpStatusBegin
						}
						return item
					})
					log.Infof("开始砸壳 %+v", dumpInfo.json)
				case DumpStatusSuccess:
					list = slice.Map(list, func(index int, item DumpAppInfo) DumpAppInfo {
						if item.info.Appid == dumpInfo.json.Appid && item.info.BundleId == dumpInfo.json.BundleId && item.info.Country == dumpInfo.json.Country {
							item.dumpStatus = DumpStatusSuccess
						}
						return item
					})
					//TODO 异步下载下一个app
					if status == EngineStatusUploading || status == EngineStatusNormal {
						unDownloads := slice.Filter(list, func(index int, item DumpAppInfo) bool {
							return item.wdaStatus == WdaStatusNormal
						})
						if len(unDownloads) > 0 {
							wdaSend <- unDownloads[0].info
						}
					}
					log.Infof("砸壳 %+v 成功", dumpInfo.json)
				case DumpStatusFail:
					list = slice.Map(list, func(index int, item DumpAppInfo) DumpAppInfo {
						if item.info.Appid == dumpInfo.json.Appid && item.info.BundleId == dumpInfo.json.BundleId && item.info.Country == dumpInfo.json.Country {
							item.dumpStatus = DumpStatusFail
						}
						return item
					})
					log.Errorf("砸壳 %+v 失败: %s", dumpInfo.json, dumpInfo.err)
				case DumpStatusUploadBegin:
					list = slice.Map(list, func(index int, item DumpAppInfo) DumpAppInfo {
						if item.info.Appid == dumpInfo.json.Appid && item.info.BundleId == dumpInfo.json.BundleId && item.info.Country == dumpInfo.json.Country {
							item.dumpStatus = DumpStatusUploadBegin
						}
						return item
					})
					status = EngineStatusUploading
					log.Infof("开始上传 %+v", dumpInfo.json)
				case DumpStatusUploadSuccess:
					status = EngineStatusNormal
					activeDumpInfo = nil
					list = slice.Filter(list, func(index int, item DumpAppInfo) bool {
						return item.dumpStatus != DumpStatusUploadSuccess
					})
					log.Infof("上传 %+v 成功", dumpInfo.json)
				case DumpStatusUploadFail:
					list = slice.Map(list, func(index int, item DumpAppInfo) DumpAppInfo {
						if item.info.Appid == dumpInfo.json.Appid && item.info.BundleId == dumpInfo.json.BundleId && item.info.Country == dumpInfo.json.Country {
							item.dumpStatus = DumpStatusUploadFail
						}
						return item
					})
					log.Errorf("上传 %+v 失败: %s", dumpInfo.json, dumpInfo.err)
				}
			case wdaInfo := <-wdaReceive:
				middleware.Logger.Infof("receive wdaInfo: %+v", wdaInfo)
				switch wdaInfo.status {
				case WdaStatusDownloadBegin:
					list = slice.Map(list, func(index int, item DumpAppInfo) DumpAppInfo {
						if item.info.Appid == wdaInfo.json.Appid && item.info.BundleId == wdaInfo.json.BundleId && item.info.Country == wdaInfo.json.Country {
							item.wdaStatus = WdaStatusDownloadBegin
						}
						return item
					})
				case WdaStatusDownloading:
					list = slice.Map(list, func(index int, item DumpAppInfo) DumpAppInfo {
						if item.info.Appid == wdaInfo.json.Appid && item.info.BundleId == wdaInfo.json.BundleId && item.info.Country == wdaInfo.json.Country {
							item.wdaStatus = WdaStatusDownloading
						}
						return item
					})
				case WdaStatusDownloadSuccess:
					//下载成功
					activeDumpInfo.wdaStatus = WdaStatusDownloadSuccess
					list = slice.Map(list, func(index int, item DumpAppInfo) DumpAppInfo {
						if item.info.Appid == wdaInfo.json.Appid && item.info.BundleId == wdaInfo.json.BundleId && item.info.Country == wdaInfo.json.Country {
							item.wdaStatus = WdaStatusDownloadSuccess
						}
						return item
					})
					//TODO 异步下载下一个app
					if status == EngineStatusUploading || status == EngineStatusNormal {
						unDownloads := slice.Filter(list, func(index int, item DumpAppInfo) bool {
							return item.wdaStatus == WdaStatusNormal
						})
						if len(unDownloads) > 0 {
							wdaSend <- unDownloads[0].info
						}
					}
				case WdaStatusDownloadFail:
					list = slice.Map(list, func(index int, item DumpAppInfo) DumpAppInfo {
						if item.info.Appid == wdaInfo.json.Appid && item.info.BundleId == wdaInfo.json.BundleId && item.info.Country == wdaInfo.json.Country {
							item.wdaStatus = WdaStatusDownloadFail
						}
						return item
					})
				}
			default:
				if activeDumpInfo == nil && len(list) > 0 {
					activeDumpInfo = &list[0]
					if activeDumpInfo.wdaStatus == WdaStatusDownloadSuccess {
						//已经下载，开始砸壳
						dumpSend <- activeDumpInfo.info
					} else {
						//下载
						wdaSend <- activeDumpInfo.info
					}
				}
			}
			time.Sleep(1 * time.Second)
		}
	}()
	return nil
}

func (s *DumpEngineService) Dumps(ctx *context.Context) ([]DumpJson, error) {
	queryUrl := `https://api.ipadump.com/dump/dumps?limit=100`
	request := s.http.Get(ctx, queryUrl)
	response, _, err := request.End()
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New("请求失败")
	}
	mockData := `{"code":0,"data":[{"appid":"6446317461","country":"cn","version":"1.0.15","name":"英雄旅途","icon":"https://is1-ssl.mzstatic.com/image/thumb/Purple126/v4/2a/f8/5b/2af85bdb-25e3-4b68-44d2-77522fd3a3ee/AppIcon-0-0-1x_U007emarketing-0-0-0-7-0-0-sRGB-0-0-0-GLES2_U002c0-512MB-85-220-0-0.png/512x512bb.jpg","link":"https://apps.apple.com/cn/app/%E8%8B%B1%E9%9B%84%E6%97%85%E9%80%94/id6446317461?uo=4","count":1,"status":1,"time":"2023-10-25 22:44:15","latest":1,"bundleId":"com.zyxk.wzgj","size":183776256,"price":0}]}`
	wrap := DumpResponsea[[]DumpJson]{}
	if err := json.Unmarshal([]byte(mockData), &wrap); err != nil {
		return nil, err
	}
	return wrap.Data, nil
}
