<script>

import { useIDStore } from '../store';

export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			store: useIDStore(),
			newusername: "",
			previewImage: null,
		}
	},
	methods: {
			async setMyUsername() {
				this.loading = true;
				this.errormsg = null;
				try {
					let response = await this.$axios.put("/settings/username", {name: this.newusername}, {headers: {Authorization: this.store.userInfo.id}});
					this.store.changeUsername(this.newusername);
					this.$router.push({name: 'homepage'});
				} catch (e) {
					if (e.toString() == "AxiosError: Request failed with status code 403"){
						this.errormsg = "Username already in use";
					} else if (e.toString() == "AxiosError: Request failed with status code 400") {
						this.errormsg = "Usernames must be between 3 and 16 alphanumeric characters; no spaces";
					} else {
						this.errormsg = e.toString();
					}
				}
				this.loading = false;
			},
			async setMyProPic() {
				this.loading = true;
				this.errormsg = null;
				try {
					await this.$axios.put("/settings/profilepicture", {image: this.previewImage}, {headers: {Authorization: this.store.userInfo.id}});
					this.$router.push({name: 'homepage'});
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
      <h1 class="h2">User settings</h1>
      <div class="btn-toolbar mb-2 mb-md-0" />
    </div>
    <p>Enter your username if you want to change it:</p>
    <input
      v-model="newusername" placeholder="Insert new username"
    >

    <div class="btn-group me-2">
      <button type="button" class="btn btn-sm btn-outline-secondary" @click.stop="setMyUsername">
        Apply changes
      </button>
    </div>

    <p class="mt-5">Insert a picture if you want to change it:</p>
    <input type="file" accept="image/png" @change="uploadImage">
    <img v-if="previewImage != null" :src="previewImage" class="image-fit">
    <div class="btn-group me-2">
      <button type="button" class="btn btn-sm btn-outline-secondary" @click.stop="setMyProPic">
        Apply changes
      </button>
      <button v-if="previewImage != null" type="button" class="btn btn-sm btn-outline-secondary" @click.stop="previewImage = null">
        Discard operation
      </button>
    </div>
    <ErrorMsg v-if="errormsg" :msg="errormsg" />
  </div>
</template>

<style>
.image-fit{
  height: 7%;
  width: 7%;
  object-fit: cover;
}
</style>
