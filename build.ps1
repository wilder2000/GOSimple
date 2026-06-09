param (
    [string]$OutputDir = "build"
)

$ErrorActionPreference = "Stop"
$rootDir = $PSScriptRoot

Write-Host "=== GOSimple Build ===" -ForegroundColor Cyan

# 1. Clean build directory
$buildDir = Join-Path $rootDir $OutputDir
if (Test-Path $buildDir) {
    Write-Host "Cleaning $OutputDir/ ..." -ForegroundColor Yellow
    Remove-Item -Path "$buildDir\*" -Recurse -Force
} else {
    New-Item -ItemType Directory -Path $buildDir -Force | Out-Null
}

# 2. Build web frontend
Write-Host "Building web frontend ..." -ForegroundColor Yellow
Push-Location (Join-Path $rootDir "web")
try {
    npm install --silent
    npm run build
    if ($LASTEXITCODE -ne 0) { throw "Frontend build failed" }
    Write-Host "  OK" -ForegroundColor Green
} finally {
    Pop-Location
}

# 2b. Copy dist to http/webdist for Go embed
Write-Host "Copying frontend dist for embed ..." -ForegroundColor Yellow
Copy-Item -Path (Join-Path $rootDir "web\dist") -Destination (Join-Path $rootDir "http\webdist") -Recurse -Force
Write-Host "  OK" -ForegroundColor Green

# 3. Build Go binary
Write-Host "Building Go binary ..." -ForegroundColor Yellow
$binaryName = "GOSimple.exe"
$outputPath = Join-Path $buildDir $binaryName
go build -o $outputPath -ldflags="-s -w" .
if ($LASTEXITCODE -ne 0) { throw "Go build failed" }
Write-Host "  -> $OutputDir/$binaryName" -ForegroundColor Green

# 4. Copy runtime config
Write-Host "Copying runtime files ..." -ForegroundColor Yellow
Copy-Item -Path (Join-Path $rootDir "conf") -Destination $buildDir -Recurse
Write-Host "  -> conf/" -ForegroundColor Green

Write-Host "=== Build complete ===" -ForegroundColor Cyan
