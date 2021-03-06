package public

import (
	"fmt"
	"github.com/go-xiaohei/pugo/src/middle"
	"github.com/go-xiaohei/pugo/src/model"
	"github.com/go-xiaohei/pugo/src/service"
	"github.com/go-xiaohei/pugo/src/utils"
	"github.com/lunny/tango"
	"net/http"
	"time"
)

var (
	readTimeCookieName = "PUGO_READ_TIME"
)

type IndexController struct {
	tango.Ctx

	middle.ThemeRender
}

func (ic *IndexController) getReadTime() int64 {
	return ic.CookieInt64(readTimeCookieName)
}

func (ic *IndexController) setReadTime() {
	ic.Cookies().Set(&http.Cookie{
		Name:     readTimeCookieName,
		Value:    fmt.Sprint(time.Now().Unix()),
		Path:     "/",
		Expires:  time.Now().Add(365 * 24 * 10 * time.Hour),
		MaxAge:   3600 * 24 * 10 * 365,
		HttpOnly: true,
	})
}

func (ic *IndexController) Get() {
	ic.Title(service.Setting.General.FullTitle())
	var (
		opt = service.ArticleListOption{
			Status:   model.ARTICLE_STATUS_PUBLISH,
			Order:    "create_time DESC",
			Page:     ic.ParamInt(":page", 1),
			Size:     service.Setting.Content.PageSize,
			IsCount:  true,
			ReadTime: ic.getReadTime(),
		}
		articles = make([]*model.Article, 0)
		pager    = new(utils.Pager)
	)
	if err := service.Call(service.Article.List, opt, &articles, pager); err != nil {
		ic.RenderError(500, err)
		return
	}
	ic.setReadTime()
	ic.Assign("Articles", articles)
	ic.Assign("Pager", pager)
	ic.Render("index.tmpl")
}
