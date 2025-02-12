<script>
import { useIDStore } from '../store';
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			store: useIDStore(),
			userquery: "",
			users: []
	}
},
	mounted() {
		this.getUsers();
	},
	methods: {
		async accessConversations() {
			this.loading = true;
			this.errormsg = null;
			try {
				this.$router.push({name: 'conversations'});
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async accessSettings() {
			this.loading = true;
			this.errormsg = null;
			try {
				this.$router.push({name: 'settings'});
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async getUsers() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/users", {headers: {Authorization: this.store.userInfo.id}, params: {username: this.userquery}});
				this.users = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
	},
}
</script>

<template>
  <div>
    <div
      class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom"
    >
      <h1 class="h2">Welcome {{ store.userInfo.username }}!</h1>
      <div class="btn-toolbar mb-2 mb-md-0">
        <div class="btn-group me-2" />
      </div>
    </div>
    <p>
      Decide what to do now:
    </p>
    <div class="btn-group me-2">
      <button type="button" class="btn btn-sm btn-outline-secondary" @click.stop="accessConversations">
        Access conversations
      </button>
    </div>
    <div class="btn-group me-2">
      <button type="button" class="btn btn-sm btn-outline-secondary" @click.stop="accessSettings">
        Settings
      </button>
    </div>
    <div class="mt-5">
      <h5 class="h5">User list</h5>
      <input v-model="userquery" placeholder="Username">
      <button type="button" class="btn btn-sm btn-outline-secondary" @click="getUsers">
        Search
      </button>
    </div>
    <div>
      <ul>
        <li v-for="u in users" :key="u">
          <img v-if="u.propic != 'NULL'" :src="u.propic" class="image-fit">
          {{ u.username }} 
        </li>
      </ul>
    </div>
    <ErrorMsg v-if="errormsg" :msg="errormsg" />
  </div>
</template>

<style>
.image-fit{
  height: 7%;
  width: 7%;
  object-fit: cover;
}
</style>
