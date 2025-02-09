<script>
import ConvBut from '../components/ConvBut.vue';
import StartConvBut from '../components/StartConvBut.vue';

export default {
	components:{
		ConvBut,
		StartConvBut
	},
	props: {token: String, username: String},
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
			newUser: "",
			groupName: "",
			showConversations: false,
		}
	},
	methods: {
		async conversations() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/conversations");
				this.conversations = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.errormsg = null;
			this.loading = false;
		},
			async callGetConvos() {
				this.loading = true;
				this.errormsg = null;
				try {this.$refs.myConvBut.getConvos();}
				catch (e) {
					this.errormsg = e.toString();
				}
				this.loading = false;
			},
			async callCreateConvo() {
				this.loading = true;
				this.errormsg = null;
				try {this.$refs.myStartConvBut.createConvo();}
				catch (e) {
					this.errormsg = e.toString();
				}
				this.loading = false;
			}
	}
}
</script>

<template>
  <div>
    <div
      class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom"
    >
      <h1 class="h2">WASAText</h1>
      <div class="btn-toolbar mb-2 mb-md-0">
        <div class="btn-group me-2">
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="conversations">
            Conversations
          </button>
        </div>
      </div>
    </div>
    <p>
      <StartConvBut ref="myStartConvBut" :token="id" :show="showConversations" @create-convo="callCreateConvo" />
    </p>
    <p>
      <ConvBut ref="myConvBut" :token="id" :show="showConversations" @get-convos="callGetConvos" />
    </p>
    <ErrorMsg v-if="errormsg" :msg="errormsg" />
  </div>
</template>

<style>
</style>
