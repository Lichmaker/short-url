package serve

import (
	"shorturl/bootstrap"
	"shorturl/pkg/app"
	"shorturl/pkg/config"
	"shorturl/pkg/console"
	"shorturl/pkg/logger"
	"strings"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/sevlyar/go-daemon"
	"github.com/spf13/cobra"
)

// CmdServe represents the available web sub-command.
var CmdServe = &cobra.Command{
	Use:   "serve",
	Short: "Start web server",
	Run:   runWeb,
	Args:  cobra.NoArgs,
}

func init() {
	// 注册子命令
	CmdServe.AddCommand(
		CmdRestart,
	)
}

// 是否以daemon形式启动,默认为false
var RunServerDaemon bool

func runWeb(cmd *cobra.Command, args []string) {

	// daemon
	cntxt := &daemon.Context{
		PidFileName: "shorturl.pid",
		PidFilePerm: 0644,
		LogFileName: "shorturl.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{app.AbsolutePath + "/" + config.GetString("app.process_name"), "serve"},
	}

	if RunServerDaemon {
		console.Warning("程序将以daemon形式启动...")

		d, err := cntxt.Reborn()
		logger.LogIf(err)
		if d != nil {
			return
		}
		defer cntxt.Release()
	}

	// 设置 gin 的运行模式，支持 debug, release, test
	// release 会屏蔽调试信息，官方建议生产环境中使用
	// 非 release 模式 gin 终端打印太多信息，干扰到我们程序中的 Log
	// 故此设置为 release，有特殊情况手动改为 debug 即可
	gin.SetMode(gin.ReleaseMode)

	// gin 实例
	router := gin.New()

	// 初始化路由绑定
	bootstrap.SetupRoute(router)

	// 运行服务器
	// err := router.Run(":" + config.Get("app.port"))
	// console.Success("server starting...")
	logger.WarnString("CMD", "serve", "starting....")
	newS := endless.NewServer(":"+config.Get("app.port"), router)

	err := newS.ListenAndServe()
	if err != nil {
		if strings.Contains(err.Error(), "use of closed network connection") {
			logger.WarnString("CMD", "serve", "server closed . "+err.Error())
		} else {
			logger.ErrorString("CMD", "serve", err.Error())
			console.Exit("Unable to start server, error:" + err.Error())
		}
	}

}
