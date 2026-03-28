# Releasing sleepoff

## Preconditions

- CI verde em `main`
- `README.md` e instalacao alinhados com o estado real do produto
- `LICENSE`, installer e assets de release revisados

## Release flow

```bash
go test ./...
go vet ./...
git tag v1.0.0
git push origin main --tags
```

## What the release workflow does

Ao receber uma tag `v*`, o workflow de release:

1. instala `goversioninfo` e `Inno Setup`
2. gera `resource.syso` com icone e metadata do executavel
3. roda `go test ./...` e `go vet ./...`
4. publica o pacote portable `sleepoff_windows_amd64.zip` via GoReleaser
5. compila o instalador `sleepoff-setup.exe`
6. atualiza `checksums.txt`
7. anexa instalador e checksums ao GitHub Release

## Release acceptance

Verifique no GitHub Release:

- `sleepoff_windows_amd64.zip`
- `sleepoff-setup.exe`
- `checksums.txt`

Teste os dois fluxos no Windows:

- installer: `sleepoff --help` em qualquer pasta
- zip portable: `.\sleepoff.exe --help` na pasta extraida
