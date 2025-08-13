<!-- 单图上传组件 -->
<template>
  <el-upload
    v-model="displayImageUrl"
    class="single-upload"
    list-type="picture-card"
    :show-file-list="false"
    :accept="props.accept"
    :before-upload="handleBeforeUpload"
    :http-request="handleUpload"
    :on-success="onSuccess"
    :on-error="onError"
  >
    <template #default>
      <el-image v-if="displayImageUrl" :src="displayImageUrl" />
      <el-icon v-if="displayImageUrl" class="single-upload__delete-btn" @click.stop="handleDelete">
        <CircleCloseFilled />
      </el-icon>
      <el-icon v-else class="single-upload__add-btn">
        <Plus />
      </el-icon>
    </template>
  </el-upload>
</template>

<script setup lang="ts">
import { UploadRawFile, UploadRequestOptions } from "element-plus";
import FileAPI, { FileInfo } from "@/api/file.api";

const props = defineProps({
  /**
   * 请求携带的额外参数
   */
  data: {
    type: Object,
    default: () => {
      return {};
    },
  },
  /**
   * 上传文件的参数名
   */
  name: {
    type: String,
    default: "file",
  },
  /**
   * 最大文件大小（单位：M）
   */
  maxFileSize: {
    type: Number,
    default: 10,
  },

  /**
   * 上传图片格式，默认支持所有图片(image/*)，指定格式示例：'.png,.jpg,.jpeg,.gif,.bmp'
   */
  accept: {
    type: String,
    default: "image/*",
  },

  /**
   * 自定义样式，用于设置组件的宽度和高度等其他样式
   */
  style: {
    type: Object,
    default: () => {
      return {
        width: "150px",
        height: "150px",
      };
    },
  },
});

// modelValue 存储 path，用于表单提交
const modelValue = defineModel("modelValue", {
  type: String,
  default: () => "",
});

// 用于显示的图片URL
const displayImageUrl = ref("");

/**
 * 监听 modelValue 变化，处理回显逻辑
 * 当 modelValue 是 path 时，需要转换为完整的 URL 用于显示
 */
watch(
  () => modelValue.value,
  (newPath: string) => {
    if (!newPath) {
      displayImageUrl.value = "";
      return;
    }

    // 如果是完整的URL（包含协议），直接使用，但确保 modelValue 不包含协议和域名
    if (newPath.startsWith("http://") || newPath.startsWith("https://")) {
      displayImageUrl.value = newPath;
      // 从 URL 中提取路径部分（去掉协议和域名），更新 modelValue
      try {
        const urlObj = new URL(newPath);
        modelValue.value = urlObj.pathname;
      } catch (error) {
        console.error("解析URL失败:", error);
      }
      return;
    }

    // 如果是相对路径，需要构建完整的URL用于显示
    // 使用配置的API URL作为基础URL
    const baseUrl = import.meta.env.VITE_APP_API_URL || "";
    displayImageUrl.value = baseUrl + newPath;
  },
  { immediate: true }
);

/**
 * 限制用户上传文件的格式和大小
 */
function handleBeforeUpload(file: UploadRawFile) {
  // 校验文件类型：虽然 accept 属性限制了用户在文件选择器中可选的文件类型，但仍需在上传时再次校验文件实际类型，确保符合 accept 的规则
  const acceptTypes = props.accept.split(",").map((type) => type.trim());

  // 检查文件格式是否符合 accept
  const isValidType = acceptTypes.some((type) => {
    if (type === "image/*") {
      // 如果是 image/*，检查 MIME 类型是否以 "image/" 开头
      return file.type.startsWith("image/");
    } else if (type.startsWith(".")) {
      // 如果是扩展名 (.png, .jpg)，检查文件名是否以指定扩展名结尾
      return file.name.toLowerCase().endsWith(type);
    } else {
      // 如果是具体的 MIME 类型 (image/png, image/jpeg)，检查是否完全匹配
      return file.type === type;
    }
  });

  if (!isValidType) {
    ElMessage.warning(`上传文件的格式不正确，仅支持：${props.accept}`);
    return false;
  }

  // 限制文件大小
  if (file.size > props.maxFileSize * 1024 * 1024) {
    ElMessage.warning("上传图片不能大于" + props.maxFileSize + "M");
    return false;
  }
  return true;
}

/*
 * 上传图片
 */
function handleUpload(options: UploadRequestOptions) {
  return new Promise((resolve, reject) => {
    const file = options.file;

    const formData = new FormData();
    formData.append(props.name, file);

    // 处理附加参数
    Object.keys(props.data).forEach((key) => {
      formData.append(key, props.data[key]);
    });

    FileAPI.upload(formData)
      .then((data) => {
        resolve(data);
      })
      .catch((error) => {
        reject(error);
      });
  });
}

/**
 * 删除图片
 */
function handleDelete() {
  modelValue.value = "";
  displayImageUrl.value = "";
}

/**
 * 上传成功回调
 *
 * @param fileInfo 上传成功后的文件信息
 */
const onSuccess = (fileInfo: FileInfo) => {
  ElMessage.success("上传成功");
  // modelValue 存储 path，用于表单提交
  modelValue.value = fileInfo.path;
  // 显示用的是 url
  displayImageUrl.value = fileInfo.url;
};

/**
 * 上传失败回调
 */
const onError = (error: any) => {
  console.log("onError");
  ElMessage.error("上传失败: " + error.message);
};
</script>

<style scoped lang="scss">
:deep(.el-upload--picture-card) {
  /*  width: var(--el-upload-picture-card-size);
  height: var(--el-upload-picture-card-size); */
  width: v-bind("props.style.width");
  height: v-bind("props.style.height");
}

.single-upload {
  position: relative;
  overflow: hidden;
  cursor: pointer;
  border: 1px var(--el-border-color) solid;
  border-radius: 5px;

  &:hover {
    border-color: var(--el-color-primary);
  }

  &__delete-btn {
    position: absolute;
    top: 1px;
    right: 1px;
    font-size: 16px;
    color: #ff7901;
    cursor: pointer;
    background: #fff;
    border-radius: 100%;

    :hover {
      color: #ff4500;
    }
  }
}
</style>
