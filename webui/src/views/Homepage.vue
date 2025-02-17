<script>
import { useIDStore } from '../store';
export default {
	data: function() {
		return {
			errormsg: null,
			store: useIDStore(),
			userquery: "",
			users: [],
			timer: '',
			tmpquery: ''
		}
	},

	mounted() {
		this.getUsers();

		// Updating page every 2000 ms
		this.timer = setInterval(this.getUsers, 2000);
	},

	unmounted() {
		// Avoiding page update when it is closed
		clearInterval(this.timer);
	},

	methods: {
		/**
		 * It redirects to the conversations page.
		 */
		async accessConversations() {
			this.errormsg = null;
			try {
				this.$router.push({name: 'conversations'});
			} catch (e) {
				this.errormsg = e.toString();
			}
		},

		/**
		 * It redirects to the user settings page.
		 */
		async accessSettings() {
			this.errormsg = null;
			try {
				this.$router.push({name: 'settings'});
			} catch (e) {
				this.errormsg = e.toString();
			}
		},

		/**
		 * It retrieves the current WASAText users.
		 */
		async getUsers() {
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/users", {headers: {Authorization: this.store.userInfo.id}, params: {username: this.userquery}});
				this.users = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
		},
	},
}
</script>

<template>
  <!--Header-->
  <div
    class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom"
  >
    <h1 class="h2">Welcome {{ store.userInfo.username }}!</h1>
  </div>

  <!--Error messages-->
  <div>
    <ErrorMsg v-if="errormsg" :msg="errormsg" />
  </div>

  <!--User pages-->
  <div>
    <p>
      Decide what to do now:
    </p>

    <div class="btn-group me-2">
      <button type="button" class="btn btn-sm btn-outline-secondary" @click="accessConversations">
        Access conversations
      </button>
    </div>

    <div class="btn-group me-2">
      <button type="button" class="btn btn-sm btn-outline-secondary" @click="accessSettings">
        Settings
      </button>
    </div>
  </div>
  
  <!--User list-->
  <div class="mt-5">
    <h5 class="h5">User list</h5>

    <input v-model="tmpquery" placeholder="Username">
    <button type="button" class="btn btn-sm btn-outline-secondary" @click="userquery = tmpquery; getUsers()">
      Search
    </button>

    <ul>
      <li v-for="u in users" :key="u">
        <img v-if="u.propic != 'NULL'" :src="u.propic" class="image-fit">
        {{ u.username }} 
      </li>
    </ul>
  </div>
</template>

<style>
.image-fit{
  height: 7%;
  width: 7%;
  object-fit: cover;
}
</style>
