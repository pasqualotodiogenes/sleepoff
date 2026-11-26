<script setup lang="ts">
import { Motion, useInView } from "motion-v";
import { computed, onMounted, onUnmounted, ref, watch } from "vue";

interface Props {
  text: string;
  class?: string;
  revealDelayMs?: number;
  charset?: string;
  flipDelayMs?: number;
  encryptedClass?: string;
  revealedClass?: string;
}

const props = withDefaults(defineProps<Props>(), {
  revealDelayMs: 50,
  charset:
    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()_+-={}[];:,.<>/?",
  flipDelayMs: 50,
});

const containerRef = ref<HTMLElement>();
const isInView = useInView(containerRef, { once: true });

const isMounted = ref(false);
const prefersReducedMotion = ref(false);
const revealCount = ref(props.text.length);
const animationFrameId = ref<number | null>(null);
const startTime = ref<number>(0);
const lastFlipTime = ref<number>(0);
const scrambleChars = ref<string[]>([]);

const textArray = computed(() => props.text.split(""));

function generateRandomCharacter(charset: string): string {
  const index = Math.floor(Math.random() * charset.length);
  return charset.charAt(index);
}

function generateGibberishPreservingSpaces(original: string, charset: string): string {
  if (!original) return "";
  let result = "";
  for (let i = 0; i < original.length; i += 1) {
    const ch = original[i];
    result += ch === " " ? " " : generateRandomCharacter(charset);
  }
  return result;
}

function getCharClass(index: number): string {
  const isRevealed = !isMounted.value || prefersReducedMotion.value || index < revealCount.value;
  const baseClass = props.class || "";
  const stateClass = isRevealed ? props.revealedClass : props.encryptedClass;
  return [baseClass, stateClass].filter(Boolean).join(" ");
}

function displayChar(index: number): string {
  const char = props.text[index];
  if (!isMounted.value || prefersReducedMotion.value || index < revealCount.value) {
    return char || "";
  }
  if (char === " ") {
    return " ";
  }
  return scrambleChars.value[index] ?? char ?? "";
}

function resetAnimation() {
  const initial = props.text ? generateGibberishPreservingSpaces(props.text, props.charset) : "";
  scrambleChars.value = initial.split("");
  startTime.value = performance.now();
  lastFlipTime.value = startTime.value;
  revealCount.value = 0;
}

function updateAnimation(now: number) {
  const elapsedMs = now - startTime.value;
  const totalLength = props.text.length;
  const currentRevealCount = Math.min(
    totalLength,
    Math.floor(elapsedMs / Math.max(1, props.revealDelayMs)),
  );

  revealCount.value = currentRevealCount;

  if (currentRevealCount >= totalLength) {
    animationFrameId.value = null;
    return;
  }

  const timeSinceLastFlip = now - lastFlipTime.value;
  if (timeSinceLastFlip >= Math.max(0, props.flipDelayMs)) {
    const newScrambleChars = [...scrambleChars.value];
    for (let index = 0; index < totalLength; index += 1) {
      if (index >= currentRevealCount) {
        newScrambleChars[index] =
          props.text[index] !== " " ? generateRandomCharacter(props.charset) : " ";
      }
    }
    scrambleChars.value = newScrambleChars;
    lastFlipTime.value = now;
  }

  animationFrameId.value = requestAnimationFrame(updateAnimation);
}

function startAnimation() {
  if (!props.text || !isMounted.value) return;
  if (prefersReducedMotion.value) {
    revealCount.value = props.text.length;
    return;
  }

  stopAnimation();
  resetAnimation();
  animationFrameId.value = requestAnimationFrame(updateAnimation);
}

function stopAnimation() {
  if (animationFrameId.value !== null) {
    cancelAnimationFrame(animationFrameId.value);
    animationFrameId.value = null;
  }
}

let mediaQuery: MediaQueryList | null = null;
function handleMotionChange(e: MediaQueryListEvent) {
  prefersReducedMotion.value = e.matches;
  if (e.matches) {
    stopAnimation();
    revealCount.value = props.text.length;
  }
}

onMounted(() => {
  mediaQuery = window.matchMedia("(prefers-reduced-motion: reduce)");
  prefersReducedMotion.value = mediaQuery.matches;
  mediaQuery.addEventListener("change", handleMotionChange);
  isMounted.value = true;

  if (!prefersReducedMotion.value && isInView.value) {
    startAnimation();
  } else {
    revealCount.value = props.text.length;
  }
});

watch(isInView, (newVal) => {
  if (!isMounted.value) return;
  if (newVal) {
    startAnimation();
  }
});

watch(
  () => props.text,
  () => {
    if (!isMounted.value || prefersReducedMotion.value) {
      revealCount.value = props.text.length;
      return;
    }
    if (isInView.value) {
      startAnimation();
    }
  },
);

onUnmounted(() => {
  stopAnimation();
  if (mediaQuery) {
    mediaQuery.removeEventListener("change", handleMotionChange);
  }
});
</script>

<template>
  <Motion
    ref="containerRef"
    as="span"
    :class="props.class"
    :aria-label="text"
    role="text"
  >
    <template
      v-for="(_, index) in textArray"
      :key="index"
    >
      <span :class="getCharClass(index)">{{ displayChar(index) }}</span>
    </template>
  </Motion>
</template>
