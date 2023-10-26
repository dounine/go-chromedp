package service

//go:generate mockgen -destination=./dump_service_mock.go -package=service go-chromedp/app/service IDumpService
import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/pkg/errors"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"go-chromedp/app/middleware"
	"go-chromedp/app/models"
	"go-chromedp/app/util"
	"gopkg.in/ini.v1"
	"os/exec"
	"regexp"
)

const (
	DumpStatusNormal DumpStatus = iota + 1
	DumpStatusBegin
	DumpStatusSuccess
	DumpStatusFail
	DumpStatusUploadBegin
	DumpStatusUploadSuccess
	DumpStatusUploadFail
)

const (
	DumpUpdateStatusDumping DumpUpdateStatus = iota + 1
	DumpUpdateStatusFinish
)

type (
	IDumpService interface {
		Config() *ini.File
		Run(ctx *context.Context, udid string, dumpPath string, receive <-chan DumpJson) <-chan DumpInfo
		QueryDumpInfo(ctx *context.Context, appid string, country string) (*DumpJson, error)
		QueryAppstoreInfo(ctx *context.Context, appid string, country string) (AppstoreInfo, error)
		AppExit(ctx *context.Context, bundleId string, udid string) bool
		MergeFileName(ctx *context.Context, dumpInfo DumpJson) string
		UploadIpa(ctx *context.Context, dumpPath string, dumpInfo DumpJson) error
		DumpApp(ctx *context.Context, dumpPath string, info DumpJson) error
	}
	TupleData interface {
		any | []any
	}
	DumpResponsea[T TupleData] struct {
		Code int `json:"code"`
		Data T   `json:"data"`
	}

	DumpService struct {
		IDumpService
		http *models.P8Http
	}
	DumpStatus int
	DumpInfo   struct {
		json   DumpJson
		status DumpStatus
		err    error
	}

	AppstoreInfoResult struct {
		artworkUrl60                       string
		artworkUrl512                      string
		artworkUrl100                      string
		artistViewUrl                      string
		screenshotUrls                     []string
		isGameCenterEnabled                bool
		features                           []string
		supportedDevices                   []string
		advisories                         []string
		ipadScreenshotUrls                 []string
		appletvScreenshotUrls              []string
		kind                               string
		trackViewUrl                       string
		minimumOsVersion                   string
		languageCodesISO2A                 []string
		fileSizeBytes                      string
		formattedPrice                     string
		contentAdvisoryRating              string
		averageUserRatingForCurrentVersion float64
		userRatingCountForCurrentVersion   int
		averageUserRating                  float64
		trackContentRating                 string
		trackCensoredName                  string
		currency                           string
		releaseNotes                       string
		artistId                           int
		artistName                         string
		genres                             []string
		price                              float64
		description                        string
		currentVersionReleaseDate          string
		isVppDeviceBasedLicensingEnabled   bool
		bundleId                           string
		genreIds                           []string
		releaseDate                        string
		primaryGenreName                   string
		primaryGenreId                     int
		version                            string
		wrapperType                        string
		sellerName                         string
		trackId                            int
		trackName                          string
		userRatingCount                    int
	}

	AppstoreInfo struct {
		resultCount int
		results     []AppstoreInfoResult
	}

	AddVersionInfo struct {
		appid    string
		version  string
		country  string
		push     int
		download int
		size     int
		official int
		des      string
		file     string
	}

	DumpUpdateStatus int

	DumpUpdateInfo struct {
		appid    string
		version  string
		country  string
		name     string
		lname    string
		icon     string
		price    string
		genres   string
		des      string
		latest   int
		status   DumpUpdateStatus
		bundleId string
	}

	DumpJson struct {
		BundleId string `json:"bundleId"`
		Appid    string `json:"appid"`
		Country  string `json:"country"`
		Version  string `json:"version"`
		Name     string `json:"name"`
		Icon     string `json:"icon"`
		Link     string `json:"link"`
		Status   int    `json:"status"`
		Time     string `json:"time"`
		Size     int    `json:"size"`
	}
)

func NewDumpService() *DumpService {
	return &DumpService{
		http: models.NewP8Http(),
	}
}

func (c DumpJson) String() string {
	str, _ := json.Marshal(c)
	return string(str)
}

func (s *DumpService) AppExit(ctx *context.Context, bundleId string, udid string) bool {
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

func (s *DumpService) Run(ctx *context.Context, udid string, dumpPath string, receive <-chan DumpJson) <-chan DumpInfo {
	fileUtil := util.NewFileUtil()
	log := middleware.Logger
	send := make(chan DumpInfo)
	go func() {
		for {
			select {
			case info := <-receive: //接收到需要dump的app信息
				log.Infof("清空dump文件夹: %+v", info)
				err := fileUtil.RemoveAllFilesInFolder(dumpPath)
				if err != nil {
					send <- DumpInfo{
						json:   info,
						status: DumpStatusFail,
						err:    err,
					}
					break
				}
				log.Info("清空dump文件夹完成")
				log.Info("开始dump")
				send <- DumpInfo{
					json:   info,
					status: DumpStatusBegin,
					err:    err,
				}
				err = s.DumpApp(ctx, dumpPath, info)
				if err != nil {
					send <- DumpInfo{
						json:   info,
						status: DumpStatusFail,
						err:    err,
					}
					break
				}
				send <- DumpInfo{
					json:   info,
					status: DumpStatusSuccess,
					err:    err,
				}
				log.Info("dump完成")
				log.Info("开始上传")
				send <- DumpInfo{
					json:   info,
					status: DumpStatusUploadBegin,
					err:    err,
				}
				err = s.UploadIpa(ctx, dumpPath, info)
				if err != nil {
					send <- DumpInfo{
						json:   info,
						status: DumpStatusUploadFail,
						err:    err,
					}
					break
				}
				send <- DumpInfo{
					json:   info,
					status: DumpStatusUploadSuccess,
					err:    err,
				}
				log.Info("上传完成")
			}
		}
	}()
	return send
}
func (s *DumpService) upsertVersion(ctx *context.Context, info AddVersionInfo) error {
	request := s.http.Post(ctx, "https://api.ipadump.com/version/upsert")
	request.Body = info
	response, data, err := request.End()
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return errors.New("服务器错误")
	}
	return json.Unmarshal(data, &info)
}
func (s *DumpService) updateDumpStatus(ctx *context.Context, info DumpUpdateInfo) error {
	request := s.http.Post(ctx, "https://api.ipadump.com/dump/update")
	request.Body = info
	response, data, err := request.End()
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return errors.New("服务器错误")
	}
	return json.Unmarshal(data, &info)
}
func (s *DumpService) QueryAppstoreInfo(ctx *context.Context, appid string, country string) (AppstoreInfo, error) {
	info := AppstoreInfo{
		resultCount: 0,
	}
	request := s.http.Post(ctx, fmt.Sprintf("https://itunes.apple.com/lookup?id=%s&country=%s", appid, country))
	response, data, err := request.End()
	if err != nil {
		return info, err
	}
	if response.StatusCode != 200 {
		return info, errors.New("服务器错误")
	}
	err = json.Unmarshal(data, &info)
	if err != nil {
		return info, err
	}
	return info, nil
}
func (s *DumpService) QueryDumpInfo(ctx *context.Context, appid string, country string) (*DumpJson, error) {
	request := s.http.Get(ctx, fmt.Sprintf("https://api.ipadump.com/app/info?appid=%s&country=%s", appid, country))
	response, data, err := request.End()
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New("服务器错误")
	}
	var appInfoWrap DumpResponsea[*DumpJson]
	err = json.Unmarshal(data, &appInfoWrap)
	if err != nil {
		return nil, err
	}
	return appInfoWrap.Data, nil
}
func (s *DumpService) MergeFileName(ctx *context.Context, dumpInfo DumpJson) string {
	//1.先查询服务器上app的应用信息
	name := dumpInfo.Name
	info, err := s.QueryDumpInfo(ctx, dumpInfo.Appid, dumpInfo.Country)
	if err != nil {
		panic(err)
	}
	if info != nil {
		name = info.Name
	} else {
		//过滤非法字符串
		nonChinesePattern := "[^\\u4e00-\\u9fa5a-zA-Z0-9-|()&+ 、：]"
		regex := regexp.MustCompile(nonChinesePattern)
		name = regex.ReplaceAllString(name, "")
	}
	return name
}
func (s *DumpService) Config() *ini.File {
	//currentDir, err := os.Getwd()
	//if err != nil {
	//	panic(err)
	//}
	cfg, err := ini.Load("../../app.ini")
	if err != nil {
		panic(err)
	}
	return cfg
}
func (s *DumpService) UploadIpa(ctx *context.Context, dumpPath string, dumpInfo DumpJson) error {
	log := middleware.Logger
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Region = &storage.ZoneHuanan
	config := s.Config().Section("qiniu")
	resumeUploader := storage.NewResumeUploaderV2(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.RputV2Extra{
		Notify: func(partNumber int64, ret *storage.UploadPartsRet) {
			log.Infof("上传进度:%d%%", partNumber)
		},
	}
	putPolicy := storage.PutPolicy{
		Scope: config.Key("bucket").String(),
	}
	mac := qbox.NewMac(config.Key("accessKey").String(), config.Key("secretKey").String())
	upToken := putPolicy.UploadToken(mac)
	ipaFile := dumpPath + "/" + dumpInfo.Appid + ".ipa"
	if !fileutil.IsExist(ipaFile) {
		panic(errors.New("ipa文件不存在"))
	}
	fileName := fmt.Sprintf("ipadump.com_%s_%s.ipa", dumpInfo.Appid, dumpInfo.Version)
	err := resumeUploader.PutFile(*ctx, &ret, upToken, fmt.Sprintf("ipas/%s/%s/%s", dumpInfo.Country, dumpInfo.Appid, fileName), ipaFile, &putExtra)
	if err != nil {
		panic(err)
	}
	log.Infof("上传成功：%s", ret.Key)
	return nil
}
func (s *DumpService) DumpApp(ctx *context.Context, dumpPath string, info DumpJson) error {
	command := exec.CommandContext(*ctx, "python3", "dump.py", info.BundleId, "-o", fmt.Sprintf("%s/%s", dumpPath, info.Appid))
	stdout, err := command.StdoutPipe()
	if err != nil {
		return err
	}
	err = command.Start()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		middleware.Logger.Info(scanner.Text())
	}
	return command.Wait()
}
