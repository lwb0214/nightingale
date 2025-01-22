package router

import (
	"net/http"
	"time"

	"github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

func (rt *Router) notifyChannelsAdd(c *gin.Context) {
	var lst []*models.NotifyChannelConfig
	ginx.BindJSON(c, &lst)
	if len(lst) == 0 {
		ginx.Bomb(http.StatusBadRequest, "input json is empty")
	}

	// me := c.MustGet("user").(*models.User)
	// gids, err := models.MyGroupIds(rt.Ctx, me.Id)
	// ginx.Dangerous(err)
	// for _, t := range lst {
	// 	if !slice.HaveIntersection(gids, t.UserGroupIds) {
	// 		ginx.Bomb(http.StatusForbidden, "no permission")
	// 	}
	// }

	username := Username(c)
	for _, tpl := range lst {
		tpl.CreateBy = username
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
	var f idsForm
	ginx.BindJSON(c, &f)
	f.Verify()

	// lst, err := models.NotifyChannelsGet(rt.Ctx, "id in (?)", f.Ids)
	// ginx.Dangerous(err)
	// // me := c.MustGet("user").(*models.User)
	// // gids, err := models.MyGroupIds(rt.Ctx, me.Id)
	// // ginx.Dangerous(err)
	// // for _, t := range lst {
	// // 	if !slice.HaveIntersection[int64](gids, t.UserGroupIds) {
	// // 		ginx.Bomb(http.StatusForbidden, "no permission")
	// // 	}
	// // }

	ginx.NewRender(c).Message(models.DB(rt.Ctx).Delete(&models.NotifyChannelConfig{},
		"id in (?)", f.Ids).Error)
}

func (rt *Router) notifyChannelPut(c *gin.Context) {
	var f models.NotifyChannelConfig
	ginx.BindJSON(c, &f)

	nc, err := models.NotifyChannelGet(rt.Ctx, "id = ?", ginx.UrlParamInt64(c, "id"))
	ginx.Dangerous(err)
	if nc == nil {
		ginx.Bomb(http.StatusNotFound, "noyify rule not found")
	}

	// me := c.MustGet("user").(*models.User)
	// gids, err := models.MyGroupIds(rt.Ctx, me.Id)
	// ginx.Dangerous(err)
	// if !slice.HaveIntersection[int64](gids, nr.UserGroupIds) {
	// 	ginx.Bomb(http.StatusForbidden, "no permission")
	// }

	f.UpdateBy = Username(c)
	ginx.NewRender(c).Message(nc.Update(rt.Ctx, f))
}

func (rt *Router) notifyChannelGet(c *gin.Context) {
	// me := c.MustGet("user").(*models.User)
	// gids, err := models.MyGroupIds(rt.Ctx, me.Id)
	// ginx.Dangerous(err)

	tid := ginx.UrlParamInt64(c, "id")
	nc, err := models.NotifyChannelGet(rt.Ctx, "id = ?", tid)
	ginx.Dangerous(err)
	if nc == nil {
		ginx.Bomb(http.StatusNotFound, "noyify channel not found")
	}
	// if !slice.HaveIntersection[int64](gids, nr.UserGroupIds) {
	// 	ginx.Bomb(http.StatusForbidden, "no permission")
	// }

	ginx.NewRender(c).Data(nc, nil)
}

func (rt *Router) notifyChannelsGet(c *gin.Context) {
	// me := c.MustGet("user").(*models.User)
	// gids, err := models.MyGroupIds(rt.Ctx, me.Id)
	// ginx.Dangerous(err)

	lst, err := models.NotifyChannelsGet(rt.Ctx, "", nil)
	ginx.Dangerous(err)

	// res := make([]*models.NotifyRule, 0)
	// for _, nr := range lst {
	// 	if slice.HaveIntersection[int64](gids, nr.UserGroupIds) {
	// 		res = append(res, nr)
	// 	}
	// }
	ginx.NewRender(c).Data(lst, nil)
}
