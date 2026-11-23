<script setup lang="ts">
import { cn } from "@inspira-ui/plugins";
import { onBeforeUnmount, onMounted, ref } from "vue";

interface Props {
  words: string[];
  duration?: number;
  class?: string;
}

const props = withDefaults(defineProps<Props>(), {
  duration: 3600,
  class: "",
});

const emit = defineEmits(["animationStart", "animationComplete"]);

const currentIndex = ref(0);
const currentWord = ref(props.words[0] || "");
let intervalTimer: any = null;

function nextWord() {
  if (!props.words.length) return;
  currentIndex.value = (currentIndex.value + 1) % props.words.length;
  currentWord.value = props.words[currentIndex.value] || "";
}

onMounted(() => {
  if (props.words.length > 1) {
    intervalTimer = setInterval(() => {
      nextWord();
    }, props.duration);
  }
});

onBeforeUnmount(() => {
  if (intervalTimer) {
    clearInterval(intervalTimer);
  }
});
</script>

<template>
  <span class="relative inline-block px-1 align-baseline">
    <Transition
      name="flip-word"
      mode="out-in"
      @after-enter="emit('animationStart')"
      @after-leave="emit('animationComplete')"
    >
      <span
        :key="currentWord"
        :class="
          cn(
            'inline-block font-extrabold tracking-tight whitespace-nowrap',
            props.class || 'text-white'
          )
        "
      >
        {{ currentWord }}
      </span>
    </Transition>
  </span>
</template>

<style scoped>
.flip-word-enter-active {
  transition: opacity 0.25s ease-out, transform 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}

.flip-word-leave-active {
  transition: opacity 0.2s ease-in, transform 0.2s cubic-bezier(0.7, 0, 0.84, 0);
}

.flip-word-enter-from {
  opacity: 0;
  transform: translateY(10px) scale(0.96);
}

.flip-word-leave-to {
  opacity: 0;
  transform: translateY(-10px) scale(0.96);
}
</style>
