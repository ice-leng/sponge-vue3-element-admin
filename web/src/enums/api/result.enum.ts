/**
 * 响应码枚举
 */
export const enum ResultEnum {
  /**
   * 成功
   */
  SUCCESS = 0,
  /**
   * 错误
   */
  ERROR = 500,

  /**
   * 访问令牌无效或过期
   */
  ACCESS_TOKEN_INVALID = 401,
}
