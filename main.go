package main

import (
	"channel"
	"common"
	_ "crypto/md5"
	_ "encoding/hex"
	"shield"
)

func main() {
	common.ReadConfiguration()

	//common.AddGitTag()

	if common.Conf.Option.Build {
		common.BuildApk()
		common.MoveApk()
	}

	if common.Conf.Option.Shield {
		shield.DoShield()
	}

	if common.Conf.Option.MulitChannel {
		channel.WriteChannels()
		channel.DoMulitChannel()
	}
}
