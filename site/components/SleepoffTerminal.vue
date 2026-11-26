<script setup lang="ts">
import { useSleepoffEmulator } from '../composables/useSleepoffEmulator'
import TerminalHeader from './terminal/TerminalHeader.vue'
import TerminalMissionBar from './terminal/TerminalMissionBar.vue'
import TerminalScreen from './terminal/TerminalScreen.vue'

const {
  isPanicActive,
  isEmulatorOpen,
  isEmulatorFocused,
  isMuted,
  missions,
  mode,
  tuiState,
  isDryRun,
  shellInput,
  shellInputRef,
  terminalContainerRef,
  terminalBodyRef,
  shellHistory,
  selectedMenuIdx,
  customMinutes,
  totalSeconds,
  remainingSeconds,
  isPaused,
  eventsLog,
  panicCountdown,
  cancelPending,
  startTimeStr,
  finishTimeStr,
  isAudioPlaying,
  activeStartingMissionId,
  completedMissionsCount,
  allMissionsCompleted,
  headerStatus,
  progressPercent,
  cancelSummaryStats,
  focusTerminal,
  toggleConsoleOpen,
  toggleConsoleWithShortcut,
  switchTab,
  toggleMute,
  startTui,
  triggerPanicDemo,
  togglePause,
  addFiveMinutes,
  subtractFiveMinutes,
  handleCancel,
  exitTuiToShell,
  handleShellCommand,
  runQuickCommand,
  selectMenuOption,
  submitCustomMinutes,
  clickMission,
  resetMissions,
} = useSleepoffEmulator()

function registerElements(els: { shellInput: HTMLInputElement | null; body: HTMLDivElement | null }) {
  shellInputRef.value = els.shellInput
  terminalBodyRef.value = els.body
}
</script>

<template>
  <div
    id="sleepoff-tactical-console"
    ref="terminalContainerRef"
    @click="focusTerminal"
    class="relative w-full max-w-5xl mx-auto rounded-none bg-[#0e1016] border transition-colors flex flex-col overflow-hidden"
    :class="
      isPanicActive
        ? 'border-[#ef4444]'
        : isEmulatorFocused
          ? 'border-[#5266eb]'
          : 'border-[#1e1e2a]'
    "
  >
    <!-- Tactical Header (1px grid cells, no macOS chrome) -->
    <TerminalHeader
      :mode="mode"
      :tui-state="tuiState"
      :is-dry-run="isDryRun"
      :is-muted="isMuted"
      :is-audio-playing="isAudioPlaying"
      :is-focused="isEmulatorFocused"
      :is-open="isEmulatorOpen"
      :header-status="headerStatus"
      @switch-tab="switchTab"
      @toggle-mute="toggleMute"
      @toggle-focus="toggleConsoleWithShortcut"
      @exit-tui="exitTuiToShell"
    />

    <!-- Accessible Framing Strip (for technical & non-technical visitors alike) -->
    <div
      class="px-3 sm:px-4 py-2 bg-[#171721] border-b border-[#1e1e2a] flex flex-wrap items-center justify-between gap-2 text-xs"
    >
      <p class="text-[#c3c3cc] leading-snug">
        Experimente ao vivo em modo seguro — nenhum comando real será executado no seu PC.
      </p>
      <div class="flex flex-wrap items-center gap-2.5">
        <span class="font-mono text-[11px] text-[#70707d]">
          App instalável no Windows ou terminal em qualquer sistema
        </span>
        <button
          v-if="isEmulatorOpen"
          type="button"
          @click.stop="toggleConsoleOpen"
          :aria-expanded="isEmulatorOpen"
          aria-label="Recolher console tático"
          class="terminal-collapse-btn min-h-[44px] min-w-[44px] px-3 py-1.5 bg-[#0e1016] text-[#c3c3cc] border border-[#272735] font-mono text-[11px] uppercase tracking-wider transition-colors cursor-pointer"
        >
          Recolher
        </button>
      </div>
    </div>

    <!-- Collapsible Body & Missions -->
    <template v-if="isEmulatorOpen">
      <TerminalScreen
        :mode="mode"
        :tui-state="tuiState"
        :is-dry-run="isDryRun"
        :is-focused="isEmulatorFocused"
        :shell-input="shellInput"
        :shell-history="shellHistory"
        :selected-menu-idx="selectedMenuIdx"
        :custom-minutes="customMinutes"
        :total-seconds="totalSeconds"
        :remaining-seconds="remainingSeconds"
        :is-paused="isPaused"
        :events-log="eventsLog"
        :panic-countdown="panicCountdown"
        :cancel-pending="cancelPending"
        :start-time-str="startTimeStr"
        :finish-time-str="finishTimeStr"
        :progress-percent="progressPercent"
        :cancel-summary-stats="cancelSummaryStats"
        @focus="focusTerminal"
        @update:shell-input="shellInput = $event"
        @update:custom-minutes="customMinutes = $event"
        @submit-shell="handleShellCommand"
        @quick-command="runQuickCommand"
        @select-menu="selectMenuOption"
        @set-tui-state="tuiState = $event"
        @submit-custom="submitCustomMinutes"
        @toggle-pause="togglePause"
        @add-five="addFiveMinutes"
        @subtract-five="subtractFiveMinutes"
        @cancel="handleCancel"
        @trigger-panic="triggerPanicDemo"
        @exit-tui="exitTuiToShell"
        @start-15m="startTui(900, true, true)"
        @register-elements="registerElements"
      />

      <TerminalMissionBar
        :missions="missions"
        :active-starting-mission-id="activeStartingMissionId"
        :completed-count="completedMissionsCount"
        :all-completed="allMissionsCompleted"
        :tui-state="tuiState"
        @click-mission="clickMission"
        @trigger-panic="triggerPanicDemo"
        @reset-missions="resetMissions"
      />
    </template>

    <div
      v-else
      class="px-4 py-6 bg-[#06070a] flex flex-wrap items-center justify-between gap-3 font-mono text-xs text-[#c3c3cc]"
    >
      <span>Console tático recolhido. Pressione <kbd class="px-1.5 py-0.5">~</kbd> para reabrir.</span>
      <button
        type="button"
        @click.stop="toggleConsoleWithShortcut"
        class="min-h-[44px] px-4 py-2 bg-[#5266eb] text-[#ededf3] font-bold cursor-pointer"
      >
        [~] Abrir Console
      </button>
    </div>
  </div>
</template>

<style scoped>
@media (hover: hover) {
  .terminal-collapse-btn:hover {
    background-color: #1e1e2a;
    color: #ededf3;
    border-color: #5266eb;
  }
}
</style>
