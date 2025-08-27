package handler

import (
	"admin/internal/ecode"
	"admin/internal/types"

	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/logger"
)

type DashboardHandler interface {
	Statistics(c *gin.Context)
	Echarts(c *gin.Context)
}

type dashboardHandler struct {
}

func NewDashboardHandler() DashboardHandler {
	return &dashboardHandler{}
}

// Statistics of data statistics
// @Summary data statistics
// @Description data statistics
// @Tags dashboard
// @accept json
// @Produce json
// @Success 200 {object} types.DashboardStatisticsReply{}
// @Router /api/v1/dashboard/statistics [get]
// @Security BearerAuth
func (d *dashboardHandler) Statistics(c *gin.Context) {
	result := []types.DashboardStatisticsItem{
		{
			Type:             "user",
			Title:            "在线用户",
			TodayCount:       1,
			TotalCount:       1,
			GrowthRate:       0,
			GranularityLabel: "日",
		},
		{
			Type:             "pv",
			Title:            "浏览量",
			TodayCount:       498,
			TotalCount:       497359,
			GrowthRate:       -0.76,
			GranularityLabel: "日",
		},
		{
			Type:             "uv",
			Title:            "访客数",
			TodayCount:       100,
			TotalCount:       2000,
			GrowthRate:       0,
			GranularityLabel: "日",
		},
		{
			Type:             "ip",
			Title:            "IP数",
			TodayCount:       96,
			TotalCount:       31048,
			GrowthRate:       -0.61,
			GranularityLabel: "日",
		},
	}
	response.Success(c, result)
}

// Echarts of data echarts
// @Summary data echarts
// @Description data echarts
// @Tags dashboard
// @accept json
// @Produce json
// @Param request query types.DashboardEchartsRequest true "query parameters"
// @Success 200 {object} types.DashboardEchartsReply{}
// @Router /api/v1/dashboard/echarts [get]
// @Security BearerAuth
func (d *dashboardHandler) Echarts(c *gin.Context) {
	request := &types.DashboardEchartsRequest{}
	err := c.ShouldBindQuery(request)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	result := types.Echarts{
		Names: []string{"浏览量(PV)", "IP"},
		Dates: []string{"2024-10-26", "2024-10-27", "2024-10-28", "2024-10-29", "2024-10-30", "2024-10-31", "2024-11-01", "2024-11-02"},
		Series: []types.EchartsSeries{
			{
				Name:      "浏览量(PV)",
				Data:      []int{1815, 1482, 3918, 3786, 3554, 4218, 4337, 498},
				AreaStyle: "rgba(64, 158, 255, 0.1)",
				LineStyle: "#4080FF",
				ItemStyle: "#4080FF",
			},
			{
				Name:      "IP",
				Data:      []int{211, 168, 475, 489, 461, 490, 457, 96},
				AreaStyle: "rgba(103, 194, 58, 0.1)",
				LineStyle: "#67C23A",
				ItemStyle: "#67C23A",
			},
		},
	}
	response.Success(c, result)
}
