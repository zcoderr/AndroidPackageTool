package shield

import (
	"bytes"
	"channel"
	"common"
	"fmt"
)

func JiaGuBaoSheild(apkPath string, targetName string) {
	var bt bytes.Buffer
	fmt.Println("开始 360 加固...")
	bt.WriteString("java -jar asset/jiagu/jiagu.jar -login ")
	bt.WriteString(common.Conf.Config360.Username)
	bt.WriteString(" ")
	bt.WriteString(common.Conf.Config360.Password)

	//cmd = exec.Command("/bin/bash", "-c", bt.String())
	result := common.ExecCommand(bt.String())

	if result {
		fmt.Println("已登录")
	}

	jiagubao_exc(apkPath, targetName)
}

func jiagubao_exc(apkPath string, targetName string) {
	command := "java -jar asset/jiagu/jiagu.jar -jiagu " + apkPath + " " + common.GetProjectPath() + "/output/temp -autosign"
	result := common.ExecCommand(command)

	if result {
		fmt.Println("已加固")
		files := common.GetFileList(common.GetProjectPath() + "/output/temp")
		// 重签名，并移动到 apk 文件夹
		channel.Resign(files[0], targetName)
		//os.Rename(files[0], targetName)
	}
}
