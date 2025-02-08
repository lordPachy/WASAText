<script>
import ConvBut from '../components/ConvBut.vue';
import StartConvBut from '../components/StartConvBut.vue';

export default {
	components:{
		ConvBut,
		StartConvBut
	},
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
			id: "",
			inputid: "",
			newUser: "",
			groupName: "",
			showConversations: false,
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
				} catch (e) {
					this.errormsg = e.toString();
				}
				this.loading = false;
			},
			async login() {
				this.loading = true;
				this.errormsg = null;
				try {
					let response = await this.$axios.post("/session", {name: this.inputid});
					this.id = response.data.identifier;
					this.showConversations = true;
				} catch (e) {
					this.errormsg = e.toString();
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
        <div class="btn-group me-2">
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="conversations">
            Conversations
          </button>
        </div>
      </div>
    </div>
    <p>Please login or create your user!</p>
    <p>Username is: {{ id }}</p>
    <input v-model="inputid" placeholder="edit me">

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
	<p>
		<ConvBut :token="id" :show="showConversations"/>
	</p>
	<p>
		<StartConvBut :token="id" :show="showConversations"/>
	</p>
    <ErrorMsg v-if="errormsg" :msg="errormsg" />
  </div>
</template>

<style>
</style>
