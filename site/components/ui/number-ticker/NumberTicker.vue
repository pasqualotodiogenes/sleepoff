<script setup lang="ts">
import { cn } from "@inspira-ui/plugins";
import { TransitionPresets, useElementVisibility, useTransition } from "@vueuse/core";
import { computed, onMounted, onUnmounted, ref, watch } from "vue";

type TransitionsPresetsKeys = keyof typeof TransitionPresets;

interface NumberTickerProps {
  value?: number;
  direction?: "up" | "down";
  duration?: number;
  delay?: number;
  decimalPlaces?: number;
  locale?: string;
  class?: string;
  transition?: TransitionsPresetsKeys;
}

const props = withDefaults(defineProps<NumberTickerProps>(), {
  value: 0,
  direction: "up",
  delay: 0,
  duration: 1000,
  decimalPlaces: 2,
  locale: "en-US",
  transition: "easeOutCubic",
});

const spanRef = ref<HTMLSpanElement>();
const prefersReducedMotion = ref(false);

const transitionValue = ref(props.direction === "down" ? props.value : 0);

const transitionOutput = useTransition(transitionValue, {
  delay: computed(() => (prefersReducedMotion.value ? 0 : props.delay)),
  duration: computed(() => (prefersReducedMotion.value ? 0 : props.duration)),
  transition: TransitionPresets[props.transition],
});

const output = computed(() => {
  const numericValue = prefersReducedMotion.value
    ? props.value
    : Number(transitionOutput.value.toFixed(props.decimalPlaces));

  return new Intl.NumberFormat(props.locale, {
    minimumFractionDigits: props.decimalPlaces,
    maximumFractionDigits: props.decimalPlaces,
  }).format(numericValue);
});

const isInView = useElementVisibility(spanRef, {
  threshold: 0,
});

const hasBeenInView = ref(false);

let mediaQuery: MediaQueryList | null = null;
function handleMotionChange(e: MediaQueryListEvent) {
  prefersReducedMotion.value = e.matches;
  if (e.matches) {
    transitionValue.value = props.value;
  }
}

onMounted(() => {
  mediaQuery = window.matchMedia("(prefers-reduced-motion: reduce)");
  prefersReducedMotion.value = mediaQuery.matches;
  mediaQuery.addEventListener("change", handleMotionChange);
  if (prefersReducedMotion.value) {
    hasBeenInView.value = true;
    transitionValue.value = props.value;
  }
});

onUnmounted(() => {
  if (mediaQuery) {
    mediaQuery.removeEventListener("change", handleMotionChange);
  }
});

const stopIsInViewWatcher = watch(
  isInView,
  (isVisible) => {
    if (isVisible && !hasBeenInView.value) {
      hasBeenInView.value = true;
      transitionValue.value =
        prefersReducedMotion.value || props.direction !== "down" ? props.value : 0;
      stopIsInViewWatcher();
    }
  },
  { immediate: true },
);

watch(
  () => props.value,
  (newVal) => {
    if (prefersReducedMotion.value) {
      transitionValue.value = newVal;
      return;
    }
    if (hasBeenInView.value) {
      transitionValue.value = props.direction === "down" ? 0 : newVal;
    }
  },
);
</script>

<template>
  <span
    ref="spanRef"
    :class="cn(`inline-block tracking-wider text-[#ededf3] tabular-nums`, props.class)"
  >
    {{ output }}
  </span>
</template>
