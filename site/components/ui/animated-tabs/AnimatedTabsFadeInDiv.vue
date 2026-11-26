<script setup lang="ts">
import type { AnimatedTab } from "./types";
import { cn } from "@inspira-ui/plugins";
import { motion } from "motion-v";
import { computed } from "vue";

const props = defineProps<{
  id: string;
  className?: string;
  tabs: AnimatedTab[];
  active: AnimatedTab | null;
  hovering?: boolean;
  reducedMotion?: boolean;
}>();

const activeTab = computed(() => props.active ?? props.tabs[0] ?? null);
</script>

<template>
  <div class="relative w-full min-w-0">
    <motion.div
      v-if="activeTab"
      :key="activeTab.value"
      :initial="reducedMotion ? false : { opacity: 0, y: 6 }"
      :animate="{ opacity: 1, y: 0 }"
      :transition="reducedMotion ? { duration: 0 } : { duration: 0.2, ease: 'easeOut' }"
      :class="cn('w-full min-w-0', className)"
    >
      <slot
        name="content"
        :tab="activeTab"
        :active="activeTab"
        :hovering="hovering"
      >
        <component
          :is="activeTab.content"
          v-if="activeTab.content"
          v-bind="activeTab.contentProps || {}"
        />
      </slot>
    </motion.div>
  </div>
</template>
