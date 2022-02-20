package serve

import (
	"fmt"
	cmdhelper "shorturl/pkg/cmd-helper"
	"shorturl/pkg/config"
	"shorturl/pkg/console"
	"shorturl/pkg/logger"
	"syscall"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var CmdRestart = &cobra.Command{
	Use:   "restart",
	Short: "restart server",
	Run:   runRestart,
}

func runRestart(cmd *cobra.Command, args []string) {
	processName := fmt.Sprintf("%s serve", config.GetString("app.process_name"))
	// kill := fmt.Sprintf("pkill -SIGHUP %s", processName)
	// _, err := cmdhelper.RunCommand(kill)
	// logger.LogIf(err)

	// 使用syscall.kill
	psCmd := `ps aux | awk '/` + processName + `/ && !/awk/ && !/restart/ {print $1}'`
	pid, err := cmdhelper.RunCommand(psCmd)
	logger.LogIf(err)
	pidint := cast.ToInt(pid)
	if pidint == 0 {
		console.Error("找不到程序pid，重启失败")
		return
	}
	err = syscall.Kill(pidint, syscall.SIGHUP)
	logger.LogIf(err)
}
