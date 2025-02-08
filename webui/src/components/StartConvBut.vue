<script>
// Component definition (JavaScript)
export default {
    // `data` defines a list of reactive variables.
    props: ['token', 'show'],
    emits: ['create-convo'],
    data: function() {
        return {
            username: ""
        }
    },
    methods: {
			async createConvo() {
				this.loading = true;
				this.errormsg = null;
                let response = await this.$axios.put("/conversations", {isgroup: false, members: [{name: this.username}]}, {headers: {Authorization: this.token}});
				this.loading = false;
			},
    }
}
</script>

<!-- template definition -->
<template>
    <div  v-if="show">
        <p>Insert the user who you want to chat with!</p>
        <input v-model="username" placeholder="Username">
        <button type="button" class="btn btn-sm btn-outline-secondary" @click="$emit('create-convo')">
            Create Conversation
        </button>
    </div>
</template>

<!-- CSS style -->
<style scoped>
</style>