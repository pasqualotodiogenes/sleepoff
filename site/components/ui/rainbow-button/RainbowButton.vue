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
        `rainbow-button group relative inline-flex h-11 cursor-pointer items-center justify-center rounded-xl border-0 bg-[linear-gradient(#0f172a,#0f172a),linear-gradient(#0f172a_50%,rgba(15,23,42,0.6)_80%,rgba(15,23,42,0)),linear-gradient(90deg,var(--color-1),var(--color-5),var(--color-3),var(--color-4),var(--color-2))] bg-size-[200%] [background-clip:padding-box,border-box,border-box] bg-origin-border px-8 py-2 font-semibold text-white transition-all duration-300 [border:calc(0.09*1rem)_solid_transparent] hover:scale-105 active:scale-95 before:absolute before:bottom-[-20%] before:left-1/2 before:z-0 before:h-1/5 before:w-3/5 before:-translate-x-1/2 before:bg-[linear-gradient(90deg,var(--color-1),var(--color-5),var(--color-3),var(--color-4),var(--color-2))] before:bg-size-[200%] before:filter-[blur(calc(0.8*1rem))] focus-visible:ring-1 focus-visible:ring-violet-400 focus-visible:outline-none disabled:pointer-events-none disabled:opacity-50`,
        props.class,
      )
    "
  >
    <slot />
  </component>
</template>

<style scoped>
.rainbow-button {
  --color-1: hsl(265 89% 66%); /* Violet #7C3AED */
  --color-2: hsl(187 92% 43%); /* Cyan #06B6D4 */
  --color-3: hsl(239 84% 67%); /* Indigo #6366F1 */
  --color-4: hsl(171 78% 46%); /* Teal */
  --color-5: hsl(280 87% 65%); /* Purple */
  --speed: v-bind(speedInSeconds);
  animation: rainbow var(--speed) infinite linear;
}

.rainbow-button:before {
  animation: rainbow var(--speed) infinite linear;
}

@keyframes rainbow {
  0% {
    background-position: 0;
  }
  100% {
    background-position: 200%;
  }
}
</style>
