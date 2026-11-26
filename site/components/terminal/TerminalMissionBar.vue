<script setup lang="ts">
import type { EmulatorMission, TuiState } from '../../composables/useSleepoffEmulator'

interface Props {
  missions: EmulatorMission[]
  activeStartingMissionId: EmulatorMission['id'] | null
  completedCount: number
  allCompleted: boolean
  tuiState: TuiState
}

defineProps<Props>()

const emit = defineEmits<{
  (e: 'click-mission', mission: EmulatorMission): void
  (e: 'trigger-panic'): void
  (e: 'reset-missions'): void
}>()
</script>

<template>
  <section aria-label="Missões interativas do emulador" class="border-t border-[#1e1e2a] bg-[#0e1016]">
    <!-- Top Bar -->
    <div class="px-3 sm:px-4 py-2.5 flex flex-wrap items-center justify-between gap-2 bg-[#171721] border-b border-[#1e1e2a]">
      <div class="flex items-center gap-2.5 min-w-0">
        <span class="font-mono-hud text-[11px] font-bold text-[#ededf3]">
          Fluxo Guiado
        </span>
        <span class="text-[#70707d] text-xs hidden sm:inline">—</span>
        <span class="text-xs text-[#c3c3cc] hidden sm:inline">
          Siga 1 → 2 → 3 → Simular Pânico → 4 sem precisar reiniciar
        </span>
      </div>

      <div class="flex flex-wrap items-center gap-2 ml-auto">
        <!-- Tactile Panic Trigger -->
        <button
          type="button"
          @click.stop="emit('trigger-panic')"
          :disabled="tuiState === 'CONFIRMATION'"
          class="min-h-[44px] min-w-[44px] px-3 py-2 font-mono text-xs font-bold border transition-colors cursor-pointer disabled:opacity-60 disabled:cursor-default"
          :class="
            tuiState === 'CONFIRMATION'
              ? 'bg-[#ef4444] text-[#ededf3] border-[#ef4444]'
              : 'bg-[#0e1016] hover:bg-[#1e1e2a] text-[#ededf3] border-[#5266eb]'
          "
        >
          Simular Pânico (10s)
        </button>

        <!-- Progress indicator -->
        <div class="flex items-center gap-2 px-2.5 py-2 bg-[#0e1016] border border-[#1e1e2a] min-h-[44px]">
          <div class="w-16 h-1.5 bg-[#171721] border border-[#272735] overflow-hidden">
            <div
              class="h-full transition-all duration-300"
              :class="allCompleted ? 'bg-[#22c55e]' : 'bg-[#5266eb]'"
              :style="{ width: `${(completedCount / Math.max(1, missions.length)) * 100}%` }"
            />
          </div>
          <span
            class="font-mono text-[11px] font-bold whitespace-nowrap"
            :class="allCompleted ? 'text-[#22c55e]' : 'text-[#c3c3cc]'"
          >
            {{ allCompleted ? '✓ FEITO (4/4)' : `${completedCount}/${missions.length}` }}
          </span>
        </div>

        <button
          v-if="completedCount > 0"
          type="button"
          @click.stop="emit('reset-missions')"
          aria-label="Reiniciar missões do emulador"
          class="min-h-[44px] min-w-[44px] px-3 py-2 bg-[#0e1016] hover:bg-[#1e1e2a] text-[#70707d] hover:text-[#ededf3] border border-[#272735] font-mono text-xs cursor-pointer"
        >
          Reiniciar
        </button>
      </div>
    </div>

    <!-- 1px Grid Mission Cells -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-px bg-[#1e1e2a]">
      <button
        v-for="mission in missions"
        :key="mission.id"
        type="button"
        @click.stop="emit('click-mission', mission)"
        class="min-h-[64px] p-3 text-left transition-colors cursor-pointer flex items-start gap-2.5 bg-[#0e1016] hover:bg-[#171721]"
        :class="
          mission.id === activeStartingMissionId && !mission.completed
            ? 'border-l-2 border-l-[#5266eb] bg-[#171721]'
            : ''
        "
      >
        <span
          class="mt-0.5 shrink-0 w-4 h-4 flex items-center justify-center border font-mono text-[10px] font-bold"
          :class="
            mission.completed
              ? 'border-[#22c55e] text-[#22c55e] bg-[#06070a]'
              : mission.id === activeStartingMissionId
                ? 'border-[#5266eb] text-[#5266eb] bg-[#06070a]'
                : 'border-[#272735] text-[#70707d] bg-[#06070a]'
          "
        >
          {{ mission.completed ? '✓' : '›' }}
        </span>

        <div class="flex-1 min-w-0">
          <div class="flex items-center justify-between gap-1.5">
            <span
              class="font-mono text-xs font-bold truncate"
              :class="mission.completed ? 'text-[#22c55e]' : 'text-[#ededf3]'"
            >
              {{ mission.title }}
            </span>
            <span
              v-if="mission.completed"
              class="text-[10px] font-mono font-bold uppercase text-[#22c55e] shrink-0"
            >
              ✓ FEITO
            </span>
            <span
              v-else-if="mission.id === activeStartingMissionId"
              class="text-[10px] font-mono font-bold uppercase text-[#5266eb] shrink-0"
            >
              PRÓXIMO
            </span>
          </div>
          <p class="text-[11px] font-mono text-[#70707d] mt-1 leading-snug">
            {{ mission.desc }}
          </p>
        </div>
      </button>
    </div>
  </section>
</template>
