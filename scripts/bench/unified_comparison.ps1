param (
    [int]$WarmupDuration = 5,
    [int]$CooldownSeconds = 2
)

$ErrorActionPreference = "Stop"

Write-Host "`n  ============================================================" -ForegroundColor Cyan
Write-Host "  ::                                                        ::" -ForegroundColor Cyan
Write-Host "  ::       TITANHTTP VS NET/HTTP PERFORMANCE LABORATORY     ::" -ForegroundColor Yellow
Write-Host "  ::                                                        ::" -ForegroundColor Cyan
Write-Host "  ============================================================`n" -ForegroundColor Cyan

# Ensure outputs directory
if (-not (Test-Path "benchmarks/data")) { New-Item -ItemType Directory -Path "benchmarks/data" -Force | Out-Null }
if (-not (Test-Path "benchmarks/logs")) { New-Item -ItemType Directory -Path "benchmarks/logs" -Force | Out-Null }
if (-not (Test-Path "benchmarks/reports")) { New-Item -ItemType Directory -Path "benchmarks/reports" -Force | Out-Null }
if (-not (Test-Path "benchmarks/profiles")) { New-Item -ItemType Directory -Path "benchmarks/profiles" -Force | Out-Null }

$payloadJson = "benchmarks/data/payload.json"
if (-not (Test-Path $payloadJson)) {
    Set-Content -Path $payloadJson -Value '{"name":"titan","value":123}' -Encoding utf8
}

$AllResults = [System.Collections.ArrayList]::new()

# Build servers
Write-Host "`n[+] Building servers..."
go build -o bin/titanhttp_bench.exe ./cmd/titanhttp
go build -o bin/nethttp_bench.exe ./cmd/nethttp_bench

# --- Helpers ---

function Get-Stats {
    param([double[]]$Values)
    if ($Values.Count -eq 0) {
        return @{ Mean = 0; StdDev = 0; Min = 0; Max = 0; Median = 0; P95 = 0; P99 = 0; Range = 0 }
    }
    
    $sum = 0
    $min = $Values[0]
    $max = $Values[0]
    foreach ($val in $Values) {
        $sum += $val
        if ($val -lt $min) { $min = $val }
        if ($val -gt $max) { $max = $val }
    }
    $mean = $sum / $Values.Count
    
    $varianceSum = 0
    foreach ($val in $Values) { $varianceSum += [math]::Pow($val - $mean, 2) }
    $variance = if ($Values.Count -gt 1) { $varianceSum / ($Values.Count - 1) } else { 0 }
    $stdDev = [math]::Sqrt($variance)
    
    $sorted = $Values | Sort-Object
    $n = $sorted.Count
    $getPercentile = {
        param([double]$p)
        $idx = ($p / 100) * ($n - 1)
        $lower = [math]::Floor($idx)
        $upper = [math]::Ceiling($idx)
        if ($lower -eq $upper) { return $sorted[$lower] }
        $fraction = $idx - $lower
        return $sorted[$lower] + $fraction * ($sorted[$upper] - $sorted[$lower])
    }
    
    return @{
        Mean = [math]::Round($mean, 2)
        StdDev = [math]::Round($stdDev, 2)
        Min = [math]::Round($min, 2)
        Max = [math]::Round($max, 2)
        Range = [math]::Round($max - $min, 2)
        Median = [math]::Round((&$getPercentile 50), 2)
        P95 = [math]::Round((&$getPercentile 95), 2)
        P99 = [math]::Round((&$getPercentile 99), 2)
    }
}

function Run-BombardierSingle {
    param ($Url, $ArgsStr)
    $cmd = "bombardier -p r -o json -l $ArgsStr `"$Url`""
    $jsonRaw = Invoke-Expression "$cmd 2> `$null"
    $res = $jsonRaw | ConvertFrom-Json
    
    $p50 = [math]::Round($res.result.latency.percentiles."50" / 1000, 2)
    $p95 = [math]::Round($res.result.latency.percentiles."95" / 1000, 2)
    $p99 = [math]::Round($res.result.latency.percentiles."99" / 1000, 2)
    
    $errors = 0
    if ($res.result.errors) {
        $errors = ($res.result.req1xx + $res.result.req2xx + $res.result.req3xx + $res.result.req4xx + $res.result.req5xx) 
        $errors = $res.result.reqsTotal - $errors
    }
    
    return @{
        Rps = [math]::Round($res.result.rps.mean, 2)
        P50 = $p50
        P95 = $p95
        P99 = $p99
        Errors = $errors
        ThroughputMB = [math]::Round($res.result.throughput / 1MB, 2)
    }
}

function Warmup {
    param($Url)
    Write-Host "  -> Warming up $Url for $WarmupDuration seconds..." -ForegroundColor DarkGray
    $cmd = "bombardier -c 50 -d $($WarmupDuration)s `"$Url`""
    Invoke-Expression "$cmd > `$null 2> `$null"
    Start-Sleep -Seconds $CooldownSeconds
}

function Run-Comparison {
    param ($Name, $Path, $ArgsStr, $Iterations, $Category, $Scheme="http")
    
    Write-Host "`n--- Testing $Name ---" -ForegroundColor Yellow
    Warmup -Url "$Scheme://localhost:8080$Path"
    Warmup -Url "$Scheme://localhost:8081$Path"
    
    $tRpsList = [System.Collections.Generic.List[double]]::new()
    $tP50List = [System.Collections.Generic.List[double]]::new()
    $tP95List = [System.Collections.Generic.List[double]]::new()
    $tP99List = [System.Collections.Generic.List[double]]::new()
    $tMbList  = [System.Collections.Generic.List[double]]::new()

    $nRpsList = [System.Collections.Generic.List[double]]::new()
    $nP50List = [System.Collections.Generic.List[double]]::new()
    $nP95List = [System.Collections.Generic.List[double]]::new()
    $nP99List = [System.Collections.Generic.List[double]]::new()
    $nMbList  = [System.Collections.Generic.List[double]]::new()
    
    for ($i = 1; $i -le $Iterations; $i++) {
        Write-Host "     [Run $i/$Iterations] $Name..." -NoNewline
        # Alternate order to avoid systemic bias
        if ($i % 2 -eq 1) {
            $t = Run-BombardierSingle -Url "$Scheme://localhost:8080$Path" -ArgsStr $ArgsStr
            $n = Run-BombardierSingle -Url "$Scheme://localhost:8081$Path" -ArgsStr $ArgsStr
        } else {
            $n = Run-BombardierSingle -Url "$Scheme://localhost:8081$Path" -ArgsStr $ArgsStr
            $t = Run-BombardierSingle -Url "$Scheme://localhost:8080$Path" -ArgsStr $ArgsStr
        }
        $tRpsList.Add($t.Rps); $tP50List.Add($t.P50); $tP95List.Add($t.P95); $tP99List.Add($t.P99); $tMbList.Add($t.ThroughputMB)
        $nRpsList.Add($n.Rps); $nP50List.Add($n.P50); $nP95List.Add($n.P95); $nP99List.Add($n.P99); $nMbList.Add($n.ThroughputMB)
        Write-Host " Done."
        if ($i -lt $Iterations) { Start-Sleep -Seconds $CooldownSeconds }
    }
    
    $tStats = @{
        RpsSamples = $tRpsList.ToArray(); RpsStats = Get-Stats $tRpsList.ToArray()
        P50Samples = $tP50List.ToArray(); P50Stats = Get-Stats $tP50List.ToArray()
        P95Samples = $tP95List.ToArray(); P95Stats = Get-Stats $tP95List.ToArray()
        P99Samples = $tP99List.ToArray(); P99Stats = Get-Stats $tP99List.ToArray()
        MbSamples  = $tMbList.ToArray();  MbStats  = Get-Stats $tMbList.ToArray()
    }
    
    $nStats = @{
        RpsSamples = $nRpsList.ToArray(); RpsStats = Get-Stats $nRpsList.ToArray()
        P50Samples = $nP50List.ToArray(); P50Stats = Get-Stats $nP50List.ToArray()
        P95Samples = $nP95List.ToArray(); P95Stats = Get-Stats $nP95List.ToArray()
        P99Samples = $nP99List.ToArray(); P99Stats = Get-Stats $nP99List.ToArray()
        MbSamples  = $nMbList.ToArray();  MbStats  = Get-Stats $nMbList.ToArray()
    }
    
    $winner = if ($tStats.RpsStats.Median -gt $nStats.RpsStats.Median) { "TitanHTTP" } else { "net/http" }
    
    $diff = 0
    if ($nStats.RpsStats.Median -gt 0) {
        $diff = [math]::Round((($tStats.RpsStats.Median - $nStats.RpsStats.Median) / $nStats.RpsStats.Median) * 100, 1)
    }
    
    $diffText = if ($diff -ge 0) { "+$($diff)%" } else { "$($diff)%" }
    Write-Host "[ WINNER ] $winner ($diffText)" -ForegroundColor Green
    
    $AllResults.Add(@{
        Name = $Name
        Category = $Category
        TitanHTTP = $tStats
        NetHttp = $nStats
        Winner = $winner
        DifferenceVal = $diff
        Difference = $diffText
    }) | Out-Null
}

# ---------------------------------------------------------
# STAGE 1: Standard Tests
# ---------------------------------------------------------
Write-Host "`n>>> [ STAGE 1 ] Starting Standard HTTP Servers..." -ForegroundColor Green
$TitanProcess = Start-Process -FilePath ".\bin\titanhttp_bench.exe" -PassThru -NoNewWindow -RedirectStandardOutput "benchmarks/logs/titan_bench.log" -RedirectStandardError "benchmarks/logs/titan_bench_err.log"
$NetHttpProcess = Start-Process -FilePath ".\bin\nethttp_bench.exe" -PassThru -NoNewWindow -RedirectStandardOutput "benchmarks/logs/nethttp_bench.log" -RedirectStandardError "benchmarks/logs/nethttp_bench_err.log"
Start-Sleep -Seconds 4

Run-Comparison -Name "Ping (Req/s)" -Path "/ping" -ArgsStr "-c 125 -d 5s" -Iterations 5 -Category "Standard"
Run-Comparison -Name "Concurrent Connections" -Path "/ping" -ArgsStr "-c 1000 -d 5s" -Iterations 5 -Category "Standard"
Run-Comparison -Name "Large Payloads" -Path "/heavy" -ArgsStr "-c 125 -d 5s" -Iterations 5 -Category "Standard"
Run-Comparison -Name "Routing Performance" -Path "/users/123" -ArgsStr "-c 125 -d 5s" -Iterations 5 -Category "Standard"
Run-Comparison -Name "Static File Serving" -Path "/static/test.txt" -ArgsStr "-c 125 -d 5s" -Iterations 5 -Category "Standard"
Run-Comparison -Name "Keep-Alive Performance" -Path "/ping" -ArgsStr "-c 125 -d 5s -a" -Iterations 5 -Category "Standard"

# ---------------------------------------------------------
# STAGE 2: Concurrency Scale (C10K Simulator)
# ---------------------------------------------------------
Write-Host "`n>>> [ STAGE 2 ] Concurrency Scale (C10K Simulator)..." -ForegroundColor Green
$ConcurrencyLevels = @(50, 100, 200)
foreach ($c in $ConcurrencyLevels) {
    Run-Comparison -Name "Concurrency $c" -Path "/ping" -ArgsStr "-c $c -d 5s" -Iterations 5 -Category "Concurrency"
}

# ---------------------------------------------------------
# STAGE 3: Payload Matrix
# ---------------------------------------------------------
Write-Host "`n>>> [ STAGE 3 ] Payload Matrix..." -ForegroundColor Green
$Payloads = @(
    @{ Name="16B"; Size=16 },
    @{ Name="1KB"; Size=1024 },
    @{ Name="10KB"; Size=10240 },
    @{ Name="100KB"; Size=102400 },
    @{ Name="1MB"; Size=1048576 },
    @{ Name="10MB"; Size=10485760 }
)
foreach ($p in $Payloads) {
    Run-Comparison -Name "Payload $($p.Name)" -Path "/payload/$($p.Size)" -ArgsStr "-c 125 -d 5s" -Iterations 5 -Category "Payload"
}

# ---------------------------------------------------------
# STAGE 4: POST JSON
# ---------------------------------------------------------
Write-Host "`n>>> [ STAGE 4 ] POST JSON..." -ForegroundColor Green
$jsonArgs = "-m POST -H `"Content-Type: application/json`" -f `"$payloadJson`" -c 125 -d 5s"
Run-Comparison -Name "POST JSON" -Path "/api/users" -ArgsStr $jsonArgs -Iterations 5 -Category "POST"

# ---------------------------------------------------------
# STAGE 5: Connection Churn
# ---------------------------------------------------------
Write-Host "`n>>> [ STAGE 5 ] Connection Churn (Connection: close)..." -ForegroundColor Green
Run-Comparison -Name "Connection Churn" -Path "/ping" -ArgsStr "-c 100 -d 5s -H `"Connection: close`"" -Iterations 3 -Category "Churn"


# ---------------------------------------------------------
# STAGE 6: Memory & CPU Profiling (Skip Warmup)
# ---------------------------------------------------------
Write-Host "`n>>> [ STAGE 6 ] Memory & CPU Leak Detector (pprof)..." -ForegroundColor Green
Write-Host "1. Forcing GC before capturing baselines..."
Invoke-WebRequest -UseBasicParsing -Uri "http://localhost:8080/debug/gc" -ErrorAction SilentlyContinue | Out-Null
Invoke-WebRequest -UseBasicParsing -Uri "http://localhost:8081/debug/gc" -ErrorAction SilentlyContinue | Out-Null

Write-Host "2. Capturing Baseline Heaps..."
Invoke-WebRequest -UseBasicParsing -Uri "http://localhost:6060/debug/pprof/heap" -OutFile "benchmarks/profiles/heap_titan_baseline.pprof"
Invoke-WebRequest -UseBasicParsing -Uri "http://localhost:6061/debug/pprof/heap" -OutFile "benchmarks/profiles/heap_nethttp_baseline.pprof"

Write-Host "3. Blasting 500k requests and capturing CPU profile simultaneously..."
# We run Bombardier and simultaneously grab a 10s CPU profile
$j1 = Start-Process -FilePath "bombardier" -ArgumentList "-c 125 -n 500000 http://localhost:8080/ping" -NoNewWindow -PassThru -RedirectStandardOutput "benchmarks/logs/bombardier_titan_mem.log"
$j2 = Start-Process -FilePath "bombardier" -ArgumentList "-c 125 -n 500000 http://localhost:8081/ping" -NoNewWindow -PassThru -RedirectStandardOutput "benchmarks/logs/bombardier_nethttp_mem.log"

# Capture 10 second CPU profile while blast runs
Invoke-WebRequest -UseBasicParsing -Uri "http://localhost:6060/debug/pprof/profile?seconds=10" -OutFile "benchmarks/profiles/cpu_titan.pprof" -ErrorAction SilentlyContinue | Out-Null
Invoke-WebRequest -UseBasicParsing -Uri "http://localhost:6061/debug/pprof/profile?seconds=10" -OutFile "benchmarks/profiles/cpu_nethttp.pprof" -ErrorAction SilentlyContinue | Out-Null

$j1.WaitForExit()
$j2.WaitForExit()

Write-Host "4. Forcing GC before capturing end profiles..."
Invoke-WebRequest -UseBasicParsing -Uri "http://localhost:8080/debug/gc" -ErrorAction SilentlyContinue | Out-Null
Invoke-WebRequest -UseBasicParsing -Uri "http://localhost:8081/debug/gc" -ErrorAction SilentlyContinue | Out-Null

Write-Host "5. Capturing End Heaps..."
Invoke-WebRequest -UseBasicParsing -Uri "http://localhost:6060/debug/pprof/heap" -OutFile "benchmarks/profiles/heap_titan_end.pprof"
Invoke-WebRequest -UseBasicParsing -Uri "http://localhost:6061/debug/pprof/heap" -OutFile "benchmarks/profiles/heap_nethttp_end.pprof"

$titanDiff = go tool pprof -text -base "benchmarks/profiles/heap_titan_baseline.pprof" "benchmarks/profiles/heap_titan_end.pprof"
$nethttpDiff = go tool pprof -text -base "benchmarks/profiles/heap_nethttp_baseline.pprof" "benchmarks/profiles/heap_nethttp_end.pprof"

function Parse-Leak { param($Output)
    $line = $Output | Where-Object { $_ -match "^Total:\s+(-?[\d\.]+MB|-?[\d\.]+kB|-?[\d\.]+B)" }
    if ($line -match "Total:\s+([\d\.]+)MB") { return [double]$matches[1] }
    return 0
}
$tLeak = Parse-Leak $titanDiff
$nLeak = Parse-Leak $nethttpDiff
Write-Host "TitanHTTP Heap Leak: $tLeak MB | net/http Heap Leak: $nLeak MB"
$MemoryResults = @{ TitanHTTP = $tLeak; NetHttp = $nLeak }

# Shutdown standard servers
Stop-Process -Id $TitanProcess.Id -Force -ErrorAction SilentlyContinue
Stop-Process -Id $NetHttpProcess.Id -Force -ErrorAction SilentlyContinue

# ---------------------------------------------------------
# STAGE 7: TLS Benchmark
# ---------------------------------------------------------
Write-Host "`n>>> [ STAGE 7 ] TLS Performance..." -ForegroundColor Green
if (-not (Test-Path "certs")) {
    Write-Host "Generating TLS certificates for benchmarks..."
    New-Item -ItemType Directory -Path "certs" -Force | Out-Null
    $goRoot = go env GOROOT
    go run "$goRoot\src\crypto\tls\generate_cert.go" --rsa-bits 2048 --host localhost
    Move-Item cert.pem certs/ -Force
    Move-Item key.pem certs/ -Force
}

$env:TLS_MODE = "1"
$env:TLS_CERT = "certs/cert.pem"
$env:TLS_KEY = "certs/key.pem"
$TitanTLS = Start-Process -FilePath ".\bin\titanhttp_bench.exe" -PassThru -NoNewWindow -RedirectStandardOutput "benchmarks/logs/titan_tls.log" -RedirectStandardError "benchmarks/logs/titan_tls_err.log"
$NetHttpTLS = Start-Process -FilePath ".\bin\nethttp_bench.exe" -PassThru -NoNewWindow -RedirectStandardOutput "benchmarks/logs/nethttp_tls.log" -RedirectStandardError "benchmarks/logs/nethttp_tls_err.log"
Start-Sleep -Seconds 4

Run-Comparison -Name "HTTPS Ping" -Path "/ping" -ArgsStr "-c 125 -d 5s --insecure" -Iterations 5 -Category "TLS" -Scheme "https"

Stop-Process -Id $TitanTLS.Id -Force -ErrorAction SilentlyContinue
Stop-Process -Id $NetHttpTLS.Id -Force -ErrorAction SilentlyContinue
Remove-Item env:TLS_MODE -ErrorAction SilentlyContinue
Remove-Item env:TLS_CERT -ErrorAction SilentlyContinue
Remove-Item env:TLS_KEY -ErrorAction SilentlyContinue


# ---------------------------------------------------------
# SCORECARD MATH
# ---------------------------------------------------------
$titanWins = 0
$nethttpWins = 0
$sumDiff = 0
$sumAbsDiff = 0
$largestWinVal = -999999
$largestWinName = ""
$closestRaceVal = 999999
$closestRaceName = ""

$categoryGains = @{}
$categoryCounts = @{}

foreach ($res in $AllResults) {
    if ($res.Winner -eq "TitanHTTP") { $titanWins++ } else { $nethttpWins++ }
    
    $d = $res.DifferenceVal
    $sumDiff += $d
    $sumAbsDiff += [math]::Abs($d)
    
    if ($d -gt $largestWinVal) {
        $largestWinVal = $d
        $largestWinName = $res.Name
    }
    
    if ([math]::Abs($d) -lt $closestRaceVal) {
        $closestRaceVal = [math]::Abs($d)
        $closestRaceName = $res.Name
    }
    
    if (-not $categoryGains.ContainsKey($res.Category)) {
        $categoryGains[$res.Category] = 0
        $categoryCounts[$res.Category] = 0
    }
    $categoryGains[$res.Category] += $d
    $categoryCounts[$res.Category]++
}

$avgGain = [math]::Round($sumDiff / $AllResults.Count, 1)
$avgSep = [math]::Round($sumAbsDiff / $AllResults.Count, 1)

# ---------------------------------------------------------
# GENERATE REPORT
# ---------------------------------------------------------
Write-Host "`n[+] Generating Unified Performance Report..." -ForegroundColor Green

$report = @()
$report += "========================================================="
$report += " TITANHTTP VS NET/HTTP PERFORMANCE REPORT"
$report += "========================================================="
$report += ""
$report += "Environment"
$report += "-----------"
$report += "CPU:        $((Get-CimInstance Win32_Processor).Name)"
$report += "OS:         $((Get-CimInstance Win32_OperatingSystem).Caption)"
$report += "Go version: $(go version)"
$report += "Date:       $(Get-Date)"
$report += "Warmup:     ${WarmupDuration}s"
$report += ""
$report += ""
$report += "========================================================="
$report += "FINAL SCORECARD"
$report += "========================================================="
$report += ""
$report += "TitanHTTP wins:            $titanWins / $($AllResults.Count)"
$report += "net/http wins:             $nethttpWins / $($AllResults.Count)"
$report += ""
$report += "Average TitanHTTP Gain:    $(if ($avgGain -ge 0) { `"+$avgGain%`" } else { `"$avgGain%`" })"
$report += "Avg. Performance Sep:      $avgSep%"
$report += ""
$report += "Largest win:               $largestWinName (+$(if($largestWinVal -gt 0){$largestWinVal}else{0})%)"
$report += "Closest race:              $closestRaceName ($closestRaceVal%)"
$report += ""
$report += "Category Averages (Signed Gain):"
foreach ($key in $categoryGains.Keys) {
    $cAvg = [math]::Round($categoryGains[$key] / $categoryCounts[$key], 1)
    $report += "  - $( '{0,-20}' -f $key ): $(if ($cAvg -ge 0) { `"+$cAvg%`" } else { `"$cAvg%`" })"
}
$report += ""
$report += ""
$report += "========================================================="
$report += "SUMMARY RANKING"
$report += "========================================================="
$report += ""
$reportFormat = "{0,-25} {1,-15} {2,-15}"
$report += $reportFormat -f "Category", "Winner", "Difference"
$report += ""
foreach ($res in $AllResults) {
    $report += $reportFormat -f $res.Name, $res.Winner, $res.Difference
}

$report += ""
$report += ""
$report += "========================================================="
$report += "RAW RESULTS"
$report += "========================================================="

foreach ($res in $AllResults) {
    $report += ""
    $report += "[$($res.Category.ToUpper()) / $($res.Name)]"
    $report += ""
    
    # TitanHTTP
    $report += "TitanHTTP:"
    $report += "----------"
    $report += "RPS:"
    $report += "Samples:    [$($res.TitanHTTP.RpsSamples -join ', ')]"
    $report += "Mean:       $($res.TitanHTTP.RpsStats.Mean)"
    $report += "Median:     $($res.TitanHTTP.RpsStats.Median)"
    $report += "StdDev:     $($res.TitanHTTP.RpsStats.StdDev)"
    $report += "Range:      $($res.TitanHTTP.RpsStats.Range)"
    $report += ""
    $report += "Throughput (MB/s):"
    $report += "Samples:    [$($res.TitanHTTP.MbSamples -join ', ')]"
    $report += "Median:     $($res.TitanHTTP.MbStats.Median)"
    $report += ""
    $report += "P50 Latency (ms):"
    $report += "Samples:    [$($res.TitanHTTP.P50Samples -join ', ')]"
    $report += "Median:     $($res.TitanHTTP.P50Stats.Median)ms"
    $report += ""
    $report += "P95 Latency (ms):"
    $report += "Samples:    [$($res.TitanHTTP.P95Samples -join ', ')]"
    $report += "Median:     $($res.TitanHTTP.P95Stats.Median)ms"
    $report += ""
    $report += "P99 Latency (ms):"
    $report += "Samples:    [$($res.TitanHTTP.P99Samples -join ', ')]"
    $report += "Median:     $($res.TitanHTTP.P99Stats.Median)ms"
    $report += ""
    
    # net/http
    $report += "net/http:"
    $report += "---------"
    $report += "RPS:"
    $report += "Samples:    [$($res.NetHttp.RpsSamples -join ', ')]"
    $report += "Mean:       $($res.NetHttp.RpsStats.Mean)"
    $report += "Median:     $($res.NetHttp.RpsStats.Median)"
    $report += "StdDev:     $($res.NetHttp.RpsStats.StdDev)"
    $report += "Range:      $($res.NetHttp.RpsStats.Range)"
    $report += ""
    $report += "Throughput (MB/s):"
    $report += "Samples:    [$($res.NetHttp.MbSamples -join ', ')]"
    $report += "Median:     $($res.NetHttp.MbStats.Median)"
    $report += ""
    $report += "P50 Latency (ms):"
    $report += "Samples:    [$($res.NetHttp.P50Samples -join ', ')]"
    $report += "Median:     $($res.NetHttp.P50Stats.Median)ms"
    $report += ""
    $report += "P95 Latency (ms):"
    $report += "Samples:    [$($res.NetHttp.P95Samples -join ', ')]"
    $report += "Median:     $($res.NetHttp.P95Stats.Median)ms"
    $report += ""
    $report += "P99 Latency (ms):"
    $report += "Samples:    [$($res.NetHttp.P99Samples -join ', ')]"
    $report += "Median:     $($res.NetHttp.P99Stats.Median)ms"
    $report += ""
    $report += "---------------------------------------------------------"
}

$report += ""
$report += "========================================================="
$report += "MEMORY & CPU PROFILING (500k requests)"
$report += "========================================================="
$report += ""
$report += "TitanHTTP Heap Leak: $($MemoryResults.TitanHTTP) MB"
$report += "net/http Heap Leak:  $($MemoryResults.NetHttp) MB"
$report += ""
$report += "CPU Profiles captured to benchmarks/profiles/cpu_*.pprof"
$report += ""

$reportPath = "benchmarks/reports/unified_comparison_report.txt"
$report | Out-File -FilePath $reportPath -Encoding utf8

$mdPath = "benchmarks/reports/unified_comparison.md"
if (Test-Path $mdPath) { Remove-Item $mdPath -Force }
$jsonPath = "benchmarks/data/unified_baseline.json"
if (Test-Path $jsonPath) { Remove-Item $jsonPath -Force }

Write-Host "`n  ============================================================" -ForegroundColor Cyan
Write-Host "  ::                                                        ::" -ForegroundColor Cyan
Write-Host "  ::            UNIFIED LABORATORY SUITE COMPLETE           ::" -ForegroundColor Yellow
Write-Host "  ::         Report generated: $reportPath" -ForegroundColor White
Write-Host "  ::                                                        ::" -ForegroundColor Cyan
Write-Host "  ============================================================`n" -ForegroundColor Cyan
