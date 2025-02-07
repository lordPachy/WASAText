<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/");
				this.some_data = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async login() {
			try {
				let response = await this.$axios.put("/session", {
					params: {
					// These values should come from the browser
					// location/GPS or some user input.
					name: "Paolo",
					}
				});
				// Axios decodes JSON automatically
				this.token = response.data;
				} catch (e) {
				alert("Error: " + e);
				}
		}
	}
}
</script>

<template>
	<!-- when this button is clicked, the list is downloaded -->
	<button @click="CreateUser">
		CreateUser
	</button>
	<!-- separator -->
	<hr />
	<p>User token:</p>
	<ul>
		Token: {{ Token }}
	</ul>
</template>
<style>
</style>
