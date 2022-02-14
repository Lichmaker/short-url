package serve

import (
	"fmt"
	cmdhelper "shorturl/pkg/cmd-helper"
	"shorturl/pkg/config"
	"shorturl/pkg/logger"

	"github.com/spf13/cobra"
)

var CmdRestart = &cobra.Command{
	Use:   "restart",
	Short: "restart server",
	Run:   runRestart,
}

func runRestart(cmd *cobra.Command, args []string) {
	processName := config.GetString("app.process_name")
	kill := fmt.Sprintf("pkill -SIGHUP %s", processName)
	_, err := cmdhelper.RunCommand(kill)
	logger.LogIf(err)

	// todo 直接使用 kill，会状态码1退出 exit status 1。 原因未知，所以暂时使用pkill
	// pid, err := cmdhelper.GetPid(processName)
	// logger.LogIf(err)
	// kill := fmt.Sprintf("kill -SIGHUP %s", pid)
	// _, err = cmdhelper.RunCommand(kill)
	// logger.LogIf(err)
}
