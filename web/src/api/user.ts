import request from "@/utils/request";

const USER_BASE_URL = "api/v1/admin/user";

class UserAPI {
  /**
   * 获取用户管理分页列表
   *
   * @param queryParams 查询参数
   * @returns 用户管理分页列表
   */
  static getPage(queryParams: UserQuery) {
    return request<any, PageResult<UserPageVO[]>>({
      url: `${USER_BASE_URL}`,
      method: "get",
      params: queryParams,
    });
  }

  /**
   * 获取用户管理表单详情
   *
   * @param id 用户管理ID
   * @returns 用户管理表单详情
   */
  static getFormData(id: number) {
    return request<any, UserForm>({
      url: `${USER_BASE_URL}/${id}`,
      method: "get",
    });
  }

  /**
   * 添加用户管理
   *
   * @param data 用户管理表单数据
   * @returns 请求结果
   */
  static add(data: UserForm) {
    return request({
      url: `${USER_BASE_URL}`,
      method: "post",
      data: data,
    });
  }

  /**
   * 修改用户管理
   *
   * @param id 用户管理ID
   * @param data 用户管理表单数据
   * @returns 请求结果
   */
  static update(id: number, data: UserForm) {
    return request({
      url: `${USER_BASE_URL}/${id}`,
      method: "put",
      data: data,
    });
  }

  /**
   * 批量删除用户管理，多个以英文逗号(,)分割
   *
   * @param ids 用户管理ID字符串，多个以英文逗号(,)分割
   * @returns 请求结果
   */
  static deleteByIds(ids: string) {
    return request({
      url: `${USER_BASE_URL}/${ids}`,
      method: "delete",
    });
  }
}

export default UserAPI;

/**
 * 查询对象类型
 */
export interface UserQuery extends PageQuery {
  startTime?: string;
  endTime?: string;
}

/**
 * 用户管理分页对象
 */
export interface UserPageVO {
  /**
   * id
   */
  id?: number;
  /**
   * 账号
   */
  username?: string;
  /**
   * 密码
   */
  password?: string;
  /**
   * 状态
   */
  status?: number;
}


/**
 * 用户管理表单类型
 */
export interface UserForm {
  /**
   * id
   */
  id?: number;
  /**
   * 账号
   */
  username?: string;
  /**
   * 密码
   */
  password?: string;
  /**
   * 状态
   */
  status?: number;
}

