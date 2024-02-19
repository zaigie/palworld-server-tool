import { defineStore } from "pinia";
import { ref } from "vue";

const whitelistStore = defineStore(
  "whitelist",
  () => {
    const whitelist = ref([]);
    // Set
    const setWhitelist = (list) => {
      whitelist.value = list;
    };
    // Get
    const getWhitelist = () => {
      return whitelist.value;
    };

    return {
      setWhitelist,
      getWhitelist,
    };
  },
  {
    persist: true,
  }
);

export default whitelistStore;
