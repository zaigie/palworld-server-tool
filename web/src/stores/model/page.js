import { defineStore } from "pinia";
import { ref } from "vue";

const pageStore = defineStore(
  "page",
  () => {
    let screenWidth = ref(0);
    // Set
    const setScreenWidth = (info) => {
      screenWidth.value = info;
    };
    // Get
    const getScreenWidth = () => {
      return screenWidth.value;
    };

    return {
      setScreenWidth,
      getScreenWidth,
    };
  },
  {
    persist: true,
  }
);

export default pageStore;
