<template>
  <div class="login-container">
    <!-- 右侧切换主题、语言按钮 -->
    <div class="flex flex-col gap-4px fixed top-40px right-40px text-lg">
      <el-tooltip :content="t('login.theneToggle')" placement="left">
        <CommonWrapper>
          <DarkModeSwitch />
        </CommonWrapper>
      </el-tooltip>
      <el-tooltip :content="t('login.languageToggle')" placement="left">
        <CommonWrapper>
          <LangSelect size="text-20px" />
        </CommonWrapper>
      </el-tooltip>
    </div>

    <!-- 登录表单 -->
    <el-card class="login-card">
      <div class="text-center relative">
        <el-image :src="logo" style="width: 84px" />
        <h2>{{ defaultSettings.title }}</h2>
        <!--        <el-tag class="ml-2 absolute-rt">{{ defaultSettings.version }}</el-tag>-->
      </div>

      <el-form ref="loginFormRef" :model="loginData" :rules="loginRules" class="login-form">
        <!-- 用户名 -->
        <el-form-item prop="username">
          <div class="input-wrapper">
            <el-icon class="mx-2">
              <User />
            </el-icon>
            <el-input
              ref="username"
              v-model="loginData.username"
              :placeholder="$t('login.username')"
              name="username"
              size="large"
              class="h-[48px]"
            />
          </div>
        </el-form-item>

        <!-- 密码 -->
        <el-tooltip :visible="isCapslock" :content="$t('login.capsLock')" placement="right">
          <el-form-item prop="password">
            <div class="input-wrapper">
              <el-icon class="mx-2">
                <Lock />
              </el-icon>
              <el-input
                v-model="loginData.password"
                :placeholder="$t('login.password')"
                type="password"
                name="password"
                size="large"
                class="h-[48px]"
                show-password
                @keyup="checkCapslock"
                @keyup.enter="handleLoginSubmit"
              />
            </div>
          </el-form-item>
        </el-tooltip>

        <!-- 验证码 -->
        <el-form-item prop="captchaCode">
          <div class="input-wrapper">
            <svg-icon icon-class="captcha" class="mx-2" />
            <el-input
              v-model="loginData.captchaCode"
              auto-complete="off"
              size="large"
              class="flex-1"
              :placeholder="$t('login.captchaCode')"
              @keyup.enter="handleLoginSubmit"
            />

            <el-image :src="captchaBase64" class="captcha-image" @click="getCaptcha" />
          </div>
        </el-form-item>

        <!-- 登录按钮 -->
        <el-button
          :loading="loading"
          type="primary"
          size="large"
          class="w-full"
          @click.prevent="handleLoginSubmit"
        >
          {{ $t("login.login") }}
        </el-button>
      </el-form>
    </el-card>

    <!-- ICP备案 -->
    <el-text size="small" class="py-2.5! fixed bottom-0 text-center">
      Copyright © 2021 - 2025 youlai.tech All Rights Reserved.
      <a href="http://beian.miit.gov.cn/" target="_blank">皖ICP备20006496号-2</a>
    </el-text>
  </div>
</template>

<script setup lang="ts">
// 外部库和依赖
import { LocationQuery, RouteLocationRaw, useRoute } from "vue-router";
import logo from "@/assets/logo.png";

// 内部依赖
import { useUserStore, useDictStore } from "@/store";
import AuthAPI, { LoginFormData } from "@/api/auth.api";
import router from "@/router";
import defaultSettings from "@/settings";
import CommonWrapper from "@/components/CommonWrapper/index.vue";
import DarkModeSwitch from "@/components/DarkModeSwitch/index.vue";

// 类型定义
import type { FormInstance } from "element-plus";
// 使用导入的依赖和库
const userStore = useUserStore();
const dictStore = useDictStore();
const route = useRoute();
// 窗口高度
const { height } = useWindowSize();
// 国际化 Internationalization
const { t } = useI18n();

// 是否显示 ICP 备案信息
const icpVisible = ref(true);
// 按钮 loading 状态
const loading = ref(false);
// 是否大写锁定
const isCapslock = ref(false);
// 验证码图片Base64字符串
const captchaBase64 = ref();
// 登录表单ref
const loginFormRef = ref<FormInstance>();

const loginData = ref<LoginFormData>({
  username: "admin",
  password: "123456",
  captchaKey: "",
  captchaCode: "",
} as LoginFormData);

const loginRules = computed(() => {
  return {
    username: [
      {
        required: true,
        trigger: "blur",
        message: t("login.message.username.required"),
      },
    ],
    password: [
      {
        required: true,
        trigger: "blur",
        message: t("login.message.password.required"),
      },
      {
        min: 6,
        message: t("login.message.password.min"),
        trigger: "blur",
      },
    ],
    captchaCode: [
      {
        required: true,
        trigger: "blur",
        message: t("login.message.captchaCode.required"),
      },
    ],
  };
});

/** 获取验证码 */
function getCaptcha() {
  AuthAPI.getCaptcha().then((data) => {
    loginData.value.captchaKey = data.captchaKey;
    captchaBase64.value = data.captchaBase64;
  });
}

/** 登录表单提交 */
function handleLoginSubmit() {
  loginFormRef.value?.validate((valid: boolean) => {
    if (valid) {
      loading.value = true;
      userStore
        .login(loginData.value)
        .then(async () => {
          await userStore.getUserInfo();
          await dictStore.loadDictItems();
          // 跳转到登录前的页面
          const redirect = resolveRedirectTarget(route.query);
          await router.push(redirect);
        })
        .catch(() => {
          getCaptcha();
        })
        .finally(() => {
          loading.value = false;
        });
    }
  });
}

/**
 * 解析重定向目标
 * @param query 路由查询参数
 * @returns 标准化后的路由地址对象
 */
function resolveRedirectTarget(query: LocationQuery): RouteLocationRaw {
  // 默认跳转路径
  const defaultPath = "/";

  // 获取原始重定向路径
  const rawRedirect = (query.redirect as string) || defaultPath;

  try {
    // 6. 使用Vue Router解析路径
    const resolved = router.resolve(rawRedirect);
    return {
      path: resolved.path,
      query: resolved.query,
    };
  } catch {
    // 7. 异常处理：返回安全路径
    return { path: defaultPath };
  }
}

/** 根据屏幕宽度切换设备模式 */
watchEffect(() => {
  if (height.value < 600) {
    icpVisible.value = false;
  } else {
    icpVisible.value = true;
  }
});

/** 检查输入大小写 */
function checkCapslock(event: KeyboardEvent) {
  // 防止浏览器密码自动填充时报错
  if (event instanceof KeyboardEvent) {
    isCapslock.value = event.getModifierState("CapsLock");
  }
}

onMounted(() => {
  getCaptcha();
});
</script>

<style lang="scss" scoped>
.login-container {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  overflow-y: auto;
  background: url("@/assets/images/login-bg.svg") no-repeat center right;

  .top-bar {
    position: absolute;
    top: 0;
    left: 0;
    display: flex;
    align-items: center;
    justify-content: flex-end;
    width: 100%;
    padding: 10px;
  }

  .login-card {
    width: 400px;
    background: transparent;
    border: none;
    border-radius: 4%;

    @media (width <= 640px) {
      width: 340px;
    }

    .input-wrapper {
      display: flex;
      align-items: center;
      width: 100%;
    }

    .captcha-image {
      height: 48px;
      cursor: pointer;
      border-top-right-radius: 6px;
      border-bottom-right-radius: 6px;
    }
  }

  .icp-info {
    position: absolute;
    bottom: 4px;
    font-size: 12px;
    text-align: center;
  }

  :deep(.el-form-item) {
    background: var(--el-input-bg-color);
    border: 1px solid var(--el-border-color);
    border-radius: 5px;
  }

  :deep(.el-input) {
    .el-input__wrapper {
      padding: 0;
      background-color: transparent;
      box-shadow: none;

      &.is-focus,
      &:hover {
        box-shadow: none !important;
      }

      input:-webkit-autofill {
        /* 通过延时渲染背景色变相去除背景颜色 */
        transition: background-color 1000s ease-in-out 0s;
      }
    }
  }
}

html.dark .login-container {
  background: url("@/assets/images/login-bg.svg") no-repeat center right;
}
</style>
