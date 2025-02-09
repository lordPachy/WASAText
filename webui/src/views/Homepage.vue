<script>

export default {
	props: {token: String, username: String},
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
			newUser: "",
			groupName: "",
			showConversations: false,
		}
	},
	methods: {
		async accessConversations() {
			this.loading = true;
			this.errormsg = null;
			try {
				this.$router.push({name: 'conversations', params: {token:this.token, username:this.username}});
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async accessSettings() {
			this.loading = true;
			this.errormsg = null;
			try {
				this.$router.push({name: 'settings', params: {token:this.token, username:this.username}});
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
      <h1 class="h2">Welcome, {{ username }}!</h1>
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
    <ErrorMsg v-if="errormsg" :msg="errormsg" />
  </div>
</template>

<style>
</style>
