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
  dotColor: "bg-[#5266eb]",
  activeTextColor: "text-[#ededf3]",
});

const buttonRef = ref<HTMLButtonElement>();
</script>

<template>
  <button
    ref="buttonRef"
    :class="
      cn(
        `ihb-root group relative flex min-h-[44px] w-auto cursor-pointer items-center justify-center overflow-hidden rounded-none border border-[#272735] bg-[#0e1016] px-6 py-2.5 text-center font-semibold text-[#ededf3] transition-colors duration-200`,
        props.class,
      )
    "
  >
    <!-- Cobalt tactical block/dot expanding only on hover-capable devices -->
    <div class="flex items-center gap-2.5">
      <div
        :class="
          cn(
            `ihb-dot size-2 scale-100 rounded-none transition-all duration-300 ease-out`,
            props.dotColor,
          )
        "
      />
      <span
        class="ihb-idle inline-flex items-center gap-2 whitespace-nowrap transition-all duration-200"
      >
        <slot name="icon" />
        {{ text }}
      </span>
    </div>

    <!-- Active content revealed on hover-capable devices -->
    <div
      :class="
        cn(
          `ihb-active pointer-events-none absolute inset-0 z-10 flex size-full translate-x-12 items-center justify-center gap-2 font-semibold opacity-0 transition-all duration-200`,
          props.activeTextColor,
        )
      "
    >
      <slot name="active-icon" />
      <span class="whitespace-nowrap">{{ text }}</span>
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width="16"
        height="16"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="square"
        stroke-linejoin="miter"
        class="ihb-arrow transition-transform duration-200"
        aria-hidden="true"
      >
        <path d="M5 12h14" />
        <path d="m12 5 7 7-7 7" />
      </svg>
    </div>
  </button>
</template>

<style scoped>
@media (hover: hover) {
  .ihb-root:hover {
    border-color: #5266eb;
  }
  .ihb-root:hover .ihb-dot {
    transform: scale(120);
  }
  .ihb-root:hover .ihb-idle {
    transform: translateX(3rem);
    opacity: 0;
  }
  .ihb-root:hover .ihb-active {
    transform: translateX(0);
    opacity: 1;
  }
  .ihb-root:hover .ihb-arrow {
    transform: translateX(0.25rem);
  }
}
</style>
