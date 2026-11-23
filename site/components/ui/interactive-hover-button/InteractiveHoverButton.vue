<script lang="ts" setup>
import { cn } from "@inspira-ui/plugins";
import { ref } from "vue";

interface Props {
  text?: string;
  class?: string;
  dotColor?: string;
  activeTextColor?: string;
}

const props = withDefaults(defineProps<Props>(), {
  text: "Button",
  dotColor: "bg-cyan-500",
  activeTextColor: "text-slate-950",
});

const buttonRef = ref<HTMLButtonElement>();
</script>

<template>
  <button
    ref="buttonRef"
    :class="
      cn(
        `group relative w-auto cursor-pointer overflow-hidden rounded-xl border border-slate-700/80 bg-slate-900/90 py-2.5 px-6 text-center font-semibold text-slate-200 transition-all duration-300 hover:border-cyan-500/60 hover:shadow-lg hover:shadow-cyan-950/40 min-h-[44px] flex items-center justify-center`,
        props.class,
      )
    "
  >
    <!-- Dot expanding background on hover -->
    <div class="flex items-center gap-2.5">
      <div
        :class="
          cn(
            `size-2 scale-100 rounded-full transition-all duration-500 ease-out group-hover:scale-[120]`,
            props.dotColor,
          )
        "
      />
      <span
        class="inline-flex items-center gap-2 whitespace-nowrap transition-all duration-300 group-hover:translate-x-12 group-hover:opacity-0"
      >
        <slot name="icon" />
        {{ text }}
      </span>
    </div>

    <!-- Active content revealed on hover -->
    <div
      :class="
        cn(
          `absolute top-0 z-10 flex size-full translate-x-12 items-center justify-center gap-2 font-bold opacity-0 transition-all duration-300 group-hover:translate-x-0 group-hover:opacity-100`,
          props.activeTextColor,
        )
      "
    >
      <slot name="active-icon" />
      <span class="whitespace-nowrap">{{ text }}</span>
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width="18"
        height="18"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2.5"
        stroke-linecap="round"
        stroke-linejoin="round"
        class="transition-transform duration-300 group-hover:translate-x-1"
      >
        <path d="M5 12h14" />
        <path d="m12 5 7 7-7 7" />
      </svg>
    </div>
  </button>
</template>

<style scoped></style>

