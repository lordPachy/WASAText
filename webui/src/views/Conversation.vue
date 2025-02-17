<script>

import { RouterLink } from 'vue-router';
import { useIDStore } from '../store';
import ConfirmedMsg from '../components/ConfirmedMsg.vue';

export default {
	props: {conversationid: String},
	data: function() {
		return {
			errormsg: null,
			confirmedmsg: null,			
			store: useIDStore(),
			data: {},
			isGroup: false,
			newuser: "",
			message: "",
			sendingpic: false,
			photo: "NULL",
			forwardingid: -1,
			adding: false,
			forwardingto: "",
			forwardables: [],
			replyingto: -1,
			chats: [],
			users: [],
			addables: [],
			timer: ''
		}
	},

	created() {
		this.refresh();
	},

	mounted() {
		this.refresh();

		// Updating page every 2000 ms
		this.timer = setInterval(this.refresh, 2000);
	},

	unmounted() {
		// Avoiding page update when it is closed
		clearInterval(this.timer);
	},

	methods: {
		/**
		 * It must be called every time a page element might be modified.
		 * 
		 * It refreshes the values of data, messages and isGroup.
		 */
		async refresh() {
			this.errormsg = null;
			try{
				this.getUsers();
				let response = await this.$axios.get("/conversations/" + this.conversationid, {headers: {Authorization: this.store.userInfo.id}});
				this.data = response.data;
				this.messages = response.data.messages;
				this.isGroup = Object.keys(response.data).length > 3;
			} catch (e) {
				this.errormsg = e.toString();
			}
		},

		/**
		 * It wraps the request using message, photo and replyingto.
		 */
		async sendMessage() {
			this.errormsg = null;
			try{
				let response = await this.$axios.post("/conversations/" + this.conversationid, {content: this.message, photo: this.photo, replyingto: this.replyingto}, {headers: {Authorization: this.store.userInfo.id}});

				// Reinitializing values
				this.message = "";
				this.photo = "NULL";
				this.replyingto = -1;
				this.sendingpic = false;

				// Updating
				this.refresh();
			} catch (e) {
				this.errormsg = e.toString();
			}
		},
		
		/**
		 * It employs newuser.
		 */
		async addMemberToGroup() {
			this.errormsg = null;
			try{
				let response = await this.$axios.put("/groups", {username: {name: this.newuser}, group: {id: parseInt(this.conversationid)}}, {headers: {Authorization: this.store.userInfo.id}});

				// Reinitializing values
				this.newuser = "";
				this.adding = false;

				// Updating
				this.refresh()
			} catch (e) {
				this.errormsg = e.toString();
			}
		},

		async leaveGroup() {
			this.errormsg = null;
			try{
				let response = await this.$axios.delete("/conversations/" + this.conversationid, {headers: {Authorization: this.store.userInfo.id}});
				this.$router.replace({name: 'conversations'});
			} catch (e) {
				this.errormsg = e.toString();
			}
		},

		async putComment(mess, reac) {
			this.errormsg = null;
			try{
				let response = await this.$axios.put("/conversations/" + this.conversationid + "/messages/" + mess.messageid, {reaction: reac}, {headers: {Authorization: this.store.userInfo.id}});

				// Updating
				this.refresh();
			} catch (e) {
				this.errormsg = e.toString();
			}
		},

		/**
		 * It returns true if the message has a comment
		 * 
		 * from the logged in user.
		 */
		hasOwnComment(mess) {
			let pres = false;
			mess.comments.forEach(c => {
				if (c.sender == this.store.userInfo.username){
					pres = true;
				}
			});

			return pres;
		},

		/**
		 * It deletes all comments to a message posted by
		 * 
		 * the logged in user (even if we are running with 
		 * 
		 * the assumption of a single comment per person).
		 */
		async deleteMyComments(mess) {
			this.errormsg = null;
			try{
				for (let i = 0; i < mess.comments.length; i++){
					if (mess.comments[i].sender == this.store.userInfo.username){
						await this.$axios.delete("/conversations/" + this.conversationid + "/messages/" + mess.messageid + "/comments/" + mess.comments[i].commentid.toString(), {headers: {Authorization: this.store.userInfo.id}});
					}
				}

				// Updating
				this.refresh();
			} catch (e) {
				this.errormsg = e.toString();
			}
		},

		/**
		 * It retrieves the content (text) of a message.
		 * 
		 * It is used only for retrieving the content of
		 * 
		 * the message the logged in user is replying to.
		 */
		repliedMessage(messid) {
			for (let i = 0; i < this.messages.length; i++){
				if (this.messages[i].messageid == messid){
					return "From " + this.messages[i].username + ": " + this.messages[i].content
				}
			}
		},

		/**
		 * It retrieves the photo of a message.
		 * 
		 * It is used only for retrieving the photo of
		 * 
		 * the message the logged in user is replying to.
		 */
		repliedPhoto(messid) {
			for (let i = 0; i < this.messages.length; i++){
				if (this.messages[i].messageid == messid){
					return this.messages[i].photo
				}
			}
		},

		/**
		 * It retrieves the photo of a user.
		 */
		userPhoto(username) {
		for (let i = 0; i < this.users.length; i++){
			if (this.users[i].username == username){
				return this.users[i].propic
			}
		}
		},

		/**
		 * It deletes the message passed as parameter.
		 * 
		 * It requires the logged in user to be the 
		 * 
		 * author of the message.
		 */
		async deleteMessage(mess) {
			this.errormsg = null;
			try{
				await this.$axios.delete("/conversations/" + this.conversationid + "/messages/" + mess.messageid, {headers: {Authorization: this.store.userInfo.id}});

				// Updating
				this.refresh();
			} catch (e) {
				this.errormsg = e.toString();
			}
		},

		/**
		 * It pushed groupsettings into the router history.
		 */
		async groupSettings() {
			this.errormsg = null;
			try {
				this.$router.push({name: 'groupsettings'}, {params: {conversationid: this.conversationid}});
			} catch (e) {
				this.errormsg = e.toString();
			}
		},

		/**
		 * It retrieves the currently active users.
		 */
		 async getUsers() {
			this.errormsg = null;
			try{
				// Retrieving WASAText active users
				let response = await this.$axios.get("/users", {headers: {Authorization: this.store.userInfo.id}, params: {username: ""}});
				this.users = response.data;

			} catch (e) {
				this.errormsg = e.toString();
			}
		},


		/**
		 * It retrieves contacts who can be added to the current groupchat.
		 * 
		 * Their usernames are saved in addables.
		 */
		async getAddables() {
			this.errormsg = null;
			try{
				this.addables = [];
				// Retrieving WASAText active users
				this.getUsers();

				// Adding all users (which are not already in the group) to the list of addables
				for (let i = 0; i < this.users.length; i++){
					for (let j = 0; j <= this.data.members.length; j++){
						if (j == this.data.members.length){
							this.addables.push(this.users[i].username);
							break;
						} else if (this.users[i].username == this.data.members[j].username){
							break;
						}
					}
				}

				// Allowing the interface to show the possible options
				this.adding = true;
			} catch (e) {
				this.errormsg = e.toString();
			}
		},

		/**
		 * It retrieves contacts whom messages can be forwarded to.
		 * 
		 * Their usernames are saved in forwardables.
		 */
		async getForwardables() {
			this.errormsg = null;
			try{
				this.forwardables = [];
				clearInterval(this.timer);

				// Retrieving user's current conversations
				let response = await this.$axios.get("/conversations", {headers: {Authorization: this.store.userInfo.id}});
				this.chats = response.data;

				// Retrieving WASAText active users
				this.getUsers();

				// Adding all user's chats to the list of forwardables
				for (let i = 0; i < this.chats.length; i++){
					this.forwardables.push(this.chats[i].name);
				}

				// Adding all users (which are not already in the started chats) to the list of forwardables
				for (let i = 0; i < this.users.length; i++){
					for (let j = 0; j <= this.forwardables.length; j++){
						if (j == this.forwardables.length && this.store.userInfo.username != this.users[i].username){
							this.forwardables.push(this.users[i].username);
							break;
						} else if (this.users[i].username == this.forwardables[j]){
							break;
						}
					}
				}
			} catch (e) {
				this.errormsg = e.toString();
			}
		},

		/**
		 * It retrieves the message id from forwardingid
		 * 
		 * and the conversation to forward to from forwardingto.
		 * 
		 * It resets their value after usage
		 */
		async forwardMessage() {
			this.errormsg = null;
			try{
				for (let i = 0; i <= this.chats.length; i++){
					// Forwarding to a non-started [private] chat
					if (i == this.chats.length){
						let response = await this.$axios.put("/conversations", {isgroup: false, members: [{name: this.forwardingto}], groupname: ""}, {headers: {Authorization: this.store.userInfo.id}});
						await this.$axios.post("/conversations/" + this.conversationid + "/messages/" + this.forwardingid, response.data, {headers: {Authorization: this.store.userInfo.id}});

					// Forwarding to a started chat
					} else if (this.forwardingto == this.chats[i].name){
						let response = await this.$axios.post("/conversations/" + this.conversationid + "/messages/" + this.forwardingid, this.chats[i].chatid, {headers: {Authorization: this.store.userInfo.id}});
						break;
					}
				}
				
				// Resetting values and refreshing page
				this.forwardingid = -1;
				this.forwardingto = "";

				// Updating
				this.refresh();
				this.timer = setInterval(this.refresh, 2000);

				// Showing confirmation message
				this.confirmedmsg = "Message forwarded successfully";
				await new Promise(resolve => setTimeout(resolve, 5000));
				this.confirmedmsg = "";
			} catch (e) {
				this.errormsg = e.toString();
			}
		},

		/**
		 * It is a support function for the image uploading
		 * 
		 * feature.
		 */
		async uploadImage(a) {
			const image = a.target.files[0];
			if (image == null){
				return;
			} else if (image.name.slice(-4) != ".png"){
				this.errormsg = "Only png images can be uploaded";
				clearInterval(this.timer);
				await new Promise(resolve => setTimeout(resolve, 5000));
				this.timer = clearInterval(this.refresh, 2000);
				this.errormsg = null;
				return;
			}
			const reader = new FileReader();
			reader.readAsDataURL(image);
			reader.onload = a =>{
			this.photo = a.target.result;
			}
		},

		/**
		 * It reinitializes values after discarding a forward.
		 */
		discardForward(){
			this.timer = setInterval(this.refresh, 2000);
			this.forwardingid = -1;
			this.forwardingto = '';
		},

		/**
		 * It returns the correct icon name to be shown according to the message 
		 * 
		 * checkmarks' status.
		 */
		checkmarks(mess){
			switch(mess.checkmarks){
				case 0:
					return "minus.png"
				case 1:
					return "single.png"
				case 2:
					return "double.png"
			}
		}


	}
}
</script>

<template>
  <!--Private chat headers-->
  <div
    v-if="!isGroup" 
    class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom"
  >
    <h1 class="h2">
      Private chat with
      <img v-if="data.user?.propic != 'NULL'" :src="data.user?.propic" class="image-big">
      {{ data.user?.username }}
    </h1>
  </div>

  <!--Group headers-->
  <div
    v-else 
    class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom"
  >
    <h1 class="h2">
      Group:
      <img v-if="data.groupphoto != 'NULL'" :src="data.groupphoto" class="image-big">
      {{ data.groupname }}
    </h1>
  </div>

  <!--Error and confirmation messages-->
  <div>
    <ErrorMsg v-if="errormsg" :msg="errormsg" />
    <ConfirmedMsg v-if="confirmedmsg" :msg="confirmedmsg" />
  </div>

  <!--Group information and options-->
  <div v-if="isGroup">
    <h5 class="h5">Group members</h5>
    <ul>
      <li v-for="u in data.members" :key="u">
        <img v-if="u.propic != 'NULL'" :src="u.propic" class="image-fit">
        {{ u.username }} 
      </li>
    </ul>

    <h5 class="h5">Group options</h5>
    <!--Member adding-->
    <div v-if="!adding">
      <button type="button" class="btn btn-sm btn-outline-secondary" @click="getAddables">
        Add new member to group
      </button>
    </div>

    <div v-else>
      <select v-model="newuser" class="mb-3">
        <option disabled value="">Please select...</option>
        <option v-for="c in addables" :key="c">{{ c }}</option>
      </select>
      <button :disabled="newuser == ''" type="button" class="btn btn-sm btn-outline-secondary" @click="addMemberToGroup">
        Add member to group
      </button>
      <button type="button" class="btn btn-sm btn-outline-secondary" @click="adding = false; newuser = ''">
        Discard
      </button>
    </div>
    
    <!--Other group settings-->
    <br>
    <button type="button" class="btn btn-sm btn-outline-secondary mt-3" @click="groupSettings">
      Group Settings
    </button>
    <button type="button" class="btn btn-sm btn-outline-secondary mt-3" @click="leaveGroup">
      Leave this group
    </button>
  </div>

  <!--Messages-->
  <div>
    <h5 class="h5 mt-4">Messages</h5>
    <ul>
      <li v-for="m in data.messages" :key="m" class="mb-4">
        <!--Content and primary options-->
        <p>
          <!--Message content-->
          ({{ m.timestamp.slice(0, 10) + " " + m.timestamp.slice(11, 19) }}) <img v-if="isGroup && userPhoto(m.username) != 'NULL' && userPhoto(m.username) != ''" :src="userPhoto(m.username)" class="image-min"> {{ m.username }}{{ (m.og_sender != "NULL") ? (" (Originally written by " + m.og_sender + ")") : ("") }}: "{{ m.content }}"
          
          <!--Message options-->
          <!--Reply-->
          <button style="background-color:#FFCCCB;" type="button" class="btn btn-sm btn-outline-secondary" @click="replyingto = m.messageid">
            Reply
          </button>
          
          <!--Forward-->
          <button v-if="forwardingid != m.messageid" type="button" class="btn btn-sm btn-outline-secondary" @click="forwardingid = m.messageid; getForwardables()">
            Forward
          </button>
          <select v-if="forwardingid == m.messageid" v-model="forwardingto" class="mb-3">
            <option disabled value="-1">Please select one</option>
            <option v-for="c in forwardables" :key="c">{{ c }}</option>
          </select>
          <button v-if="forwardingid == m.messageid" type="button" class="btn btn-sm btn-outline-secondary" @click="forwardMessage">
            Send forwarded message
          </button>
          <button v-if="forwardingid == m.messageid" type="button" class="btn btn-sm btn-outline-secondary" @click="discardForward">
            Discard
          </button>
          
          <!--Delete-->
          <button v-if="m.username == store.userInfo.username" style="background-color:#ccccc8;" type="button" class="btn btn-sm btn-outline-secondary" @click="deleteMessage(m)">
            Delete message
          </button>
        </p>

        <!--Checkmarks-->
        <div v-if="m.username == store.userInfo.username">
          Status: <img :src="checkmarks(m)" class="image-min">
        </div>
		
        <!--Picture showing-->
        <div>
          <span class="bigspan" />
          <img v-if="m.photo != 'NULL' && m.photo != ''" :src="m.photo" class="image-big">
        </div>

        <!--Replying to-->
        <div>
          <p v-if="m.replyingto != -1">
            <span class="bigspan" />This message is replying to:<br>
            <span class="bigspan" />{{ repliedMessage(m.replyingto) }}
            <img v-if="repliedPhoto(m.replyingto) != 'NULL' && repliedPhoto(m.replyingto) != ''" :src="repliedPhoto(m.replyingto)" class="image-min">
          </p>
        </div>

        <!--Put a comment-->
        <div v-if="m.username != store.userInfo.username && !hasOwnComment(m)">
          <span class="bigspan" />Put a comment:
          <br>
          <span class="bigspan" />
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="putComment(m, 'laugh')"> 
            <img src="/laugh.png" class="image-min">
          </button>
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="putComment(m, 'sad')"> 
            <img src="/sad.png" class="image-min">
          </button>
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="putComment(m, 'love')"> 
            <img src="/love.png" class="image-min">
          </button>
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="putComment(m, 'pray')"> 
            <img src="/pray.png" class="image-min">
          </button>
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="putComment(m, 'thumbs_up')"> 
            <img src="/thumbs_up.png" class="image-min">
          </button>
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="putComment(m, 'surprised')"> 
            <img src="/surprised.png" class="image-min">
          </button>
        </div>

        <!--Delete a comment-->
        <div v-if="m.username != store.userInfo.username && hasOwnComment(m)">
          <span class="bigspan" />
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="deleteMyComments(m)">
            Delete my comment
          </button>
        </div>

        <!--Show comments-->
        <div>
          <span class="bigspan" />Comments:
          <p v-for="r in m.comments" :key="r">
            <span class="bigspan" />
            <img :src="r.reaction + '.png'" class="image-min">
            by {{ r.sender }}
          </p>
        </div>
      </li>
    </ul>
  </div>

  <!--Message reply -->
  <div v-if="replyingto != -1">
    Currently replying to the following message:<br>
    {{ repliedMessage(replyingto) }}
    <div>
      <span class="bigspan" />
      <img v-if="repliedPhoto(replyingto) != 'NULL' && repliedPhoto(replyingto) != ''" :src="repliedPhoto(replyingto)" class="image-min">
    </div>
    <button type="button" class="btn btn-sm btn-outline-secondary" @click="replyingto = -1">
      Abort reply
    </button>
  </div>

  <!--Sending picture-->
  <div v-if="sendingpic">
    <p class="mt-5">Select a picture:</p>
    <input type="file" accept="image/png" @change="uploadImage">
    <img v-if="photo != 'NULL'" :src="photo" class="image-fit">
    <div class="btn-group me-2">
      <span class="smallspan" />
      <button type="button" class="btn btn-sm btn-outline-secondary" @click="photo = 'NULL'; sendingpic = false">
        Discard
      </button>
    </div>
  </div>

  <br>

  <!--New message box-->
  <div> 
    <textarea v-model="message" class="bottom" placeholder="New message" />
    <div class="btn-group me-2">
      <button type="button" :disabled="photo == 'NULL' && message.length < 1" class="btn btn-sm btn-outline-secondary" @click="sendMessage">
        Send
      </button>
      <button v-if="!sendingpic" type="button" class="btn btn-sm btn-outline-secondary" @click="sendingpic = true">
        Picture selection
      </button>
    </div>
  </div>
</template>

<style>
.image-big{
  width: 2cm;
  object-fit: fit;
}
.image-min{
  width: 0.5cm;
  object-fit: fit;
}
</style>
