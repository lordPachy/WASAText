<script>
import { useIDStore } from '../store';
export default {
	data: function() {
		return {
			errormsg: null,
			store: useIDStore(),
			input: "",
		}
	},

	mounted () {
		// Clearing up cached data when authenticating
		this.store.changeID("");
		this.store.changeUsername("");
	},

	methods: {
		/**
		 * It creates a user (if the username is valid) and
		 * 
		 * redirects to the homepage.
		 */
		async createUser() {
			this.errormsg = null;
			try {
				let response = await this.$axios.put("/session", {name: this.input});
				
				// Updating stored data
				this.store.changeID(response.data.identifier);
				this.store.changeUsername(this.input);

				// Redirecting
				this.$router.push({name: 'homepage'});
			} catch (e) {
				if (e.toString() == "AxiosError: Request failed with status code 403"){
					this.errormsg = "Username already in use";
				} else if (e.toString() == "AxiosError: Request failed with status code 400") {
					this.errormsg = "Usernames must be between 3 and 16 alphanumeric characters; no spaces";
				} else {
					this.errormsg = e.toString();
				}

				await new Promise(resolve => setTimeout(resolve, 7000));
				this.errormsg = null;
			}
		},

		/**
		 * It logins a user (if they exist) and
		 * 
		 * redirects to the homepage.
		 */
		async login() {
			this.errormsg = null;
			try {
				let response = await this.$axios.post("/session", {name: this.input});

				// Updating stored data
				this.store.changeID(response.data.identifier);
				this.store.changeUsername(this.input);

				// Redirecting
				this.$router.push({name: 'homepage'});
			} catch (e) {
				if (e.toString() == "AxiosError: Request failed with status code 404"){
					this.errormsg = "User not found";
				} else {
					this.errormsg = e.toString();
				}

				await new Promise(resolve => setTimeout(resolve, 7000));
				this.errormsg = null;
			}
		}
	}
}
</script>

<template>
  <!--Header-->
  <div
    class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom"
  >
    <h1 class="h2">Login</h1>
  </div>

  <!--Error messages-->
  <ErrorMsg v-if="errormsg" :msg="errormsg" />

  <!--User input and login/creation-->
  <div>
    <p>Please login or create your user!<br> Note that usernames are between 3 and 16 alphanumeric characters long, with no blank spaces.</p>
    <p>Username is:</p>
    <input v-model="input" placeholder="Insert here">

    <div class="btn-group me-2">
      <button type="button" :disabled="input.length < 3 || input.length > 16" class="btn btn-sm btn-outline-secondary" @click="createUser">
        Create User
      </button>
    </div>

    <div class="btn-group me-2">
      <button type="button" :disabled="input.length < 3 || input.length > 16" class="btn btn-sm btn-outline-secondary" @click="login">
        Login
      </button>
    </div>
  </div>
</template>

