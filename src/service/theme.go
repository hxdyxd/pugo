package service

import (
	"errors"
	"github.com/go-xiaohei/pugo/src/core"
	"github.com/go-xiaohei/pugo/src/model"
)

var (
	Theme = newThemeService()

	ErrThemeInvalid = errors.New("theme-invalid")
)

type ThemeService struct {
	cacheThemes  map[string]*model.Theme
	themeSlice   []*model.Theme
	currentTheme *model.Theme
}

func newThemeService() *ThemeService {
	return &ThemeService{
		cacheThemes: make(map[string]*model.Theme),
	}
}

func (ts *ThemeService) Load(_ interface{}) (*Result, error) {
	themes := make([]*model.Theme, 0)
	if err := core.Db.OrderBy("status ASC").Find(&themes); err != nil {
		return nil, err
	}
	for _, t := range themes {
		ts.cacheThemes[t.Name] = t
		if t.IsCurrent() {
			ts.currentTheme = t
		}
	}
	ts.themeSlice = themes
	res := newResult(ts.Load, &themes)
	return res, nil
}

func (ts *ThemeService) Current(_ interface{}) (*Result, error) {
	if ts.currentTheme == nil {
		if _, err := ts.Load(nil); err != nil {
			return nil, err
		}
	}
	if ts.currentTheme == nil || !ts.currentTheme.IsValid() {
		return nil, ErrThemeInvalid
	}
	return newResult(ts.Current, ts.currentTheme), nil
}

func (ts *ThemeService) Admin(_ interface{}) (*Result, error) {
	if len(ts.cacheThemes) == 0 {
		if _, err := ts.Load(nil); err != nil {
			return nil, err
		}
	}
	theme := ts.cacheThemes["admin"]
	if theme == nil || !theme.IsValid() {
		return nil, ErrThemeInvalid
	}
	return newResult(ts.Admin, theme), nil
}

func (ts *ThemeService) All(_ interface{}) (*Result, error) {
	if len(ts.themeSlice) == 0 {
		if _, err := ts.Load(nil); err != nil {
			return nil, err
		}
	}
	return newResult(ts.All, &ts.themeSlice), nil
}
