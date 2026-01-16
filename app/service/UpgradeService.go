package service

type UpgradeService struct {
}

// 简化版本，用于PostgreSQL迁移
func (this *UpgradeService) UpgradeBlog() bool {
	return true
}

func (this *UpgradeService) UpgradeBetaToBeta2(userId string) (ok bool, msg string) {
	return true, "success"
}

func (this *UpgradeService) Api(userId string) (ok bool, msg string) {
	return true, ""
}
