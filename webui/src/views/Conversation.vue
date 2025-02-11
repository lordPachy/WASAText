<script>

import { RouterLink } from 'vue-router';

export default {
	props: {conversationid: String},
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
			data: {},
			isGroup: false,
			newuser: "",
			message: "",
			photo: "",
			replyingto: -1,
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
					let response = await this.$axios.get("/conversations/" + this.conversationid, {headers: {Authorization: this.$router.id}});
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
					let response = await this.$axios.post("/conversations/" + this.conversationid, {content: this.message, replyingto: this.replyingto}, {headers: {Authorization: this.$router.id}});
					this.message = "";
					this.photo = "";
					this.replyingto = -1;
					this.refresh();
				} catch (e) {
					this.errormsg = e.toString();
				}
				this.loading = false;
			},
			
			async addMemberToGroup() {
				this.loading = true;
				this.errormsg = null;
				try{
                	let response = await this.$axios.put("/groups", {username: {name: this.newuser}, group: {id: parseInt(this.conversationid)}}, {headers: {Authorization: this.$router.id}});
				} catch (e) {
					this.errormsg = e.toString();
				}
				this.loading = false;
			},
			async leaveGroup() {
				this.loading = true;
				this.errormsg = null;
				try{
                	let response = await this.$axios.delete("/conversations/" + this.conversationid, {headers: {Authorization: this.$router.id}});
					this.$router.push({name: 'conversations'});
				} catch (e) {
					this.errormsg = e.toString();
				}
				this.loading = false;
			},
			async putComment(mess, reac) {
				this.loading = true;
				this.errormsg = null;
				try{
                	let response = await this.$axios.put("/conversations/" + this.conversationid + "/messages/" + mess.messageid, {reaction: reac}, {headers: {Authorization: this.$router.id}});
					this.refresh();
				} catch (e) {
					this.errormsg = e.toString();
				}
				this.loading = false;
			},
			hasOwnComment(mess) {
				let pres = false;
				mess.comments.forEach(c => {
					if (c.sender == this.$router.username){
						pres = true;
					}
				});

				return pres;
			},

			async deleteMyComments(mess) {
				this.loading = true;
				this.errormsg = null;
				try{
					for (let i = 0; i < mess.comments.length; i++){
						if (mess.comments[i].sender == this.$router.username){
							console.log(mess.comments[i].commentid)
							await this.$axios.delete("/conversations/" + this.conversationid + "/messages/" + mess.messageid + "/comments/" + mess.comments[i].commentid.toString(), {headers: {Authorization: this.$router.id}});
						}
					}

					this.refresh();
				} catch (e) {
					this.errormsg = e.toString();
				}
				this.loading = false;
			},

			repliedMessage(messid) {
				for (let i = 0; i < this.messages.length; i++){
					if (this.messages[i].messageid == messid){
						return "From " + this.messages[i].username + ": " + this.messages[i].content
					}
				}
			},
			async deleteMessage(mess) {
				this.loading = true;
				this.errormsg = null;
				try{
					await this.$axios.delete("/conversations/" + this.conversationid + "/messages/" + mess.messageid, {headers: {Authorization: this.$router.id}});

					this.refresh();
				} catch (e) {
					this.errormsg = e.toString();
				}
				this.loading = false;
			},

			async groupSettings() {
				this.loading = true;
				this.errormsg = null;
				try {
					this.$router.push({name: 'groupsettings'}, {params: {conversationid: this.conversationid}});
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
      v-if="!isGroup" 
      class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom"
    >
      <h1 class="h2">Private chat with {{ data.user.username }}</h1>
      <div class="btn-toolbar mb-2 mb-md-0" />
    </div>
    <div
      v-else 
      class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom"
    >
      <h1 class="h2">Group: {{ data.groupname }}</h1>
      <div class="btn-toolbar mb-2 mb-md-0" />
    </div>

    <h5 class="h5">Group options</h5>
    <div v-if="isGroup">
      <input v-model="newuser" placeholder="New group member">
      <button type="button" class="btn btn-sm btn-outline-secondary" @click.stop="addMemberToGroup">
        Add member to group
      </button>
    </div>
    <div v-if="isGroup">
      <button type="button" class="btn btn-sm btn-outline-secondary mt-3" @click.stop="groupSettings">
        Group Settings
      </button>
      <button type="button" class="btn btn-sm btn-outline-secondary mt-3" @click.stop="leaveGroup">
        Leave this group
      </button>
    </div>
    <h5 class="h5 mt-4">Messages</h5>
    <div>
      <ul>
        <li v-for="m in data.messages" :key="m" class="mb-4">
          <p>
            ({{ m.timestamp.slice(0, 10) + " " + m.timestamp.slice(11, 19) }}) {{ m.username }}: "{{ m.content }}"
            <button type="button" class="btn btn-sm btn-outline-secondary" @click.stop="replyingto = m.messageid">
              Reply
            </button>
            <button v-if="m.username == $router.username" type="button" class="btn btn-sm btn-outline-secondary" @click.stop="deleteMessage(m)">
              Delete message
            </button>
          </p> 
          <div v-if="m.username === $router.username" class="mb-3"><span />Checkmarks: {{ m.checkmarks }}</div>
          <div>
            <p v-if="m.replyingto != -1">
              <span />This message is replying to:<br>
              <span />{{ repliedMessage(m.replyingto) }}
            </p>
          </div> 
          <div v-if="m.username != $router.username && !hasOwnComment(m)">
            <span />Put a comment:
          </div>
          <div v-if="m.username != $router.username && hasOwnComment(m)">
            <span /><button type="button" class="btn btn-sm btn-outline-secondary" @click.stop="deleteMyComments(m)">
              Delete my comment
            </button>
          </div>
          <div v-if="m.username != $router.username && !hasOwnComment(m)">
            <span />
            <select v-model="selected" class="mb-3" @click.stop="putComment(m, selected)">
              <option disabled value="None">Please select one</option>
              <option>laugh</option>
              <option>sad</option>
              <option>thumbs_up</option>
              <option>surprised</option>
              <option>love</option>
              <option>pray</option>
            </select>
          </div>
          <span />Comments:
          <p v-for="r in m.comments" :key="r"><span />{{ r.reaction }} by {{ r.sender }} </p>
        </li>
      </ul>
    </div>
	
    <p v-if="replyingto != -1">
      Currently replying to the following message:<br>
      {{ repliedMessage(replyingto) }}
      <button type="button" class="btn btn-sm btn-outline-secondary" @click.stop="replyingto = -1">
        Abort reply
      </button>
    </p>
    
    <textarea v-model="message" class="bottom" placeholder="New message" />
    <div class="btn-group me-2">
      <button type="button" class="btn btn-sm btn-outline-secondary" @click.stop="sendMessage">
        Send
      </button>
    </div>
  </div>
</template>

<style>
span { margin-left:11.8em }
</style>
