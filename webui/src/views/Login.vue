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
					this.games = response.data;
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
    <p>Message is: {{ message }}</p>
    <input v-model="message" placeholder="edit me">
    <ErrorMsg v-if="errormsg" :msg="errormsg" />
  </div>
</template>

<style>
</style>
