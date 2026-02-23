$ErrorActionPreference = "Stop"

$repo = "indigo-net/Brf.it"
$installDir = "$env:USERPROFILE\.brfit"

function Write-Info($msg) { Write-Host "==> $msg" -ForegroundColor Green }
function Write-Warn($msg) { Write-Host "==> $msg" -ForegroundColor Yellow }
function Write-Err($msg) { Write-Host "Error: $msg" -ForegroundColor Red; exit 1 }

Write-Info "Fetching latest version..."
try {
    $release = Invoke-RestMethod "https://api.github.com/repos/$repo/releases/latest"
    $version = $release.tag_name
} catch {
    Write-Err "Failed to fetch release info: $_"
}
Write-Info "Latest version: $version"

$filename = "brfit_$($version.TrimStart('v'))_windows_amd64.zip"
$url = "https://github.com/$repo/releases/download/$version/$filename"
$checksumUrl = "https://github.com/$repo/releases/download/$version/checksums.txt"

Write-Info "Downloading $filename..."
$tmpDir = [System.IO.Path]::GetTempPath() + [System.Guid]::NewGuid().ToString()
New-Item -ItemType Directory -Force -Path $tmpDir | Out-Null

try {
    Invoke-WebRequest -Uri $url -OutFile "$tmpDir\$filename"
    Invoke-WebRequest -Uri $checksumUrl -OutFile "$tmpDir\checksums.txt"
} catch {
    Write-Err "Download failed: $_"
}

Write-Info "Verifying checksum..."
$checksums = Get-Content "$tmpDir\checksums.txt"
$expected = ($checksums | Where-Object { $_ -match $filename }) -split '\s+' | Select-Object -First 1
$actual = (Get-FileHash "$tmpDir\$filename" -Algorithm SHA256).Hash.ToLower()

if ($expected -ne $actual) {
    Write-Err "Checksum mismatch! Expected: $expected, Got: $actual"
}
Write-Info "Checksum verified"

Write-Info "Installing to $installDir..."
New-Item -ItemType Directory -Force -Path $installDir | Out-Null
Expand-Archive -Path "$tmpDir\$filename" -DestinationPath $installDir -Force

# Add to PATH
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -notlike "*$installDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$userPath;$installDir", "User")
    Write-Warn "Added $installDir to PATH"
}

# Cleanup
Remove-Item -Recurse -Force $tmpDir

Write-Info "brfit $version installed successfully!"
Write-Host ""
Write-Warn "IMPORTANT: Please restart your terminal for PATH changes to take effect."
Write-Host ""
Write-Host "Run 'brfit --help' to get started."
