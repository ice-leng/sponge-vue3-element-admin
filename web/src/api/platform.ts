import request from "@/utils/request";

const PLATFORM_BASE_URL = "/api/v1/admin/platform";

class PlatformAPI {
  /**
   * 登录成功后获取用户信息（昵称、头像、权限集合和角色集合）
   *
   * @returns 用户信息
   */
  static getInfo() {
    return request<any, PlatformInfo>({
      url: `${PLATFORM_BASE_URL}/me`,
      method: "get",
    });
  }

  /**
   * 获取用户分页列表
   *
   * @param queryParams 查询参数
   * @returns 用户分页列表
   */
  static getPage(queryParams: PlatformQuery) {
    return request<any, PageResult<PlatformPageVO[]>>({
      url: `${PLATFORM_BASE_URL}`,
      method: "get",
      params: queryParams,
    });
  }

  /**
   * 获取用户表单详情
   *
   * @param id 用户ID
   * @returns 用户表单详情
   */
  static getFormData(id: number) {
    return request<any, PlatformForm>({
      url: `${PLATFORM_BASE_URL}/${id}`,
      method: "get",
    });
  }

  /**
   * 添加用户
   *
   * @param data 用户表单数据
   * @returns 请求结果
   */
  static add(data: PlatformForm) {
    return request({
      url: `${PLATFORM_BASE_URL}`,
      method: "post",
      data: data,
    });
  }

  /**
   * 修改用户
   *
   * @param id 用户ID
   * @param data 用户表单数据
   * @returns 请求结果
   */
  static update(id: number, data: PlatformForm) {
    return request({
      url: `${PLATFORM_BASE_URL}/${id}`,
      method: "put",
      data: data,
    });
  }

  /**
   * 批量删除用户，多个以英文逗号(,)分割
   *
   * @param ids 用户ID字符串，多个以英文逗号(,)分割
   * @returns 请求结果
   */
  static deleteByIds(ids: string) {
    return request({
      url: `${PLATFORM_BASE_URL}/${ids}`,
      method: "delete",
    });
  }
}

export default PlatformAPI;

/**
 * 登录用户信息
 */
export interface PlatformInfo {
  id?: number;
  username?: string;
  avatar?: string;
  roles: string[];
  perms: string[];
}

/**
 * 用户查询对象类型
 */
export interface PlatformQuery extends PageQuery {
  username?: string;
  status?: number;
  startTime?: string;
  endTime?: string;
}

/**
 * 用户分页对象
 */
export interface PlatformPageVO {
  /**
   * 创建时间
   */
  createdAt?: Date;
  /**
   * 用户ID
   */
  id?: number;
  /**
   * 用户昵称
   */
  nickname?: string;
  /**
   * 角色名称，多个使用英文逗号(,)分割
   */
  roleNames?: string;
  /**
   * 用户状态(1:启用;0:禁用)
   */
  status?: number;
  /**
   * 用户名
   */
  username?: string;
}

/**
 * 用户表单类型
 */
export interface PlatformForm {
  /**
   * 用户头像
   */
  avatar?: string;
  /**
   * ID
   */
  id?: number;
  /**
   * 角色ID集合
   */
  roleId?: number[];
  /**
   * 用户状态(1:正常;0:禁用)
   */
  status?: number;
  /**
   * 账号
   */
  username?: string;
  /**
   * 密码
   */
  password?: string;
}
