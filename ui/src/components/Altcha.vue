<script setup lang="ts">
import { ref, useTemplateRef, onMounted, onUnmounted, watch } from 'vue';
import   'altcha';

const altchaWidget2 = useTemplateRef<HTMLElement>("altchaWidget2");
const props = defineProps({
  payload: {
    type: String,
    required: false,
  },
  label: {
    type: Object,
    required: false,
    default: () => ({})
  }
});
const emit = defineEmits(['update:payload', 'update:start']);
const internalValue = ref(props.payload);

watch(internalValue, (v) => {
  emit('update:payload', v || '');
});

const onStateChange = (ev: Event) => {
  if ('detail' in ev) {
    const { payload, state } = (<CustomEvent>ev).detail;
    if (state === 'verified' && payload) {
      internalValue.value = payload;
    } else {
      if (state === 'serververification') {
        emit('update:start');
      }
      internalValue.value = '';
    }
  }
};

onMounted(() => {
  if (altchaWidget2.value) {
    altchaWidget2.value.addEventListener('statechange', onStateChange);
  }
});

onUnmounted(() => {
  if (altchaWidget2.value) {
    altchaWidget2.value.removeEventListener('statechange', onStateChange);
  }
});
</script>

<template>
  <!-- Configure your `challengeurl` and remove the `test` attribute, see docs: https://altcha.org/docs/website-integration/#using-altcha-widget -->
  <altcha-widget
    ref="altchaWidget2"
    challengeurl="/ddos/challenge"
    style="--altcha-max-width:100%;--altcha-border-width: 0px;"
    hidelogo
    hidefooter
    :strings="JSON.stringify(props.label)"
  ></altcha-widget>
</template>