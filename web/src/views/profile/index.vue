<template>
  <div class="app-container">
    <el-tabs tab-position="left">
      <!-- 基本设置 Tab Pane -->
      <el-tab-pane label="账号信息">
        <div class="w-full">
          <el-card>
            <!-- 头像和昵称部分 -->
            <div class="relative w-100px h-100px flex-center">
              <el-avatar :src="userProfile.avatar" :size="100" />
              <el-button
                type="info"
                class="absolute bottom-0 right-0 cursor-pointer"
                circle
                :icon="Camera"
                size="small"
                @click="triggerFileUpload"
              />
              <input
                ref="fileInput"
                type="file"
                style="display: none"
                @change="handleFileChange"
              />
            </div>
            <!-- 用户信息描述 -->
            <el-descriptions :column="1" class="mt-10">
              <!-- 用户名 -->
              <el-descriptions-item>
                <template #label>
                  <el-icon class="align-middle"><User /></el-icon>
                  账号
                </template>
                {{ userProfile.username }}
              </el-descriptions-item>
              <el-descriptions-item>
                <template #label>
                  <SvgIcon icon-class="role" />
                  角色
                </template>
                {{ userProfile.roleNames }}
              </el-descriptions-item>

              <el-descriptions-item>
                <template #label>
                  <el-icon class="align-middle"><Timer /></el-icon>
                  创建日期
                </template>
                {{ userProfile.createdAt }}
              </el-descriptions-item>
            </el-descriptions>
          </el-card>
        </div>
      </el-tab-pane>

      <!-- 安全设置  -->
      <el-tab-pane label="安全设置">
        <el-card>
          <!-- 账户密码 -->
          <el-row>
            <el-col :span="16">
              <div class="font-bold">账户密码</div>
              <div class="text-14px mt-2">
                定期修改密码有助于保护账户安全
                <el-button
                  type="primary"
                  plain
                  size="small"
                  class="ml-5"
                  @click="() => handleOpenDialog(DialogType.PASSWORD)"
                >
                  修改
                </el-button>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-tab-pane>
    </el-tabs>

    <!-- 弹窗 -->
    <el-dialog v-model="dialog.visible" :title="dialog.title" :width="500">
      <!-- 修改密码 -->
      <el-form
        v-if="dialog.type === DialogType.PASSWORD"
        ref="passwordChangeFormRef"
        :model="passwordChangeForm"
        :rules="passwordChangeRules"
        :label-width="100"
      >
        <el-form-item label="原密码" prop="oldPassword">
          <el-input
            v-model="passwordChangeForm.oldPassword"
            type="password"
            show-password
          />
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input
            v-model="passwordChangeForm.newPassword"
            type="password"
            show-password
          />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input
            v-model="passwordChangeForm.confirmPassword"
            type="password"
            show-password
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialog.visible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import PlatformAPI, {
  PlatformProfileVO,
  PasswordChangeForm,
  PlatformProfileForm,
} from "@/api/system/platform";

import FileAPI from "@/api/file";

import { Camera } from "@element-plus/icons-vue";

const userProfile = ref<PlatformProfileVO>({});

enum DialogType {
  PASSWORD = "password",
}

const dialog = reactive({
  visible: false,
  title: "",
  type: "" as DialogType, // 修改账号资料,修改密码、绑定手机、绑定邮箱
});

const userProfileForm = reactive<PlatformProfileForm>({});
const passwordChangeForm = reactive<PasswordChangeForm>({});

// 修改密码校验规则
const passwordChangeRules = {
  oldPassword: [{ required: true, message: "请输入原密码", trigger: "blur" }],
  newPassword: [{ required: true, message: "请输入新密码", trigger: "blur" }],
  confirmPassword: [
    { required: true, message: "请再次输入新密码", trigger: "blur" },
  ],
};

/**
 * 打开弹窗
 * @param type 弹窗类型 ACCOUNT: 账号资料 PASSWORD: 修改密码 MOBILE: 绑定手机 EMAIL: 绑定邮箱
 */
const handleOpenDialog = (type: DialogType) => {
  dialog.type = type;
  dialog.visible = true;
  switch (type) {
    case DialogType.PASSWORD:
      dialog.title = "修改密码";
      break;
  }
};

/**
 * 提交表单
 */
const handleSubmit = async () => {
  if (dialog.type === DialogType.PASSWORD) {
    if (passwordChangeForm.newPassword !== passwordChangeForm.confirmPassword) {
      ElMessage.error("两次输入的密码不一致");
      return;
    }
    PlatformAPI.changePassword(passwordChangeForm).then(() => {
      ElMessage.success("密码修改成功");
      dialog.visible = false;
    });
  }
};

const fileInput = ref<HTMLInputElement | null>(null);

const triggerFileUpload = () => {
  fileInput.value?.click();
};

const handleFileChange = async (event: Event) => {
  const target = event.target as HTMLInputElement;
  const file = target.files ? target.files[0] : null;
  if (file) {
    // 调用文件上传API
    try {
      const data = await FileAPI.upload(file);
      // 更新用户头像
      userProfile.value.avatar = data.url;
      // 更新用户信息
      await PlatformAPI.updateProfile({
        avatar: data.path,
      });
    } catch (error) {
      ElMessage.error("头像上传失败");
    }
  }
};

/** 加载用户信息 */
const loadUserProfile = async () => {
  const data = await PlatformAPI.getProfile();
  userProfile.value = data;
};

onMounted(async () => {
  await loadUserProfile();
});
</script>

<style lang="scss" scoped>
/** 关闭tag标签  */
.app-container {
  /* 50px = navbar = 50px */
  height: calc(100vh - 50px);
  background: var(--el-fill-color-blank);
}

/** 开启tag标签  */
.hasTagsView {
  .app-container {
    /* 84px = navbar + tags-view = 50px + 34px */
    height: calc(100vh - 84px);
  }
}
</style>
