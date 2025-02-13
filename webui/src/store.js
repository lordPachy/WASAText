import { defineStore } from 'pinia'

const getDefaultUserInfo = () => ({
    id: "",
    username: "",
})

const getUserInfo = () => {
    const userInfo = localStorage.getItem('store');
    return userInfo ? JSON.parse(userInfo) : getDefaultUserInfo();
}

export const useIDStore = defineStore('store', {
  state: () => ({userInfo: {
    id: "",
    username: "",
}}),
  actions: {
    changeID(newid) {
      this.userInfo.id = newid;
    },

    changeUsername(newusername) {
        this.userInfo.username = newusername;
    },
  },
  persist: {
    enabled: true
  }
})