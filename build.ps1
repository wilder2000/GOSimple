param (
    [string]$OutputDir = "build",
    [string]$PublishDir = "",
    [switch]$OverwriteConf = $false
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

# 2b. Copy dist contents to http/webdist for Go embed (clean first to remove stale hash files)
Write-Host "Copying frontend dist for embed ..." -ForegroundColor Yellow
$webdistDir = Join-Path $rootDir "http\webdist"
if (Test-Path $webdistDir) {
    Remove-Item -Path "$webdistDir\*" -Recurse -Force
}
# Copy-Item with trailing "\*" copies contents, not the source folder itself
Copy-Item -Path "$(Join-Path $rootDir 'web\dist')\*" -Destination $webdistDir -Recurse
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

# 5. Publish to custom directory (optional)
if ($PublishDir) {
    $pubDir = $PublishDir
    if (-not [System.IO.Path]::IsPathRooted($pubDir)) {
        $pubDir = Join-Path $rootDir $PublishDir
    }
    Write-Host "Publishing to $pubDir ..." -ForegroundColor Yellow
    if (-not (Test-Path $pubDir)) {
        New-Item -ItemType Directory -Path $pubDir -Force | Out-Null
    }
    Copy-Item -Path $outputPath -Destination (Join-Path $pubDir $binaryName) -Force
    $pubConf = Join-Path $pubDir "conf"
    if ($OverwriteConf -or -not (Test-Path $pubConf)) {
        Copy-Item -Path (Join-Path $rootDir "conf") -Destination $pubDir -Recurse -Force
        Write-Host "  -> conf/ (overwritten)" -ForegroundColor Green
    } else {
        Write-Host "  -> conf/ (skipped, use -OverwriteConf to replace)" -ForegroundColor DarkYellow
    }
    Write-Host "  OK" -ForegroundColor Green
}

Write-Host "=== Build complete ===" -ForegroundColor Cyan
