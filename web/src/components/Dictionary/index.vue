<template>
  <el-select
    v-model="selectedValue"
    :placeholder="placeholder"
    :disabled="disabled"
    clearable
    @change="handleChange"
  >
    <el-option
      v-for="option in options"
      :key="option.value"
      :label="option.label"
      :value="option.value"
    />
  </el-select>
</template>

<script setup lang="ts">

const props = defineProps({
  /**
   * 字典编码(eg: 性别-gender)
   */
  options: {
    type: Array as PropType<OptionType[]>,
    default: () => [],
    required: true,
  },
  modelValue: {
    type: [String, Number],
  },
  placeholder: {
    type: String,
    default: "请选择",
  },
  disabled: {
    type: Boolean,
    default: false,
  },
});

const emits = defineEmits(["update:modelValue"]);
const selectedValue = ref<string | number | undefined>();

watch(
  () => [props.modelValue],
  (newModelValue) => {
    if (newModelValue[0] == undefined) {
      selectedValue.value = undefined;
      return;
    }
    if (typeof props.options[0].value === "number") {
      selectedValue.value = Number(newModelValue);
    } else if (typeof props.options[0].value === "string") {
      selectedValue.value = String(newModelValue);
    } else {
      selectedValue.value = String(newModelValue);
    }
  }
);

function handleChange(val?: string | number | undefined) {
  emits("update:modelValue", val);
}
</script>
