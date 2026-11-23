<script setup lang="ts">
import { cn } from "@inspira-ui/plugins";
import { computed } from "vue";

interface BorderBeamProps {
  class?: string;
  size?: number;
  duration?: number;
  borderWidth?: number;
  borderRadius?: number;
  anchor?: number;
  colorFrom?: string;
  colorTo?: string;
  delay?: number;
}

const props = withDefaults(defineProps<BorderBeamProps>(), {
  size: 260,
  duration: 12,
  anchor: 90,
  borderWidth: 1.5,
  borderRadius: 16,
  colorFrom: "#ffaa40",
  colorTo: "#9c40ff",
  delay: 0,
});

const durationInSeconds = computed(() => {
  const d = props.duration;
  return d > 100 ? `${d / 1000}s` : `${d}s`;
});
const delayInSeconds = computed(() => {
  const d = props.delay;
  return d > 100 ? `${d / 1000}s` : `${d}s`;
});
</script>

<template>
  <div
    :class="
      cn(
        `border-beam animate-border-beam pointer-events-none absolute inset-0 z-20 rounded-[inherit] after:absolute after:aspect-square after:w-[calc(var(--size)*1px)] after:[animation-delay:var(--delay)] after:[background:linear-gradient(to_left,var(--color-from),var(--color-to),transparent)] after:[offset-anchor:calc(var(--anchor)*1%)_50%] after:[offset-path:rect(0_auto_auto_0_round_calc(var(--border-radius)*1px))]`,
        props.class,
      )
    "
  />
</template>

<style scoped>
.border-beam {
  --size: v-bind(size);
  --duration: v-bind(durationInSeconds);
  --anchor: v-bind(anchor);
  --border-width: v-bind(borderWidth);
  --border-radius: v-bind(borderRadius);
  --color-from: v-bind(colorFrom);
  --color-to: v-bind(colorTo);
  --delay: v-bind(delayInSeconds);

  border: calc(var(--border-width) * 1px) solid transparent;
  mask: linear-gradient(#fff 0 0) padding-box, linear-gradient(#fff 0 0) border-box;
  mask-composite: exclude;
  -webkit-mask: linear-gradient(#fff 0 0) padding-box, linear-gradient(#fff 0 0) border-box;
  -webkit-mask-composite: destination-out;
}

.animate-border-beam::after {
  animation: border-beam-anim var(--duration) infinite linear;
}

@keyframes border-beam-anim {
  to {
    offset-distance: 100%;
  }
}
</style>
