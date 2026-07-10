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
    Write-Host "    [1] Run Micro & Component Benchmarks (Go Tests)" -ForegroundColor White
    Write-Host "    [2] Run Framework Comparisons (net/http, Gin, Fiber, Chi)" -ForegroundColor White
    Write-Host "    [3] Run Unified Cinematic Suite (TitanHTTP vs net/http)" -ForegroundColor White
    Write-Host "    [4] HTTP Compliance Tests" -ForegroundColor White
    Write-Host "    [5] Unified Security Suite`n" -ForegroundColor White
    Write-Host "    [A] Run Full Cinematic Engine Suite" -ForegroundColor Magenta
    Write-Host "    [Q] Quit`n" -ForegroundColor DarkGray
}

function Format-GoBench {
    param([string[]]$Lines)
    Write-Host "      $('-'*85)" -ForegroundColor DarkGray
    Write-Host ("      {0,-40} | {1,12} | {2,10} | {3,10}" -f "Benchmark Name", "Time (ns/op)", "Memory (B/op)", "Allocations") -ForegroundColor Yellow
    Write-Host "      $('-'*85)" -ForegroundColor DarkGray
    foreach ($line in $Lines) {
        if ($line -match "^(Benchmark\S+)\s+\d+\s+([0-9.]+)\s+ns/op\s+(\d+)\s+B/op\s+(\d+)\s+allocs/op") {
            Write-Host ("      {0,-40} | {1,12} | {2,10} | {3,10}" -f $matches[1], $matches[2], $matches[3], $matches[4]) -ForegroundColor Cyan
        } elseif ($line -match "^FAIL") {
            Write-Host "      $line" -ForegroundColor Red
        }
    }
    Write-Host "      $('-'*85)" -ForegroundColor DarkGray
}

function Export-GoBenchJson {
    param([string[]]$Lines, [string]$OutPath)
    $arr = @()
    foreach ($line in $Lines) {
        if ($line -match "^(Benchmark\S+)\s+\d+\s+([0-9.]+)\s+ns/op\s+(\d+)\s+B/op\s+(\d+)\s+allocs/op") {
            $arr += [PSCustomObject]@{
                Name = $matches[1]
                TimeNs = [double]$matches[2]
                MemoryBytes = [int]$matches[3]
                Allocs = [int]$matches[4]
            }
        }
    }
    $arr | ConvertTo-Json | Out-File -FilePath $OutPath -Encoding utf8
}

function Run-Micro {
    Write-Host "`n==========================================================" -ForegroundColor Cyan
    Write-Host "         Micro & Component Benchmarks (Go Tests)          " -ForegroundColor Yellow
    Write-Host "==========================================================" -ForegroundColor Cyan
    
    Write-Host "`n[+] Running Micro Benchmarks..." -ForegroundColor Magenta
    
    $rawMicro = go test "-bench=." -benchmem ./internal/...
    Export-GoBenchJson -Lines $rawMicro -OutPath "benchmarks/data/micro_latest.json"
    Format-GoBench -Lines $rawMicro
    
    Write-Host "`n[+] Running Component Benchmarks..." -ForegroundColor Magenta
    
    $rawComp = go test "-bench=." -benchmem "-run=^$" ./internal/server/...
    Export-GoBenchJson -Lines $rawComp -OutPath "benchmarks/data/component_latest.json"
    Format-GoBench -Lines $rawComp
    
    # Generate micro benchmarks markdown report
    $reportPath = "benchmarks/reports/micro_benchmarks.md"
    $mdLines = @(
        "# Micro & Component Benchmarks Report"
        "Generated on $(Get-Date)"
        ""
        "## Micro Benchmarks"
        "| Benchmark Name | Time (ns/op) | Memory (B/op) | Allocations |"
        "| :--- | :--- | :--- | :--- |"
    )
    foreach ($line in $rawMicro) {
        if ($line -match "^(Benchmark\S+)\s+\d+\s+([0-9.]+)\s+ns/op\s+(\d+)\s+B/op\s+(\d+)\s+allocs/op") {
            $mdLines += "| $($matches[1]) | $($matches[2]) | $($matches[3]) | $($matches[4]) |"
        }
    }
    
    $mdLines += ""
    $mdLines += "## Component Benchmarks"
    $mdLines += "| Benchmark Name | Time (ns/op) | Memory (B/op) | Allocations |"
    $mdLines += "| :--- | :--- | :--- | :--- |"
    
    foreach ($line in $rawComp) {
        if ($line -match "^(Benchmark\S+)\s+\d+\s+([0-9.]+)\s+ns/op\s+(\d+)\s+B/op\s+(\d+)\s+allocs/op") {
            $mdLines += "| $($matches[1]) | $($matches[2]) | $($matches[3]) | $($matches[4]) |"
        }
    }
    
    ($mdLines -join "`r`n") | Out-File -FilePath $reportPath -Encoding utf8
    Write-Host "`nCompleted Micro & Component Tests.`n" -ForegroundColor Green
}

function Run-Frameworks {
    Write-Host "`n>>> Running Framework Comparisons..." -ForegroundColor Green
    & .\scripts\bench\compare_frameworks.ps1
    Write-Host "Completed Framework Comparisons.`n" -ForegroundColor Green
}

function Run-UnifiedSuite {
    Write-Host "`n>>> Running Unified Comprehensive Suite..." -ForegroundColor Green
    & .\scripts\bench\unified_comparison.ps1
    Write-Host "Completed Unified Suite.`n" -ForegroundColor Green
}

function Run-Compliance {
    Write-Host "`n>>> Running HTTP Compliance Tests..." -ForegroundColor Green
    $TitanProcess = Start-Process -FilePath ".\bin\titanhttp_bench.exe" -PassThru -NoNewWindow -RedirectStandardOutput "benchmarks/logs/titan_compliance.log" -RedirectStandardError "benchmarks/logs/titan_compliance_err.log"
    Start-Sleep -Seconds 2
    
    $compOutput = go run ./cmd/compliance/main.go localhost:8080 2>&1
    $compOutput | ForEach-Object { Write-Host $_ }
    
    Stop-Process -Id $TitanProcess.Id -Force -ErrorAction SilentlyContinue
    
    # Write report
    $reportPath = "benchmarks/reports/compliance_report.md"
    $mdLines = @(
        "# HTTP Compliance Tests Report"
        "Generated on $(Get-Date)"
        ""
        "## Test Results"
        ""
        '```text'
        $compOutput
        '```'
    )
    ($mdLines -join "`r`n") | Out-File -FilePath $reportPath -Encoding utf8
    Write-Host "Completed Compliance Tests.`n" -ForegroundColor Green
}

function Run-SecuritySuite {
    Write-Host "`n>>> Running Unified Security Suite..." -ForegroundColor Green
    $env:IDLE_TIMEOUT = "5s"
    $TitanProcess = Start-Process -FilePath ".\bin\titanhttp_bench.exe" -PassThru -NoNewWindow -RedirectStandardOutput "benchmarks/logs/titan_security.log" -RedirectStandardError "benchmarks/logs/titan_security_err.log"
    Start-Sleep -Seconds 2
    
    $secOutput = go run ./cmd/securitysuite/main.go localhost:8080 2>&1
    $secOutput | ForEach-Object { Write-Host $_ }
    
    Remove-Item env:IDLE_TIMEOUT -ErrorAction SilentlyContinue
    Stop-Process -Id $TitanProcess.Id -Force -ErrorAction SilentlyContinue
    
    # Write report
    $reportPath = "benchmarks/reports/security_report.md"
    $mdLines = @(
        "# Unified Security Suite Report"
        "Generated on $(Get-Date)"
        ""
        "## Test Results"
        ""
        '```text'
        $secOutput
        '```'
    )
    ($mdLines -join "`r`n") | Out-File -FilePath $reportPath -Encoding utf8
    Write-Host "Completed Security Suite.`n" -ForegroundColor Green
}

Write-Host "Building TitanHTTP for benchmarking..." -ForegroundColor Cyan
if (-not (Test-Path "bin")) { New-Item -ItemType Directory -Path "bin" | Out-Null }
if (-not (Test-Path "benchmarks/data")) { New-Item -ItemType Directory -Path "benchmarks/data" -Force | Out-Null }
if (-not (Test-Path "benchmarks/logs")) { New-Item -ItemType Directory -Path "benchmarks/logs" -Force | Out-Null }
if (-not (Test-Path "benchmarks/reports")) { New-Item -ItemType Directory -Path "benchmarks/reports" -Force | Out-Null }

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
            Write-Host "Full Suite Complete!" -ForegroundColor Cyan
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
            Write-Host "Full Suite Complete!" -ForegroundColor Cyan
            pause
        }
        'Q' { Write-Host "Exiting orchestrator." -ForegroundColor Yellow; break }
        default { Write-Host "Invalid option. Please try again." -ForegroundColor Red; Start-Sleep -Seconds 1 }
    }
} while ($choice -ne 'Q')
