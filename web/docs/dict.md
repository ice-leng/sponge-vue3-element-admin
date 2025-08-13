# dict 使用方法

## 列表使用
```vue

<el-table-column label="状态" width="120">
  <template #default="scope">
    <DictLabel v-model="scope.row.status" code="order_status" />
  </template>
</el-table-column>

```

## 表单使用

```vue
  ....
    <Dict v-model="queryParams.status" code="order_status" />
```
