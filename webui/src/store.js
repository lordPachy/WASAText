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
  state: () => ({userInfo: getUserInfo()}),
  // could also be defined as
  // state: () => ({ count: 0 })
  actions: {
    changeID(newid) {
      this.userInfo.id = newid;
      localStorage.setItem('store', JSON.stringify(this.userInfo));
    },

    changeUsername(newusername) {
        this.userInfo.username = newusername;
        localStorage.setItem('store', JSON.stringify(this.userInfo));
    },
  },
  persist: true
})