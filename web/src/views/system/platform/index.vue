<template>
  <div class="app-container">
    <!-- 搜索 -->
    <div class="search-container">
      <el-form ref="queryFormRef" :model="queryParams" :inline="true">
        <el-form-item label="用户名" prop="keywords">
          <el-input
            v-model="queryParams.username"
            placeholder="用户名"
            clearable
            style="width: 200px"
            @keyup.enter="handleQuery"
          />
        </el-form-item>

        <el-form-item label="状态" prop="status">
          <el-select
            v-model="queryParams.status"
            placeholder="全部"
            clearable
            class="!w-[100px]"
          >
            <el-option label="启用" value="1" />
            <el-option label="禁用" value="0" />
          </el-select>
        </el-form-item>

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
            ><i-ep-search />搜索</el-button
          >
          <el-button @click="handleResetQuery"><i-ep-refresh />重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 操作按钮 -->
    <el-card shadow="never" class="table-container">
      <template #header>
        <el-button
          v-hasPerm="['sys:platform:add']"
          type="success"
          @click="handleOpenDialog()"
          ><i-ep-plus />新增</el-button
        >
        <el-button
          v-hasPerm="['sys:platform:delete']"
          type="danger"
          :disabled="removeIds.length === 0"
          @click="handleDelete()"
          ><i-ep-delete />删除</el-button
        >
      </template>

      <el-table
        v-loading="loading"
        :data="pageData"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="50" align="center" />
        <el-table-column
          key="id"
          label="编号"
          align="center"
          prop="id"
          width="100"
        />
        <el-table-column
          key="username"
          label="账号"
          align="center"
          prop="username"
        />

        <el-table-column
          key="roleNames"
          label="角色"
          align="center"
          prop="roleNames"
        />

        <el-table-column label="状态" align="center" prop="status">
          <template #default="scope">
            <el-tag :type="scope.row.status == 1 ? 'success' : 'info'">{{
              scope.row.status == 1 ? "启用" : "禁用"
            }}</el-tag>
          </template>
        </el-table-column>

        <el-table-column
          label="最新登录时间"
          align="center"
          prop="lastTime"
          width="180"
        />

        <el-table-column
          label="创建时间"
          align="center"
          prop="createdAt"
          width="180"
        />

        <el-table-column label="操作" fixed="right" width="220">
          <template #default="scope">
            <el-button
              v-hasPerm="['sys:platform:edit']"
              type="primary"
              link
              size="small"
              @click="handleOpenDialog(scope.row.id)"
              ><i-ep-edit />编辑</el-button
            >
            <el-button
              v-hasPerm="['sys:platform:delete']"
              type="danger"
              link
              size="small"
              @click="handleDelete(scope.row.id)"
              ><i-ep-delete />删除</el-button
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

    <!-- 管理员表单弹窗 -->
    <el-dialog
      v-model="dialog.visible"
      :title="dialog.title"
      width="500px"
      @close="handleCloseDialog"
    >
      <el-form
        ref="platformFormRef"
        :model="formData"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="账号" prop="username">
          <el-input
            v-model="formData.username"
            :readonly="!!formData.id"
            placeholder="请输入账号"
          />
        </el-form-item>

        <el-form-item label="密码" prop="password">
          <el-input v-model="formData.password" placeholder="请输入密码" />
        </el-form-item>

        <el-form-item label="角色" prop="roleId">
          <el-select v-model="formData.roleId" multiple placeholder="请选择">
            <el-option
              v-for="item in roleOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio :label="1">正常</el-radio>
            <el-radio :label="0">停用</el-radio>
          </el-radio-group>
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
  name: "Platform",
  inheritAttrs: false,
});

import PlatformAPI, {
  PlatformQuery,
  PlatformPageVO,
  PlatformForm,
} from "@/api/platform";
import RoleAPI from "@/api/role";

const queryFormRef = ref(ElForm);
const platformFormRef = ref(ElForm);

const loading = ref(false);
const removeIds = ref([]);
const total = ref(0);

const pageData = ref<PlatformPageVO[]>();
/** 角色下拉选项 */
const roleOptions = ref<OptionType[]>();
/** 用户查询参数  */
const queryParams = reactive<PlatformQuery>({
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

// 弹窗
const dialog = reactive({
  title: "",
  visible: false,
});
// 管理员表单数据
const formData = reactive<PlatformForm>({
  status: 1,
});

const rules = reactive({
  username: [{ required: true, message: "账号不能为空", trigger: "blur" }],
  roleId: [{ required: true, message: "角色不能为空", trigger: "blur" }],
});

/** 查询 */
function handleQuery() {
  loading.value = true;
  PlatformAPI.getPage(queryParams)
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
  queryParams.username = undefined;
  handleQuery();
}

/** 行复选框选中记录选中ID集合 */
function handleSelectionChange(selection: any) {
  removeIds.value = selection.map((item: any) => item.id);
}

/** 打开管理员弹窗 */
async function handleOpenDialog(id?: number) {
  dialog.visible = true;
  // 加载角色下拉数据源
  roleOptions.value = await RoleAPI.getOptions();
  if (id) {
    dialog.title = "修改管理员";
    PlatformAPI.getFormData(id).then((data) => {
      Object.assign(formData, { ...data });
    });
  } else {
    dialog.title = "新增管理员";
  }
}

/** 关闭管理员弹窗 */
function handleCloseDialog() {
  dialog.visible = false;

  platformFormRef.value.resetFields();
  platformFormRef.value.clearValidate();

  formData.id = undefined;
  formData.status = 1;
}

/** 提交角色表单 */
function handleSubmit() {
  platformFormRef.value.validate((valid: any) => {
    if (valid) {
      loading.value = true;
      const platformId = formData.id;
      if (platformId) {
        PlatformAPI.update(platformId, formData)
          .then(() => {
            ElMessage.success("修改管理员成功");
            handleCloseDialog();
            handleResetQuery();
          })
          .finally(() => (loading.value = false));
      } else {
        PlatformAPI.add(formData)
          .then(() => {
            ElMessage.success("新增管理员成功");
            handleCloseDialog();
            handleResetQuery();
          })
          .finally(() => (loading.value = false));
      }
    }
  });
}

/** 删除角色 */
function handleDelete(id?: number) {
  const platformIds = [id || removeIds.value].join(",");
  if (!platformIds) {
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
      PlatformAPI.deleteByIds(platformIds)
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
