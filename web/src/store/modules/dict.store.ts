import { store } from "@/store";
import ConfigAPI from "@/api/system/config.api";

export const useDictStore = defineStore("dict", () => {
  // 字典数据缓存
  const dictCache = useStorage<Record<string, OptionType[]>>("dict_cache", {});

  /**
   * 缓存字典数据
   * @param dictCode 字典编码
   * @param data 字典项列表
   */
  const cacheDictItems = (dictCode: string, data: OptionType[]) => {
    dictCache.value[dictCode] = data;
  };

  /**
   * 加载字典数据（如果缓存中没有则请求）
   */
  const loadDictItems = async () => {
    await ConfigAPI.dict().then((data) => {
      for (let dictCode in data) {
        cacheDictItems(dictCode, data[dictCode]);
      }
    });
  };

  /**
   * 获取字典项列表
   * @param dictCode 字典编码
   * @returns 字典项列表
   */
  const getDictItems = (dictCode: string): OptionType[] => {
    return dictCache.value[dictCode] || [];
  };

  /**
   * 清空字典缓存
   */
  const clearDictCache = () => {
    dictCache.value = {};
  };

  return {
    loadDictItems,
    getDictItems,
    clearDictCache,
  };
});

export function useDictStoreHook() {
  return useDictStore(store);
}
