<script>
// Component definition (JavaScript)
export default {
    // `data` defines a list of reactive variables
    props: ['token', 'show'],
    emits: ['get-convos'],
    data: function() {
        return {
            chats: {},
            showchats: false,
        }
    },
    methods: {
			async getConvos() {
				this.loading = true;
				this.errormsg = null;
                let response = await this.$axios.get("/conversations", {headers: {Authorization: this.token}});
                this.chats = response.data;
                this.showchats = true;
				this.loading = false;
			},
    }
}
</script>

<!-- template definition -->
<template>
    <div  v-if="show">
        <button type="button" class="btn btn-sm btn-outline-secondary" @click="$emit('get-convos')"> 
            Get conversations
        </button>
        <div>
            <p>Chats are: {{ chats }}</p>
        </div>
        
    </div>
</template>

<!-- CSS style -->
<style scoped>
</style>