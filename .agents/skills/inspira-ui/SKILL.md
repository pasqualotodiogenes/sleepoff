---
name: inspira-ui
description: Uso do registry Inspira UI com Nuxt, Tailwind CSS v4, Motion-V e ClientOnly
---

# Inspira UI - Guia de Integração e Melhores Práticas

## 1. Registry & Instalação
A Inspira UI utiliza um registry distribuído compatível com shadcn.
Documentação oficial: [https://inspira-ui.com/docs/en/getting-started](https://inspira-ui.com/docs/en/getting-started)

No `components.json`:
```json
{
  "$schema": "https://shadcn-vue.com/schema.json",
  "style": "default",
  "typescript": true,
  "tailwind": {
    "config": "",
    "css": "assets/css/main.css",
    "baseColor": "slate",
    "cssVariables": true
  },
  "framework": "nuxt",
  "aliases": {
    "components": "~/components",
    "utils": "~/lib/utils",
    "ui": "~/components/ui",
    "lib": "~/lib",
    "composables": "~/composables"
  },
  "registries": {
    "@inspira-ui": "https://registry.inspira-ui.com/{name}.json"
  }
}
```

Para adicionar componentes:
```bash
npx shadcn-vue@latest add @inspira-ui/<componente>
# ou baixar direto do registry:
# curl -s https://registry.inspira-ui.com/<componente>.json
```

## 2. Regras Estritas de SSR e `<ClientOnly>`
Componentes que tocam APIs de DOM, Browser (`window`, `document`, `requestAnimationFrame`), Canvas ou WebGL **NUNCA** devem rodar no SSR (Server-Side Rendering / Prerender).

- Envolva componentes visuais com efeitos canvas/partículas ou interatividade do browser com `<ClientOnly>`.
- Forneça sempre um `#fallback` limpo com dimensões preservadas para evitar CLS (Cumulative Layout Shift).
- Em composables ou scripts, proteja acessos imperativos com `if (import.meta.client)` ou dentro do hook `onMounted()`.

Exemplo:
```vue
<ClientOnly>
  <ParticlesBg :quantity="50" class="absolute inset-0" />
  <template #fallback>
    <div class="absolute inset-0 bg-slate-950" />
  </template>
</ClientOnly>
```

## 3. Performance & Acessibilidade
- Animações CSS e Motion-V devem utilizar estritamente propriedades aceleradas por GPU (`transform`, `opacity`). Evite animar `top`, `left`, `width`, `height` ou `filter` excessivo.
- Respeitar sempre as preferências do usuário para redução de movimento:
```css
@media (prefers-reduced-motion: reduce) {
  *, ::before, ::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
    scroll-behavior: auto !important;
  }
}
```
- Limitar a contagem de nós dinâmicos em fundos animados (partículas, estrelas, grids) para manter 60 FPS estáveis mesmo em dispositivos modestos.
