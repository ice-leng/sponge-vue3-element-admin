package types

// DashboardStatisticsReply statistics
type DashboardStatisticsReply struct {
	Code int                       `json:"code"` // return code
	Msg  string                    `json:"msg"`  // return information description
	Data []DashboardStatisticsItem `json:"data"` // return data
}

type DashboardStatisticsItem struct {
	Type             string  `json:"type"`             // 类型 "pv" | "uv" | "ip"
	Title            string  `json:"title"`            // 标题
	TodayCount       int     `json:"todayCount"`       // 今日数量
	TotalCount       int     `json:"totalCount"`       // 总数量
	GrowthRate       float64 `json:"growthRate"`       // 增长率
	GranularityLabel string  `json:"granularityLabel"` // 粒度标签
}

// DashboardEchartsReply echarts
type DashboardEchartsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Dates  []string `json:"dates"`  // 日期
		IpList []int    `json:"ipList"` // ip 列表
		PvList []int    `json:"pvList"` // pv 列表
	} `json:"data"` // return data
}
