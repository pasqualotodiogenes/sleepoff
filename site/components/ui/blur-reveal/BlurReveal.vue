<script setup lang="ts">
import { Motion } from "motion-v";
import { computed, onMounted, onUnmounted, ref } from "vue";

interface Props {
  duration?: number;
  delay?: number;
  blur?: string;
  yOffset?: number;
  once?: boolean;
  class?: string;
}

const props = withDefaults(defineProps<Props>(), {
  duration: 0.45,
  delay: 0.05,
  blur: "8px",
  yOffset: 12,
  once: true,
  class: "",
});

const prefersReducedMotion = ref(false);
let mediaQuery: MediaQueryList | null = null;

function handleMotionChange(e: MediaQueryListEvent) {
  prefersReducedMotion.value = e.matches;
}

onMounted(() => {
  mediaQuery = window.matchMedia("(prefers-reduced-motion: reduce)");
  prefersReducedMotion.value = mediaQuery.matches;
  mediaQuery.addEventListener("change", handleMotionChange);
});

onUnmounted(() => {
  if (mediaQuery) {
    mediaQuery.removeEventListener("change", handleMotionChange);
  }
});

const initialState = computed(() => {
  if (prefersReducedMotion.value) {
    return false as const;
  }
  return {
    opacity: 0,
    filter: `blur(${props.blur})`,
    y: props.yOffset,
  };
});

const animateState = computed(() => ({
  opacity: 1,
  filter: "blur(0px)",
  y: 0,
}));

const transitionConfig = computed(() => {
  if (prefersReducedMotion.value) {
    return { duration: 0, delay: 0 };
  }
  return {
    duration: props.duration,
    delay: props.delay,
    ease: "easeOut",
  };
});
</script>

<template>
  <div :class="props.class">
    <Motion
      as="div"
      :initial="initialState"
      :while-in-view="animateState"
      :in-view-options="{ once: props.once }"
      :transition="transitionConfig"
    >
      <slot />
    </Motion>
  </div>
</template>
