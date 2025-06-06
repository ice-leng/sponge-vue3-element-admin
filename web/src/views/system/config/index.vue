<!-- 系统配置 -->
<template>
  <div class="app-container">
    <div class="search-bar">
      <el-form ref="queryFormRef" :model="queryParams" :inline="true">
        <el-form-item label="关键字" prop="keywords">
          <el-input
            v-model="queryParams.name"
            placeholder="请输入配置名称"
            clearable
            @keyup.enter="handleQuery"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="handleQuery">搜索</el-button>
          <el-button icon="refresh" @click="handleResetQuery">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-card shadow="never">
      <div class="mb-10px">
        <el-button
          v-hasPerm="['sys:config:add']"
          type="success"
          icon="plus"
          @click="handleOpenDialog()"
        >
          新增
        </el-button>
        <el-button
          v-hasPerm="['sys:config:refresh']"
          color="#626aef"
          icon="RefreshLeft"
          @click="handleRefreshCache"
        >
          刷新缓存
        </el-button>
      </div>

      <el-table
        ref="dataTableRef"
        v-loading="loading"
        :data="pageData"
        highlight-current-row
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="index" label="序号" width="60" />
        <el-table-column key="name" label="配置名称" prop="name" min-width="100" />
        <el-table-column key="key" label="配置键" prop="key" min-width="100" />
        <el-table-column key="value" label="配置值" prop="value" min-width="100" />
        <el-table-column key="description" label="描述" prop="description" min-width="100" />
        <el-table-column fixed="right" label="操作" width="220">
          <template #default="scope">
            <el-button
              v-hasPerm="['sys:config:edit']"
              type="primary"
              size="small"
              link
              icon="edit"
              @click="handleOpenDialog(scope.row.id)"
            >
              编辑
            </el-button>
            <el-button
              v-hasPerm="['sys:config:delete']"
              type="danger"
              size="small"
              link
              icon="delete"
              @click="handleDelete(scope.row.id)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <pagination
        v-if="total > 0"
        v-model:total="total"
        v-model:page="queryParams.page"
        v-model:limit="queryParams.pageSize"
        @pagination="handleQuery"
      />
    </el-card>

    <!-- 系统配置表单弹窗 -->
    <el-dialog
      v-model="dialog.visible"
      :title="dialog.title"
      width="500px"
      @close="handleCloseDialog"
    >
      <el-form
        ref="dataFormRef"
        :model="formData"
        :rules="rules"
        label-suffix=":"
        label-width="100px"
      >
        <el-form-item label="配置名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入配置名称" :maxlength="50" />
        </el-form-item>
        <el-form-item label="配置键" prop="key">
          <el-input v-model="formData.key" placeholder="请输入配置键" :maxlength="50" />
        </el-form-item>
        <el-form-item label="配置值" prop="value">
          <el-input v-model="formData.value" placeholder="请输入配置值" :maxlength="100" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="formData.description"
            :rows="4"
            :maxlength="100"
            show-word-limit
            type="textarea"
            placeholder="请输入描述"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button type="primary" @click="handleSubmit">确定</el-button>
          <el-button @click="handleCloseDialog">取消</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
defineOptions({
  name: "Config",
  inheritAttrs: false,
});

import ConfigAPI, { ConfigPageVO, ConfigForm, ConfigPageQuery } from "@/api/system/config.api";

const queryFormRef = ref();
const dataFormRef = ref();

const loading = ref(false);
const selectIds = ref<number[]>([]);
const total = ref(0);

const queryParams = reactive<ConfigPageQuery>({
  page: 1,
  pageSize: 10,
  name: "",
});

// 系统配置表格数据
const pageData = ref<ConfigPageVO[]>([]);

const dialog = reactive({
  title: "",
  visible: false,
});

const formData = reactive<ConfigForm>({
  id: undefined,
  name: "",
  key: "",
  value: "",
  description: "",
});

const rules = reactive({
  name: [{ required: true, message: "请输入系统配置名称", trigger: "blur" }],
  key: [{ required: true, message: "请输入系统配置编码", trigger: "blur" }],
  value: [{ required: true, message: "请输入系统配置值", trigger: "blur" }],
});

// 查询系统配置
function handleQuery() {
  loading.value = true;
  ConfigAPI.getPage(queryParams)
    .then((data) => {
      pageData.value = data.list;
      total.value = data.total;
    })
    .finally(() => {
      loading.value = false;
    });
}

// 重置查询
function handleResetQuery() {
  queryFormRef.value.resetFields();
  queryParams.page = 1;
  handleQuery();
}

// 行复选框选中项变化
function handleSelectionChange(selection: any) {
  selectIds.value = selection.map((item: any) => item.id);
}

// 打开系统配置弹窗
function handleOpenDialog(id?: string) {
  dialog.visible = true;
  if (id) {
    dialog.title = "修改系统配置";
    ConfigAPI.getFormData(id).then((data) => {
      Object.assign(formData, data);
    });
  } else {
    dialog.title = "新增系统配置";
    formData.id = undefined;
  }
}

// 刷新缓存(防抖)
const handleRefreshCache = useDebounceFn(() => {
  ConfigAPI.refreshCache().then(() => {
    ElMessage.success("刷新成功");
  });
}, 1000);

// 系统配置表单提交
function handleSubmit() {
  dataFormRef.value.validate((valid: any) => {
    if (valid) {
      loading.value = true;
      const id = formData.id;
      if (id) {
        ConfigAPI.update(id, formData)
          .then(() => {
            ElMessage.success("修改成功");
            handleCloseDialog();
            handleResetQuery();
          })
          .finally(() => (loading.value = false));
      } else {
        ConfigAPI.create(formData)
          .then(() => {
            ElMessage.success("新增成功");
            handleCloseDialog();
            handleResetQuery();
          })
          .finally(() => (loading.value = false));
      }
    }
  });
}

// 重置表单
function resetForm() {
  dataFormRef.value.resetFields();
  dataFormRef.value.clearValidate();
  formData.id = undefined;
}

// 关闭系统配置弹窗
function handleCloseDialog() {
  dialog.visible = false;
  resetForm();
}

// 删除系统配置
function handleDelete(id: string) {
  ElMessageBox.confirm("确认删除该项配置?", "警告", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(() => {
    loading.value = true;
    ConfigAPI.deleteById(id)
      .then(() => {
        ElMessage.success("删除成功");
        handleResetQuery();
      })
      .finally(() => (loading.value = false));
  });
}

onMounted(() => {
  handleQuery();
});
</script>
