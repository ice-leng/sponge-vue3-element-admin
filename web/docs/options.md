# options 使用方法

## 加载 options
```js
  ...
  import xxxAPI from "@/api/xxx";

  const xxxxOptions = ref<OptionType[]>([]);
  async function options() {
    xxxxOptions.value = await xxxAPI.getOptions();
  }
  options();

```

## 列表使用
```vue

  import { showDictLabel } from "@/utils/index";
  ....
  <el-table-column
    key="xxxx"
    label="xxxx"
    align="center"
    prop="xxxx"
    width="180"
  >
    <template #default="scope">
        {{ showDictLabel(xxxxOptions, scope.row.xxxx) }}
    </template>
  </el-table-column>

```

## 表单使用

```vue
  ....
    <dictionary v-model="formData.xxx" :options="xxxxOptions" />
```
