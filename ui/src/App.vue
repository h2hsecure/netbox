<script setup lang="ts">
import { ref, computed } from 'vue'
import UserValidate from './views/UserValidate.vue'
import Queue from './views/Queue.vue'
import Forbiden from './views/Forbiden.vue'

let routes = new Map<string, any> ()
routes.set('/', UserValidate)
routes.set('queue', Queue)
routes.set('forbiden', Forbiden)


const currentPath = ref(window.location.hash)

window.addEventListener('hashchange', () => {
  currentPath.value = window.location.hash
})

const currentView = computed(() => {
  return routes.get(currentPath.value.slice(1) || '/') || UserValidate
})

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
    <component :is="currentView" />
  </main>  
  <footer>
    <hr />
    <p><small>{{ system }}: {{ systemId }}</small></p>
    <p><small v-html="per"></small></p>
  </footer>
</template>
