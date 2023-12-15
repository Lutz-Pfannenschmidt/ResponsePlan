<script setup lang="ts">
import SunIcon from './icons/SunIcon.vue';
import MoonIcon from './icons/MoonIcon.vue';
import { onMounted } from 'vue';

function setTheme(theme: "business" | "light") {
  console.log("switch")
  const switcher = document.querySelector("input#theme-switch") as HTMLInputElement
  switcher.checked = theme === "business"
  document.documentElement.dataset.theme = theme
  window.localStorage.setItem("theme", theme)
}

function switchTheme() {
  const curr = document.documentElement.dataset.theme || "business"
  setTheme(curr === "light" ? "business" : "light" )
}

onMounted(() => {
  const switcher = document.querySelector("input#theme-switch") as HTMLInputElement
  switcher.addEventListener("click", switchTheme)

  const localStorageTheme = window.localStorage.getItem("theme")
  const prefersbusinessMode = window.matchMedia && window.matchMedia('(prefers-color-scheme: business)').matches;
  switcher.checked = localStorageTheme ? localStorageTheme === "business" : prefersbusinessMode
  document.documentElement.dataset.theme = switcher.checked ? "business" : "light"
})

</script>

<template>

  <label class="swap swap-rotate px-4">
    <input type="checkbox" id="theme-switch" />
    
    <div class="swap-on w-10 h-10"><SunIcon /></div>
    <div class="swap-off w-10 h-10"><MoonIcon /></div>

  </label>
</template>