package common

type Configuration struct {
	Option Option `json:"option"`
	ConfigKeyStore ConfigKeyStore   `json:"keystore"`
	ConfigBuild    ConfBuild        `json:"build"`
	ConfigLeGu     ConfigLeGu       `json:"yingyongbao"`
	Config360      Config360        `json:"jiagubao"`
	ConfigShield   ConfigShield     `json:"shield"`
	ConfigChannel  [] ConfigChannel `json:"channel"`
	ConfigChannels ConfigChannels   `json:"channels"`
}

type Option struct {
	Build bool `json:"build"`
	Shield bool `json:"shield"`
	MulitChannel bool `json:"mulit_channel"`
} 

type ConfigKeyStore struct {
	Path         string `json:"path"`
	KeyAlias     string `json:"key_alias"`
	KeyStorePass string `json:"key_store_pass"`
	KeyPass      string `json:"key_pass"`
}

type ConfBuild struct {
	ProjectPath   string          `json:"project_path"`
	BuildVariants []BuildVariants `json:"build_variants"`
}

type BuildVariants struct {
	Task       string `json:"task"`
	OutputPath string `json:"output_path"`
	TargetPath string `json:"target_path"`
}

type Config360 struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ConfigLeGu struct {
	CosSecretId  string `json:"cos_secret_id"`
	CosSecretKey string `json:"cos_secret_key"`
	CosBucketUrl string `json:"cos_bucket_url"`
}

type StudioOutputJson struct {
	Path    string  `json:path`
	ApkInfo ApkInfo `json:apkInfo`
}

type ConfigShield struct {
	Legu     [] Shield `json:"legu"`
	Jiagubao [] Shield `json:"jiagubao"`
}

type Shield struct {
	ApkPath    string `json:"apk_path"`
	TargetPath string `json:"target_path"`
}

type ConfigChannel struct {
	ApkPath     string `json:"apk_path"`
	ChannelName string `json:"channel_name"`
}

type ConfigChannels struct {
	ChannelFilePath string `json:"channel_file_path"`
	ApkPath         string `json:"apk_path"`
}
