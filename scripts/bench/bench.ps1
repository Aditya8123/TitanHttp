[CmdletBinding()]
param (
    [Parameter(Mandatory=$false)]
    [string]$RunMode
)

$ErrorActionPreference = "Stop"

function Show-Menu {
    Clear-Host
    Write-Host "`n  ============================================================" -ForegroundColor Cyan
    Write-Host "  ::                                                        ::" -ForegroundColor Cyan
    Write-Host "  ::           TITANHTTP PERFORMANCE ORCHESTRATOR           ::" -ForegroundColor Yellow
    Write-Host "  ::                                                        ::" -ForegroundColor Cyan
    Write-Host "  ============================================================`n" -ForegroundColor Cyan
    Write-Host "    [1] Run Micro & Component Benchmarks" -ForegroundColor White
    Write-Host "    [2] Run Framework Comparisons" -ForegroundColor White
    Write-Host "    [3] Run Unified Performance Laboratory" -ForegroundColor White
    Write-Host "    [4] Run HTTP Compliance Tests" -ForegroundColor White
    Write-Host "    [5] Run Advanced Security Suite`n" -ForegroundColor White
    Write-Host "    [A] Run Full Cinematic Engine Suite" -ForegroundColor Magenta
    Write-Host "    [Q] Quit`n" -ForegroundColor DarkGray
}

function Run-Micro {
    Write-Host "`n>>> [ STAGE 1 ] Running Micro & Component Benchmarks..." -ForegroundColor Green
    & .\scripts\bench\stats.ps1 -Internal
    Write-Host "Completed Micro & Component Tests.`n" -ForegroundColor Green
}

function Run-Frameworks {
    Write-Host "`n>>> [ STAGE 2 ] Running Framework Comparisons..." -ForegroundColor Green
    & .\scripts\bench\stats.ps1 -Frameworks
    Write-Host "Completed Framework Comparisons.`n" -ForegroundColor Green
}

function Run-UnifiedSuite {
    Write-Host "`n>>> [ STAGE 3 ] Running Unified Performance Laboratory..." -ForegroundColor Green
    & .\scripts\bench\unified_comparison.ps1
    Write-Host "Completed Unified Performance Laboratory.`n" -ForegroundColor Green
}

function Run-Compliance {
    Write-Host "`n>>> [ STAGE 4 ] Running HTTP Compliance Tests..." -ForegroundColor Green
    $TitanProcess = Start-Process -FilePath ".\bin\titanhttp_bench.exe" -PassThru -NoNewWindow -RedirectStandardOutput "benchmarks/logs/titan_compliance.log" -RedirectStandardError "benchmarks/logs/titan_compliance_err.log"
    Start-Sleep -Seconds 2
    
    $compOutput = go run ./cmd/compliance/main.go localhost:8080 2>&1
    $compOutput | ForEach-Object { Write-Host $_ }
    
    Stop-Process -Id $TitanProcess.Id -Force -ErrorAction SilentlyContinue
    
    $reportPath = "benchmarks/reports/compliance_report.txt"
    $mdLines = @(
        "# HTTP Compliance Tests Report"
        "Generated on $(Get-Date)"
        ""
        "## Test Results"
        ""
        '```text'
        ($compOutput -replace "\x1B\[[0-9;]*[a-zA-Z]", "")
        '```'
    )
    ($mdLines -join "`r`n") | Out-File -FilePath $reportPath -Encoding utf8
    Write-Host "Completed Compliance Tests.`n" -ForegroundColor Green
}

function Run-SecuritySuite {
    Write-Host "`n>>> [ STAGE 5 ] Running Advanced Security Suite..." -ForegroundColor Green
    & .\scripts\bench\advanced_security.ps1
    Write-Host "Completed Security Suite.`n" -ForegroundColor Green
}

Write-Host "Initializing Performance Orchestrator..." -ForegroundColor Cyan
if (-not (Test-Path "bin")) { New-Item -ItemType Directory -Path "bin" | Out-Null }
if (-not (Test-Path "benchmarks/data")) { New-Item -ItemType Directory -Path "benchmarks/data" -Force | Out-Null }
if (-not (Test-Path "benchmarks/logs")) { New-Item -ItemType Directory -Path "benchmarks/logs" -Force | Out-Null }
if (-not (Test-Path "benchmarks/reports")) { New-Item -ItemType Directory -Path "benchmarks/reports" -Force | Out-Null }

Write-Host "Building Test Binaries..." -ForegroundColor Gray
go build -o bin/titanhttp_bench.exe ./cmd/titanhttp
go build -o bin/nethttp_bench.exe ./cmd/nethttp_bench

if ($RunMode) {
    switch ($RunMode.ToUpper()) {
        '1' { Run-Micro }
        '2' { Run-Frameworks }
        '3' { Run-UnifiedSuite }
        '4' { Run-Compliance }
        '5' { Run-SecuritySuite }
        'A' {
            Run-Micro
            Run-Frameworks
            Run-UnifiedSuite
            Run-Compliance
            Run-SecuritySuite
            Write-Host "Full Engine Suite Complete!" -ForegroundColor Cyan
        }
        default { Write-Host "Invalid RunMode specified." -ForegroundColor Red }
    }
    exit
}

do {
    Show-Menu
    $choice = Read-Host "Select an option"

    switch ($choice) {
        '1' { Run-Micro; pause }
        '2' { Run-Frameworks; pause }
        '3' { Run-UnifiedSuite; pause }
        '4' { Run-Compliance; pause }
        '5' { Run-SecuritySuite; pause }
        'A' {
            Run-Micro
            Run-Frameworks
            Run-UnifiedSuite
            Run-Compliance
            Run-SecuritySuite
            Write-Host "Full Engine Suite Complete!" -ForegroundColor Cyan
            pause
        }
        'Q' { Write-Host "Exiting orchestrator." -ForegroundColor Yellow; break }
        default { Write-Host "Invalid option. Please try again." -ForegroundColor Red; Start-Sleep -Seconds 1 }
    }
} while ($choice -ne 'Q')
