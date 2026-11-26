<script setup lang="ts">
import { cn } from "@inspira-ui/plugins";
import { computed } from "vue";

interface Props {
  max?: number;
  value?: number;
  min?: number;
  gaugePrimaryColor?: string;
  gaugeSecondaryColor?: string;
  class?: string;
  circleStrokeWidth?: number;
  showPercentage?: boolean;
  duration?: number;
}

const props = withDefaults(defineProps<Props>(), {
  max: 100,
  min: 0,
  value: 0,
  gaugePrimaryColor: "#5266eb",
  gaugeSecondaryColor: "#1e1e2a",
  circleStrokeWidth: 8,
  showPercentage: true,
  duration: 0.35,
});

const circumference = 2 * Math.PI * 45;
const percentPx = circumference / 100;

const currentPercent = computed(() => {
  const span = props.max - props.min;
  if (span <= 0) return 0;
  const raw = ((props.value - props.min) / span) * 100;
  return Math.max(0, Math.min(100, Math.round(raw)));
});
const percentageInPx = computed(() => `${percentPx}px`);
const durationInSeconds = computed(() => `${props.duration}s`);
</script>

<template>
  <div :class="cn(`progress-circle-base relative size-12 font-mono text-xs font-semibold text-[#ededf3]`, props.class)">
    <svg
      fill="none"
      class="size-full"
      stroke-width="2"
      viewBox="0 0 100 100"
      aria-hidden="true"
    >
      <circle
        cx="50"
        cy="50"
        r="45"
        :stroke-width="circleStrokeWidth"
        stroke-dashoffset="0"
        stroke-linecap="butt"
        stroke-linejoin="miter"
        class="gauge-secondary-stroke opacity-100"
      />
      <circle
        cx="50"
        cy="50"
        r="45"
        :stroke-width="circleStrokeWidth"
        stroke-dashoffset="0"
        stroke-linecap="butt"
        stroke-linejoin="miter"
        class="gauge-primary-stroke opacity-100"
      />
    </svg>
    <span
      v-if="showPercentage"
      :data-current-value="currentPercent"
      class="absolute inset-0 m-auto size-fit tabular-nums"
    >
      {{ currentPercent }}%
    </span>
  </div>
</template>

<style scoped lang="css">
.progress-circle-base {
  --circle-size: 100px;
  --circumference: v-bind(circumference);
  --percent-to-px: v-bind(percentageInPx);
  transform: translateZ(0);
}

.gauge-primary-stroke {
  stroke: v-bind(gaugePrimaryColor);
  --stroke-percent: v-bind(currentPercent);
  stroke-dasharray: calc(var(--stroke-percent) * var(--percent-to-px)) var(--circumference);
  transition: stroke-dasharray v-bind(durationInSeconds) ease, stroke v-bind(durationInSeconds) ease;
  transform: rotate(-90deg);
  transform-origin: calc(var(--circle-size) / 2) calc(var(--circle-size) / 2);
}

.gauge-secondary-stroke {
  stroke: v-bind(gaugeSecondaryColor);
  stroke-dasharray: var(--circumference) var(--circumference);
}

@media (prefers-reduced-motion: reduce) {
  .gauge-primary-stroke {
    transition: none;
  }
}
</style>
