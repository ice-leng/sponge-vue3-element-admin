<!-- 用户管理 -->
<template>
  <div class="app-container">
    <!-- 用户列表 -->
    <div class="search-bar">
      <el-form ref="queryFormRef" :model="queryParams" :inline="true">
        <el-form-item label="关键词" prop="keywords">
          <el-input
            v-model="queryParams.keyword"
            placeholder="用户名/昵称"
            clearable
            style="width: 200px"
            @keyup.enter="handleQuery"
          />
        </el-form-item>

        <el-form-item label="手机号" prop="mobile">
          <el-input
            v-model="queryParams.mobile"
            placeholder="手机号"
            clearable
            style="width: 200px"
            @keyup.enter="handleQuery"
          />
        </el-form-item>

        <el-form-item label="状态" prop="status">
          <el-select v-model="queryParams.status" placeholder="全部" clearable class="!w-[100px]">
            <el-option label="正常" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>

        <el-form-item label="创建时间">
          <el-date-picker
            v-model="dateTimeRange"
            :editable="false"
            class="!w-[240px]"
            type="daterange"
            range-separator="~"
            start-placeholder="开始时间"
            end-placeholder="截止时间"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleQuery">
            <template #icon><Search /></template>
            搜索
          </el-button>
          <el-button @click="handleResetQuery">
            <template #icon><Refresh /></template>
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-card shadow="never">
      <div class="flex-x-between mb-10px">
        <div>
          <el-button v-hasPerm="['sys:platform:edit']" type="success" @click="handleOpenDialog()">
            <template #icon><Plus /></template>
            新增
          </el-button>
          <el-button
            v-hasPerm="['sys:platform:delete']"
            type="danger"
            :disabled="removeIds.length === 0"
            @click="handleDelete()"
          >
            <template #icon><Delete /></template>
            删除
          </el-button>
        </div>
      </div>

      <el-table v-loading="loading" :data="pageData" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="50" align="center" />
        <el-table-column key="id" label="编号" align="center" prop="id" width="100" />

        <el-table-column key="username" label="用户名" align="center" prop="username" />

        <el-table-column key="nickname" label="昵称" align="center" prop="nickname" />

        <el-table-column key="mobile" label="手机号" align="center" prop="mobile" />

        <el-table-column key="gender" label="性别" align="center" prop="gender">
          <template #default="scope">
            <DictLabel v-model="scope.row.gender" code="gender" />
          </template>
        </el-table-column>

        <el-table-column key="roleNames" label="角色" align="center" prop="roleNames" />

        <el-table-column label="状态" align="center" prop="status">
          <template #default="scope">
            <el-tag :type="scope.row.status == 1 ? 'success' : 'info'">
              {{ scope.row.status == 1 ? "正常" : "禁用" }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="最新登录时间" align="center" prop="lastTime" width="180" />

        <el-table-column label="创建时间" align="center" prop="createdAt" width="180" />
        <el-table-column label="操作" fixed="right" width="220">
          <template #default="scope">
            <el-button
              v-hasPerm="['sys:platform:password:reset']"
              type="primary"
              size="small"
              link
              @click="handleResetPassword(scope.row)"
            >
              <template #icon><RefreshLeft /></template>
              重置密码
            </el-button>
            <el-button
              v-hasPerm="['sys:platform:edit']"
              type="primary"
              link
              size="small"
              @click="handleOpenDialog(scope.row.id)"
            >
              <template #icon><Edit /></template>
              编辑
            </el-button>
            <el-button
              v-hasPerm="['sys:platform:delete']"
              type="danger"
              link
              size="small"
              @click="handleDelete(scope.row.id)"
            >
              <template #icon><Delete /></template>
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

    <!-- 用户表单弹窗 -->
    <el-dialog
      v-model="dialog.visible"
      :title="dialog.title"
      width="500px"
      @close="handleCloseDialog"
    >
      <el-form ref="platformFormRef" :model="formData" :rules="rules" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input
            v-model="formData.username"
            :readonly="!!formData.id"
            placeholder="请输入用户名"
          />
        </el-form-item>

        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="formData.nickname" placeholder="请输入昵称" />
        </el-form-item>

        <el-form-item label="手机号" prop="mobile">
          <el-input v-model="formData.mobile" placeholder="手机号" />
        </el-form-item>

        <el-form-item label="性别" prop="gender">
          <Dict v-model="formData.gender" code="gender" />
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
          <el-switch
            v-model="formData.status"
            inline-prompt
            active-text="正常"
            inactive-text="禁用"
            :active-value="1"
            :inactive-value="0"
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
  name: "Platform",
  inheritAttrs: false,
});

import PlatformAPI, {
  PlatformPageQuery,
  PlatformPageVO,
  PlatformForm,
} from "@/api/system/platform.api";

import RoleAPI from "@/api/system/role.api";

const queryFormRef = ref(ElForm);
const platformFormRef = ref(ElForm);

const loading = ref(false);
const removeIds = ref([]);
const total = ref(0);
const pageData = ref<PlatformPageVO[]>();
// 角色下拉数据源
const roleOptions = ref<OptionType[]>();
// 用户查询参数
const queryParams = reactive<PlatformPageQuery>({
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

// 用户弹窗
const dialog = reactive({
  visible: false,
  title: "",
});

// 用户表单数据
const formData = reactive<PlatformForm>({
  status: 1,
});

// 用户表单校验规则
const rules = reactive({
  username: [{ required: true, message: "用户名不能为空", trigger: "blur" }],
  roleId: [{ required: true, message: "用户角色不能为空", trigger: "blur" }],
  mobile: [
    { message: "请输入手机号", trigger: "blur" },
    {
      pattern: /^1[3|4|5|6|7|8|9][0-9]\d{8}$/,
      message: "请输入正确的手机号码",
      trigger: "blur",
    },
  ],
});

//查询
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

// 重置查询
function handleResetQuery() {
  queryFormRef.value.resetFields();
  queryParams.page = 1;
  dateTimeRange.value = "";
  queryParams.startTime = undefined;
  queryParams.endTime = undefined;
  handleQuery();
}

// 行复选框选中记录选中ID集合
function handleSelectionChange(selection: any) {
  removeIds.value = selection.map((item: any) => item.id);
}

// 重置密码
function handleResetPassword(row: { [key: string]: any }) {
  ElMessageBox.prompt("请输入用户「" + row.username + "」的新密码", "重置密码", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
  }).then(
    ({ value }) => {
      if (!value || value.length < 6) {
        // 检查密码是否为空或少于6位
        ElMessage.warning("密码至少需要6位字符，请重新输入");
        return false;
      }
      PlatformAPI.resetPassword(row.id, value).then(() => {
        ElMessage.success("密码重置成功，新密码是：" + value);
      });
    },
    () => {
      ElMessage.info("已取消重置密码");
    }
  );
}

/**
 * 打开弹窗
 *
 * @param id 用户ID
 */
async function handleOpenDialog(id?: number) {
  dialog.visible = true;
  // 加载角色下拉数据源
  roleOptions.value = await RoleAPI.getOptions();

  if (id) {
    dialog.title = "修改用户";
    PlatformAPI.getFormData(id).then((data) => {
      Object.assign(formData, { ...data });
    });
  } else {
    dialog.title = "新增用户";
  }
}

// 关闭弹窗
function handleCloseDialog() {
  dialog.visible = false;
  platformFormRef.value.resetFields();
  platformFormRef.value.clearValidate();

  formData.id = undefined;
  formData.status = 1;
}

// 表单提交
const handleSubmit = useThrottleFn(() => {
  platformFormRef.value.validate((valid: any) => {
    if (valid) {
      const userId = formData.id;
      loading.value = true;
      if (userId) {
        PlatformAPI.update(userId, formData)
          .then(() => {
            ElMessage.success("修改用户成功");
            handleCloseDialog();
            handleResetQuery();
          })
          .finally(() => (loading.value = false));
      } else {
        PlatformAPI.add(formData)
          .then(() => {
            ElMessage.success("新增用户成功");
            handleCloseDialog();
            handleResetQuery();
          })
          .finally(() => (loading.value = false));
      }
    }
  });
}, 3000);

// 删除用户
function handleDelete(id?: number) {
  const userIds = [id || removeIds.value].join(",");
  if (!userIds) {
    ElMessage.warning("请勾选删除项");
    return;
  }

  ElMessageBox.confirm("确认删除用户?", "警告", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(
    function () {
      loading.value = true;
      PlatformAPI.deleteByIds(userIds)
        .then(() => {
          ElMessage.success("删除成功");
          handleResetQuery();
        })
        .finally(() => (loading.value = false));
    },
    function () {
      ElMessage.info("已取消删除");
    }
  );
}

onMounted(() => {
  handleQuery();
});
</script>
