<script setup lang="ts">
import { ref, computed } from 'vue'
import Atacha from '../components/Altcha.vue'

const altchaPayload = ref('');
const favico = ref(window._ntb_dds.l)
const count = ref(false)
const done = ref(true)

const checkCallback = () => {
  var bodyFormData = new FormData();
  bodyFormData.append('challenge', altchaPayload.value)
  fetch('/ddos/accept', { method: "POST", body: bodyFormData})
  .then(() => {
        done.value = false
        setTimeout(() => {  window.location.assign(window._ntb_dds.r);  } , 500);
    })
}

const startCallback = () => {
    setTimeout(() => { count.value = true;  } , 3000);
}

const hostname = computed(() => {
return window._ntb_dds.h
});

const h1 = computed(() => {
return window._ntb_dds.t.main.h1
});

const h2 = computed(() => {
return window._ntb_dds.t.main.h2
});

const timeout = computed(() => {
return window._ntb_dds.t.main.timeout
});

const lbl = window._ntb_dds.t.main.at
</script>

<template>
    <div style="display: flex">
        <div>
            <img :src="favico" width="64"/>
        </div>
        <div>&nbsp;</div>
        <div>
            <h2>{{ hostname }} </h2>
        </div>
    </div>
    <h5>{{ h1 }}</h5>    
    <br />
    <br />
    <Atacha v-model:payload="altchaPayload" id="red-checkbox" @update:payload="checkCallback" @update:start="startCallback"
    :label="lbl"
    />
    <br />
    <br />
    <br />
    <p v-if="done">{{ hostname }} {{ h2 }}</p>

    <p v-if="count"><small v-html="timeout"></small></p>
</template>
