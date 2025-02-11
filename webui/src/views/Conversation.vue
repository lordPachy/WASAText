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
			sendingpic: false,
			photo: "NULL",
			forwarding: false,
			forwardingto: -1,
			replyingto: -1,
			chats: []
		}
	},
	created() {
		this.refresh();
		this.getConvos();
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
					let response = await this.$axios.post("/conversations/" + this.conversationid, {content: this.message, photo: this.photo, replyingto: this.replyingto}, {headers: {Authorization: this.$router.id}});
					this.message = "";
					this.photo = "NULL";
					this.replyingto = -1;
					this.sendingpic = false;
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
			repliedPhoto(messid) {
				for (let i = 0; i < this.messages.length; i++){
					if (this.messages[i].messageid == messid){
						return this.messages[i].photo
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
			async getConvos() {
				this.loading = true;
				this.errormsg = null;
				try{
					let response = await this.$axios.get("/conversations", {headers: {Authorization: this.$router.id}});
					this.chats = response.data;
				} catch (e) {
					this.errormsg = e.toString();
				}
				this.loading = false;
			},
			async forwardMessage(mess) {
				this.loading = true;
				this.errormsg = null;
				try{
					let response = await this.$axios.post("/conversations/" + this.conversationid + "/messages/" + mess.messageid, {id: this.forwardingto}, {headers: {Authorization: this.$router.id}});
					this.forwarding = false;
					this.forwardingto = -1;
					this.refresh();
				} catch (e) {
					this.errormsg = e.toString();
				}
				this.loading = false;
			},
			async uploadImage(a) {
				this.loading = true;
				this.errormsg = null;
				try {
					const image = a.target.files[0];
					const reader = new FileReader();
					reader.readAsDataURL(image);
					reader.onload = a =>{
						this.photo = a.target.result;
					};
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
      <h1 class="h2">
        Private chat with
        <img v-if="data.user.propic != 'NULL'" :src="data.user.propic" class="image-big">
        {{ data.user.username }}
      </h1>
      <div class="btn-toolbar mb-2 mb-md-0" />
    </div>
    <div
      v-else 
      class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom"
    >
      <h1 class="h2">
        Group:
        <img v-if="data.groupphoto != 'NULL'" :src="data.groupphoto" class="image-big">
        {{ data.groupname }}
      </h1>
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
            ({{ m.timestamp.slice(0, 10) + " " + m.timestamp.slice(11, 19) }}) {{ m.username }}{{ (m.og_sender != "NULL") ? (" (Originally written by " + m.og_sender + ")") : ("") }}: "{{ m.content }}"
            <button type="button" class="btn btn-sm btn-outline-secondary" @click.stop="replyingto = m.messageid">
              Reply
            </button>
            <button type="button" class="btn btn-sm btn-outline-secondary" @click.stop="forwarding = !forwarding">
              Forward
            </button>
            <select v-if="forwarding" v-model="forwardingto" class="mb-3" @click.stop="forwardMessage(m)">
              <option disabled value="-1">Please select one</option>
              <option v-for="c in chats" :key="c" :value="c.chatid.id">{{ c.name }}</option>
            </select>
            <button v-if="m.username == $router.username" type="button" class="btn btn-sm btn-outline-secondary" @click.stop="deleteMessage(m)">
              Delete message
            </button>
          </p>
          <div>
            <bigspan />
            <img v-if="m.photo != 'NULL' && m.photo != ''" :src="m.photo" class="image-big">
          </div>
          <div v-if="m.username === $router.username" class="mb-3"><bigspan />Checkmarks: {{ m.checkmarks }}</div>
          <div>
            <p v-if="m.replyingto != -1">
              <bigspan />This message is replying to:<br>
              <bigspan />{{ repliedMessage(m.replyingto) }}
              <img v-if="repliedPhoto(m.replyingto) != 'NULL' && repliedPhoto(m.replyingto) != ''" :src="repliedPhoto(m.replyingto)" class="image-min">
            </p>
          </div> 
          <div v-if="m.username != $router.username && !hasOwnComment(m)">
            <bigspan />Put a comment:
          </div>
          <div v-if="m.username != $router.username && hasOwnComment(m)">
            <bigspan /><button type="button" class="btn btn-sm btn-outline-secondary" @click.stop="deleteMyComments(m)">
              Delete my comment
            </button>
          </div>
          <div v-if="m.username != $router.username && !hasOwnComment(m)">
            <bigspan />
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
          <bigspan />Comments:
          <p v-for="r in m.comments" :key="r"><bigspan />{{ r.reaction }} by {{ r.sender }} </p>
        </li>
      </ul>
    </div>
	
    <div v-if="replyingto != -1">
      Currently replying to the following message:<br>
      {{ repliedMessage(replyingto) }}
      <div>
        <bigspan />
        <img v-if="repliedPhoto(replyingto) != 'NULL' && repliedPhoto(replyingto) != ''" :src="repliedPhoto(replyingto)" class="image-min">
      </div>
      <button type="button" class="btn btn-sm btn-outline-secondary" @click.stop="replyingto = -1">
        Abort reply
      </button>
    </div>
    
    <div v-if="sendingpic">
      <p class="mt-5">Select a picture:</p>
      <input type="file" accept="image/png" @change="uploadImage">
      <img v-if="photo != 'NULL'" :src="photo" class="image-fit">
      <div class="btn-group me-2">
        <smallspan />
        <button type="button" class="btn btn-sm btn-outline-secondary" @click.stop="photo = 'NULL'; sendingpic = false">
          Discard
        </button>
      </div>
    </div>
    
    <br>
    <textarea v-model="message" class="bottom" placeholder="New message" />
    <div class="btn-group me-2">
      <button type="button" class="btn btn-sm btn-outline-secondary" @click.stop="sendMessage">
        Send
      </button>
      <button v-if="!sendingpic" type="button" class="btn btn-sm btn-outline-secondary" @click.stop="sendingpic = true">
        Picture selection
      </button>
    </div>
  </div>
</template>

<style>
bigspan { margin-left:11.8em }
.image-big{
  width: 2cm;
  object-fit: fit;
}
.image-min{
  width: 0.5cm;
  object-fit: fit;
}
</style>
