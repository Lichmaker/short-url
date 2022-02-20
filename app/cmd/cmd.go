// Package cmd 存放程序的所有子命令
package cmd

import (
	"os"
	"shorturl/app/cmd/consumer"
	"shorturl/app/cmd/serve"
	"shorturl/pkg/helpers"

	"github.com/spf13/cobra"
)

// Env 存储全局选项 --env 的值
var Env string

// RegisterGlobalFlags 注册全局选项（flag）
func RegisterGlobalFlags(rootCmd *cobra.Command) {
	rootCmd.PersistentFlags().StringVarP(&Env, "env", "e", "", "load .env file, example: --env=testing will use .env.testing file")
}

// 注册其他命令的flag
func RegisterCommandFlags() {
	serve.CmdServe.Flags().BoolVarP(&serve.RunServerDaemon, "daemon", "d", false, "run server as a daemon")
	consumer.CmdConsumer.Flags().StringVarP(&consumer.RunConsumerCount, "count", "c", "", "传入需要启动的消费者数量")
	consumer.CmdConsumer.Flags().BoolVarP(&consumer.RunInstance, "instance", "i", false, "直接启动一个实例")
}

// RegisterDefaultCmd 注册默认命令
func RegisterDefaultCmd(rootCmd *cobra.Command, subCmd *cobra.Command) {
	cmd, _, err := rootCmd.Find(os.Args[1:])
	firtArg := helpers.FirstElement(os.Args[1:])
	if err == nil && cmd.Use == rootCmd.Use && firtArg != "-h" && firtArg != "--help" {
		args := append([]string{subCmd.Use}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}
}
