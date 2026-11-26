<script lang="ts" setup>
import { cn } from "@inspira-ui/plugins";
import { computed, onMounted, onUnmounted, ref } from "vue";
import ScrollWord from "./ScrollWord.vue";

interface Props {
  class?: string;
  text: string;
}

const props = defineProps<Props>();

const textScrollRevealRef = ref<HTMLElement | null>(null);
const words = computed(() => props.text.split(" "));
const scrollYProgress = ref(0);
const prefersReducedMotion = ref(false);

function updateScrollYProgress() {
  if (prefersReducedMotion.value) {
    scrollYProgress.value = 1;
    return;
  }
  if (textScrollRevealRef.value) {
    const boundingRect = textScrollRevealRef.value.getBoundingClientRect();
    const windowHeight = window.innerHeight || 1;
    const raw = (boundingRect.y / windowHeight) * -1;
    scrollYProgress.value = Math.max(0, Math.min(1, raw));
  }
}

let mediaQuery: MediaQueryList | null = null;
function handleMotionChange(e: MediaQueryListEvent) {
  prefersReducedMotion.value = e.matches;
  updateScrollYProgress();
}

onMounted(() => {
  mediaQuery = window.matchMedia("(prefers-reduced-motion: reduce)");
  prefersReducedMotion.value = mediaQuery.matches;
  mediaQuery.addEventListener("change", handleMotionChange);

  window.addEventListener("scroll", updateScrollYProgress, { passive: true });
  window.addEventListener("resize", updateScrollYProgress, { passive: true });
  updateScrollYProgress();
});

onUnmounted(() => {
  window.removeEventListener("scroll", updateScrollYProgress);
  window.removeEventListener("resize", updateScrollYProgress);
  if (mediaQuery) {
    mediaQuery.removeEventListener("change", handleMotionChange);
  }
});
</script>

<template>
  <div
    ref="textScrollRevealRef"
    :class="cn(`relative z-0 h-[140vh]`, $props.class)"
  >
    <div class="sticky top-0 mx-auto flex h-1/2 max-w-4xl items-center bg-transparent px-4 py-12">
      <p
        class="font-display flex flex-wrap text-2xl font-bold text-[#70707d]/20 md:text-3xl lg:text-4xl xl:text-5xl"
      >
        <ScrollWord
          v-for="(word, i) in words"
          :key="i"
          :word="word"
          :progress="scrollYProgress"
          :range="[i / Math.max(1, words.length), (i + 1) / Math.max(1, words.length)]"
        />
      </p>
    </div>
  </div>
</template>
