<script>

import { RouterLink } from 'vue-router';
import { useIDStore } from '../store';

export default {
	data: function() {
		return {
			errormsg: null,
			store: useIDStore(),
			isGroup: false,
			newGroupName: "",
			username: "",
			creatingChat: false,
			availableUsers: [],
			newGroupMembers: [],
			newGroupMembersReq: [],
			chats: null,
			timer: ''
		}
	},

	mounted() {
		this.getConvos();

		// Updating page every 2000 ms
		this.timer = setInterval(this.getConvos, 2000);
	},

	unmounted() {
		// Avoiding page update when it is closed
		clearInterval(this.timer);
	},

	methods: {
		/**
		 * It retrieves the user's started conversations.
		 */
		async getConvos() {
			this.errormsg = null;
			try{
				let response = await this.$axios.get("/conversations", {headers: {Authorization: this.store.userInfo.id}});
				this.chats = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
		},

		/**
		 * It creates a private conversation with a single user.
		*/
		async createPrivateConvo() {
			this.errormsg = null;
			try{
				let response = await this.$axios.put("/conversations", {isgroup: this.isGroup, members: [{name: this.username}], groupname: this.groupName}, {headers: {Authorization: this.store.userInfo.id}});

				// Updating
				this.getConvos();
				this.resetGroupMembers();
				this.creatingChat = false;
			} catch (e) {
				// Capturing error
				if (e.toString() == "AxiosError: Request failed with status code 400" || e.toString() == "AxiosError: Request failed with status code 404") {
					this.errormsg = "User not found";
				} else {
					this.errormsg = e.toString();
				}

				// Ensuring the message can be seen
				await new Promise(resolve => setTimeout(resolve, 7000));
				this.errormsg = null;
			}
		},

		/**
		 * It creates a group chat.
		 */
		async createGroup() {
			this.errormsg = null;
			try{
				let response = await this.$axios.put("/conversations", {isgroup: true, members: this.newGroupMembersReq, groupname: this.newGroupName}, {headers: {Authorization: this.store.userInfo.id}});
				
				// Reinitializing variables
				this.newGroupMembers = [];
				this.newGroupMembersReq = [];
				this.newGroupName = "";

				// Updating
				this.getConvos();
				this.resetGroupMembers();
				this.creatingChat = false;
			} catch (e) {
				if (e.toString() == "AxiosError: Request failed with status code 400") {
					this.errormsg = "Error: user(s) might not exist, or groupname is not valid (it must be between 3 and 16 alphanumeric characters; no spaces)";
				} else {
					this.errormsg = e.toString();
				}

				// Ensuring the message can be seen
				await new Promise(resolve => setTimeout(resolve, 7000));
				this.errormsg = null;
			}
		},

		/**
		 * It adds "username" to newGroupMembers and newGroupMembersReq.
		 */
		addToGroupMembers() {
			this.newGroupMembers.push(this.username);
			this.newGroupMembersReq.push({name: this.username});
			this.username = "";
		},

		/**
		 * It cleans all temporary structures for the creation
		 * of a new group.
		 */
		resetGroupMembers() {
			this.newGroupMembers = [];
			this.newGroupMembersObj = [];
			this.username = "";
			this.newGroupName = "";
		},

		/**
		 * It calculates users that can be selected either for a
		 * private chat or a group.
		 */
		async computeAvailableUsers() {
			this.errormsg = null;
			try{
				this.availableUsers = [];

				// Retrieving WASAText active users
				let users = await this.$axios.get("/users", {headers: {Authorization: this.store.userInfo.id}, params: {username: ""}});
				users = users.data;

				// Retrieving chats in a moment in time
				let nowchats = await this.$axios.get("/conversations", {headers: {Authorization: this.store.userInfo.id}});
				nowchats = nowchats.data;
				// Private chat case
				if (!this.isGroup){
					// Adding all users (which are not already in a private chat) to the list of available ones
					for (let i = 0; i < users.length; i++){
						if (nowchats != null){
							for (let j = 0; j <= nowchats.length; j++){
								if (j == nowchats.length){
									if (users[i].username != this.store.userInfo.username){
										this.availableUsers.push(users[i].username);
									}
									
									break;
								} else if (nowchats[j].chatid.id < 5000 && users[i].username == nowchats[j].name){
									break;
								}
							}
						} else {
							if (users[i].username != this.store.userInfo.username){
								this.availableUsers.push(users[i].username);
							}
						}
					}
				// Group chat case 
				} else {
					// Adding all users (which are not already selected for the group) to the list of available ones
					for (let i = 0; i < users.length; i++){
						for (let j = 0; j <= this.newGroupMembers.length; j++){
							if (j == this.newGroupMembers.length) {
								if (users[i].username != this.store.userInfo.username){
									this.availableUsers.push(users[i].username);
								}

								break;
							} else if (users[i].username == this.newGroupMembers[j]){
								break;
							}
						}
					}
				}
			} catch (e) {
				this.errormsg = e.toString();
			}
		},
	}
}
</script>

<template>
  <!--Header-->
  <div
    class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom"
  >
    <h1 class="h2">Conversations</h1>
  </div>

  <!--Error messages-->
  <div>
    <ErrorMsg v-if="errormsg" :msg="errormsg" />
  </div>

  <!--Conversation creation-->
  <div>
    <h6>Conversation creation</h6>
    <div v-if="!creatingChat">
      <button type="button" class="btn btn-sm btn-outline-secondary" @click="creatingChat = true; isGroup = false; computeAvailableUsers()">
        Private conversation
      </button>
      <button type="button" class="btn btn-sm btn-outline-secondary" @click="creatingChat = true; isGroup = true; computeAvailableUsers()">
        Groupchat
      </button>
    </div>
    <div v-else>
      <button type="button" class="btn btn-sm btn-outline-secondary" @click="creatingChat = false; resetGroupMembers()">
        Discard
      </button>
    </div>

    <br>

    <div v-if="creatingChat">
      <select v-model="username" class="mb-3">
        <option disabled value="">Please select...</option>
        <option v-for="c in availableUsers" :key="c">{{ c }}</option>
      </select>

      <!--Private conversation creation-->
      <button v-if="!isGroup" type="button" class="btn btn-sm btn-outline-secondary" @click="createPrivateConvo">
        Create Conversation
      </button>
	
      <!--Group conversation creation-->
      <button v-if="isGroup" type="button" class="btn btn-sm btn-outline-secondary" @click="addToGroupMembers(); computeAvailableUsers();">
        Add to new group
      </button>
      <button v-if="isGroup" type="button" class="btn btn-sm btn-outline-secondary" @click="resetGroupMembers(); computeAvailableUsers();">
        Reset group
      </button>

      <div v-if="isGroup">
        <br> Note that group names are between 3 and 16 alphanumeric characters long.
        <br>
        <input v-model="newGroupName" placeholder="New group name">
        <button type="button" :disabled="newGroupName.length < 3 || newGroupName.length > 16 || newGroupMembers.length < 1" class="btn btn-sm btn-outline-secondary" @click="createGroup">
          Create new group
        </button>
      </div>
      <div v-if="isGroup">
        <p>Current group members are (beside you): {{ newGroupMembers }}</p>
      </div>
    </div>
  </div>

  <!--Started conversations-->
  <div class="mt-5">
    <h6>Started conversations</h6>

    <ul>
      <li v-for="f in chats" :key="f">
        <img v-if="f.photo != 'NULL'" :src="f.photo" class="image-fit"><span v-if="f.photo != 'NULL'" class="smallspan" />
        {{ f.name }}: 
        <img v-if="f.lastmessage.photo != 'NULL' && f.lastmessage.photo != ''" :src="f.lastmessage.photo" class="image-min">
        "{{ f.lastmessage.content }}" ({{ f.lastmessage.timestamp.slice(0, 10) + " " + f.lastmessage.timestamp.slice(11, 19) }})

        <RouterLink :to="/conversations/+f.chatid.id" class="nav-link">
          <button type="button" class="btn btn-sm btn-outline-secondary"> 
            Open
          </button>
        </RouterLink>
        <br><br>
      </li>
    </ul>
  </div>
</template>

