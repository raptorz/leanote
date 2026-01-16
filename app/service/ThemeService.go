package service

import (
	"github.com/leanote/leanote/app/info"
)

type ThemeService struct {
}

// 简化版本，让编译通过
func (this *ThemeService) GetDefaultTheme(style string) info.Theme {
	return info.Theme{}
}

func (this *ThemeService) GetTheme(themeId string) info.Theme {
	return info.Theme{}
}

func (this *ThemeService) GetThemes() []info.Theme {
	return []info.Theme{}
}

func (this *ThemeService) AddTheme(theme info.Theme) bool {
	return false
}

func (this *ThemeService) UpdateTheme(themeId string, theme info.Theme) bool {
	return false
}

func (this *ThemeService) DeleteTheme(themeId string) bool {
	return false
}

func (this *ThemeService) SetDefaultTheme(themeId string) bool {
	return false
}
