package channel

import (
	"bytes"
	"common"
	"fmt"
	"os"
	"strings"
)

func Resign(apkPath string, targetPath string) {
	// 4k 对齐
	common.ExecCommand("./asset/zipalign -v 4 " + apkPath + " " + common.GetProjectPath() + "./output/temp/temp_sign.apk")

	// 重签名
	var bt bytes.Buffer
	bt.WriteString("java -jar " + common.GetProjectPath() + "/asset/apksigner.jar")
	bt.WriteString(" sign --ks " + common.Conf.ConfigKeyStore.Path)
	bt.WriteString(" --ks-key-alias " + common.Conf.ConfigKeyStore.KeyAlias)
	bt.WriteString(" --ks-pass pass:" + common.Conf.ConfigKeyStore.KeyStorePass)
	bt.WriteString(" --key-pass pass:" + common.Conf.ConfigKeyStore.KeyPass)
	bt.WriteString(" --out " + common.GetProjectPath() + "/output/temp/temp_sign.apk " + targetPath)
	if (common.ExecCommand(bt.String())) {
		fmt.Println(targetPath + "应用宝重签名成功")
	} else {
		fmt.Println(targetPath + "应用宝重签名失败")
	}
}

func DoMulitChannel() {
	for _, channel := range common.Conf.ConfigChannel {
		writeChannelTag(channel.ApkPath, channel.ChannelName)
	}
}

// 对加固和变种包写渠道
func writeChannelTag(apkPath string, channelName string) {
	var bt bytes.Buffer
	bt.WriteString("java -jar " + common.GetProjectPath() + "/asset/VasDolly.jar")
	bt.WriteString(" put -c")
	bt.WriteString(" " + channelName + " " + apkPath + " ./output/channel")
	common.ExecCommand(bt.String())

	files := common.GetFileList("./output/channel")
	for _, file := range files {
		if strings.Contains(file,"-") {
			print(file)
			// 重命名
			os.Rename(file, "./output/channel/"+channelName+".apk")
		}
	}
}

func WriteChannels() {
	var bt bytes.Buffer
	bt.WriteString("java -jar " + common.GetProjectPath() + "/asset/VasDolly.jar")
	bt.WriteString(" put -c")
	bt.WriteString(" " + common.Conf.ConfigChannels.ChannelFilePath + " " + common.Conf.ConfigChannels.ApkPath + " " + common.GetProjectPath() + "/output/channel")
	common.ExecCommand(bt.String())
	fmt.Println("写入渠道完成")
}
