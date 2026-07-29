package admin

import (
	"admin/internal/ecode"
	logic "admin/internal/logic/admin"
	types "admin/internal/types/admin"
	"admin/pkg/gin/validator"

	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
)

var _ DashboardHandler = (*dashboardHandler)(nil)

// DashboardHandler defining the handler interface
type DashboardHandler interface {
	Statistics(c *gin.Context)
	Echarts(c *gin.Context)
}

type dashboardHandler struct {
	logic logic.DashboardLogic
}

// NewDashboardHandler creating the handler interface
func NewDashboardHandler() DashboardHandler {
	return &dashboardHandler{
		logic: logic.NewDashboardLogic(),
	}
}

// Statistics of data statistics
// @Summary data statistics
// @Description data statistics
// @Tags admin/dashboard
// @accept json
// @Produce json
// @Success 200 {object} types.DashboardStatisticsReply{}
// @Router /admin/v1/dashboard/statistics [get]
// @Security BearerAuth
func (d *dashboardHandler) Statistics(c *gin.Context) {
	ctx := middleware.WrapCtx(c)
	result := d.logic.Statistics(ctx)
	response.Success(c, result)
}

// Echarts of data echarts
// @Summary data echarts
// @Description data echarts
// @Tags admin/dashboard
// @accept json
// @Produce json
// @Param request query types.DashboardEchartsRequest true "query parameters"
// @Success 200 {object} types.DashboardEchartsReply{}
// @Router /admin/v1/dashboard/echarts [get]
// @Security BearerAuth
func (d *dashboardHandler) Echarts(c *gin.Context) {
	request := &types.DashboardEchartsRequest{}
	err := c.ShouldBindQuery(request)
	if err != nil {
		response.Error(c, ecode.InvalidParams.RewriteMsg(validator.GetValidatorErrorMsg(err)))
		return
	}
	ctx := middleware.WrapCtx(c)
	result := d.logic.Echarts(ctx, request)
	response.Success(c, result)
}
