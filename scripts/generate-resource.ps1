param(
    [Parameter(Mandatory = $true)]
    [string]$Version,
    [string]$OutFile = "resource.syso"
)

$ErrorActionPreference = "Stop"

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
$templatePath = Join-Path $repoRoot "packaging\\windows\\versioninfo.template.json"
$iconPath = Join-Path $repoRoot "packaging\\windows\\icon.ico"

if (-not (Test-Path $templatePath)) {
    throw "versioninfo template not found: $templatePath"
}

if (-not (Test-Path $iconPath)) {
    throw "icon not found: $iconPath"
}

$goversioninfo = Get-Command goversioninfo -ErrorAction SilentlyContinue
if (-not $goversioninfo) {
    throw "goversioninfo not found in PATH. Install with: go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest"
}

$normalizedVersion = $Version.Trim()
if ($normalizedVersion.StartsWith("v")) {
    $normalizedVersion = $normalizedVersion.Substring(1)
}

$parts = $normalizedVersion.Split(".")
while ($parts.Count -lt 3) {
    $parts += "0"
}

$rendered = Get-Content -Path $templatePath -Raw
$rendered = $rendered.Replace("{{VERSION}}", $normalizedVersion)
$rendered = $rendered.Replace("{{MAJOR}}", $parts[0])
$rendered = $rendered.Replace("{{MINOR}}", $parts[1])
$rendered = $rendered.Replace("{{PATCH}}", $parts[2])

$tempDir = Join-Path $env:TEMP ("sleepoff-versioninfo-" + [guid]::NewGuid().ToString("N"))
New-Item -ItemType Directory -Path $tempDir | Out-Null

try {
    Set-Content -Path (Join-Path $tempDir "versioninfo.json") -Value $rendered -Encoding UTF8

    if ([System.IO.Path]::IsPathRooted($OutFile)) {
        $outputPath = $OutFile
    } else {
        $outputPath = Join-Path $repoRoot $OutFile
    }

    Push-Location $tempDir
    try {
        & $goversioninfo.Source "-icon=$iconPath" "-o=$outputPath"
    } finally {
        Pop-Location
    }
} finally {
    Remove-Item -Path $tempDir -Recurse -Force -ErrorAction SilentlyContinue
}
