import { defineStore } from "pinia";
import { ref } from "vue";

const userStore = defineStore(
  "user",
  () => {
    const isLogin = ref(false);
    const token = ref(null);

    const setIsLogin = (status, key) => {
      if (status) {
        isLogin.value = true;
        token.value = key;
      } else {
        isLogin.value = false;
        token.value = null;
      }
    };
    const getLoginInfo = () => {
      return {
        isLogin: isLogin.value,
        token: token,
      };
    };

    return {
      setIsLogin,
      getLoginInfo,
    };
  },
  {
    persist: true,
  }
);

export default userStore;
