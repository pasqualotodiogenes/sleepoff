<script setup lang="ts">
import type { AnimatedTab } from "./types";
import { cn } from "@inspira-ui/plugins";
import { motion } from "motion-v";
import { computed, onMounted, onScopeDispose, ref, useId, watch } from "vue";
import AnimatedTabsFadeInDiv from "./AnimatedTabsFadeInDiv.vue";

const props = withDefaults(
  defineProps<{
    tabs: AnimatedTab[];
    modelValue?: string;
    containerClassName?: string;
    activeTabClassName?: string;
    tabClassName?: string;
    contentClassName?: string;
  }>(),
  {
    containerClassName: "",
    activeTabClassName: "",
    tabClassName: "",
    contentClassName: "",
  },
);

const emit = defineEmits<{
  "update:modelValue": [value: string];
}>();

const isControlled = computed(() => props.modelValue !== undefined);
const uncontrolledValue = ref<string | null>(props.tabs[0]?.value ?? null);
const selectedValue = computed(() =>
  isControlled.value ? props.modelValue : uncontrolledValue.value,
);
const id = useId();
const pillLayoutId = computed(() => `animated-tabs-pill-${id}`);
const prefersReducedMotion = ref(false);

const reorderTabs = (tabs: AnimatedTab[], value: string | null | undefined) => {
  const next = tabs.slice();
  if (!value) return next;

  const idx = next.findIndex((t) => t.value === value);
  if (idx <= 0) return next;

  const selected = next.splice(idx, 1)[0];
  if (!selected) return next;

  next.unshift(selected);
  return next;
};

const stackedTabs = ref<AnimatedTab[]>(reorderTabs(props.tabs, selectedValue.value));
const active = ref<AnimatedTab | null>(stackedTabs.value[0] ?? null);
const pointerHovering = ref(false);
const programmaticHovering = ref(false);
const hovering = computed(() => pointerHovering.value || programmaticHovering.value);

let programmaticHoverTimer: ReturnType<typeof setTimeout> | null = null;
const triggerProgrammaticHover = () => {
  if (pointerHovering.value) return;

  programmaticHovering.value = true;
  if (programmaticHoverTimer) clearTimeout(programmaticHoverTimer);
  programmaticHoverTimer = setTimeout(() => {
    programmaticHoverTimer = null;
    if (!pointerHovering.value) programmaticHovering.value = false;
  }, 240);
};

let mediaQuery: MediaQueryList | null = null;
function handleMotionChange(e: MediaQueryListEvent) {
  prefersReducedMotion.value = e.matches;
}

onMounted(() => {
  mediaQuery = window.matchMedia("(prefers-reduced-motion: reduce)");
  prefersReducedMotion.value = mediaQuery.matches;
  mediaQuery.addEventListener("change", handleMotionChange);
});

onScopeDispose(() => {
  if (programmaticHoverTimer) clearTimeout(programmaticHoverTimer);
  if (mediaQuery) mediaQuery.removeEventListener("change", handleMotionChange);
});

watch(
  [() => props.tabs, () => selectedValue.value],
  ([nextTabs, value]) => {
    const nextStack = reorderTabs(nextTabs, value);
    stackedTabs.value = nextStack;
    active.value = nextStack[0] ?? null;

    if (!isControlled.value) {
      if (!value || !nextTabs.some((t) => t.value === value)) {
        uncontrolledValue.value = nextStack[0]?.value ?? null;
      }
    }
  },
  { deep: false, immediate: true },
);

watch(
  () => props.modelValue,
  (value, oldValue) => {
    if (!isControlled.value) return;
    if (value === undefined || oldValue === undefined || value === oldValue) return;
    triggerProgrammaticHover();
  },
);

const selectTab = (idx: number) => {
  const selected = props.tabs[idx];
  if (!selected) return;

  if (isControlled.value) {
    emit("update:modelValue", selected.value);
    return;
  }

  uncontrolledValue.value = selected.value;
};
</script>

<template>
  <div
    class="flex w-full min-w-0 flex-col"
    @mouseenter="pointerHovering = true"
    @mouseleave="pointerHovering = false"
  >
    <div
      role="tablist"
      :class="
        cn(
          'relative flex w-full max-w-full shrink-0 flex-row items-stretch justify-start overflow-x-auto border border-[#1e1e2a] bg-[#0e1016] [-ms-overflow-style:none] [-webkit-overflow-scrolling:touch] [scrollbar-width:none] [&::-webkit-scrollbar]:hidden',
          containerClassName,
        )
      "
    >
      <button
        v-for="(tab, idx) in tabs"
        :key="tab.value"
        type="button"
        role="tab"
        :aria-selected="active?.value === tab.value"
        :class="
          cn(
            'relative min-h-[44px] shrink-0 cursor-pointer rounded-none border-r border-[#1e1e2a] px-4 py-2.5 text-xs font-mono uppercase tracking-[0.08em] transition-colors last:border-r-0',
            tabClassName,
          )
        "
        @click="selectTab(idx)"
      >
        <motion.div
          v-if="active?.value === tab.value"
          :layout-id="pillLayoutId"
          :transition="
            prefersReducedMotion
              ? { duration: 0 }
              : { type: 'spring', bounce: 0.15, duration: 0.35 }
          "
          :class="
            cn(
              'pointer-events-none absolute inset-0 rounded-none bg-[#5266eb]',
              activeTabClassName,
            )
          "
        />

        <span
          class="relative z-10 block whitespace-nowrap"
          :class="active?.value === tab.value ? 'font-semibold text-[#ededf3]' : 'text-[#c3c3cc]'"
        >
          {{ tab.title }}
        </span>
      </button>
    </div>

    <div class="relative mt-4 min-w-0 flex-1">
      <AnimatedTabsFadeInDiv
        :id="id"
        :tabs="stackedTabs"
        :active="active"
        :hovering="hovering"
        :reduced-motion="prefersReducedMotion"
        :class-name="contentClassName"
      >
        <template #content="{ tab, active: slotActive, hovering: slotHovering }">
          <slot
            name="content"
            :tab="tab"
            :active="slotActive"
            :hovering="slotHovering"
          >
            <component
              :is="tab.content"
              v-if="tab.content"
              v-bind="tab.contentProps || {}"
              :tab="tab"
              :active="slotActive"
              :hovering="slotHovering"
            />
          </slot>
        </template>
      </AnimatedTabsFadeInDiv>
    </div>
  </div>
</template>
