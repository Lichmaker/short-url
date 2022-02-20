package consumer

import (
	"fmt"
	cmdhelper "shorturl/pkg/cmd-helper"
	"shorturl/pkg/config"
	"shorturl/pkg/console"
	"shorturl/pkg/logger"
	"strings"
	"syscall"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var CmdShutdown = &cobra.Command{
	Use:   "shutdown",
	Short: "shutdown all cunsumer",
	Run:   runShutdown,
}

func runShutdown(cmd *cobra.Command, args []string) {
	processName := fmt.Sprintf("%s consumer", config.GetString("app.process_name"))

	// var i int = 0

	// 使用syscall.kill
	psCmd := `ps aux | awk '/` + processName + `/ && !/awk/ && !/shutdown/ {print $1}'`
	pid, err := cmdhelper.RunCommand(psCmd)
	fmt.Printf("拿到的pid %s \n", pid)
	logger.LogIf(err)

	pidGroup := strings.Split(pid, "\n")

	for _, eachPid := range pidGroup {
		err = syscall.Kill(cast.ToInt(eachPid), syscall.SIGQUIT)
		logger.LogIf(err)
	}

	console.Success(fmt.Sprintf("总共关闭%d个consumer", len(pidGroup)))
}
