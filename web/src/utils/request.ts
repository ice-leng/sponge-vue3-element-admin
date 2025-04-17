import axios, { type InternalAxiosRequestConfig, type AxiosResponse } from "axios";
import qs from "qs";
import { useUserStoreHook } from "@/store/modules/user.store";
import { ResultEnum } from "@/enums/api/result.enum";
import { getAccessToken, setAccessToken } from "@/utils/auth";
import router from "@/router";

// 创建 axios 实例
const service = axios.create({
  baseURL: import.meta.env.VITE_APP_BASE_API,
  timeout: 50000,
  headers: { "Content-Type": "application/json;charset=utf-8" },
  paramsSerializer: (params) => qs.stringify(params),
});
// 请求拦截器
service.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const accessToken = getAccessToken();
    // 如果 Authorization 设置为 no-auth，则不携带 Token
    if (config.headers.Authorization !== "no-auth" && accessToken) {
      config.headers.Authorization = `Bearer ${accessToken}`;
    } else {
      delete config.headers.Authorization;
    }
    return config;
  },
  (error) => Promise.reject(error)
);
// 响应拦截器
service.interceptors.response.use(
  (response: AxiosResponse) => {
    // 如果响应是二进制流，则直接返回，用于下载文件、Excel 导出等
    if (response.config.responseType === "blob") {
      return response;
    }
    const { code, data, msg } = response.data;
    if (code === ResultEnum.SUCCESS) {
      if (response.headers["x-renewed-token"]) {
        setAccessToken(response.headers["x-renewed-token"]);
      }
      return data;
    }
    ElMessage.error(msg || "系统出错");
    return Promise.reject(new Error(msg || "Error"));
  },
  async (error) => {
    console.error("request error", error); // for debug
    const { config, response } = error;
    if (response) {
      const { code, msg } = response.data;
      if (code === ResultEnum.ACCESS_TOKEN_INVALID) {
        // 刷新 Token 过期，跳转登录页
        await handleSessionExpired();
        return Promise.reject(new Error(msg || "Error"));
      } else {
        ElMessage.error(msg || "系统出错");
      }
    }
    return Promise.reject(error.message);
  }
);
export default service;
// 处理会话过期
async function handleSessionExpired() {
  ElNotification({
    title: "提示",
    message: "您的会话已过期，请重新登录",
    type: "info",
  });
  await useUserStoreHook().clearSessionAndCache();
  router.push("/login");
}
