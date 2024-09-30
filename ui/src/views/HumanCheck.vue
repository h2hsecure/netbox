<script setup>
import { reactive } from 'vue';
import axios from 'axios';

var data = reactive({
  checked: false,
  referer: "www.pdaccess.com",
  fronterid: "1231312"
})

function checkCallback() {
  axios.get('/ddos/check')
    .then(function (response) {
      console.log(response);
      window.location = "/"
    })
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
            Why did I see this Page? <a href="#" class="font-semibold text-indigo-600"><span class="absolute inset-0"
                aria-hidden="true"></span>Read more <span aria-hidden="true">&rarr;</span></a>
          </div>
        </div>
        <div class="text-center">
          <h1 class="text-4xl font-bold tracking-tight text-gray-900 sm:text-6xl">Frontier wants to check you are a
            human
          </h1>
          <p class="mt-6 text-lg leading-8 text-gray-600">After checking your request you will be continue on {{
            data.referer }}</p>
          <div class="mt-10 flex items-center justify-center gap-x-6">
            <div class="flex items-center me-4">
              <input checked id="red-checkbox" type="checkbox" v-model="data.checked" @click="checkCallback"
                class="w-4 h-4 text-red-600 bg-gray-100 border-gray-300 rounded focus:ring-red-500 dark:focus:ring-red-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600">
              <label for="red-checkbox" class="ms-2 text-sm font-bold"> Please click that you are a human</label>
            </div>
          </div>
        </div>
        <div class="hidden sm:mt-8 sm:flex sm:justify-center">
          <div
            class="relative px-3 py-1 text-sm leading-6 text-gray-600 ring-1 ring-gray-900/10 hover:ring-gray-900/20">
            front Id: {{ data.fronterid }} <span aria-hidden="true">&rarr;</span> Provided by: H2H Secure
            www.h2hsecure.com
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


</template>
<style>
.chat-notification {
  display: flex;
  max-width: 24rem;
  margin: 0 auto;
  padding: 1.5rem;
  border-radius: 0.5rem;
  background-color: #fff;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
}

.chat-notification-logo-wrapper {
  flex-shrink: 0;
}

.chat-notification-logo {
  height: 3rem;
  width: 3rem;
}

.chat-notification-content {
  margin-left: 1.5rem;
  padding-top: 0.25rem;
}

.chat-notification-title {
  color: #1a202c;
  font-size: 1.25rem;
  line-height: 1.25;
}

.chat-notification-message {
  color: #718096;
  font-size: 1rem;
  line-height: 1.5;
}
</style>
