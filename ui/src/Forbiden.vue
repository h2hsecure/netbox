<script setup lang="ts">
import { computed } from 'vue'
import Forbiden from './views/Forbiden.vue'


const systemId = computed(() => {
return window._ntb_dds.id
});

const system = computed(() => {
return window._ntb_dds.t.main.footer.system
});

const per = computed(() => {
return window._ntb_dds.t.main.footer.per
});

</script>

<template>
  <header>
    <br /><br /><br />
  </header>
  <main>
    <Forbiden />
  </main>  
  <footer>
    <hr />
    <p><small>{{ system }}: {{ systemId }}</small></p>
    <p><small v-html="per"></small></p>
  </footer>
</template>
