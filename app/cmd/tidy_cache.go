package cmd

import (
	"fmt"
	"shorturl/pkg/console"
	"shorturl/pkg/short"

	"github.com/spf13/cobra"
)

var CmdTidyCache = &cobra.Command{
	Use:   "tidycache",
	Short: "整理缓存，把长时间没有使用的数据清理掉",
	Run:   runTidyCache,
}

func runTidyCache(cmd *cobra.Command, args []string) {
	console.Info("开始执行缓存整理...")
	c := short.TidyCache()
	console.Success(fmt.Sprintf("缓存整理完成，共处理%d个数据", c))
}
