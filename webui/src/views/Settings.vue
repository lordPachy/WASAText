<script>

export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
			newusername: "",
			showimage: false,
			previewImage: null,
		}
	},
	methods: {
			async setMyUsername() {
				this.loading = true;
				this.errormsg = null;
				try {
					let response = await this.$axios.put("/settings/username", {name: this.newusername}, {headers: {Authorization: this.$router.id}});
					this.$router.username = this.newusername;
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
					const reader = new FileReader();
					reader.readAsDataURL(image);
					reader.onload = a =>{
						this.previewImage = a.target.result;
						console.log(this.previewImage);
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
    <img :src="previewImage" class="uploading-image">
    <div class="btn-group me-2">
      <button type="button" class="btn btn-sm btn-outline-secondary">
        Apply changes
      </button>
    </div>
    <ErrorMsg v-if="errormsg" :msg="errormsg" />
  </div>
</template>

<style>
</style>
