import request from "@/utils/request";

const CONFIG_BASE_URL = "api/v1/admin/config";

class ConfigAPI {
  /**
   * 获取配置分页列表
   *
   * @param queryParams 查询参数
   * @returns 配置分页列表
   */
  static getPage(queryParams: ConfigQuery) {
    return request<any, PageResult<ConfigPageVO[]>>({
      url: `${CONFIG_BASE_URL}`,
      method: "get",
      params: queryParams,
    });
  }

  /**
   * 获取配置表单详情
   *
   * @param id 配置ID
   * @returns 配置表单详情
   */
  static getFormData(id: number) {
    return request<any, ConfigForm>({
      url: `${CONFIG_BASE_URL}/${id}`,
      method: "get",
    });
  }

  /**
   * 添加配置
   *
   * @param data 配置表单数据
   * @returns 请求结果
   */
  static add(data: ConfigForm) {
    return request({
      url: `${CONFIG_BASE_URL}`,
      method: "post",
      data: data,
    });
  }

  /**
   * 修改配置
   *
   * @param id 配置ID
   * @param data 配置表单数据
   * @returns 请求结果
   */
  static update(id: number, data: ConfigForm) {
    return request({
      url: `${CONFIG_BASE_URL}/${id}`,
      method: "put",
      data: data,
    });
  }

  /**
   * 批量删除配置，多个以英文逗号(,)分割
   *
   * @param ids 配置ID字符串，多个以英文逗号(,)分割
   * @returns 请求结果
   */
  static deleteByIds(ids: string) {
    return request({
      url: `${CONFIG_BASE_URL}/${ids}`,
      method: "delete",
    });
  }
}

export default ConfigAPI;

/**
 * 查询对象类型
 */
export interface ConfigQuery extends PageQuery {
  startTime?: string;
  endTime?: string;
}

/**
 * 配置分页对象
 */
export interface ConfigPageVO {
  /**
   * id
   */
  id?: number;
  /**
   * 配置标题
   */
  name?: string;
  /**
   * 配置项描述
   */
  description?: string;
  /**
   * 配置项键
   */
  key?: string;
  /**
   * 配置值
   */
  value?: string;
}


/**
 * 配置表单类型
 */
export interface ConfigForm {
  /**
   * id
   */
  id?: number;
  /**
   * 配置标题
   */
  name?: string;
  /**
   * 配置项描述
   */
  description?: string;
  /**
   * 配置项键
   */
  key?: string;
  /**
   * 配置值
   */
  value?: string;
}

