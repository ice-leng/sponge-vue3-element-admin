# 图片 使用方法

## 列表使用
```vue
  ....
  <el-table-column
    key="imageUrl"
    label="图片"
    align="center"
    prop="imageUrl"
    width="180"
  >
    <template #default="scope">
      <el-image
        :src="scope.row.imageUrl"
        :preview-src-list="[scope.row.imageUrl]"
        :preview-teleported="true"
        :style="`width: 40px;`"
      />
    </template>
  </el-table-column>

```

## 表单使用

```vue
  ....

  <el-form-item label="封面图" prop="cover">
    <SingleImageUpload v-model="formData.cover" />
  </el-form-item>

```
