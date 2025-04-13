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
  startDate: string;
  /** 结束日期 */
  endDate: string;
}

/**  访问趋势视图对象 */
export interface EchartsVO {
  /** 日期列表 */
  dates: string[];
  /** 浏览量(PV) */
  pvList: number[];
  /** 访客数(UV) */
  uvList: number[];
  /** IP数 */
  ipList: number[];
}
