import { useState } from '#imports'

export type SleepoffSectionId = '01' | '02' | '03' | '04'

export interface QuickCommandRequest {
  id: number
  command: string
}

export function useSleepoffSiteState() {
  const isPanicActive = useState<boolean>('sleepoff:isPanicActive', () => false)
  const panicSecondsRemaining = useState<number>('sleepoff:panicSecondsRemaining', () => 0)
  const isEmulatorOpen = useState<boolean>('sleepoff:isEmulatorOpen', () => true)
  const isEmulatorFocused = useState<boolean>('sleepoff:isEmulatorFocused', () => false)
  const isCommandPaletteOpen = useState<boolean>('sleepoff:isCommandPaletteOpen', () => false)
  const activeSectionId = useState<SleepoffSectionId>('sleepoff:activeSectionId', () => '01')
  const scrollProgress = useState<number>('sleepoff:scrollProgress', () => 0)
  const requestedQuickCommand = useState<QuickCommandRequest | null>(
    'sleepoff:requestedQuickCommand',
    () => null,
  )
  const isMuted = useState<boolean>('sleepoff:isMuted', () => false)
  const quickCommandSequence = useState<number>('sleepoff:quickCommandSeq', () => 0)

  function setPanicState(active: boolean, seconds?: number) {
    isPanicActive.value = active
    if (typeof seconds === 'number' && Number.isFinite(seconds)) {
      panicSecondsRemaining.value = Math.max(0, Math.min(10, Math.round(seconds)))
    } else if (!active) {
      panicSecondsRemaining.value = 0
    }
  }

  function triggerCommandInEmulator(command: string) {
    const normalized = command.trim()
    if (!normalized) return
    quickCommandSequence.value += 1
    isEmulatorOpen.value = true
    isEmulatorFocused.value = true
    requestedQuickCommand.value = {
      id: quickCommandSequence.value,
      command: normalized,
    }
  }

  function toggleEmulatorFocus(force?: boolean) {
    const next = typeof force === 'boolean' ? force : !isEmulatorFocused.value
    if (next) {
      isEmulatorOpen.value = true
    }
    isEmulatorFocused.value = next
  }

  function toggleCommandPalette(force?: boolean) {
    isCommandPaletteOpen.value =
      typeof force === 'boolean' ? force : !isCommandPaletteOpen.value
  }

  return {
    isPanicActive,
    panicSecondsRemaining,
    isEmulatorOpen,
    isEmulatorFocused,
    isCommandPaletteOpen,
    activeSectionId,
    scrollProgress,
    requestedQuickCommand,
    isMuted,
    setPanicState,
    triggerCommandInEmulator,
    toggleEmulatorFocus,
    toggleCommandPalette,
  }
}
