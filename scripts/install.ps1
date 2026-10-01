# Instala o CV-Craft no Windows a partir da última release do GitHub.
#
#   irm https://raw.githubusercontent.com/fabiodrneles/cv-craft/main/scripts/install.ps1 | iex
#
# Instala em %LOCALAPPDATA%\Programs\cv-craft e acrescenta essa pasta ao PATH
# do usuário (não precisa de administrador). Variáveis opcionais:
#   CV_CRAFT_VERSION      versão a instalar (ex.: 1.0.0); padrão: a mais recente
#   CV_CRAFT_INSTALL_DIR  pasta de destino
$ErrorActionPreference = 'Stop'
$ProgressPreference = 'SilentlyContinue' # o Invoke-WebRequest fica muito mais rápido

$repo = 'fabiodrneles/cv-craft'
$dir = if ($env:CV_CRAFT_INSTALL_DIR) { $env:CV_CRAFT_INSTALL_DIR } else { Join-Path $env:LOCALAPPDATA 'Programs\cv-craft' }
$arch = if ($env:PROCESSOR_ARCHITECTURE -eq 'ARM64') { 'arm64' } else { 'amd64' }

$version = $env:CV_CRAFT_VERSION
if (-not $version) {
    $version = (Invoke-RestMethod "https://api.github.com/repos/$repo/releases/latest").tag_name
}
$version = $version.TrimStart('v')

$archive = "cv-craft_${version}_windows_$arch.zip"
$base = "https://github.com/$repo/releases/download/v$version"
$tmp = Join-Path ([IO.Path]::GetTempPath()) ([Guid]::NewGuid())
New-Item -ItemType Directory -Path $tmp | Out-Null
try {
    Write-Host "Baixando o CV-Craft $version (windows/$arch)..."
    Invoke-WebRequest "$base/$archive" -OutFile (Join-Path $tmp $archive) -UseBasicParsing
    Invoke-WebRequest "$base/checksums.txt" -OutFile (Join-Path $tmp 'checksums.txt') -UseBasicParsing

    $line = Get-Content (Join-Path $tmp 'checksums.txt') | Where-Object { $_ -match " $([regex]::Escape($archive))$" }
    $expected = ($line -split ' ')[0]
    $actual = (Get-FileHash (Join-Path $tmp $archive) -Algorithm SHA256).Hash.ToLower()
    if (-not $expected -or $expected -ne $actual) { throw 'o checksum do arquivo baixado não confere' }

    Expand-Archive (Join-Path $tmp $archive) -DestinationPath (Join-Path $tmp 'x') -Force
    New-Item -ItemType Directory -Force -Path $dir | Out-Null
    Copy-Item (Join-Path $tmp 'x\cv-craft.exe') (Join-Path $dir 'cv-craft.exe') -Force
} finally {
    Remove-Item $tmp -Recurse -Force -ErrorAction SilentlyContinue
}

$userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
if (($userPath -split ';') -notcontains $dir) {
    $newPath = if ($userPath) { "$userPath;$dir" } else { $dir }
    [Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
    Write-Host "A pasta $dir foi adicionada ao seu PATH."
}
if (($env:Path -split ';') -notcontains $dir) { $env:Path = "$env:Path;$dir" }

Write-Host "Instalado em $dir\cv-craft.exe"
& (Join-Path $dir 'cv-craft.exe') version
Write-Host ''
Write-Host 'Pronto! Neste terminal o comando cv-craft já funciona. Nos terminais que já estavam abertos, feche e abra de novo.'
