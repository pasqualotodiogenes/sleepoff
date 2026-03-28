# TODO

## Concluido nesta entrega

- [x] Fluxo de timer estabilizado (ticks ativos apos splash/menu)
- [x] Tela de panico corrigida para 10 segundos reais
- [x] Cancelamento na tela de panico volta para fluxo de "cancelado" (nao "sucesso")
- [x] Tratamento explicito de erro ao executar shutdown
- [x] Validacao de duracao (`> 0`) no modo CLI
- [x] Ajuste de `+/-` com timer pausado mantendo `Remaining` consistente
- [x] Execucao por duplo clique no Windows sem bloqueio do Cobra (`MousetrapHelpText`)
- [x] Build cross-platform (stubs de shutdown com build tags)
- [x] Testes automatizados adicionados para parsing e state machine
- [x] Modulo publico alinhado com `github.com/pasqualotodiogenes/sleepoff`
- [x] Pipeline de release com GoReleaser e GitHub Releases
- [x] Instalador Windows por usuario com Inno Setup
- [x] Script de instalacao local no `PATH` para validar a experiencia final
- [x] README reescrito com instalacao real, uso local vs global e artefatos de release

## Melhorias Futuras

### Icone do App
- [x] Criar icone `.ico`
- [x] Embed de metadata/icon no build de release com `goversioninfo`

### Distribuicao
- [ ] Publicar o primeiro release `v1.0.0` no GitHub
- [ ] Submeter o app ao `winget`
- [ ] Avaliar `.msi` ou instalador alternativo alem do Inno Setup

### Features
- [ ] Configuracao persistente (ultima escolha, tema, etc.)
- [ ] Sons customizaveis
- [ ] Temas de cores
- [ ] Modo "apenas avisar" (sem desligar, so notificacao)
