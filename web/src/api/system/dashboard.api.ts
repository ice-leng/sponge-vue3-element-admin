import request from "@/utils/request";

const DASHBOARD_BASE_URL = "/api/v1/dashboard";

const DashboardAPI = {
  getStatistics() {
    return request<any, StatisticsVO[]>({
      url: `${DASHBOARD_BASE_URL}/statistics`,
      method: "get",
    });
  },

  getEcharts(queryParams: EchartsQuery) {
    return request<any, EchartsVO>({
      url: `${DASHBOARD_BASE_URL}/echarts`,
      method: "get",
      params: queryParams,
    });
  },
};

export default DashboardAPI;

/**  访问统计 */
export interface StatisticsVO {
  /** 标题 */
  title: string;
  /** 类型 */
  type: "pv" | "uv" | "ip";

  /** 今日访问量 */
  todayCount: number;
  /** 总访问量 */
  totalCount: number;
  /** 同比增长率（相对于昨天同一时间段的增长率） */
  growthRate: number;
}

/** 访问趋势查询参数 */
export interface EchartsQuery {
  /** 开始日期 */
  startTime: string;
  /** 结束日期 */
  endTime: string;
}

/**  访问趋势视图对象 */
export interface EchartsVO {
  /** 名称列表 */
  names: string[];
  /** 日期列表 */
  dates: string[];
  /** 数据列表 */
  series: EchartsSeries[];
}

/**  访问趋势数据项 */
export interface EchartsSeries {
  /** 名称 */
  name: string;
  /** 数据 */
  data: any[];
  /** 面积样式 */
  areaStyle: string;
  /** 项目样式 */
  itemStyle: string;
  /** 线样式 */
  lineStyle: string;
}
