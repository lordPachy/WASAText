<script>

export default {
	props: {conversationid: String},
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
			newgroupname: "",
			showimage: false,
			previewImage: null,
		}
	},
	methods: {
			async setGroupName() {
				this.loading = true;
				this.errormsg = null;
				try {
					let response = await this.$axios.put("/conversations/" + this.conversationid + "/settings/groupname", {value: this.newgroupname}, {headers: {Authorization: this.$router.id}});
					this.$router.push({name: "conversation", params: {conversationid: this.conversationid}});
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
					};
				} catch (e) {
					this.errormsg = e.toString();
				}
				this.loading = false;
			},

			async setGroupPic() {
				this.loading = true;
				this.errormsg = null;
				try {
					console.log("hello");
					await this.$axios.put("/conversations/" + this.conversationid + "/settings/grouphoto", {image: this.previewImage}, {headers: {Authorization: this.$router.id}});
					this.$router.push({name: "conversation", params: {conversationid: this.conversationid}});
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
      <h1 class="h2">Group settings</h1>
      <div class="btn-toolbar mb-2 mb-md-0" />
    </div>
    <p>Enter the new groupname if you want to change it:</p>
    <input
      v-model="newgroupname" placeholder="Insert new group name"
    >

    <div class="btn-group me-2">
      <button type="button" class="btn btn-sm btn-outline-secondary" @click.stop="setGroupName">
        Apply changes
      </button>
    </div>

    <p class="mt-5">Insert a picture if you want to change it:</p>
    <input type="file" accept="image/png" @change="uploadImage">
    <img v-if="previewImage != null" :src="previewImage" class="image-fit">
    <div class="btn-group me-2">
      <button type="button" class="btn btn-sm btn-outline-secondary" @click.stop="setGroupPic">
        Apply changes
      </button>
      <button type="button" class="btn btn-sm btn-outline-secondary" @click.stop="previewImage = null">
        Reset image
      </button>
    </div>
    <ErrorMsg v-if="errormsg" :msg="errormsg" />
  </div>
</template>

<style>
.image-fit{
  width: 1cm;
  object-fit: fit;
}
</style>
