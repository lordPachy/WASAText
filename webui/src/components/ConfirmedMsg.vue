<script>
export default {
	props: {'msg': String}
}
</script>

<template>
  <div id="novin-alert-container">
    <div class="alert novin-alert">
      <img src="/single.png" class="image-min"> {{ msg }}
    </div>
  </div>
</template>

<style>
</style>
