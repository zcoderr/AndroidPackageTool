package common

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var Conf Configuration

var err error
var cmd *exec.Cmd

func GetProjectPath() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func ReadConfiguration() {
	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&Conf)
	if err != nil {
		fmt.Println("配置文件解析错误:", err)
	} else {
		fmt.Println("配置文件解析完成")
	}

	// 把配置文件中的相对路径重新赋值为绝对路径
	Conf.ConfigKeyStore.Path = GetProjectPath() + Conf.ConfigKeyStore.Path
	for i, item := range Conf.ConfigShield.Legu {
		item.ApkPath = GetProjectPath() + item.ApkPath
		item.TargetPath = GetProjectPath() + item.TargetPath
		Conf.ConfigShield.Legu[i] = item
	}

	for i, item := range Conf.ConfigShield.Jiagubao {
		item.ApkPath = GetProjectPath() + item.ApkPath
		item.TargetPath = GetProjectPath() + item.TargetPath
		Conf.ConfigShield.Jiagubao[i] = item
	}

	for i, item := range Conf.ConfigChannel {
		item.ApkPath = GetProjectPath() + item.ApkPath
		Conf.ConfigChannel[i] = item
	}

	Conf.ConfigChannels.ChannelFilePath = GetProjectPath() + Conf.ConfigChannels.ChannelFilePath
	Conf.ConfigChannels.ApkPath = GetProjectPath() + Conf.ConfigChannels.ApkPath
}

func ExecCommand(commandStr string) bool {
	cmd = exec.Command("/bin/bash", "-c", commandStr)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("Cmd：" + commandStr)
	cmd.Start()

	reader := bufio.NewReader(stdout)
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Println(line)
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func GetFileList(path string) (files []string) {
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		println(path)
		files = append(files, path)
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
	return files
}

func AddGitTag() {
	f, err := os.Open(Conf.ConfigBuild.ProjectPath + "/app/build.gradle")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')

		if err != nil || io.EOF == err {
			break
		}
		if strings.Contains(line, "versionName \"") {
			ver := strings.Split(line, "\"")[1]
			var bt bytes.Buffer
			bt.WriteString("cd " + Conf.ConfigBuild.ProjectPath)
			bt.WriteString(";")
			bt.WriteString("git tag " + ver)
			if (ExecCommand(bt.String())) {
				fmt.Println("已添加 git tag " + ver)
			} else {
				fmt.Println("tag" + ver + " 已存在!!!!!!!!")
			}
			break
		}
	}
}
