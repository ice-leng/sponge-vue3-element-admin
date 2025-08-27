package types

// DashboardStatisticsReply statistics
type DashboardStatisticsReply struct {
	Code int                       `json:"code"` // return code
	Msg  string                    `json:"msg"`  // return information description
	Data []DashboardStatisticsItem `json:"data"` // return data
}

type DashboardStatisticsItem struct {
	Type             string      `json:"type"`             // 类型 "pv" | "uv" | "ip"
	Title            string      `json:"title"`            // 标题
	TodayCount       interface{} `json:"todayCount"`       // 今日数量
	TotalCount       interface{} `json:"totalCount"`       // 总数量
	GrowthRate       interface{} `json:"growthRate"`       // 增长率
	GranularityLabel string      `json:"granularityLabel"` // 粒度标签
}

// DashboardEchartsRequest request params
type DashboardEchartsRequest struct {
	StartTime string `json:"startTime,omitempty" form:"startTime" binding:""` // 开始时间
	EndTime   string `json:"endTime,omitempty" form:"endTime" binding:""`     // 结束时间
}

type Echarts struct {
	Names  []string        `json:"names"`
	Dates  []string        `json:"dates"` // 日期
	Series []EchartsSeries `json:"series"`
}

type EchartsSeries struct {
	Name      string `json:"name"`
	Data      []int  `json:"data"`
	AreaStyle string `json:"areaStyle"`
	ItemStyle string `json:"itemStyle"`
	LineStyle string `json:"lineStyle"`
}

// DashboardEchartsReply echarts
type DashboardEchartsReply struct {
	Code int     `json:"code"` // return code
	Msg  string  `json:"msg"`  // return information description
	Data Echarts `json:"data"` // return data
}
