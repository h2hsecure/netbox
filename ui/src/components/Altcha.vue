<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue';

// Importing altcha package will introduce a new element <altcha-widget>
import   'altcha';

const altchaWidget2 = ref(null);
const props = defineProps({
  payload: {
    type: String,
    required: false,
  },
  label: {
    type: String,
    required: false,
  }
});
const emit = defineEmits(['update:payload']);
const internalValue = ref(props.payload);

watch(internalValue, (v) => {
  emit('update:payload', v || '');
});

const onStateChange = (ev) => {
  if ('detail' in ev) {
    const { payload, state } = ev.detail;
    if (state === 'verified' && payload) {
      internalValue.value = payload;
    } else {
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
    auto="onload"
    strings="{&quot;label&quot;:&quot;{{props.label}}&quot;}"
  ></altcha-widget>
</template>