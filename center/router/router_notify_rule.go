package router

import (
	"net/http"
	"time"

	"github.com/ccfos/nightingale/v6/models"
	"github.com/ccfos/nightingale/v6/pkg/slice"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

func (rt *Router) notifyRulesAdd(c *gin.Context) {
	var lst []*models.NotifyRule
	ginx.BindJSON(c, &lst)
	if len(lst) == 0 {
		ginx.Bomb(http.StatusBadRequest, "input json is empty")
	}

	me := c.MustGet("user").(*models.User)
	gids, err := models.MyGroupIds(rt.Ctx, me.Id)
	ginx.Dangerous(err)
	for _, t := range lst {
		if !slice.HaveIntersection(gids, t.UserGroupIds) {
			ginx.Bomb(http.StatusForbidden, "no permission")
		}
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

func (rt *Router) notifyRulesDel(c *gin.Context) {
	var f idsForm
	ginx.BindJSON(c, &f)
	f.Verify()

	lst, err := models.NotifyRulesGet(rt.Ctx, "id in (?)", f.Ids)
	ginx.Dangerous(err)
	me := c.MustGet("user").(*models.User)
	gids, err := models.MyGroupIds(rt.Ctx, me.Id)
	ginx.Dangerous(err)
	for _, t := range lst {
		if !slice.HaveIntersection[int64](gids, t.UserGroupIds) {
			ginx.Bomb(http.StatusForbidden, "no permission")
		}
	}

	ginx.NewRender(c).Message(models.DB(rt.Ctx).
		Delete(&models.NotifyRule{}, "id in (?)", f.Ids).Error)
}

func (rt *Router) notifyRulePut(c *gin.Context) {
	var f models.NotifyRule
	ginx.BindJSON(c, &f)

	nr, err := models.NotifyRuleGet(rt.Ctx, "id = ?", ginx.UrlParamInt64(c, "id"))
	ginx.Dangerous(err)
	if nr == nil {
		ginx.Bomb(http.StatusNotFound, "notify rule not found")
	}

	me := c.MustGet("user").(*models.User)
	gids, err := models.MyGroupIds(rt.Ctx, me.Id)
	ginx.Dangerous(err)
	if !slice.HaveIntersection[int64](gids, nr.UserGroupIds) {
		ginx.Bomb(http.StatusForbidden, "no permission")
	}

	f.UpdateBy = me.Username
	ginx.NewRender(c).Message(nr.Update(rt.Ctx, f))
}

func (rt *Router) notifyRuleGet(c *gin.Context) {
	me := c.MustGet("user").(*models.User)
	gids, err := models.MyGroupIds(rt.Ctx, me.Id)
	ginx.Dangerous(err)

	tid := ginx.UrlParamInt64(c, "id")
	nr, err := models.NotifyRuleGet(rt.Ctx, "id = ?", tid)
	ginx.Dangerous(err)
	if nr == nil {
		ginx.Bomb(http.StatusNotFound, "notify rule not found")
	}
	if !slice.HaveIntersection[int64](gids, nr.UserGroupIds) {
		ginx.Bomb(http.StatusForbidden, "no permission")
	}

	ginx.NewRender(c).Data(nr, nil)
}

func (rt *Router) notifyRulesGet(c *gin.Context) {
	me := c.MustGet("user").(*models.User)
	gids, err := models.MyGroupIds(rt.Ctx, me.Id)
	ginx.Dangerous(err)

	lst, err := models.NotifyRulesGet(rt.Ctx, "", nil)
	ginx.Dangerous(err)

	res := make([]*models.NotifyRule, 0)
	for _, nr := range lst {
		if slice.HaveIntersection[int64](gids, nr.UserGroupIds) {
			res = append(res, nr)
		}
	}
	ginx.NewRender(c).Data(res, nil)
}
