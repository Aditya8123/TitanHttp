# stats.ps1 - TitanHTTP Statistical Benchmarking Suite
# Runs benchmarks 10 times and computes statistical data (average, variance, std dev, min/max, CoV)

param (
    [switch]$Frameworks,
    [switch]$Internal,
    [switch]$All,
    [int]$Count = 10
)

# If no switch specified, run all
if (-not $Frameworks -and -not $Internal -and -not $All) {
    $All = $true
}

$ErrorActionPreference = "Stop"

# Ensure output directories exist
$baseDir = "$PSScriptRoot/../../benchmarks"
$dataDir = "$baseDir/data"
$reportsDir = "$baseDir/reports"

if (-not (Test-Path $dataDir)) { New-Item -ItemType Directory -Path $dataDir -Force | Out-Null }
if (-not (Test-Path $reportsDir)) { New-Item -ItemType Directory -Path $reportsDir -Force | Out-Null }

Write-Host "============================================================" -ForegroundColor Cyan
Write-Host "::      TITANHTTP STATISTICAL BENCHMARK ENGINE (Count=$Count)      ::" -ForegroundColor Yellow
Write-Host "============================================================" -ForegroundColor Cyan

# Stats computation helper function
function Get-Stats {
    param(
        [double[]]$Values
    )
    if ($Values.Count -eq 0) {
        return @{ Avg = 0; Var = 0; StdDev = 0; Min = 0; Max = 0; CoV = 0; Median = 0; P50 = 0; P90 = 0; P95 = 0; P99 = 0 }
    }
    
    $sum = 0
    $min = $Values[0]
    $max = $Values[0]
    foreach ($val in $Values) {
        $sum += $val
        if ($val -lt $min) { $min = $val }
        if ($val -gt $max) { $max = $val }
    }
    $avg = $sum / $Values.Count
    
    $varianceSum = 0
    foreach ($val in $Values) {
        $varianceSum += [math]::Pow($val - $avg, 2)
    }
    $variance = if ($Values.Count -gt 1) { $varianceSum / ($Values.Count - 1) } else { 0 }
    $stdDev = [math]::Sqrt($variance)
    $cov = if ($avg -gt 0) { ($stdDev / $avg) * 100 } else { 0 }
    
    # Calculate Percentiles
    $sorted = $Values | Sort-Object
    $n = $sorted.Count
    $getPercentile = {
        param([double]$p)
        $idx = ($p / 100) * ($n - 1)
        $lower = [math]::Floor($idx)
        $upper = [math]::Ceiling($idx)
        if ($lower -eq $upper) {
            return $sorted[$lower]
        }
        $fraction = $idx - $lower
        return $sorted[$lower] + $fraction * ($sorted[$upper] - $sorted[$lower])
    }
    
    $median = &$getPercentile 50
    $p50 = $median
    $p90 = &$getPercentile 90
    $p95 = &$getPercentile 95
    $p99 = &$getPercentile 99
    
    return @{
        Avg = [math]::Round($avg, 2)
        Var = [math]::Round($variance, 4)
        StdDev = [math]::Round($stdDev, 2)
        Min = [math]::Round($min, 2)
        Max = [math]::Round($max, 2)
        CoV = [math]::Round($cov, 2)
        Median = [math]::Round($median, 2)
        P50 = [math]::Round($p50, 2)
        P90 = [math]::Round($p90, 2)
        P95 = [math]::Round($p95, 2)
        P99 = [math]::Round($p99, 2)
    }
}

$rawLines = @()

# ---------------------------------------------------------
# STAGE 1: Run Framework Benchmarks
# ---------------------------------------------------------
if ($Frameworks -or $All) {
    Write-Host "`n[+] Executing Framework Benchmarks..." -ForegroundColor Green
    $fwDir = "$PSScriptRoot/../../benchmarks/frameworks"
    $fwRawFile = "$dataDir/raw_frameworks.bench"
    
    # Run benchmarks using cmd to ensure clean UTF-8 formatting and avoid PowerShell formatting issues
    Push-Location $fwDir
    Write-Host "  -> Running 'go test' ($Count iterations)..." -ForegroundColor Gray
    cmd /c "go test -bench=Benchmark -benchmem -count=$Count -run=^$ > `"$fwRawFile`""
    Pop-Location
    
    if (Test-Path $fwRawFile) {
        $rawLines += Get-Content $fwRawFile
    }
}

# ---------------------------------------------------------
# STAGE 2: Run Internal Microbenchmarks
# ---------------------------------------------------------
if ($Internal -or $All) {
    Write-Host "`n[+] Executing Internal Microbenchmarks..." -ForegroundColor Green
    $internalRawFile = "$dataDir/raw_internal.bench"
    
    Push-Location "$PSScriptRoot/../.."
    Write-Host "  -> Running 'go test ./internal/...' ($Count iterations)..." -ForegroundColor Gray
    cmd /c "go test -bench=Benchmark -benchmem -count=$Count -run=^$ ./internal/... > `"$internalRawFile`""
    Pop-Location
    
    if (Test-Path $internalRawFile) {
        $rawLines += Get-Content $internalRawFile
    }
}

# ---------------------------------------------------------
# STAGE 3: Parse Benchmark Output
# ---------------------------------------------------------
Write-Host "`n[+] Parsing raw benchmark outputs..." -ForegroundColor Green

$benchmarkGroups = @{}

foreach ($line in $rawLines) {
    # Match standard Go benchmark output lines
    # e.g., BenchmarkRouterStatic-16   13410976   95.39 ns/op   4 B/op   1 allocs/op
    if ($line -match "^\s*(Benchmark\S+)\s+\d+\s+([0-9.]+)\s+ns/op(?:\s+(\d+)\s+B/op)?(?:\s+(\d+)\s+allocs/op)?") {
        $name = $matches[1]
        $nsOp = [double]$matches[2]
        $bOp = if ($matches[3]) { [double]$matches[3] } else { 0.0 }
        $allocsOp = if ($matches[4]) { [double]$matches[4] } else { 0.0 }
        
        if (-not $benchmarkGroups.ContainsKey($name)) {
            $benchmarkGroups[$name] = @{
                Times = [System.Collections.Generic.List[double]]::new()
                Bytes = [System.Collections.Generic.List[double]]::new()
                Allocs = [System.Collections.Generic.List[double]]::new()
            }
        }
        
        $benchmarkGroups[$name].Times.Add($nsOp)
        $benchmarkGroups[$name].Bytes.Add($bOp)
        $benchmarkGroups[$name].Allocs.Add($allocsOp)
    }
}

if ($benchmarkGroups.Count -eq 0) {
    Write-Error "No benchmarks were successfully parsed from the test runs."
    exit 1
}

# ---------------------------------------------------------
# STAGE 4: Calculate Descriptive Statistics & Ranking
# ---------------------------------------------------------
Write-Host "[+] Computing statistical metrics and ranking..." -ForegroundColor Green

$summaries = @()

foreach ($entry in $benchmarkGroups.GetEnumerator()) {
    $name = $entry.Key
    $data = $entry.Value
    
    $timeStats = Get-Stats -Values $data.Times.ToArray()
    $byteStats = Get-Stats -Values $data.Bytes.ToArray()
    $allocStats = Get-Stats -Values $data.Allocs.ToArray()
    
    $summaries += [PSCustomObject]@{
        Name      = $name
        Time      = $timeStats
        Memory    = $byteStats
        Allocs    = $allocStats
        RawTimes  = $data.Times.ToArray()
        RawBytes  = $data.Bytes.ToArray()
        RawAllocs = $data.Allocs.ToArray()
    }
}

# Rank benchmarks by Median Time (P50) (fastest first)
$rankedSummaries = $summaries | Sort-Object { $_.Time.Median }
for ($i = 0; $i -lt $rankedSummaries.Count; $i++) {
    $rankedSummaries[$i] | Add-Member -MemberType NoteProperty -Name "Rank" -Value ($i + 1)
}

$reportTXT = @()
$reportTXT += "============================================================"
$reportTXT += "TITANHTTP BENCHMARK STATISTICAL REPORT"
$reportTXT += "Generated on: $(Get-Date)"
$reportTXT += "Count: $Count iterations"
$reportTXT += "============================================================"
$reportTXT += ""
$reportTXT += "BENCHMARK PERFORMANCE RANKING (Fastest to Slowest)"
$reportTXT += ("=" * 100)

$rankingFormat = "{0,-6} | {1,-45} | {2,-15} | {3,-15} | {4,-10}"
$reportTXT += $rankingFormat -f "Rank", "Benchmark Name", "Median Time", "Memory", "Allocs"
$reportTXT += ("-" * 100)
foreach ($s in $rankedSummaries) {
    $timeStr = "$($s.Time.Median) ns"
    $byteStr = "$($s.Memory.Avg) B"
    $allocStr = "$($s.Allocs.Avg) allocs"
    $reportTXT += $rankingFormat -f "#$($s.Rank)", $s.Name, $timeStr, $byteStr, $allocStr
}
$reportTXT += ""
$reportTXT += ""

# Format console table headers
Write-Host "`nBENCHMARK PERFORMANCE RANKING (Fastest to Slowest)" -ForegroundColor Yellow
Write-Host ("-" * 100) -ForegroundColor DarkGray
Write-Host ($rankingFormat -f "Rank", "Benchmark Name", "Median Time", "Memory", "Allocs") -ForegroundColor Cyan
Write-Host ("-" * 100) -ForegroundColor DarkGray
foreach ($s in $rankedSummaries) {
    $timeStr = "$($s.Time.Median) ns"
    $byteStr = "$($s.Memory.Avg) B"
    $allocStr = "$($s.Allocs.Avg) allocs"
    # Highlight TitanHTTP in green, others in white
    $color = if ($s.Name -like "*TitanHTTP*") { "Green" } else { "White" }
    Write-Host ($rankingFormat -f "#$($s.Rank)", $s.Name, $timeStr, $byteStr, $allocStr) -ForegroundColor $color
}
Write-Host ("-" * 100) -ForegroundColor DarkGray
Write-Host "`n"

foreach ($s in $rankedSummaries) {
    # Save statistics and raw data in text report
    $reportTXT += "Benchmark: $($s.Name) (Rank: #$($s.Rank))"
    $reportTXT += ("-" * 60)
    $reportTXT += "Raw Data:"
    $reportTXT += "  Time (ns/op):   [$($s.RawTimes -join ', ')]"
    $reportTXT += "  Memory (B/op):  [$($s.RawBytes -join ', ')]"
    $reportTXT += "  Allocs (op):    [$($s.RawAllocs -join ', ')]"
    $reportTXT += ""
    $reportTXT += "Calculated Statistics:"
    $reportTXT += "  * Time (ns/op):"
    $reportTXT += "    - Average:         $($s.Time.Avg) ns/op"
    $reportTXT += "    - Median / P50:    $($s.Time.Median) ns/op"
    $reportTXT += "    - P90:             $($s.Time.P90) ns/op"
    $reportTXT += "    - P95:             $($s.Time.P95) ns/op"
    $reportTXT += "    - P99:             $($s.Time.P99) ns/op"
    $reportTXT += "    - Variance:        $($s.Time.Var)"
    $reportTXT += "    - Std Dev:         $($s.Time.StdDev) ns/op"
    $reportTXT += "    - Min:             $($s.Time.Min) ns/op"
    $reportTXT += "    - Max:             $($s.Time.Max) ns/op"
    $reportTXT += "    - CoV (Std Dev %): $($s.Time.CoV)%"
    $reportTXT += "  * Memory (B/op):"
    $reportTXT += "    - Average:         $($s.Memory.Avg) B/op"
    $reportTXT += "    - Median / P50:    $($s.Memory.Median) B/op"
    $reportTXT += "    - Variance:        $($s.Memory.Var)"
    $reportTXT += "    - Std Dev:         $($s.Memory.StdDev) B/op"
    $reportTXT += "    - Min:             $($s.Memory.Min) B/op"
    $reportTXT += "    - Max:             $($s.Memory.Max) B/op"
    $reportTXT += "    - CoV (Std Dev %): $($s.Memory.CoV)%"
    $reportTXT += "  * Allocs (op):"
    $sAvg = $s.Allocs.Avg
    $sMedian = $s.Allocs.Median
    $sVar = $s.Allocs.Var
    $sStdDev = $s.Allocs.StdDev
    $sMin = $s.Allocs.Min
    $sMax = $s.Allocs.Max
    $sCoV = $s.Allocs.CoV
    $reportTXT += "    - Average:         $sAvg allocs/op"
    $reportTXT += "    - Median / P50:    $sMedian allocs/op"
    $reportTXT += "    - Variance:        $sVar"
    $reportTXT += "    - Std Dev:         $sStdDev allocs/op"
    $reportTXT += "    - Min:             $sMin allocs/op"
    $reportTXT += "    - Max:             $sMax allocs/op"
    $reportTXT += "    - CoV (Std Dev %): $sCoV%"
    $reportTXT += ""
    $reportTXT += ""
}

# Save text report
$txtPath = "$reportsDir/stats_report.txt"
$reportTXT -join "`r`n" | Out-File -FilePath $txtPath -Encoding utf8

# Clean up any leftover stats JSON/MD files if they exist to keep workspace clean
$jsonPath = "$dataDir/stats_report.json"
$mdPath = "$reportsDir/stats_report.md"
if (Test-Path $jsonPath) { Remove-Item $jsonPath -Force }
if (Test-Path $mdPath) { Remove-Item $mdPath -Force }

Write-Host "[+] Statistical analysis report saved to:" -ForegroundColor Green
Write-Host "  -> Text Report: $txtPath" -ForegroundColor White
Write-Host "============================================================" -ForegroundColor Cyan
