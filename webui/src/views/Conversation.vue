<script>

import { RouterLink } from 'vue-router';

export default {
	props: {cid: String},
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
			data: {},
			isGroup: false,
			message: "",
			photo: "",
		}
	},
	created() {
		this.refresh();
	},
	methods: {
			async refresh() {
				this.loading = true;
				this.errormsg = null;
				try{
					let response = await this.$axios.get("/conversations/" + this.cid, {headers: {Authorization: this.$router.id}});
					this.data = response.data;
					this.messages = response.data.messages;
					this.isGroup = Object.keys(response.data).length > 3;
				} catch (e) {
					this.errormsg = e.toString();
				}
				this.loading = false;
			},
			async sendMessage() {
				this.loading = true;
				this.errormsg = null;
				try{
					let response = await this.$axios.post("/conversations/" + this.cid, {content: this.message, replyingto: 0}, {headers: {Authorization: this.$router.id}});
					this.message = "";
					this.photo = "";
					this.refresh();
				} catch (e) {
					this.errormsg = e.toString();
				}
				this.loading = false;
			},
	}
}
</script>

<template>
  <div v-if="!isGroup">
    <div
      class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom"
    >
      <h1 class="h2">Private chat with {{ data.user.username }}</h1>
      <div class="btn-toolbar mb-2 mb-md-0" />
    </div>

    <div>
      <ul>
        <li v-for="m in data.messages" :key="m">
          ({{ m.timestamp.slice(0, 10) + " " + m.timestamp.slice(11, 19) }}) {{ m.username }}: "{{ m.content }}" 
        </li>
      </ul>
    </div>
	
    <textarea v-model="message" class="bottom" placeholder="New message" />
    <div class="btn-group me-2">
      <button type="button" class="btn btn-sm btn-outline-secondary" @click.stop="sendMessage">
        Send
      </button>
    </div>
    <ErrorMsg v-if="errormsg" :msg="errormsg" />
  </div>
</template>

<style>
</style>
