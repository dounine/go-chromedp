package main

import (
	"context"
	"encoding/base64"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/duke-git/lancet/v2/fileutil"
	"strings"
	"time"
)

func main() {
	//cmd.Execute()
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)

	allocExt, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(allocExt)
	defer cancel()
	//qrcodeService := service.NewQrcodeService()

	js := `new Promise((resolve, reject) => {
				setTimeout(async () => {
					let response = await fetch('https://mp.weixin.qq.com/cgi-bin/scanloginqrcode?action=getqrcode&random=1698639598411')
					let blob = await response.blob()
					let reader = new FileReader()
					reader.readAsDataURL(blob)
					reader.onloadend = function () {
						let base64data = reader.result
						resolve(base64data)
					}
				}, 1000); // 模拟异步操作
			});
	`
	//qrcodeJs := `(await fetch('https://mp.weixin.qq.com/cgi-bin/scanloginqrcode?action=getqrcode&random=1698639598411')).text()`
	//var qrcodeUrl *string
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
		//截图
		//chromedp.ActionFunc(func(ctx context.Context) error {
		//	return retry.Retry(func() error {
		//		log := middleware.Logger
		//		var buf []byte
		//		if err := chromedp.Screenshot("img.login__type__container__scan__qrcode", &buf, chromedp.NodeVisible).Do(ctx); err != nil {
		//			log.Errorf("chromedp.Screenshot err: %s", err.Error())
		//			return err
		//		}
		//		err := fileutil.WriteBytesToFile("./qrcode.png", buf)
		//		if err != nil {
		//			log.Errorf("fileutil.WriteBytesToFile err: %s", err.Error())
		//			return err
		//		}
		//		defer os.Remove("./qrcode.png")
		//		qrcodeUrl, err := qrcodeService.Parser("./qrcode.png")
		//		if err != nil {
		//			log.Errorf("qrcodeService.Parser err: %s", err.Error())
		//			return err
		//		}
		//		if qrcodeUrl != nil && *qrcodeUrl != "" {
		//			log.Errorf("qrcodeUrl: %s", *qrcodeUrl)
		//			return nil
		//		}
		//		return errors.New("qrcode not found")
		//	}, retry.RetryTimes(3))
		//}),
	)
	arrayBufferBytes, err := base64.StdEncoding.DecodeString(strings.Replace(result, "data:image/jpg;base64,", "", -1))
	if err != nil {
		panic(err)
	}
	fileutil.WriteBytesToFile("./qrcode.png", arrayBufferBytes)
	//fmt.Println(strings.Replace(result, "data:image/jpg;base64,", "", -1))
	//middleware.Logger.Infof("qrcodeUrl: %s", *qrcodeUrl)

	time.Sleep(10 * time.Minute)
	if err != nil {
		panic(err)
	}
}
