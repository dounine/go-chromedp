package cmd

import (
	"fmt"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go-chromedp/app/middleware"
	"go-chromedp/app/service"
	"go-chromedp/app/views"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func quit() {
	fmt.Println("程序退出了")
}

func quitListen() {
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGHUP)
	go func() {
		for s := range sig {
			switch s {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGHUP:
				quit()
				if i, ok := s.(syscall.Signal); ok {
					os.Exit(int(i))
				} else {
					os.Exit(0)
				}
			}
		}
	}()
}

var (
	host string
	port string
	udid string
)
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "运行程序",
	Long:  `运行程序`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		imagePath, err := filepath.Abs("./images")
		if err != nil {
			fmt.Println(err)
			return
		}
		if !fileutil.IsExist(imagePath) {
			fmt.Println("images目录不存在")
			return
		}
		middleware.Logger.Infof("imagePath: %s", imagePath)
		middleware.Logger.Infof("udid: %s", udid)
		err = service.NewDumpEngineService().Run(&ctx, udid, imagePath)
		if err != nil {
			fmt.Println(err)
			return
		}
		quitListen()
		//store.InitSchemas() //创建数据库表
		r := gin.Default()
		r.SetTrustedProxies([]string{host})
		r.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Token"},
			AllowCredentials: true,
		}))
		r.Use(middleware.IPLimiter())
		r.Use(middleware.LoggerMiddleware())
		views.NewFileRouter().Init(r)
		views.NewDumpRouter().Init(r)

		r.NoRoute(func(c *gin.Context) {
			c.JSON(200, gin.H{
				"msg": "404",
			})
		})
		r.Run(host + ":" + port)
	},
}

func init() {
	runCmd.PersistentFlags().StringVarP(&host, "host", "H", "127.0.0.1", "绑定的服务器地址")
	runCmd.PersistentFlags().StringVarP(&port, "port", "P", "8080", "绑定的端口")
	runCmd.PersistentFlags().StringVarP(&udid, "udid", "U", "00008030-001D24223E20802E", "绑定的手机udid")
	rootCmd.AddCommand(runCmd)
}
