import type { Component } from "vue";

export interface AnimatedTab {
  title: string;
  value: string;
  panelClassName?: string;
  content?: Component;
  contentProps?: Record<string, unknown>;
}
