package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/duke-git/lancet/v2/fileutil"
	assert2 "github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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
func TestQrodeServie_Parser2(t *testing.T) {
	assert := assert2.New(t)
	qrcodeService := NewQrcodeService()
	url, err := qrcodeService.Parser("./file/aa.png")
	assert.NoError(err)
	fmt.Println(*url)
}
func TestQrcodeService_QrcodeGet(t *testing.T) {
	assert := assert2.New(t)
	qrcodeService := NewQrcodeService()
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)
	allocExt, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()
	ctx, cancel = chromedp.NewContext(allocExt)
	defer cancel()
	js := `new Promise((resolve, reject) => {
				setTimeout(async () => {
					let response = await fetch('https://mp.weixin.qq.com/cgi-bin/scanloginqrcode?action=getqrcode&random='+new Date().getTime())
					let blob = await response.blob()
					let reader = new FileReader()
					reader.readAsDataURL(blob)
					reader.onloadend = function () {
						let base64data = reader.result
						resolve(base64data)
					}
					reader.onerror = function () {
						reject('blob转base64失败，请联系开发人员')
					}
				}, 1000); // 模拟异步操作
			});
	`
	var result string
	err := chromedp.Run(
		ctx,
		chromedp.Navigate("https://mp.weixin.qq.com"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			cookies, err := network.GetCookies().Do(ctx)
			if err != nil {
				panic(err)
			}
			for _, cookie := range cookies {
				println(cookie.Name + ":" + cookie.Value)
			}
			return nil
		}),
		chromedp.Evaluate(js, &result, func(ep *runtime.EvaluateParams) *runtime.EvaluateParams {
			return ep.WithAwaitPromise(true)
		}),
	)
	regexImage := regexp.MustCompile(`^data:image/(jpg|png);base64,`)
	//正则替换jpg或者png
	arrayBufferBytes, err := base64.StdEncoding.DecodeString(regexImage.ReplaceAllString(result, ""))
	if err != nil {
		panic(err)
	}
	fileutil.WriteBytesToFile("./qrcode.png", arrayBufferBytes)
	defer os.Remove("./qrcode.png")
	context, err := qrcodeService.Parser("./qrcode.png")
	assert.NoError(err)
	assert.NotEmptyf(context, "context is empty")
	assert.Equal(strings.HasPrefix(*context, "http://mp.weixin.qq.com"), true)
}
