$ErrorActionPreference = "Stop"

Write-Host "`n  ============================================================" -ForegroundColor Cyan
Write-Host "  ::                                                        ::" -ForegroundColor Cyan
Write-Host "  ::         TITANHTTP VS FRAMEWORKS CINEMATIC BRAWL        ::" -ForegroundColor Yellow
Write-Host "  ::                                                        ::" -ForegroundColor Cyan
Write-Host "  ============================================================`n" -ForegroundColor Cyan
Write-Host "Setting up isolated environment for Framework Comparisons..."
Push-Location ".\benchmarks\frameworks"

Write-Host "Downloading dependencies (Gin, Fiber, Echo, Chi)..."
go mod tidy

Write-Host "`nRunning Framework Benchmarks (10 iterations)..."
$benchFile = "../../benchmarks/data/frameworks_latest.bench"
if (Test-Path $benchFile) { Remove-Item $benchFile }
go test -bench="." -benchmem -benchtime=5s -count=10 | ForEach-Object {
    Write-Host $_
    Add-Content -Path $benchFile -Value $_ -Encoding Ascii
}

$jsonArr = @()
foreach ($line in (Get-Content $benchFile)) {
    if ($line -match "^(Benchmark\S+)\s+\d+\s+([0-9.]+)\s+ns/op\s+(\d+)\s+B/op\s+(\d+)\s+allocs/op") {
        $jsonArr += [PSCustomObject]@{
            Name = $matches[1]
            TimeNs = [double]$matches[2]
            MemoryBytes = [int]$matches[3]
            Allocs = [int]$matches[4]
        }
    }
}
$jsonArr | ConvertTo-Json | Out-File -FilePath "../../benchmarks/data/frameworks_latest.json" -Encoding utf8


Write-Host "`nAnalyzing Results with benchstat..."
$benchstatOutput = benchstat $benchFile
$benchstatOutput | ForEach-Object { Write-Host $_ }

$reportPath = "../../benchmarks/reports/framework_comparison.md"
$mdLines = @(
    "# Framework Comparison Report"
    "Generated on $(Get-Date)"
    ""
    "## Performance Metrics (benchstat)"
    ""
    '```text'
    $benchstatOutput
    '```'
)
$reportContent = $mdLines -join "`r`n"

$reportContent | Out-File -FilePath $reportPath -Encoding utf8

Pop-Location
Write-Host "`nFramework Comparison Complete."
