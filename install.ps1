# central-set installer for Windows
#irm https://realdatadriven.github.io/central-set-go/install.ps1 | iex

$ErrorActionPreference = "Stop"

$Repo = "realdatadriven/central-set-go"
$ApiUrl = "https://api.github.com/repos/$Repo/releases/latest"

function Die {
    param([string]$Message)
    Write-Error $Message
    exit 1
}

function Detect-Platform {
    $arch = $env:PROCESSOR_ARCHITECTURE

    switch ($arch) {
        "AMD64" { return "windows-amd64" }
        "ARM64" { return "windows-arm64" }
        default { Die "Unsupported Windows architecture: $arch" }
    }
}

function Fetch-Release {
    param([string]$Platform)

    Write-Host "Fetching latest release information..."

    try {
        $release = Invoke-RestMethod `
            -Uri $ApiUrl `
            -Headers @{
                Accept = "application/vnd.github+json"
                "User-Agent" = "PowerShell"
            }
    }
    catch {
        Die "Failed to fetch release information."
    }

    $version = $release.tag_name

    if ([string]::IsNullOrWhiteSpace($version)) {
        Die "Could not determine latest version."
    }

    $assetName = "central-set-$Platform.zip"
    $assetVersionedName = "central-set-$Platform-$version.zip"
    $downloadUrl = "https://github.com/$Repo/releases/download/$version/$assetName"

    Write-Host "Version : $version"
    Write-Host "Asset   : $assetName"
    Write-Host "URL     : $downloadUrl"

    return @{
        Version = $version
        AssetName = $assetName
        AssetVersionedName = $assetVersionedName
        DownloadUrl = $downloadUrl
    }
}

function Install-Release {
    param(
        [string]$Platform,
        [hashtable]$ReleaseInfo
    )

    $tmpZip = Join-Path $env:TEMP $ReleaseInfo.AssetVersionedName

    Write-Host ""
    Write-Host "Downloading..."

    if (-not (Test-Path $tmpZip)) {
        Invoke-WebRequest `
            -Uri $ReleaseInfo.DownloadUrl `
            -OutFile $tmpZip
    }

    Write-Host ""
    Write-Host "Extracting into $(Get-Location)..."

    Expand-Archive `
        -Path $tmpZip `
        -DestinationPath "." `
        -Force

    $sourceExe = "central-set-$Platform.exe"

    if (-not (Test-Path "c7.exe") -and (Test-Path $sourceExe)) {
        Rename-Item $sourceExe "c7.exe"
        Write-Host "Renamed $sourceExe to c7.exe"
    }

    if (-not (Test-Path ".env") -and (Test-Path "dot-env-example.txt")) {
        Rename-Item "dot-env-example.txt" ".env"
        Write-Host "Created .env from dot-env-example.txt"
    }

    if (Test-Path "c7.exe") {
        & ".\c7.exe" --init --model "admin_model.md"
        Write-Host "Set up admin model"
    }

    if (Test-Path "c7.exe") {
        & ".\c7.exe" --init --model "etlx_model.md"
        Write-Host "Set up etlx model"
    }

    Write-Host ""
    Write-Host "Installation completed."
    Write-Host ""
    Write-Host "Files extracted to:"
    Write-Host "  $(Get-Location)"
    Write-Host ""
}

try {
    $platform = Detect-Platform
    $releaseInfo = Fetch-Release -Platform $platform
    Install-Release -Platform $platform -ReleaseInfo $releaseInfo
}
catch {
    Die $_.Exception.Message
}