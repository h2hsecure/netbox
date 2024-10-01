<script setup>
import { ref } from 'vue';
import axios from 'axios';

const data = ref({
  checked: false,
  referer: window.location.host,
  fronterid: "1231312",
  showReadme: false
})

function checkCallback() {
  axios.get('/ddos/check')
    .then(function (response) {
      console.log(response);
      window.location = "/"
    })
}

function showHideDrawer() {
  data.value.showReadme = !data.value.showReadme
}
</script>

<template>
  <div class="bg-white">
    <div class="relative isolate px-6 pt-4 lg:px-8">
      <div class="absolute inset-x-0 -top-40 -z-10 transform-gpu overflow-hidden blur-3xl sm:-top-80"
        aria-hidden="true">
        <div
          class="relative left-[calc(50%-11rem)] aspect-[1155/678] w-[36.125rem] -translate-x-1/2 rotate-[30deg] bg-gradient-to-tr from-[#ff80b5] to-[#9089fc] opacity-30 sm:left-[calc(50%-30rem)] sm:w-[72.1875rem]"
          style="clip-path: polygon(74.1% 44.1%, 100% 61.6%, 97.5% 26.9%, 85.5% 0.1%, 80.7% 2%, 72.5% 32.5%, 60.2% 62.4%, 52.4% 68.1%, 47.5% 58.3%, 45.2% 34.5%, 27.5% 76.7%, 0.1% 64.9%, 17.9% 100%, 27.6% 76.8%, 76.1% 97.7%, 74.1% 44.1%)">
        </div>
      </div>
      <div class="mx-auto max-w-2xl py-32 sm:py-48 lg:py-56">
        <div class="hidden sm:mb-8 sm:flex sm:justify-center">
          <div
            class="relative rounded-full px-3 py-1 text-sm leading-6 text-gray-600 ring-1 ring-gray-900/10 hover:ring-gray-900/20">
            {{ $t("main.why") }} <a href="#" class="font-semibold text-indigo-600" @click.stop="showHideDrawer"><span
                class="absolute inset-0" aria-hidden="true"></span>{{ $t("main.more") }} <span
                aria-hidden="true">&rarr;</span></a>
          </div>
        </div>
        <div class="text-center">
          <h1 class="text-4xl font-bold tracking-tight text-gray-900 sm:text-6xl">{{ $t("main.h1") }}
          </h1>
          <p class="mt-6 text-lg leading-8 text-gray-600"> {{ $t("main.h2") }} {{
            data.referer }}</p>
          <div class="mt-10 flex items-center justify-center gap-x-6">
            <div class="flex items-center me-4">
              <input checked id="red-checkbox" type="checkbox" v-model="data.checked" @click="checkCallback"
                class="w-4 h-4 text-red-600 bg-gray-100 border-gray-300 rounded focus:ring-red-500 dark:focus:ring-red-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600">
              <label for="red-checkbox" class="ms-2 text-sm font-bold"> {{ $t("main.check") }}</label>
            </div>
          </div>
        </div>
        <div class="hidden sm:mt-8 sm:flex sm:justify-center">
          <div
            class="relative px-3 py-1 text-sm leading-6 text-gray-600 ring-1 ring-gray-900/10 hover:ring-gray-900/20">
            {{ $t("main.id") }}: {{ data.fronterid }} <span aria-hidden="true">&rarr;</span> {{ $t("main.provided",
              ["www.h2hsecure.com"]) }}

          </div>
        </div>
      </div>
      <div
        class="absolute inset-x-0 top-[calc(100%-13rem)] -z-10 transform-gpu overflow-hidden blur-3xl sm:top-[calc(100%-30rem)]"
        aria-hidden="true">
        <div
          class="relative left-[calc(50%+3rem)] aspect-[1155/678] w-[36.125rem] -translate-x-1/2 bg-gradient-to-tr from-[#ff80b5] to-[#9089fc] opacity-30 sm:left-[calc(50%+36rem)] sm:w-[72.1875rem]"
          style="clip-path: polygon(74.1% 44.1%, 100% 61.6%, 97.5% 26.9%, 85.5% 0.1%, 80.7% 2%, 72.5% 32.5%, 60.2% 62.4%, 52.4% 68.1%, 47.5% 58.3%, 45.2% 34.5%, 27.5% 76.7%, 0.1% 64.9%, 17.9% 100%, 27.6% 76.8%, 76.1% 97.7%, 74.1% 44.1%)">

        </div>
      </div>
    </div>
  </div>

  <div id="drawer-example"
    class="fixed top-0 left-0 right-0 z-40 w-full p-4 transition-transform bg-white light:bg-gray-800"
    :class="{ 'translate- none': data.showReadme, '-translate-y-full': !data.showReadme }" tabindex="-1"
    aria-labelledby="drawer-label">
    <h5 id="drawer-label"
      class="inline-flex items-center mb-4 text-base font-semibold text-gray-500 dark:text-gray-600"><svg
        class="w-4 h-4 me-2.5" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="currentColor"
        viewBox="0 0 20 20">
        <path
          d="M10 .5a9.5 9.5 0 1 0 9.5 9.5A9.51 9.51 0 0 0 10 .5ZM9.5 4a1.5 1.5 0 1 1 0 3 1.5 1.5 0 0 1 0-3ZM12 15H8a1 1 0 0 1 0-2h1v-3H8a1 1 0 0 1 0-2h2a1 1 0 0 1 1 1v4h1a1 1 0 0 1 0 2Z" />
      </svg> {{ $t("drawer.title") }}</h5>
    <button type="button" data-drawer-hide="drawer-example" aria-controls="drawer-example"
      class="text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm w-8 h-8 absolute top-2.5 end-2.5 flex items-center justify-center dark:hover:bg-gray-600 dark:hover:text-white"
      @click="showHideDrawer">
      <svg class="w-3 h-3" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 14 14">
        <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6" />
      </svg>
      <span class="sr-only">{{ $t("drawer.close") }}</span>
    </button>

    <p class="mb-6 text-sm text-gray-500 dark:text-gray-400">{{ $t("drawer.content") }}</p>
  </div>

</template>
<style></style>
