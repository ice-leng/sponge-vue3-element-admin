<template>
  <div class="app-container">
    <!-- 搜索 -->
    <div class="search-container">
      <el-form ref="queryFormRef" :model="queryParams" :inline="true">
        <el-form-item label="创建时间">
          <el-date-picker
            class="!w-[240px]"
            v-model="dateTimeRange"
            type="daterange"
            range-separator="~"
            start-placeholder="开始时间"
            end-placeholder="截止时间"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleQuery"
          >
            <i-ep-search/>
            搜索
          </el-button
          >
          <el-button @click="handleResetQuery">
            <i-ep-refresh/>
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 操作按钮 -->
    <el-card shadow="never" class="table-container">
      <template #header>
        <el-button
          v-hasPerm="['sys:config:add']"
          type="success"
          @click="handleOpenDialog()"
        >
          <i-ep-plus/>
          新增
        </el-button
        >
        <el-button
          v-hasPerm="['sys:config:delete']"
          type="danger"
          :disabled="removeIds.length === 0"
          @click="handleDelete()"
        >
          <i-ep-delete/>
          删除
        </el-button
        >
      </template>

      <el-table
        v-loading="loading"
        :data="pageData"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="50" align="center"/>
        <el-table-column key="id" label="序号" align="center" prop="id" width="180"/>
        <el-table-column key="createdAt" label="创建时间" align="center" prop="createdAt" width="180"/>
        <el-table-column key="name" label="配置名称" align="center" prop="name" width="180"/>
        <el-table-column key="key" label="配置项键" align="center" prop="key" width="180"/>
        <el-table-column key="value" label="配置值" align="center" prop="value"/>
        <el-table-column key="description" label="描述" align="center" prop="description"/>

        <el-table-column label="操作" fixed="right" width="220">
          <template #default="scope">
            <el-button
              v-hasPerm="['sys:config:edit']"
              type="primary"
              link
              size="small"
              @click="handleOpenDialog(scope.row.id)"
            >
              <i-ep-edit/>
              编辑
            </el-button
            >
            <el-button
              v-hasPerm="['sys:config:delete']"
              type="danger"
              link
              size="small"
              @click="handleDelete(scope.row.id)"
            >
              <i-ep-delete/>
              删除
            </el-button
            >
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

    <!-- 表单弹窗 -->
    <el-dialog
      v-model="dialog.visible"
      :title="dialog.title"
      width="500px"
      @close="handleCloseDialog"
    >
      <el-form
        ref="configFormRef"
        :model="formData"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="配置名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入配置名称"/>
        </el-form-item>
        <el-form-item label="配置项键" prop="key">
          <el-input v-model="formData.key" placeholder="请输入配置项键"/>
        </el-form-item>
        <el-form-item label="配置值" prop="value">
          <el-input v-model="formData.value" placeholder="请输入配置值"/>
        </el-form-item>
        <el-form-item label="配置项描述" prop="description">
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
          <el-button type="primary" @click="handleSubmit">确 定</el-button>
          <el-button @click="handleCloseDialog">取 消</el-button>
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

import ConfigAPI, {
  ConfigQuery,
  ConfigPageVO,
  ConfigForm,
} from "@/api/config";

const queryFormRef = ref(ElForm);
const configFormRef = ref(ElForm);

const loading = ref(false);
const removeIds = ref([]);
const total = ref(0);

const pageData = ref<ConfigPageVO[]>();

/** 查询参数  */
const queryParams = reactive<ConfigQuery>({
  page: 1,
  pageSize: 10,
});

const dateTimeRange = ref("");
watch(dateTimeRange, (newVal) => {
  if (newVal) {
    queryParams.startTime = newVal[0];
    queryParams.endTime = newVal[1];
  }
});

// 弹窗配置
const dialog = reactive({
  title: "",
  visible: false,
});
// 配置表单数据
const formData = reactive<ConfigForm>({
  id: undefined,
});

const rules = reactive({
  id: [{required: true, message: "id不能为空", trigger: "blur"}],
  name: [{required: true, message: "配置名称不能为空", trigger: "blur"}],
  description: [{required: true, message: "配置项描述不能为空", trigger: "blur"}],
  key: [{required: true, message: "配置项键不能为空", trigger: "blur"}],
  value: [{required: true, message: "配置值不能为空", trigger: "blur"}],

});

/** 查询 */
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

/** 重置查询 */
function handleResetQuery() {
  queryFormRef.value.resetFields();
  dateTimeRange.value = "";
  queryParams.page = 1;
  queryParams.startTime = undefined;
  queryParams.endTime = undefined;
  handleQuery();
}

/** 行复选框选中记录选中ID集合 */
function handleSelectionChange(selection: any) {
  removeIds.value = selection.map((item: any) => item.id);
}

/** 打开配置弹窗 */
async function handleOpenDialog(id?: number) {
  dialog.visible = true;
  if (id) {
    dialog.title = "修改配置";
    ConfigAPI.getFormData(id).then((data) => {
      Object.assign(formData, {...data});
    });
  } else {
    dialog.title = "新增配置";
  }
}

/** 关闭配置弹窗 */
function handleCloseDialog() {
  dialog.visible = false;

  configFormRef.value.resetFields();
  configFormRef.value.clearValidate();

  formData.id = undefined;
}

/** 提交配置表单 */
function handleSubmit() {
  configFormRef.value.validate((valid: any) => {
    if (valid) {
      loading.value = true;
      const id = formData.id;
      if (id) {
        ConfigAPI.update(id, formData)
          .then(() => {
            ElMessage.success("修改配置成功");
            handleCloseDialog();
            handleResetQuery();
          })
          .finally(() => (loading.value = false));
      } else {
        ConfigAPI.add(formData)
          .then(() => {
            ElMessage.success("新增配置成功");
            handleCloseDialog();
            handleResetQuery();
          })
          .finally(() => (loading.value = false));
      }
    }
  });
}

/** 删除配置 */
function handleDelete(id?: number) {
  const ids = [id || removeIds.value].join(",");
  if (!ids) {
    ElMessage.warning("请勾选删除项");
    return;
  }

  ElMessageBox.confirm("确认删除已选中的数据项?", "警告", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(
    () => {
      loading.value = true;
      ConfigAPI.deleteByIds(ids)
        .then(() => {
          ElMessage.success("删除成功");
          handleResetQuery();
        })
        .finally(() => (loading.value = false));
    },
    () => {
      ElMessage.info("已取消删除");
    }
  );
}

onMounted(() => {
  handleQuery();
});
</script>
