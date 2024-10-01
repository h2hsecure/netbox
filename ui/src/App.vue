<script setup>
import { RouterView } from 'vue-router'
import { onMounted } from 'vue'


onMounted(() => {
  const availibleLocales = this.$i18n.availableLocales;
  const usersLanguage = window.navigator.language

  if (availibleLocales.includes(usersLanguage)) {
    this.$i18n.locale = usersLanguage;
  } else {
    this.$i18n.locale = 'en';
  }
})
</script>

<template>


  <RouterView />
</template>
