<script setup lang="ts">
import { cn } from "@inspira-ui/plugins";
import { computed, onMounted, onUnmounted, ref, useTemplateRef, watch } from "vue";

export type DitheringMode = "bayer" | "halftone" | "noise" | "crosshatch";
export type ColorMode = "original" | "grayscale" | "duotone" | "custom";

export interface Props {
  /** Source image URL */
  src: string;
  /** Size of the dithering grid cells */
  gridSize?: number;
  /** Type of dithering pattern */
  ditherMode?: DitheringMode;
  /** Color processing mode */
  colorMode?: ColorMode;
  /** Invert the dithered output colors */
  invert?: boolean;
  /** Pixelation multiplier (1 = no pixelation, higher = more pixelated) */
  pixelRatio?: number;
  /** Primary color for duotone mode */
  primaryColor?: string;
  /** Secondary color for duotone mode */
  secondaryColor?: string;
  /** Custom color palette array for custom mode */
  customPalette?: string[];
  /** Brightness adjustment (-1 to 1) */
  brightness?: number;
  /** Contrast adjustment (0 to 2, 1 = normal) */
  contrast?: number;
  /** Background color behind the dithered image */
  backgroundColor?: string;
  /** Object fit behavior */
  objectFit?: "cover" | "contain" | "fill" | "none";
  /** Threshold bias for dithering (0 to 1) */
  threshold?: number;
  /** Enable animation effect */
  animated?: boolean;
  /** Animation speed (lower = slower) */
  animationSpeed?: number;
  /** Additional CSS classes for the container */
  class?: string;
}

const props = withDefaults(defineProps<Props>(), {
  gridSize: 4,
  ditherMode: "bayer",
  colorMode: "original",
  invert: false,
  pixelRatio: 1,
  primaryColor: "#06070a",
  secondaryColor: "#5266eb",
  customPalette: () => ["#06070a", "#5266eb"],
  brightness: 0,
  contrast: 1,
  backgroundColor: "transparent",
  objectFit: "cover",
  threshold: 0.5,
  animated: false,
  animationSpeed: 0.02,
});

const BAYER_MATRIX_4x4 = [
  [0, 8, 2, 10],
  [12, 4, 14, 6],
  [3, 11, 1, 9],
  [15, 7, 13, 5],
];

const BAYER_MATRIX_8x8 = [
  [0, 32, 8, 40, 2, 34, 10, 42],
  [48, 16, 56, 24, 50, 18, 58, 26],
  [12, 44, 4, 36, 14, 46, 6, 38],
  [60, 28, 52, 20, 62, 30, 54, 22],
  [3, 35, 11, 43, 1, 33, 9, 41],
  [51, 19, 59, 27, 49, 17, 57, 25],
  [15, 47, 7, 39, 13, 45, 5, 37],
  [63, 31, 55, 23, 61, 29, 53, 21],
];

const containerRef = useTemplateRef<HTMLDivElement>("containerRef");
const canvasRef = useTemplateRef<HTMLCanvasElement>("canvasRef");
const animationRef = ref<number | null>(null);
const timeRef = ref<number>(0);
const imageRef = ref<HTMLImageElement | null>(null);
const loadedSrc = ref<string>("");
const imageDataRef = ref<ImageData | null>(null);
const dimensions = ref({ width: 0, height: 0 });
const prefersReducedMotion = ref(false);

function parseColor(color: string): [number, number, number] {
  if (color.startsWith("#")) {
    const hex = color.slice(1);
    if (hex.length === 3) {
      return [
        Number.parseInt((hex[0] ?? "0") + (hex[0] ?? "0"), 16),
        Number.parseInt((hex[1] ?? "0") + (hex[1] ?? "0"), 16),
        Number.parseInt((hex[2] ?? "0") + (hex[2] ?? "0"), 16),
      ];
    }
    return [
      Number.parseInt(hex.slice(0, 2), 16),
      Number.parseInt(hex.slice(2, 4), 16),
      Number.parseInt(hex.slice(4, 6), 16),
    ];
  }
  const match = color.match(/rgb\((\d+)\s*,\s*(\d+)\s*,\s*(\d+)\)/i);
  if (match) {
    return [
      Number.parseInt(match[1] ?? "0", 10),
      Number.parseInt(match[2] ?? "0", 10),
      Number.parseInt(match[3] ?? "0", 10),
    ];
  }
  return [0, 0, 0];
}

function getLuminance(r: number, g: number, b: number): number {
  return 0.299 * r + 0.587 * g + 0.114 * b;
}

function clamp(value: number, min: number, max: number): number {
  return Math.max(min, Math.min(max, value));
}

const parsedPrimaryColor = computed(() => parseColor(props.primaryColor));
const parsedSecondaryColor = computed(() => parseColor(props.secondaryColor));
const parsedCustomPalette = computed(() => props.customPalette.map(parseColor));

function applyDithering(
  ctx: CanvasRenderingContext2D,
  displayWidth: number,
  displayHeight: number,
  time = 0,
) {
  const canvas = canvasRef.value;
  if (!canvas || !imageDataRef.value) return;

  if (props.backgroundColor !== "transparent") {
    ctx.fillStyle = props.backgroundColor;
    ctx.fillRect(0, 0, displayWidth, displayHeight);
  } else {
    ctx.clearRect(0, 0, displayWidth, displayHeight);
  }

  const sourceData = imageDataRef.value.data;
  const sourceWidth = imageDataRef.value.width;
  const sourceHeight = imageDataRef.value.height;

  const effectivePixelSize = Math.max(1, Math.floor(props.gridSize * props.pixelRatio));
  const matrixSize = props.gridSize <= 4 ? 4 : 8;
  const bayerMatrix = props.gridSize <= 4 ? BAYER_MATRIX_4x4 : BAYER_MATRIX_8x8;
  const matrixScale = matrixSize === 4 ? 16 : 64;

  for (let y = 0; y < displayHeight; y += effectivePixelSize) {
    for (let x = 0; x < displayWidth; x += effectivePixelSize) {
      const srcX = Math.floor((x / displayWidth) * sourceWidth);
      const srcY = Math.floor((y / displayHeight) * sourceHeight);
      const srcIdx = (srcY * sourceWidth + srcX) * 4;

      let r = sourceData[srcIdx] || 0;
      let g = sourceData[srcIdx + 1] || 0;
      let b = sourceData[srcIdx + 2] || 0;
      const a = sourceData[srcIdx + 3] || 0;

      if (a < 10) continue;

      r = clamp((r - 128) * props.contrast + 128 + props.brightness * 255, 0, 255);
      g = clamp((g - 128) * props.contrast + 128 + props.brightness * 255, 0, 255);
      b = clamp((b - 128) * props.contrast + 128 + props.brightness * 255, 0, 255);

      const luminance = getLuminance(r, g, b) / 255;

      let ditherThreshold: number;
      const matrixX = Math.floor(x / props.gridSize) % matrixSize;
      const matrixY = Math.floor(y / props.gridSize) % matrixSize;

      switch (props.ditherMode) {
        case "bayer":
          ditherThreshold = (bayerMatrix[matrixY]?.[matrixX] ?? 0) / matrixScale;
          break;
        case "halftone": {
          const angle = Math.PI / 4;
          const scale = props.gridSize * 2;
          const rotX = x * Math.cos(angle) + y * Math.sin(angle);
          const rotY = -x * Math.sin(angle) + y * Math.cos(angle);
          ditherThreshold = (Math.sin(rotX / scale) + Math.sin(rotY / scale) + 2) / 4;
          break;
        }
        case "noise": {
          const noiseVal = Math.sin(x * 12.9898 + y * 78.233 + time * 100) * 43758.5453;
          ditherThreshold = noiseVal - Math.floor(noiseVal);
          break;
        }
        case "crosshatch": {
          const line1 = (x + y) % (props.gridSize * 2) < props.gridSize ? 1 : 0;
          const line2 =
            (x - y + props.gridSize * 4) % (props.gridSize * 2) < props.gridSize ? 1 : 0;
          ditherThreshold = (line1 + line2) / 2;
          break;
        }
        default:
          ditherThreshold = (bayerMatrix[matrixY]?.[matrixX] ?? 0) / matrixScale;
      }

      ditherThreshold = ditherThreshold * (1 - props.threshold) + props.threshold * 0.5;

      let outputColor: [number, number, number];

      switch (props.colorMode) {
        case "grayscale": {
          const shouldBeDark = luminance < ditherThreshold;
          outputColor = shouldBeDark ? [0, 0, 0] : [255, 255, 255];
          break;
        }
        case "duotone": {
          const shouldBeDark = luminance < ditherThreshold;
          outputColor = shouldBeDark ? parsedPrimaryColor.value : parsedSecondaryColor.value;
          break;
        }
        case "custom": {
          const palette = parsedCustomPalette.value;
          if (palette.length <= 2) {
            const shouldBeDark = luminance < ditherThreshold;
            outputColor = shouldBeDark ? (palette[0] ?? [0, 0, 0]) : (palette[1] ?? [255, 255, 255]);
          } else {
            const adjustedLuminance = luminance + (ditherThreshold - 0.5) * 0.5;
            const paletteIndex = Math.floor(
              clamp(adjustedLuminance, 0, 1) * (palette.length - 1),
            );
            outputColor = palette[paletteIndex] ?? [0, 0, 0];
          }
          break;
        }
        case "original":
        default: {
          const ditherAmount = ditherThreshold - 0.5;
          const adjustedR = clamp(r + ditherAmount * 64, 0, 255);
          const adjustedG = clamp(g + ditherAmount * 64, 0, 255);
          const adjustedB = clamp(b + ditherAmount * 64, 0, 255);
          const levels = 4;
          outputColor = [
            Math.round(adjustedR / (255 / levels)) * (255 / levels),
            Math.round(adjustedG / (255 / levels)) * (255 / levels),
            Math.round(adjustedB / (255 / levels)) * (255 / levels),
          ];
          break;
        }
      }

      if (props.invert) {
        outputColor = [255 - outputColor[0], 255 - outputColor[1], 255 - outputColor[2]];
      }

      ctx.fillStyle = `rgb(${outputColor[0]}, ${outputColor[1]}, ${outputColor[2]})`;
      ctx.fillRect(x, y, effectivePixelSize, effectivePixelSize);
    }
  }
}

function processImage(img: HTMLImageElement) {
  const canvas = canvasRef.value;
  if (!canvas || dimensions.value.width === 0 || dimensions.value.height === 0) return;

  const dpr = typeof window !== "undefined" ? Math.min(window.devicePixelRatio || 1, 2) : 1;
  const displayWidth = dimensions.value.width;
  const displayHeight = dimensions.value.height;

  canvas.width = Math.floor(displayWidth * dpr);
  canvas.height = Math.floor(displayHeight * dpr);

  const ctx = canvas.getContext("2d");
  if (!ctx) return;
  ctx.resetTransform();
  ctx.scale(dpr, dpr);

  const offscreen = document.createElement("canvas");
  const iw = img.naturalWidth || displayWidth;
  const ih = img.naturalHeight || displayHeight;

  let dw = displayWidth;
  let dh = displayHeight;
  let dx = 0;
  let dy = 0;

  if (props.objectFit === "cover") {
    const scale = Math.max(displayWidth / iw, displayHeight / ih);
    dw = Math.ceil(iw * scale);
    dh = Math.ceil(ih * scale);
    dx = Math.floor((displayWidth - dw) / 2);
    dy = Math.floor((displayHeight - dh) / 2);
  } else if (props.objectFit === "contain") {
    const scale = Math.min(displayWidth / iw, displayHeight / ih);
    dw = Math.ceil(iw * scale);
    dh = Math.ceil(ih * scale);
    dx = Math.floor((displayWidth - dw) / 2);
    dy = Math.floor((displayHeight - dh) / 2);
  } else if (props.objectFit === "fill") {
    dw = displayWidth;
    dh = displayHeight;
  } else {
    dw = iw;
    dh = ih;
    dx = Math.floor((displayWidth - dw) / 2);
    dy = Math.floor((displayHeight - dh) / 2);
  }

  offscreen.width = displayWidth;
  offscreen.height = displayHeight;
  const offCtx = offscreen.getContext("2d");
  if (!offCtx) return;

  offCtx.drawImage(img, dx, dy, dw, dh);

  try {
    imageDataRef.value = offCtx.getImageData(0, 0, displayWidth, displayHeight);
  } catch {
    return;
  }

  applyDithering(ctx, displayWidth, displayHeight, timeRef.value);

  if (props.animated && !prefersReducedMotion.value) {
    const animate = () => {
      timeRef.value += props.animationSpeed;
      applyDithering(ctx, displayWidth, displayHeight, timeRef.value);
      animationRef.value = requestAnimationFrame(animate);
    };
    animationRef.value = requestAnimationFrame(animate);
  }
}

function loadImage() {
  if (!import.meta.client) return;
  if (animationRef.value !== null) {
    cancelAnimationFrame(animationRef.value);
    animationRef.value = null;
  }

  if (imageRef.value && imageRef.value.complete && loadedSrc.value === props.src) {
    processImage(imageRef.value);
    return;
  }

  const img = new Image();
  img.crossOrigin = "anonymous";
  img.src = props.src;

  img.onload = () => {
    imageRef.value = img;
    loadedSrc.value = props.src;
    processImage(img);
  };
}

let resizeObserver: ResizeObserver | null = null;
let mediaQuery: MediaQueryList | null = null;

function handleMotionChange(e: MediaQueryListEvent) {
  prefersReducedMotion.value = e.matches;
  loadImage();
}

onMounted(() => {
  const container = containerRef.value;
  if (!container) return;

  mediaQuery = window.matchMedia("(prefers-reduced-motion: reduce)");
  prefersReducedMotion.value = mediaQuery.matches;
  mediaQuery.addEventListener("change", handleMotionChange);

  resizeObserver = new ResizeObserver((entries) => {
    for (const entry of entries) {
      const { width, height } = entry.contentRect;
      if (width > 0 && height > 0) {
        dimensions.value = { width: Math.floor(width), height: Math.floor(height) };
      }
    }
  });

  resizeObserver.observe(container);
});

onUnmounted(() => {
  if (resizeObserver) {
    resizeObserver.disconnect();
  }
  if (mediaQuery) {
    mediaQuery.removeEventListener("change", handleMotionChange);
  }
  if (animationRef.value !== null) {
    cancelAnimationFrame(animationRef.value);
    animationRef.value = null;
  }
});

watch(
  [
    () => props.src,
    () => dimensions.value,
    () => props.objectFit,
    () => props.animated,
    () => props.animationSpeed,
    () => props.gridSize,
    () => props.ditherMode,
    () => props.colorMode,
    () => props.invert,
    () => props.pixelRatio,
    () => props.primaryColor,
    () => props.secondaryColor,
    () => props.customPalette,
    () => props.brightness,
    () => props.contrast,
    () => props.backgroundColor,
    () => props.threshold,
  ],
  () => {
    loadImage();
  },
  { deep: true },
);
</script>

<template>
  <div
    ref="containerRef"
    :class="cn(`relative h-full w-full`, props.class)"
  >
    <canvas
      ref="canvasRef"
      class="absolute inset-0 h-full w-full"
      style="image-rendering: pixelated"
      aria-hidden="true"
    />
  </div>
</template>
