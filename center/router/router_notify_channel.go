package router

import (
	"net/http"
	"time"

	"github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

func (rt *Router) notifyChannelsAdd(c *gin.Context) {
	me := c.MustGet("user").(*models.User)
	if !me.IsAdmin() {
		ginx.Bomb(http.StatusForbidden, "no permission")
	}

	var lst []*models.NotifyChannelConfig
	ginx.BindJSON(c, &lst)
	if len(lst) == 0 {
		ginx.Bomb(http.StatusBadRequest, "input json is empty")
	}

	for _, tpl := range lst {
		tpl.CreateBy = me.Username
		tpl.CreateAt = time.Now().Unix()
	}

	ginx.Dangerous(models.DB(rt.Ctx).CreateInBatches(lst, 100).Error)
	ids := make([]uint, 0, len(lst))
	for _, tpl := range lst {
		ids = append(ids, tpl.ID)
	}
	ginx.NewRender(c).Data(ids, nil)
}

func (rt *Router) notifyChannelsDel(c *gin.Context) {
	me := c.MustGet("user").(*models.User)
	if !me.IsAdmin() {
		ginx.Bomb(http.StatusForbidden, "no permission")
	}

	var f idsForm
	ginx.BindJSON(c, &f)
	f.Verify()

	ginx.NewRender(c).Message(models.DB(rt.Ctx).
		Delete(&models.NotifyChannelConfig{}, "id in (?)", f.Ids).Error)
}

func (rt *Router) notifyChannelPut(c *gin.Context) {
	me := c.MustGet("user").(*models.User)
	if !me.IsAdmin() {
		ginx.Bomb(http.StatusForbidden, "no permission")
	}

	var f models.NotifyChannelConfig
	ginx.BindJSON(c, &f)

	nc, err := models.NotifyChannelGet(rt.Ctx, "id = ?", ginx.UrlParamInt64(c, "id"))
	ginx.Dangerous(err)
	if nc == nil {
		ginx.Bomb(http.StatusNotFound, "notify channel not found")
	}

	f.UpdateBy = me.Username
	ginx.NewRender(c).Message(nc.Update(rt.Ctx, f))
}

func (rt *Router) notifyChannelGet(c *gin.Context) {
	tid := ginx.UrlParamInt64(c, "id")
	nc, err := models.NotifyChannelGet(rt.Ctx, "id = ?", tid)
	ginx.Dangerous(err)
	if nc == nil {
		ginx.Bomb(http.StatusNotFound, "notify channel not found")
	}

	ginx.NewRender(c).Data(nc, nil)
}

func (rt *Router) notifyChannelsGet(c *gin.Context) {
	lst, err := models.NotifyChannelsGet(rt.Ctx, "", nil)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(lst, nil)
}
