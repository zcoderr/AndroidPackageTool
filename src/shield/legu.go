package shield

import (
	"channel"
	"common"
	"context"
	"crypto/md5"
	_ "crypto/md5"
	"encoding/hex"
	_ "encoding/hex"
	"encoding/json"
	"fmt"
	tencent_common "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ms/v20180408"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

type ShieldData struct {
	AppInfo     AppInfo     `json:"AppInfo"`
	ServiceInfo ServiceInfo `json:"ServiceInfo"`
}

type AppInfo struct {
	AppUrl string `json:"AppUrl"`
	AppMd5 string `json:"AppMd5"`
}

type ServiceInfo struct {
	ServiceEdition string `json:"ServiceEdition"`
	CallbackUrl    string `json:"CallbackUrl"`
	SubmitSource   string `json:"SubmitSource"`
}

func LeGuShield(sourceApkPath string, targetApkPath string) {
	println("开始应用宝加固...")
	apkUrl := uploadToTencentCos(sourceApkPath)
	itemId := createShield(apkUrl, MD5File(sourceApkPath))
	apkDownloadUrl := checkShieldStatus(itemId)
	tempPath := common.GetProjectPath() + "/output/temp/yingyongbaotemp.apk"
	downloadApk(apkDownloadUrl, tempPath)
	channel.Resign(tempPath, targetApkPath)
} /**/

func createShield(apkUrl string, apkMd5 string) string {
	println("创建加固资源...")
	credential := tencent_common.NewCredential(
		common.Conf.ConfigLeGu.CosSecretId,
		common.Conf.ConfigLeGu.CosSecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "ms.tencentcloudapi.com"
	client, _ := ms.NewClient(credential, "ap-beijing", cpf)

	request := ms.NewCreateShieldInstanceRequest()

	params := ShieldData{
		AppInfo: AppInfo{
			AppUrl: apkUrl,
			AppMd5: apkMd5,
		},
		ServiceInfo: ServiceInfo{
			ServiceEdition: "basic",
			CallbackUrl:    "",
			SubmitSource:   "MC",
		},
	}
	jsonStr, _ := json.Marshal(params)

	err := request.FromJsonString(string(jsonStr))
	if err != nil {
		panic(err)
	}
	response, err := client.CreateShieldInstance(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return ""
	}
	if err != nil {
		panic(err)
	}

	fmt.Println("加固资源创建成功，加固任务 ID:" + *response.Response.ItemId)
	return *response.Response.ItemId
}

func checkShieldStatus(itemId string) string {
	println("正在加固...")
	credential := tencent_common.NewCredential(
		common.Conf.ConfigLeGu.CosSecretId,
		common.Conf.ConfigLeGu.CosSecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "ms.tencentcloudapi.com"
	client, _ := ms.NewClient(credential, "ap-beijing", cpf)

	request := ms.NewDescribeShieldResultRequest()

	params := "{\"ItemId\":\"" + itemId + "\"}"
	err := request.FromJsonString(params)
	if err != nil {
		panic(err)
	}

	for {
		response, err := client.DescribeShieldResult(request)
		if _, ok := err.(*errors.TencentCloudSDKError); ok {
			fmt.Printf("An API error has returned: %s", err)
			return ""
		}
		if err != nil {
			panic(err)
		}

		if *response.Response.TaskStatus != 2 {
			if *response.Response.TaskStatus == 1 {
				println("加固完成")
				return *response.Response.ShieldInfo.AppUrl
			} else if *response.Response.TaskStatus == 3 {
				println("加固出错，任务中止。")
				os.Exit(0)
			} else if *response.Response.TaskStatus == 4 {
				println("加固超时，任务中止。")
				os.Exit(0)
			}
		}
		println("应用宝加固中,请等待...")
		time.Sleep(time.Duration(10) * time.Second)
	}
}

func downloadApk(url string, targetPath string) {
	println("Apk 下载中...")
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	f, err := os.Create(targetPath)
	if err != nil {
		panic(err)
	}
	io.Copy(f, res.Body)
	println("下载完成")
}

func uploadToTencentCos(filePath string) string {
	u, _ := url.Parse(common.Conf.ConfigLeGu.CosBucketUrl)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  common.Conf.ConfigLeGu.CosSecretId,
			SecretKey: common.Conf.ConfigLeGu.CosSecretKey,
		},
	})

	println("上传 Apk 至腾讯云...")

	name := "shield/yingyongbao.apk"
	_, err := c.Object.PutFromFile(context.Background(), name, filePath, nil)
	if err != nil {
		panic(err)
	}
	println("上传完成")
	return common.Conf.ConfigLeGu.CosBucketUrl + "/" + name
}

func MD5File(file string) string {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return ""
	}
	ret := md5.Sum(data)
	return hex.EncodeToString(ret[:])
}
