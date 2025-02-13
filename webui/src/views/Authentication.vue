<script>
import { useIDStore } from '../store';
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			store: useIDStore(),
			id: "",
			inputid: "",
			newUser: "",
		}
	},
	methods: {
			async refresh() {
				this.loading = true;
				this.errormsg = null;
				try {
					let response = await this.$axios.get("/");
				} catch (e) {
					this.errormsg = e.toString();
				}
				this.loading = false;
			},
			async createUser() {
				this.loading = true;
				this.errormsg = null;
				try {
					let response = await this.$axios.put("/session", {name: this.inputid});
					this.id = response.data.identifier;
					this.store.changeID(this.id);
					this.store.changeUsername(this.inputid);
					this.$router.push({name: 'homepage'});
				} catch (e) {
					if (e.toString() == "AxiosError: Request failed with status code 403"){
						this.errormsg = "Username already in use";
					} else if (e.toString() == "AxiosError: Request failed with status code 400") {
						this.errormsg = "Usernames must be between 3 and 16 alphanumeric characters; no spaces";
					} else {
						this.errormsg = e.toString();
					}
				}
				this.loading = false;
			},
			async login() {
				this.loading = true;
				this.errormsg = null;
				try {
					let response = await this.$axios.post("/session", {name: this.inputid});
					this.id = response.data.identifier;
					this.store.changeID(this.id);
					this.store.changeUsername(this.inputid);
					this.$router.push({name: 'homepage'});
					this.showConversations = true;
				} catch (e) {
					if (e.toString() == "AxiosError: Request failed with status code 404"){
						this.errormsg = "User not found";
					} else {
						this.errormsg = e.toString();
					}
				}
				this.loading = false;
			},
	}
}
</script>

<template>
  <div>
    <div
      class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom"
    >
      <h1 class="h2">Login</h1>
      <div class="btn-toolbar mb-2 mb-md-0">
        <div class="btn-group me-2" />
      </div>
    </div>
    <p>Please login or create your user!</p>
    <p>Username is: {{ store.userInfo.username }}</p>
    <input v-model="inputid" placeholder="Insert here">

    <div class="btn-group me-2">
      <button type="button" class="btn btn-sm btn-outline-secondary" @click.stop="createUser">
        Create User
      </button>
    </div>

    <div class="btn-group me-2">
      <button type="button" class="btn btn-sm btn-outline-secondary" @click.stop="login">
        Login
      </button>
    </div>
    <ErrorMsg v-if="errormsg" :msg="errormsg" />
  </div>
</template>

<style>
</style>
