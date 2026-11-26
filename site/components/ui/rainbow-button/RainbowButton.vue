<script setup lang="ts">
import { cn } from "@inspira-ui/plugins";
import { computed } from "vue";

interface RainbowButtonProps {
  class?: string;
  is?: string;
  speed?: number;
}

const props = withDefaults(defineProps<RainbowButtonProps>(), {
  speed: 2.5,
  is: "button",
});

const speedInSeconds = computed(() => `${props.speed}s`);
</script>

<template>
  <component
    :is="is"
    :class="
      cn(
        `tactical-cobalt-btn group relative inline-flex min-h-[44px] cursor-pointer items-center justify-center rounded-none border border-[#5266eb] bg-[#0e1016] px-6 py-2.5 font-semibold text-[#ededf3] transition-colors duration-200 focus-visible:ring-1 focus-visible:ring-[#5266eb] focus-visible:outline-none disabled:pointer-events-none disabled:opacity-50`,
        props.class,
      )
    "
  >
    <slot />
  </component>
</template>

<style scoped>
.tactical-cobalt-btn {
  --speed: v-bind(speedInSeconds);
}

@media (hover: hover) {
  .tactical-cobalt-btn:hover {
    background-color: #5266eb;
    color: #ededf3;
  }
}
</style>
