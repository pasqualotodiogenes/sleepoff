<script setup lang="ts">
import { useDevicePixelRatio, useMouse } from "@vueuse/core";
import { computed, onBeforeUnmount, onMounted, reactive, ref, watch } from "vue";

interface Circle {
  x: number;
  y: number;
  translateX: number;
  translateY: number;
  size: number;
  alpha: number;
  targetAlpha: number;
  dx: number;
  dy: number;
  magnetism: number;
}

interface Props {
  color?: string;
  quantity?: number;
  staticity?: number;
  ease?: number;
  connectLines?: boolean;
  lineColor?: string;
  lineDistance?: number;
  lineOpacity?: number;
  class?: string;
}

const props = withDefaults(defineProps<Props>(), {
  color: "#ededf3",
  quantity: 65,
  staticity: 50,
  ease: 50,
  connectLines: false,
  lineColor: "#5266eb",
  lineDistance: 120,
  lineOpacity: 0.18,
  class: "",
});

const canvasRef = ref<HTMLCanvasElement | null>(null);
const canvasContainerRef = ref<HTMLDivElement | null>(null);
const context = ref<CanvasRenderingContext2D | null>(null);
const circles = ref<Circle[]>([]);
const mouse = reactive<{ x: number; y: number }>({ x: 0, y: 0 });
const canvasSize = reactive<{ w: number; h: number }>({ w: 0, h: 0 });
const { x: mouseX, y: mouseY } = useMouse();
const { pixelRatio } = useDevicePixelRatio();

const rafId = ref<number | null>(null);
const prefersReducedMotion = ref(false);
let mediaQuery: MediaQueryList | null = null;

function parseHexToRgbTuple(hexInput: string): string {
  let hex = hexInput.replace(/^#/, "");
  if (hex.length === 3) {
    hex = hex
      .split("")
      .map((char) => char + char)
      .join("");
  }
  const bigint = Number.parseInt(hex, 16);
  if (Number.isNaN(bigint)) return "237, 237, 243";
  const r = (bigint >> 16) & 255;
  const g = (bigint >> 8) & 255;
  const b = bigint & 255;
  return `${r}, ${g}, ${b}`;
}

const colorRgb = computed(() => parseHexToRgbTuple(props.color));
const lineColorRgb = computed(() => parseHexToRgbTuple(props.lineColor));

const effectiveQuantity = computed(() => {
  if (typeof window !== "undefined" && window.innerWidth < 768) {
    return Math.min(props.quantity, Math.max(10, Math.round(props.quantity * 0.4)));
  }
  return props.quantity;
});

function stopAnimation() {
  if (rafId.value !== null) {
    window.cancelAnimationFrame(rafId.value);
    rafId.value = null;
  }
}

function handleMotionChange(e: MediaQueryListEvent) {
  prefersReducedMotion.value = e.matches;
  stopAnimation();
  renderStaticOrStartLoop();
}

onMounted(() => {
  if (canvasRef.value) {
    context.value = canvasRef.value.getContext("2d");
  }

  mediaQuery = window.matchMedia("(prefers-reduced-motion: reduce)");
  prefersReducedMotion.value = mediaQuery.matches;
  mediaQuery.addEventListener("change", handleMotionChange);

  initCanvas();
  renderStaticOrStartLoop();
  window.addEventListener("resize", handleResize);
});

onBeforeUnmount(() => {
  stopAnimation();
  window.removeEventListener("resize", handleResize);
  if (mediaQuery) {
    mediaQuery.removeEventListener("change", handleMotionChange);
  }
});

watch([mouseX, mouseY], () => {
  if (!prefersReducedMotion.value) {
    onMouseMove();
  }
});

watch(
  [() => props.color, () => props.lineColor, () => props.connectLines, () => props.quantity],
  () => {
    if (!import.meta.client) return;
    if (prefersReducedMotion.value) {
      drawStaticFrame();
    }
  },
);

function handleResize() {
  initCanvas();
  if (prefersReducedMotion.value) {
    drawStaticFrame();
  }
}

function renderStaticOrStartLoop() {
  if (prefersReducedMotion.value) {
    drawStaticFrame();
    return;
  }
  if (rafId.value === null) {
    rafId.value = window.requestAnimationFrame(animate);
  }
}

function initCanvas() {
  resizeCanvas();
  drawParticles();
}

function onMouseMove() {
  if (canvasRef.value) {
    const rect = canvasRef.value.getBoundingClientRect();
    const { w, h } = canvasSize;
    const x = mouseX.value - rect.left - w / 2;
    const y = mouseY.value - rect.top - h / 2;

    const inside = x < w / 2 && x > -w / 2 && y < h / 2 && y > -h / 2;
    if (inside) {
      mouse.x = x;
      mouse.y = y;
    }
  }
}

function resizeCanvas() {
  if (canvasContainerRef.value && canvasRef.value && context.value) {
    circles.value.length = 0;
    canvasSize.w = canvasContainerRef.value.offsetWidth;
    canvasSize.h = canvasContainerRef.value.offsetHeight;
    const dpr = Math.min(pixelRatio.value || 1, 2);
    canvasRef.value.width = canvasSize.w * dpr;
    canvasRef.value.height = canvasSize.h * dpr;
    canvasRef.value.style.width = `${canvasSize.w}px`;
    canvasRef.value.style.height = `${canvasSize.h}px`;
    context.value.setTransform(dpr, 0, 0, dpr, 0, 0);
  }
}

function circleParams(): Circle {
  const x = Math.floor(Math.random() * Math.max(1, canvasSize.w));
  const y = Math.floor(Math.random() * Math.max(1, canvasSize.h));
  const size = Math.floor(Math.random() * 2) + 1;
  const targetAlpha = Number.parseFloat((Math.random() * 0.5 + 0.15).toFixed(2));
  const alpha = prefersReducedMotion.value ? targetAlpha : 0;
  const dx = (Math.random() - 0.5) * 0.2;
  const dy = (Math.random() - 0.5) * 0.2;
  const magnetism = 0.1 + Math.random() * 4;
  return {
    x,
    y,
    translateX: 0,
    translateY: 0,
    size,
    alpha,
    targetAlpha,
    dx,
    dy,
    magnetism,
  };
}

function drawCircle(circle: Circle, update = false) {
  if (!context.value) return;
  const dpr = Math.min(pixelRatio.value || 1, 2);
  const { x, y, translateX, translateY, size, alpha } = circle;
  context.value.translate(translateX, translateY);
  context.value.beginPath();
  context.value.arc(x, y, size, 0, 2 * Math.PI);
  context.value.fillStyle = `rgba(${colorRgb.value}, ${alpha})`;
  context.value.fill();
  context.value.setTransform(dpr, 0, 0, dpr, 0, 0);

  if (!update) {
    circles.value.push(circle);
  }
}

function drawConstellationLines() {
  if (!props.connectLines || !context.value) return;
  const ctx = context.value;
  const maxDist = props.lineDistance;
  const maxDistSq = maxDist * maxDist;
  const list = circles.value;
  const len = list.length;

  ctx.lineWidth = 1;
  for (let i = 0; i < len; i++) {
    const a = list[i];
    if (!a) continue;
    const ax = a.x + a.translateX;
    const ay = a.y + a.translateY;
    for (let j = i + 1; j < len; j++) {
      const b = list[j];
      if (!b) continue;
      const bx = b.x + b.translateX;
      const by = b.y + b.translateY;
      const dx = ax - bx;
      const dy = ay - by;
      const distSq = dx * dx + dy * dy;
      if (distSq < maxDistSq) {
        const distRatio = 1 - Math.sqrt(distSq) / maxDist;
        const alpha = (distRatio * props.lineOpacity).toFixed(3);
        ctx.strokeStyle = `rgba(${lineColorRgb.value}, ${alpha})`;
        ctx.beginPath();
        ctx.moveTo(ax, ay);
        ctx.lineTo(bx, by);
        ctx.stroke();
      }
    }
  }
}

function clearContext() {
  if (context.value) {
    context.value.clearRect(0, 0, canvasSize.w, canvasSize.h);
  }
}

function drawParticles() {
  clearContext();
  const particleCount = effectiveQuantity.value;
  for (let i = 0; i < particleCount; i++) {
    const circle = circleParams();
    drawCircle(circle);
  }
  drawConstellationLines();
}

function drawStaticFrame() {
  clearContext();
  for (const circle of circles.value) {
    circle.alpha = circle.targetAlpha;
    drawCircle(circle, true);
  }
  drawConstellationLines();
}

function remapValue(
  value: number,
  start1: number,
  end1: number,
  start2: number,
  end2: number,
): number {
  const remapped = ((value - start1) * (end2 - start2)) / (end1 - start1) + start2;
  return remapped > 0 ? remapped : 0;
}

function animate() {
  if (prefersReducedMotion.value) {
    rafId.value = null;
    drawStaticFrame();
    return;
  }

  clearContext();
  circles.value.forEach((circle, i) => {
    const edge = [
      circle.x + circle.translateX - circle.size,
      canvasSize.w - circle.x - circle.translateX - circle.size,
      circle.y + circle.translateY - circle.size,
      canvasSize.h - circle.y - circle.translateY - circle.size,
    ];

    const closestEdge = edge.reduce((a, b) => Math.min(a, b));
    const remapClosestEdge = Number.parseFloat(remapValue(closestEdge, 0, 20, 0, 1).toFixed(2));

    if (remapClosestEdge > 1) {
      circle.alpha += 0.02;
      if (circle.alpha > circle.targetAlpha) circle.alpha = circle.targetAlpha;
    } else {
      circle.alpha = circle.targetAlpha * remapClosestEdge;
    }

    circle.x += circle.dx;
    circle.y += circle.dy;
    circle.translateX +=
      (mouse.x / (props.staticity / circle.magnetism) - circle.translateX) / props.ease;
    circle.translateY +=
      (mouse.y / (props.staticity / circle.magnetism) - circle.translateY) / props.ease;

    if (
      circle.x < -circle.size ||
      circle.x > canvasSize.w + circle.size ||
      circle.y < -circle.size ||
      circle.y > canvasSize.h + circle.size
    ) {
      circles.value.splice(i, 1);
      const newCircle = circleParams();
      drawCircle(newCircle);
    } else {
      drawCircle(circle, true);
    }
  });

  drawConstellationLines();
  rafId.value = window.requestAnimationFrame(animate);
}
</script>

<template>
  <div
    ref="canvasContainerRef"
    :class="$props.class"
    aria-hidden="true"
  >
    <canvas ref="canvasRef" />
  </div>
</template>
