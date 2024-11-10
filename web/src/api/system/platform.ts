import request from "@/utils/request";

const PLATFORM_BASE_URL = "/api/v1/platform";

const PlatformAPI = {
  /**
   * 获取当前登录用户信息
   *
   * @returns 登录用户昵称、头像信息，包括角色和权限
   */
  getInfo() {
    return request<any, PlatformInfo>({
      url: `${PLATFORM_BASE_URL}/me`,
      method: "get",
    });
  },

  /**
   * 获取用户分页列表
   *
   * @param queryParams 查询参数
   */
  getPage(queryParams: PlatformPageQuery) {
    return request<any, PageResult<PlatformPageVO[]>>({
      url: `${PLATFORM_BASE_URL}`,
      method: "get",
      params: queryParams,
    });
  },

  /**
   * 获取用户表单详情
   *
   * @param platformId 用户ID
   * @returns 用户表单详情
   */
  getFormData(platformId: number) {
    return request<any, PlatformForm>({
      url: `${PLATFORM_BASE_URL}/${platformId}`,
      method: "get",
    });
  },

  /**
   * 添加用户
   *
   * @param data 用户表单数据
   */
  add(data: PlatformForm) {
    return request({
      url: `${PLATFORM_BASE_URL}`,
      method: "post",
      data: data,
    });
  },

  /**
   * 修改用户
   *
   * @param id 用户ID
   * @param data 用户表单数据
   */
  update(id: number, data: PlatformForm) {
    return request({
      url: `${PLATFORM_BASE_URL}/${id}`,
      method: "put",
      data: data,
    });
  },

  /**
   * 修改用户密码
   *
   * @param id 用户ID
   * @param password 新密码
   */
  resetPassword(id: number, password: string) {
    return request({
      url: `${PLATFORM_BASE_URL}/password/reset`,
      method: "put",
      data: { id: id, password: password },
    });
  },

  /**
   * 批量删除用户，多个以英文逗号(,)分割
   *
   * @param ids 用户ID字符串，多个以英文逗号(,)分割
   */
  deleteByIds(ids: string) {
    return request({
      url: `${PLATFORM_BASE_URL}/${ids}`,
      method: "delete",
    });
  },

  /** 获取个人中心用户信息 */
  getProfile() {
    return request<any, PlatformProfileVO>({
      url: `${PLATFORM_BASE_URL}/profile`,
      method: "get",
    });
  },

  /** 修改个人中心用户信息 */
  updateProfile(data: PlatformProfileForm) {
    return request({
      url: `${PLATFORM_BASE_URL}/profile`,
      method: "put",
      data: data,
    });
  },

  /** 修改个人中心用户密码 */
  changePassword(data: PasswordChangeForm) {
    return request({
      url: `${PLATFORM_BASE_URL}/password`,
      method: "put",
      data: data,
    });
  },
};

export default PlatformAPI;

/** 登录用户信息 */
export interface PlatformInfo {
  /** 用户ID */
  id?: number;

  /** 用户名 */
  username?: string;

  /** 头像URL */
  avatar?: string;

  /** 角色 */
  roles: string[];

  /** 权限 */
  perms: string[];
}

/**
 * 用户分页查询对象
 */
export interface PlatformPageQuery extends PageQuery {
  /** 搜索关键字 */
  username?: string;

  /** 用户状态 */
  status?: number;

  /** 开始时间 */
  startTime?: string;
  endTime?: string;
}

/** 用户分页对象 */
export interface PlatformPageVO {
  /** 用户头像URL */
  avatar?: string;
  /** 创建时间 */
  createdAt?: Date;
  /** 用户ID */
  id?: number;
  /** 角色名称，多个使用英文逗号(,)分割 */
  roleNames?: string;
  /** 用户状态(1:启用;0:禁用) */
  status?: number;
  /** 用户名 */
  username?: string;
}

/** 用户表单类型 */
export interface PlatformForm {
  /** 用户头像 */
  avatar?: string;
  /** 用户ID */
  id?: number;
  /** 角色ID集合 */
  roleId?: number[];
  /** 用户状态(1:正常;0:禁用) */
  status?: number;
  /** 用户名 */
  username?: string;
}

/** 个人中心用户信息 */
export interface PlatformProfileVO {
  /** 用户ID */
  id?: number;

  /** 用户名 */
  username?: string;

  /** 头像URL */
  avatar?: string;

  /** 角色名称，多个使用英文逗号(,)分割 */
  roleNames?: string;

  /** 创建时间 */
  createdAt?: Date;
}

/** 个人中心用户信息表单 */
export interface PlatformProfileForm {
  /** 用户ID */
  id?: number;

  /** 用户名 */
  username?: string;

  /** 头像URL */
  avatar?: string;
}

/** 修改密码表单 */
export interface PasswordChangeForm {
  /** 原密码 */
  oldPassword?: string;
  /** 新密码 */
  newPassword?: string;
  /** 确认新密码 */
  confirmPassword?: string;
}
