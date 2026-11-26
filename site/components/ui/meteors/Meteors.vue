<script lang="ts" setup>
import { cn } from "@inspira-ui/plugins";
import { computed } from "vue";

interface Props {
  count?: number;
  number?: number;
  class?: string;
}

const props = withDefaults(defineProps<Props>(), {
  count: 20,
  number: undefined,
  class: "",
});

const meteorCount = computed(() => props.number ?? props.count ?? 20);

function deterministicOffset(index: number): number {
  return ((index * 137) % 800) - 400;
}

function deterministicDelay(index: number): string {
  return `${(0.2 + ((index * 19) % 60) / 100).toFixed(2)}s`;
}

function deterministicDuration(index: number): string {
  return `${2 + ((index * 3) % 8)}s`;
}
</script>

<template>
  <span
    v-for="index in meteorCount"
    :key="`meteor-${index}`"
    :style="{
      top: 0,
      left: `${deterministicOffset(index)}px`,
      animationDelay: deterministicDelay(index),
      animationDuration: deterministicDuration(index),
    }"
    :class="
      cn(
        `animate-meteor absolute h-0.5 w-0.5 rotate-45 rounded-none bg-[#c3c3cc] before:absolute before:top-1/2 before:h-px before:w-[50px] before:-translate-y-[50%] before:transform before:bg-gradient-to-r before:from-[#5266eb] before:to-transparent before:content-['']`,
        $props.class,
      )
    "
  />
</template>

<style scoped>
@keyframes meteor {
  0% {
    transform: rotate(215deg) translateX(0);
    opacity: 1;
  }
  70% {
    opacity: 1;
  }
  100% {
    transform: rotate(215deg) translateX(-500px);
    opacity: 0;
  }
}

.animate-meteor {
  animation: meteor 5s linear infinite;
}
</style>
