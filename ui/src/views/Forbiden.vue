<script setup lang="ts">
import { ref, onMounted } from 'vue'

const count = ref(false)

onMounted(() => {
  setTimeout(() => { count.value = true;  } , 5000);
})

</script>

<template>
    <h2>Forbiden</h2>
    <p v-show="count"><small> takes too much time for validating <a href="javascript:window.location.reload()"> please click</a> to refresh page</small></p>
</template>