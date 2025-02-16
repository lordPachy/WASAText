<script>

import { useIDStore } from '../store';

export default {
	props: {conversationid: String},
	data: function() {
		return {
			errormsg: null,
			store: useIDStore(),
			newgroupname: "",
			previewImage: null,
		}
	},
	methods: {
		/**
		 * It sets the current group name and returns
		 * to the group chat page.
		 */
		async setGroupName() {
			this.errormsg = null;
			try {
				let response = await this.$axios.put("/conversations/" + this.conversationid + "/settings/groupname", {value: this.newgroupname}, {headers: {Authorization: this.store.userInfo.id}});
				this.$router.push({name: "conversation", params: {conversationid: this.conversationid}});
			} catch (e) {
				if (e.toString() == "AxiosError: Request failed with status code 400") {
					this.errormsg = "Groupnames must be between 3 and 16 alphanumeric characters; no spaces";
				} else {
					this.errormsg = e.toString();
				}
			}
		},

		/**
		 * It is a support function for the image uploading
		 * 
		 * feature.
		 */
		async uploadImage(a) {
			this.errormsg = null;
			try {
				const image = a.target.files[0];
				if (image == null){
					return;
				} else if (image.name.slice(-4) != ".png"){
					this.errormsg = "Only png images can be uploaded";
					return;
				}
				const reader = new FileReader();
				reader.readAsDataURL(image);
				reader.onload = a =>{
					this.previewImage = a.target.result;
				};
			} catch (e) {
				this.errormsg = e.toString();
			}
		},

		/**
		 * It sets the group profile picture and returns
		 * 
		 * to the group chat page.
		 */
		async setGroupPic() {
			this.errormsg = null;
			try {
				await this.$axios.put("/conversations/" + this.conversationid + "/settings/grouphoto", {image: this.previewImage}, {headers: {Authorization: this.store.userInfo.id}});
				this.$router.push({name: "conversation", params: {conversationid: this.conversationid}});
			} catch (e) {
				this.errormsg = e.toString();
			}
		},

		/**
		 * A utility function to clean the entire page.
		 */
		cleanPage() {
			this.previewImage = null;
			window.location.reload();
		}
	}
}
</script>

<template>
  <!--Header-->
  <div
    class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom"
  >
    <h1 class="h2">Group settings</h1>
  </div>

  <!--Error messages-->
  <div>
    <ErrorMsg v-if="errormsg" :msg="errormsg" />
  </div>

  <!--Group name changing-->
  <div>
    <p>Enter the new groupname if you want to change it:</p>
    <input
      v-model="newgroupname" placeholder="Insert new group name"
    >

    <div class="btn-group me-2">
      <button type="button" class="btn btn-sm btn-outline-secondary" @click="setGroupName">
        Apply changes
      </button>
    </div>
  </div>
  
  <!--Group picture changing-->
  <div>
    <p class="mt-5">Insert a picture if you want to change it:</p>

    <input type="file" accept="image/png" @change="uploadImage">
    <img v-if="previewImage != null" :src="previewImage" class="image-fit">

    <div class="btn-group me-2">
      <button type="button" class="btn btn-sm btn-outline-secondary" @click="setGroupPic">
        Apply changes
      </button>

      <button v-if="previewImage != null" type="button" class="btn btn-sm btn-outline-secondary" @click="cleanPage">
        Discard operation
      </button>
    </div>
  </div>
</template>

<style>
.image-fit{
  width: 1cm;
  object-fit: fit;
}
</style>
