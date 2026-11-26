<script setup lang="ts">
import type { HTMLAttributes } from "vue";
import { cn } from "@inspira-ui/plugins";
import { computed, onMounted, onUnmounted, ref } from "vue";

const props = withDefaults(
  defineProps<{
    class?: HTMLAttributes["class"];
    slotClass?: HTMLAttributes["class"];
    gradientSize?: number;
    gradientColor?: string;
    gradientOpacity?: number;
  }>(),
  {
    class: "",
    slotClass: "",
    gradientSize: 220,
    gradientColor: "#5266eb",
    gradientOpacity: 0.14,
  },
);

const mouseX = ref(-props.gradientSize * 10);
const mouseY = ref(-props.gradientSize * 10);
const canHover = ref(false);
let hoverQuery: MediaQueryList | null = null;

function handleHoverChange(e: MediaQueryListEvent) {
  canHover.value = e.matches;
  if (!e.matches) {
    handleMouseLeave();
  }
}

function handleMouseMove(e: MouseEvent) {
  if (!canHover.value) return;
  const target = e.currentTarget as HTMLElement;
  const rect = target.getBoundingClientRect();
  mouseX.value = e.clientX - rect.left;
  mouseY.value = e.clientY - rect.top;
}

function handleMouseLeave() {
  mouseX.value = -props.gradientSize * 10;
  mouseY.value = -props.gradientSize * 10;
}

onMounted(() => {
  hoverQuery = window.matchMedia("(hover: hover)");
  canHover.value = hoverQuery.matches;
  hoverQuery.addEventListener("change", handleHoverChange);
  mouseX.value = -props.gradientSize * 10;
  mouseY.value = -props.gradientSize * 10;
});

onUnmounted(() => {
  if (hoverQuery) {
    hoverQuery.removeEventListener("change", handleHoverChange);
  }
});

const backgroundStyle = computed(() => {
  return `radial-gradient(
    ${props.gradientSize}px circle at ${mouseX.value}px ${mouseY.value}px,
    ${props.gradientColor} 0%,
    rgba(0, 0, 0, 0) 70%
  )`;
});
</script>

<template>
  <div
    :class="
      cn(
        `spotlight-card group relative flex size-full w-full min-w-0 overflow-hidden rounded-none border border-[#1e1e2a] bg-[#0e1016] text-[#ededf3]`,
        $props.class,
      )
    "
    @mousemove="handleMouseMove"
    @mouseleave="handleMouseLeave"
  >
    <div :class="cn(`relative z-10 w-full min-w-0`, props.slotClass)">
      <slot />
    </div>
    <div
      v-if="canHover"
      class="spotlight-overlay pointer-events-none absolute inset-0 rounded-none opacity-0 transition-opacity duration-300"
      :style="{
        background: backgroundStyle,
        '--spotlight-opacity': gradientOpacity,
      }"
      aria-hidden="true"
    />
  </div>
</template>

<style scoped>
@media (hover: hover) {
  .spotlight-card:hover .spotlight-overlay {
    opacity: var(--spotlight-opacity, 0.14);
  }
}
</style>
