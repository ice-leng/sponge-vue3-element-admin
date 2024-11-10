package handler

import (
	"admin/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/sponge/pkg/gin/response"
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
// @Success 200 {object} types.DashboardEchartsReply{}
// @Router /api/v1/dashboard/echarts [get]
// @Security BearerAuth
func (d *dashboardHandler) Echarts(c *gin.Context) {

	result := map[string]interface{}{
		"dates":  []string{"2024-10-26", "2024-10-27", "2024-10-28", "2024-10-29", "2024-10-30", "2024-10-31", "2024-11-01", "2024-11-02"},
		"pvList": []int{1815, 1482, 3918, 3786, 3554, 4218, 4337, 498},
		"ipList": []int{211, 168, 475, 489, 461, 490, 457, 96},
	}
	response.Success(c, result)
}
