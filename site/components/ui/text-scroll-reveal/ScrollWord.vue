<script lang="ts" setup>
import { computed } from "vue";

interface Props {
  word: string;
  progress: number;
  range: Array<number>;
}

const props = defineProps<Props>();

const computedOpacity = computed(() => {
  const [start = 0, end = 1] = props.range;
  const progress = props.progress;

  if (progress <= start) return 0;
  if (progress >= end) return 1;

  const span = Math.max(0.0001, end - start);
  return (progress - start) / span;
});
</script>

<template>
  <span class="relative mx-1 inline-block lg:mx-2">
    <span class="text-[#70707d] opacity-20 select-none" aria-hidden="true">{{ word }}</span>
    <span
      :style="{ opacity: computedOpacity }"
      class="absolute inset-0 text-[#ededf3]"
    >
      {{ word }}
    </span>
  </span>
</template>
