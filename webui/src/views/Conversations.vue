<script>

export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
			newUser: "",
			isGroup: false,
			newGroupName: "",
			username: "",
			showConversations: false,
			newGroupMembers: [],
			newGroupMembersReq: [],
			chats: null
		}
	},
	methods: {
			async getConvos() {
				this.loading = true;
				this.errormsg = null;
				try{
					let response = await this.$axios.get("/conversations", {headers: {Authorization: this.$router.id}});
					this.chats = response.data;
					this.showConversations = true;
				} catch (e) {
					this.errormsg = e.toString();
				}
				this.loading = false;
			},
			async createPrivateConvo() {
				this.loading = true;
				this.errormsg = null;
				try{
                	let response = await this.$axios.put("/conversations", {isgroup: this.isGroup, members: [{name: this.username}], groupname: this.groupName}, {headers: {Authorization: this.$router.id}});
				} catch (e) {
					this.errormsg = e.toString();
				}
				this.loading = false;
			},
			async createGroup() {
				this.loading = true;
				this.errormsg = null;
				try{
                	let response = await this.$axios.put("/conversations", {isgroup: true, members: this.newGroupMembersReq, groupname: this.newGroupName}, {headers: {Authorization: this.$router.id}});
					this.newGroupMembers = [];
					this.newGroupMembersReq = [];
					this.newGroupName = "";
				} catch (e) {
					this.errormsg = e.toString();
				}
				this.loading = false;
			},
			async addToGroupMembers() {
				this.loading = true;
				this.errormsg = null;
				try{
                	this.newGroupMembers.push(this.username);
					this.newGroupMembersReq.push({name: this.username});
					this.username = "";
				} catch (e) {
					this.errormsg = e.toString();
				}
				this.loading = false;
			},
			async resetGroupMembers() {
				this.loading = true;
				this.errormsg = null;
				try{
                	this.newGroupMembers = [];
					this.newGroupMembersObj = [];
					this.username = "";
					this.newGroupName = "";
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
      <h1 class="h2">Conversations</h1>
      <div class="btn-toolbar mb-2 mb-md-0" />
    </div>
    <div>
      <h6>Create a new conversation</h6>
      <div>
        <input id="checkbox" v-model="isGroup" type="checkbox">
        <label for="checkbox">Is is a group?</label>
      </div>
      <input v-model="username" placeholder="Username">
      <button v-if="!isGroup" type="button" class="btn btn-sm btn-outline-secondary" @click="createPrivateConvo">
        Create Conversation
      </button>

      <button v-if="isGroup" type="button" class="btn btn-sm btn-outline-secondary" @click="addToGroupMembers">
        Add to new group
      </button>
      <button v-if="isGroup" type="button" class="btn btn-sm btn-outline-secondary" @click="resetGroupMembers">
        Reset group
      </button>
      <div v-if="isGroup">
        <input v-model="newGroupName" placeholder="New group name">
        <button type="button" class="btn btn-sm btn-outline-secondary" @click="createGroup">
          Create new group
        </button>
      </div>
      <div v-if="isGroup">
        <p>Current group members are (beside you): {{ newGroupMembers }}</p>
      </div>
    </div>





    <div class="mt-5">
      <h6>Started conversations</h6>
      <button type="button" class="btn btn-sm btn-outline-secondary" @click="getConvos"> 
        Get conversations
      </button>
      <button type="button" class="btn btn-sm btn-outline-secondary" @click="showConversations = false"> 
        Hide conversations
      </button>
      <div>
        <p v-if="showConversations">Chats are: {{ chats }}</p>
      </div>
    </div>

	
    <ErrorMsg v-if="errormsg" :msg="errormsg" />
  </div>
</template>

<style>
</style>
