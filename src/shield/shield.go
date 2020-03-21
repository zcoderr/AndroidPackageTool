package shield

import (
	"common"
)

func DoShield(){
	DoLeguSheild()
	DoJiagubaoShield()
}

func DoLeguSheild() {
	for _, shield := range common.Conf.ConfigShield.Legu {
		LeGuShield(shield.ApkPath, shield.TargetPath)
	}
}

func DoJiagubaoShield() {
	for _, shield := range common.Conf.ConfigShield.Jiagubao {
		JiaGuBaoSheild(shield.ApkPath, shield.TargetPath)
	}
}
