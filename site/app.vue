<script setup lang="ts">
import { NuxtPage } from '#components'
</script>

<template>
  <div
    id="top"
    class="relative min-h-screen overflow-x-clip bg-[#06070a] text-[#ededf3] selection:bg-[#5266eb] selection:text-[#ededf3]"
  >
    <div
      class="film-grain pointer-events-none fixed inset-0 z-50"
      aria-hidden="true"
    />
    <NuxtPage />
  </div>
</template>
