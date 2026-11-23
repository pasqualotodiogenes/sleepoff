<script lang="ts" setup>
import { computed } from "vue";
import { cn } from "@inspira-ui/plugins";
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
</script>

<template>
  <span
    v-for="index in meteorCount"
    :key="`meteor ${index}`"
    :style="{
      top: 0,
      left: `${Math.floor(Math.random() * (400 - -400) + -400)}px`,
      animationDelay: `${Math.random() * (0.8 - 0.2) + 0.2}s`,
      animationDuration: `${Math.floor(Math.random() * (10 - 2) + 2)}s`,
    }"
    :class="
      cn(
        `animate-meteor absolute h-0.5 w-0.5 rotate-45 rounded-[9999px] bg-slate-400 shadow-[0_0_0_1px_#ffffff10] before:absolute before:top-1/2 before:h-px before:w-[50px] before:-translate-y-[50%] before:transform before:bg-gradient-to-r before:from-violet-400 before:to-transparent before:content-['']`,
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
