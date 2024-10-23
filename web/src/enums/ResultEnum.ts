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
   * 令牌无效或过期
   */
  TOKEN_INVALID = 401,
}
