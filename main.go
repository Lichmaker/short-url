package main

import (
	"fmt"
	"shorturl/app/cmd"
	"shorturl/app/cmd/make"
	"shorturl/app/cmd/serve"

	"os"
	"shorturl/bootstrap"
	btsConfig "shorturl/config"
	"shorturl/pkg/app"
	"shorturl/pkg/config"
	"shorturl/pkg/console"
	"shorturl/pkg/traceid"

	"github.com/spf13/cobra"
)

func init() {
	// 初始化所有的config
	btsConfig.Initialize()
}

func main() {
	app.AbsolutePath, _ = os.Getwd()

	// 主入口，调用 cmd.CmdServe 命令
	var rootCmd = &cobra.Command{
		Use:   config.Get("app.name"),
		Short: "my project describe",
		Long:  `默认会使用 "serve" 命令，输入 -h 查看帮助`,

		// 定义所有自命令都会执行的func
		PersistentPreRun: func(command *cobra.Command, args []string) {
			traceid.Boot("")

			config.InitConfig(cmd.Env)

			bootstrap.SetupLogger()
			bootstrap.SetupDB()
			bootstrap.SetupRedis()
			bootstrap.SetupCache()
		},
	}

	// 注册子命令
	rootCmd.AddCommand(
		serve.CmdServe,
		cmd.CmdKey,
		cmd.CmdPlay,
		cmd.CmdMigrate,
		make.CmdMake,
		cmd.CmdTidyCache,
		// cmd.CmdRestart,
	)

	// 配置默认使用的子命令
	// cmd.RegisterDefaultCmd(rootCmd, cmd.CmdServe)

	// 注册指定命令，非全局的flag
	cmd.RegisterCommandFlags()

	// 注册全局参数
	cmd.RegisterGlobalFlags(rootCmd)

	// 执行主命令
	if err := rootCmd.Execute(); err != nil {
		console.Exit(fmt.Sprintf("启动失败 %v: %s", os.Args, err.Error()))
	}

	// var env string
	// flag.StringVar(&env, "env", "", "加载 .env 文件， 如 --env=testing 则加载 .env.testing 文件")
	// flag.Parse()
	// config.InitConfig(env)
	// bootstrap.SetupLogger()

	// r := gin.New()
	// bootstrap.SetupRoute(r)
	// bootstrap.SetupDB()
	// bootstrap.SetupRedis()

	// // 设置 gin 的运行模式，支持 debug, release, test
	// // release 会屏蔽调试信息，官方建议生产环境中使用
	// // 非 release 模式 gin 终端打印太多信息，干扰到我们程序中的 Log
	// // 故此设置为 release，有特殊情况手动改为 debug 即可
	// gin.SetMode(gin.ReleaseMode)

	// err := r.Run(":" + config.GetString("app.port"))
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
}

// 获取当前执行文件绝对路径（go run）
// func getCurrentAbPathByCaller() string {
// 	var abPath string
// 	_, filename, _, ok := runtime.Caller(0)
// 	if ok {
// 		abPath = path.Dir(filename)
// 	}
// 	return abPath
// }
