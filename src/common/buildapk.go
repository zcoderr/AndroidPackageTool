package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type ApkInfo struct {
	FullName string `json:fullName`
}

func ReadStudioOutputJson(filePath string) (apkName string) {
	var dataArr []StudioOutputJson

	byteArray, _ := ioutil.ReadFile(filePath)
	err := json.Unmarshal(byteArray, &dataArr)
	if err != nil {
		fmt.Println("配置文件解析错误:", err)
	}

	return dataArr[0].Path
}

/**
构建多个变种的 apk
*/
func BuildApk() {
	var bt bytes.Buffer
	bt.WriteString("cd ")
	bt.WriteString(Conf.ConfigBuild.ProjectPath)
	bt.WriteString(";")
	bt.WriteString("./gradlew ")
	for _, variant := range Conf.ConfigBuild.BuildVariants {
		bt.WriteString(" " + variant.Task)
	}
	result := ExecCommand(bt.String())
	if result {
		fmt.Println("编译完成...")
	}
}

/*
把 AndroidStudio 构建输出的 Apk 移动到目标位置
*/
func MoveApk() {
	for _, variant := range Conf.ConfigBuild.BuildVariants {
		apkFileName := ReadStudioOutputJson(Conf.ConfigBuild.ProjectPath + variant.OutputPath + "/output.json")
		os.Rename(Conf.ConfigBuild.ProjectPath+variant.OutputPath+"/"+apkFileName, GetProjectPath()+variant.TargetPath)
	}
	fmt.Println("移动APK到指定目录...")
}
