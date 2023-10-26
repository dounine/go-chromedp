package main

import (
	"context"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/duke-git/lancet/v2/fileutil"
)

func main() {
	//cmd.Execute()
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

	err := chromedp.Run(
		ctx,
		chromedp.Navigate("https://baidu.com"),
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
		//截图
		chromedp.ActionFunc(func(ctx context.Context) error {
			var buf []byte
			if err := chromedp.CaptureScreenshot(&buf).Do(ctx); err != nil {
				return err
			}
			err := fileutil.WriteBytesToFile("baidu.png", buf)
			return err
		}),
	)
	if err != nil {
		panic(err)
	}
}
