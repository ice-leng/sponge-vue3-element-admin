import request from "@/utils/request";

const ROLE_BASE_URL = "/api/v1/admin/role";

class RoleAPI {
  /** 获取角色分页数据 */
  static getPage(queryParams?: RolePageQuery) {
    return request<any, PageResult<RolePageVO[]>>({
      url: `${ROLE_BASE_URL}`,
      method: "get",
      params: queryParams,
    });
  }

  /** 获取角色下拉数据源 */
  static getOptions() {
    return request<any, OptionType[]>({
      url: `${ROLE_BASE_URL}/options`,
      method: "get",
    });
  }

  /**
   * 获取角色的菜单ID集合
   *
   * @param roleId 角色ID
   * @returns 角色的菜单ID集合
   */
  static getRoleMenuIds(roleId: number) {
    return request<any, number[]>({
      url: `${ROLE_BASE_URL}/menuId/${roleId}`,
      method: "get",
    });
  }

  /**
   * 分配菜单权限
   *
   * @param roleId 角色ID
   * @param data 菜单ID集合
   */
  static updateRoleMenus(roleId: number, data: number[]) {
    return request({
      url: `${ROLE_BASE_URL}/menu/${roleId}`,
      method: "put",
      data: data,
    });
  }

  /**
   * 获取角色表单数据
   *
   * @param id 角色ID
   * @returns 角色表单数据
   */
  static getFormData(id: number) {
    return request<any, RoleForm>({
      url: `${ROLE_BASE_URL}/${id}`,
      method: "get",
    });
  }

  /** 添加角色 */
  static add(data: RoleForm) {
    return request({
      url: `${ROLE_BASE_URL}`,
      method: "post",
      data: data,
    });
  }

  /**
   * 更新角色
   *
   * @param id 角色ID
   * @param data 角色表单数据
   */
  static update(id: number, data: RoleForm) {
    return request({
      url: `${ROLE_BASE_URL}/${id}`,
      method: "put",
      data: data,
    });
  }

  /**
   * 批量删除角色，多个以英文逗号(,)分割
   *
   * @param ids 角色ID字符串，多个以英文逗号(,)分割
   */
  static deleteByIds(ids: string) {
    return request({
      url: `${ROLE_BASE_URL}/${ids}`,
      method: "delete",
    });
  }
}

export default RoleAPI;

/** 角色分页查询参数 */
export interface RolePageQuery extends PageQuery {
  /** 搜索关键字 */
  name?: string;
}

/** 角色分页对象 */
export interface RolePageVO {
  /** 角色编码 */
  code?: string;
  /** 角色ID */
  id?: number;
  /** 角色名称 */
  name?: string;
  /** 排序 */
  sort?: number;
  /** 角色状态 */
  status?: number;
  /** 创建时间 */
  createdAt?: Date;
  /** 修改时间 */
  updatedAt?: Date;
}

/** 角色表单对象 */
export interface RoleForm {
  /** 角色ID */
  id?: number;
  /** 角色编码 */
  code: string;
  /** 角色名称 */
  name: string;
  /** 排序 */
  sort?: number;
  /** 角色状态(1-正常；0-停用) */
  status?: number;
}
