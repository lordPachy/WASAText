<script>
// Component definition (JavaScript)
export default {
    // `data` defines a list of reactive variables.
    props: ['token', 'show'],
    emits: ['chats'],
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
				try {
					let response = await this.$axios.get("/converations", {identifier: token});
                    this.chats = response.data;
                    this.showchats = true;
				} catch (e) {
					this.errormsg = e.toString();
				}
				this.loading = false;
			},
    }
}
</script>

<!-- template definition -->
<template>
    <div  v-if="show">
        <button type="button" class="btn btn-sm btn-outline-secondary" @click.stop="getConvos">
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