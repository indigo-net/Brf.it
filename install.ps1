# Brf.it installer for Windows
# Usage: .\install.ps1 [-Version <version>]
# Example: .\install.ps1 -Version v0.5.0
#
# Environment variables:
#   BRFIT_VERSION - Specify version to install (e.g., v0.5.0)
#   BRFIT_INSTALL_DIR - Specify custom install directory

param(
    [string]$Version = $env:BRFIT_VERSION
)

$ErrorActionPreference = "Stop"

$repo = "indigo-net/Brf.it"
$installDir = if ($env:BRFIT_INSTALL_DIR) { $env:BRFIT_INSTALL_DIR } else { "$env:ProgramFiles\brfit" }
$binaryName = "brfit.exe"

function Write-Info($msg) {
    Write-Host "==> " -NoNewline -ForegroundColor Green
    Write-Host $msg
}

function Write-Warn($msg) {
    Write-Host "==> " -NoNewline -ForegroundColor Yellow
    Write-Host $msg
}

function Write-Err($msg) {
    Write-Host "Error: " -NoNewline -ForegroundColor Red
    Write-Host $msg
    exit 1
}

# Check architecture (only amd64 supported for Windows)
$arch = [System.Environment]::GetEnvironmentVariable("PROCESSOR_ARCHITECTURE")
if ($arch -ne "AMD64") {
    Write-Err "Unsupported architecture: $arch. Only AMD64 (x86_64) is supported on Windows."
}

Write-Info "Detected: windows/amd64"

# Get version
if (-not $Version) {
    Write-Info "Fetching latest version..."
    try {
        $release = Invoke-RestMethod "https://api.github.com/repos/$repo/releases/latest"
        $Version = $release.tag_name
    } catch {
        Write-Err "Failed to fetch release info: $_"
    }
}

# Ensure version starts with 'v'
if (-not $Version.StartsWith("v")) {
    $Version = "v$Version"
}

Write-Info "Installing brfit $Version"

# Build download URLs
$versionNum = $Version.TrimStart('v')
$filename = "brfit_${versionNum}_windows_amd64.zip"
$url = "https://github.com/$repo/releases/download/$Version/$filename"
$checksumUrl = "https://github.com/$repo/releases/download/$Version/checksums.txt"

# Create temp directory
$tmpDir = Join-Path ([System.IO.Path]::GetTempPath()) ([System.Guid]::NewGuid().ToString())
New-Item -ItemType Directory -Force -Path $tmpDir | Out-Null

try {
    # Download files
    Write-Info "Downloading $filename..."
    try {
        Invoke-WebRequest -Uri $url -OutFile "$tmpDir\$filename" -UseBasicParsing
    } catch {
        Write-Err "Download failed: $_"
    }

    Write-Info "Downloading checksums..."
    try {
        Invoke-WebRequest -Uri $checksumUrl -OutFile "$tmpDir\checksums.txt" -UseBasicParsing
    } catch {
        Write-Err "Checksum download failed: $_"
    }

    # Verify checksum
    Write-Info "Verifying checksum..."
    $checksums = Get-Content "$tmpDir\checksums.txt"
    $expectedLine = $checksums | Where-Object { $_ -match [regex]::Escape($filename) }
    if (-not $expectedLine) {
        Write-Err "Checksum not found for $filename"
    }
    $expected = ($expectedLine -split '\s+')[0].ToLower()
    $actual = (Get-FileHash "$tmpDir\$filename" -Algorithm SHA256).Hash.ToLower()

    if ($expected -ne $actual) {
        Write-Err "Checksum mismatch!`n  Expected: $expected`n  Actual:   $actual"
    }
    Write-Info "Checksum verified"

    # Extract archive
    Write-Info "Extracting..."
    Expand-Archive -Path "$tmpDir\$filename" -DestinationPath $tmpDir -Force

    # Install binary
    Write-Info "Installing to $installDir..."

    # Check if we need admin rights
    $needsAdmin = $false
    try {
        if (-not (Test-Path $installDir)) {
            New-Item -ItemType Directory -Force -Path $installDir | Out-Null
        }
        # Test write permission
        $testFile = Join-Path $installDir ".write_test"
        [System.IO.File]::Create($testFile).Close()
        Remove-Item $testFile
    } catch {
        $needsAdmin = $true
    }

    if ($needsAdmin) {
        Write-Warn "Administrator privileges required to install to $installDir"
        Write-Host "Please run this script as Administrator, or set BRFIT_INSTALL_DIR to a writable location."
        exit 1
    }

    Copy-Item "$tmpDir\brfit.exe" "$installDir\$binaryName" -Force

    # Add to PATH if not already present
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    $pathAdded = $false
    if ($userPath -notlike "*$installDir*") {
        [Environment]::SetEnvironmentVariable("Path", "$userPath;$installDir", "User")
        $pathAdded = $true
    }

    # Success message
    Write-Info "brfit $Version installed successfully!"
    Write-Host ""

    if ($pathAdded) {
        Write-Host ([char]0x2713) -NoNewline -ForegroundColor Green
        Write-Host " Added $installDir to your PATH."
        Write-Host ""
        Write-Host "Please restart your terminal for the changes to take effect."
    } else {
        Write-Host ([char]0x2713) -NoNewline -ForegroundColor Green
        Write-Host " $installDir is already in your PATH."
    }

    Write-Host ""
    Write-Host "Run 'brfit --help' to get started."

} finally {
    # Cleanup
    if (Test-Path $tmpDir) {
        Remove-Item -Recurse -Force $tmpDir -ErrorAction SilentlyContinue
    }
}
