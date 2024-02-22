import { defineStore } from "pinia";
import { ref } from "vue";

const playerToGuildStore = defineStore(
  "playerToGuild",
  () => {
    const currentUid = ref(null);
    const updateStatus = ref("players");
    // Set
    const setCurrentUid = (uid) => {
      currentUid.value = uid;
    };
    const setUpdateStatus = (status) => {
      updateStatus.value = status;
    };
    // Get
    const getCurrentUid = () => {
      return currentUid.value;
    };
    const getUpdateStatus = () => {
      return updateStatus.value;
    };

    return {
      setCurrentUid,
      setUpdateStatus,
      getCurrentUid,
      getUpdateStatus,
    };
  },
  {
    persist: true,
  }
);

export default playerToGuildStore;
