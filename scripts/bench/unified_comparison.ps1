$ErrorActionPreference = "Stop"

Write-Host "`n  ============================================================" -ForegroundColor Cyan
Write-Host "  ::                                                        ::" -ForegroundColor Cyan
Write-Host "  ::       TITANHTTP VS NET/HTTP UNIFIED CINEMATIC SUITE    ::" -ForegroundColor Yellow
Write-Host "  ::                                                        ::" -ForegroundColor Cyan
Write-Host "  ============================================================`n" -ForegroundColor Cyan

# Ensure outputs directory
if (-not (Test-Path "benchmarks/data")) { New-Item -ItemType Directory -Path "benchmarks/data" -Force | Out-Null }
if (-not (Test-Path "benchmarks/logs")) { New-Item -ItemType Directory -Path "benchmarks/logs" -Force | Out-Null }
if (-not (Test-Path "benchmarks/reports")) { New-Item -ItemType Directory -Path "benchmarks/reports" -Force | Out-Null }
if (-not (Test-Path "benchmarks/profiles")) { New-Item -ItemType Directory -Path "benchmarks/profiles" -Force | Out-Null }

$AllResults = @{}

# Build servers
Write-Host "`n[+] Building servers..."
go build -o bin/titanhttp_bench.exe ./cmd/titanhttp
go build -o bin/nethttp_bench.exe ./cmd/nethttp_bench

# Helper to run bombardier and parse stats
function Run-Bombardier {
    param ($Url, $ArgsStr)
    $cmd = "bombardier -p r -o json -l $ArgsStr $Url"
    $jsonRaw = Invoke-Expression "$cmd 2> `$null"
    $res = $jsonRaw | ConvertFrom-Json
    
    $p50 = [math]::Round($res.result.latency.percentiles."50" / 1000, 2)
    $p99 = [math]::Round($res.result.latency.percentiles."99" / 1000, 2)
    
    $errors = 0
    if ($res.result.errors) {
        $errors = ($res.result.req1xx + $res.result.req2xx + $res.result.req3xx + $res.result.req4xx + $res.result.req5xx) 
        $errors = $res.result.reqsTotal - $errors
    }
    
    return @{
        Rps = [math]::Round($res.result.rps.mean, 2)
        P50 = $p50
        P99 = $p99
        Errors = $errors
        ThroughputMB = [math]::Round($res.result.throughput / 1MB, 2)
    }
}

# ---------------------------------------------------------
# STAGE 1: Standard Tests
# ---------------------------------------------------------
Write-Host "`n>>> [ STAGE 1 ] Starting Standard HTTP Servers..." -ForegroundColor Green
$TitanProcess = Start-Process -FilePath ".\bin\titanhttp_bench.exe" -PassThru -NoNewWindow -RedirectStandardOutput "benchmarks/logs/titan_bench.log" -RedirectStandardError "benchmarks/logs/titan_bench_err.log"
$NetHttpProcess = Start-Process -FilePath ".\bin\nethttp_bench.exe" -PassThru -NoNewWindow -RedirectStandardOutput "benchmarks/logs/nethttp_bench.log" -RedirectStandardError "benchmarks/logs/nethttp_bench_err.log"
Start-Sleep -Seconds 4

$StandardTests = @(
    @{ Name = "Throughput (Req/s)"; Path = "/ping"; Args = "-c 125 -d 5s" },
    @{ Name = "Concurrent Connections"; Path = "/ping"; Args = "-c 1000 -d 5s" },
    @{ Name = "Large Payloads"; Path = "/heavy"; Args = "-c 125 -d 5s" },
    @{ Name = "Routing Performance"; Path = "/users/123"; Args = "-c 125 -d 5s" },
    @{ Name = "Static File Serving"; Path = "/static/test.txt"; Args = "-c 125 -d 5s" },
    @{ Name = "Keep-Alive Performance"; Path = "/ping"; Args = "-c 125 -d 5s -a" }
)

$Stage1Results = @{}
foreach ($test in $StandardTests) {
    Write-Host "`n--- Testing $($test.Name) ---" -ForegroundColor Yellow
    Write-Host "  -> Running against TitanHTTP..."
    $tRes = Run-Bombardier -Url "http://localhost:8080$($test.Path)" -ArgsStr $test.Args
    Write-Host "  -> Running against net/http..."
    $nRes = Run-Bombardier -Url "http://localhost:8081$($test.Path)" -ArgsStr $test.Args
    
    $Stage1Results[$test.Name] = @{ TitanHTTP = $tRes; NetHttp = $nRes }
    
    $winner = if ($tRes.Rps -gt $nRes.Rps) { "TitanHTTP" } else { "net/http" }
    Write-Host ("{0,-20} | {1,-18} | {2,-18}" -f "Metric", "TitanHTTP", "net/http") -ForegroundColor Cyan
    Write-Host ("-" * 64) -ForegroundColor DarkGray
    Write-Host ("{0,-20} | {1,-18} | {2,-18}" -f "RPS", $tRes.Rps, $nRes.Rps) -ForegroundColor White
    Write-Host ("{0,-20} | {1,-18} | {2,-18}" -f "P50 Latency (ms)", $tRes.P50, $nRes.P50) -ForegroundColor Gray
    Write-Host ("{0,-20} | {1,-18} | {2,-18}" -f "P99 Latency (ms)", $tRes.P99, $nRes.P99) -ForegroundColor Gray
    Write-Host ("{0,-20} | {1,-18} | {2,-18}" -f "Throughput (MB/s)", $tRes.ThroughputMB, $nRes.ThroughputMB) -ForegroundColor Gray
    Write-Host ("{0,-20} | {1,-18} | {2,-18}" -f "Errors", $tRes.Errors, $nRes.Errors) -ForegroundColor Gray
    Write-Host "[ WINNER ] $winner`n" -ForegroundColor Green
}
$AllResults["Standard"] = $Stage1Results


# ---------------------------------------------------------
# STAGE 2: Concurrency Scale (C10K)
# ---------------------------------------------------------
Write-Host "`n>>> [ STAGE 2 ] Concurrency Scale (C10K Simulator)..." -ForegroundColor Green
$ConcurrencyLevels = @(50, 100, 200)
$Stage2Results = @()

foreach ($c in $ConcurrencyLevels) {
    Write-Host "`n[ IGNITION ] Ramping to $c concurrent connections..." -ForegroundColor Magenta
    $tRes = Run-Bombardier -Url "http://localhost:8080/ping" -ArgsStr "-c $c -d 5s"
    $nRes = Run-Bombardier -Url "http://localhost:8081/ping" -ArgsStr "-c $c -d 5s"
    
    $Stage2Results += @{ Connections = $c; TitanHTTP = $tRes; NetHttp = $nRes }
    Write-Host "TitanHTTP RPS: $($tRes.Rps) | P50: $($tRes.P50)ms | P99: $($tRes.P99)ms | Err: $($tRes.Errors)"
    Write-Host "net/http  RPS: $($nRes.Rps) | P50: $($nRes.P50)ms | P99: $($nRes.P99)ms | Err: $($nRes.Errors)`n"
}
$AllResults["Concurrency"] = $Stage2Results

# ---------------------------------------------------------
# STAGE 3: Connection Churn
# ---------------------------------------------------------
Write-Host "`n>>> [ STAGE 3 ] Connection Churn (Connection: close)..." -ForegroundColor Green
Write-Host "  -> TitanHTTP..."
$tChurn = Run-Bombardier -Url "http://localhost:8080/ping" -ArgsStr "-c 100 -d 5s -H `"Connection: close`""
Write-Host "  -> net/http..."
$nChurn = Run-Bombardier -Url "http://localhost:8081/ping" -ArgsStr "-c 100 -d 5s -H `"Connection: close`""
$AllResults["ConnectionChurn"] = @{ TitanHTTP = $tChurn; NetHttp = $nChurn }

Write-Host "TitanHTTP RPS: $($tChurn.Rps) | P50: $($tChurn.P50)ms | P99: $($tChurn.P99)ms | Err: $($tChurn.Errors)"
Write-Host "net/http  RPS: $($nChurn.Rps) | P50: $($nChurn.P50)ms | P99: $($nChurn.P99)ms | Err: $($nChurn.Errors)`n"


# ---------------------------------------------------------
# STAGE 4: Memory Leak Detector
# ---------------------------------------------------------
Write-Host "`n>>> [ STAGE 4 ] Memory Leak Detector (pprof)..." -ForegroundColor Green
Write-Host "1. Capturing Baseline Heaps..."
Invoke-WebRequest -Uri "http://localhost:6060/debug/pprof/heap" -OutFile "benchmarks/profiles/titan_baseline.pprof"
Invoke-WebRequest -Uri "http://localhost:6061/debug/pprof/heap" -OutFile "benchmarks/profiles/nethttp_baseline.pprof"

Write-Host "2. Blasting both servers with 500,000 requests each..."
$j1 = Start-Process -FilePath "bombardier" -ArgumentList "-c 125 -n 500000 http://localhost:8080/ping" -NoNewWindow -PassThru -RedirectStandardOutput "benchmarks/logs/bombardier_titan.log" -RedirectStandardError "benchmarks/logs/bombardier_titan_err.log"
$j2 = Start-Process -FilePath "bombardier" -ArgumentList "-c 125 -n 500000 http://localhost:8081/ping" -NoNewWindow -PassThru -RedirectStandardOutput "benchmarks/logs/bombardier_nethttp.log" -RedirectStandardError "benchmarks/logs/bombardier_nethttp_err.log"
$j1.WaitForExit()
$j2.WaitForExit()

Write-Host "3. Capturing End Heaps..."
Invoke-WebRequest -Uri "http://localhost:6060/debug/pprof/heap" -OutFile "benchmarks/profiles/titan_end.pprof"
Invoke-WebRequest -Uri "http://localhost:6061/debug/pprof/heap" -OutFile "benchmarks/profiles/nethttp_end.pprof"

$titanDiff = go tool pprof -text -base "benchmarks/profiles/titan_baseline.pprof" "benchmarks/profiles/titan_end.pprof"
$nethttpDiff = go tool pprof -text -base "benchmarks/profiles/nethttp_baseline.pprof" "benchmarks/profiles/nethttp_end.pprof"

function Parse-Leak { param($Output)
    $line = $Output | Where-Object { $_ -match "^Total:\s+(-?[\d\.]+MB|-?[\d\.]+kB|-?[\d\.]+B)" }
    if ($line -match "Total:\s+([\d\.]+)MB") { return [double]$matches[1] }
    return 0
}

$tLeak = Parse-Leak $titanDiff
$nLeak = Parse-Leak $nethttpDiff
Write-Host "TitanHTTP Leak: $tLeak MB | net/http Leak: $nLeak MB"
$AllResults["MemoryLeakMB"] = @{ TitanHTTP = $tLeak; NetHttp = $nLeak }


# Shutdown standard servers
Stop-Process -Id $TitanProcess.Id -Force -ErrorAction SilentlyContinue
Stop-Process -Id $NetHttpProcess.Id -Force -ErrorAction SilentlyContinue


# ---------------------------------------------------------
# STAGE 5: TLS Benchmark
# ---------------------------------------------------------
Write-Host "`n>>> [ STAGE 5 ] TLS Performance..." -ForegroundColor Green
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

Write-Host "  -> TitanHTTP (HTTPS)..."
$tTLS = Run-Bombardier -Url "https://localhost:8443/ping" -ArgsStr "-c 125 -d 5s --insecure"
Write-Host "  -> net/http (HTTPS)..."
$nTLS = Run-Bombardier -Url "https://localhost:8444/ping" -ArgsStr "-c 125 -d 5s --insecure"

$AllResults["TLS"] = @{ TitanHTTP = $tTLS; NetHttp = $nTLS }
Write-Host "TitanHTTP TLS RPS: $($tTLS.Rps) | P50: $($tTLS.P50)ms | P99: $($tTLS.P99)ms"
Write-Host "net/http  TLS RPS: $($nTLS.Rps) | P50: $($nTLS.P50)ms | P99: $($nTLS.P99)ms"

Stop-Process -Id $TitanTLS.Id -Force -ErrorAction SilentlyContinue
Stop-Process -Id $NetHttpTLS.Id -Force -ErrorAction SilentlyContinue
Remove-Item env:TLS_MODE -ErrorAction SilentlyContinue
Remove-Item env:TLS_CERT -ErrorAction SilentlyContinue
Remove-Item env:TLS_KEY -ErrorAction SilentlyContinue


# ---------------------------------------------------------
# Final Export
# ---------------------------------------------------------
Write-Host "`n[+] Exporting Unified Results to TXT format..." -ForegroundColor Green
$exportData = @()

foreach ($testName in $StandardTests.Name) {
    $t = $AllResults["Standard"][$testName].TitanHTTP
    $n = $AllResults["Standard"][$testName].NetHttp
    $exportData += [PSCustomObject]@{ Stage="1-Standard"; Test=$testName; TRps=$t.Rps; TP50=$t.P50; TP99=$t.P99; NRps=$n.Rps; NP50=$n.P50; NP99=$n.P99 }
}

foreach ($item in $AllResults["Concurrency"]) {
    $c = $item.Connections
    $t = $item.TitanHTTP
    $n = $item.NetHttp
    $exportData += [PSCustomObject]@{ Stage="2-Concurrency"; Test="$c Conns"; TRps=$t.Rps; TP50=$t.P50; TP99=$t.P99; NRps=$n.Rps; NP50=$n.P50; NP99=$n.P99 }
}

$tChurn = $AllResults["ConnectionChurn"].TitanHTTP
$nChurn = $AllResults["ConnectionChurn"].NetHttp
$exportData += [PSCustomObject]@{ Stage="3-Churn"; Test="Conn: close"; TRps=$tChurn.Rps; TP50=$tChurn.P50; TP99=$tChurn.P99; NRps=$nChurn.Rps; NP50=$nChurn.P50; NP99=$nChurn.P99 }

$tTLS = $AllResults["TLS"].TitanHTTP
$nTLS = $AllResults["TLS"].NetHttp
$exportData += [PSCustomObject]@{ Stage="5-TLS"; Test="HTTPS Ping"; TRps=$tTLS.Rps; TP50=$tTLS.P50; TP99=$tTLS.P99; NRps=$nTLS.Rps; NP50=$nTLS.P50; NP99=$nTLS.P99 }

$exportData | Format-Table -AutoSize | Out-File -FilePath "benchmarks/data/unified_comparison.txt" -Encoding utf8

Write-Host "[+] Generating Unified Markdown Report..." -ForegroundColor Green
$reportPath = "benchmarks/reports/unified_comparison.md"
$mdLines = @(
    "# TitanHTTP vs net/http Unified Comparison Report"
    "Generated on $(Get-Date)"
    ""
    "## 1. Concurrency Scaling Performance"
    "This benchmark measures how throughput (RPS) and latency scale as concurrent connections increase."
    ""
    "| Concurrency (Connections) | Server | Throughput (Req/s) | P50 Latency (ms) | P99 Latency (ms) | TitanHTTP Speedup | P99 Latency Reduction |"
    "| :--- | :--- | :--- | :--- | :--- | :--- | :--- |"
)
$md = $mdLines -join "`r`n"

foreach ($item in $AllResults["Concurrency"]) {
    $c = $item.Connections
    $t = $item.TitanHTTP
    $n = $item.NetHttp
    
    $speedup = 0
    if ($n.Rps -gt 0) {
        $speedup = [math]::Round((($t.Rps - $n.Rps) / $n.Rps) * 100, 1)
    }
    $p99Red = 0
    if ($n.P99 -gt 0) {
        $p99Red = [math]::Round((($n.P99 - $t.P99) / $n.P99) * 100, 1)
    }
    
    $speedupText = if ($speedup -ge 0) { "+$($speedup)%" } else { "$($speedup)%" }
    $p99RedText = if ($p99Red -ge 0) { "-$($p99Red)%" } else { "+$([math]::Abs($p99Red))%" }
    
    $md += "`n| **$c Connections** | **TitanHTTP** | **$($t.Rps)** | **$($t.P50)** | **$($t.P99)** | **$speedupText** | **$p99RedText** |"
    $md += "`n| | net/http | $($n.Rps) | $($n.P50) | $($n.P99) | | |"
}

$mdLines2 = @(
    ""
    ""
    "## 2. Standard Workload Profiles"
    "A side-by-side comparison of specific networking patterns under standard load."
    ""
    "| Workload | Server | Throughput (Req/s) | P50 Latency (ms) | P99 Latency (ms) | RPS Difference | P99 Difference |"
    "| :--- | :--- | :--- | :--- | :--- | :--- | :--- |"
)
$md += "`r`n" + ($mdLines2 -join "`r`n")

foreach ($testName in $AllResults["Standard"].Keys) {
    $t = $AllResults["Standard"][$testName].TitanHTTP
    $n = $AllResults["Standard"][$testName].NetHttp
    
    $speedup = 0
    if ($n.Rps -gt 0) {
        $speedup = [math]::Round((($t.Rps - $n.Rps) / $n.Rps) * 100, 1)
    }
    $p99Red = 0
    if ($n.P99 -gt 0) {
        $p99Red = [math]::Round((($n.P99 - $t.P99) / $n.P99) * 100, 1)
    }
    
    $speedupText = if ($speedup -ge 0) { "+$($speedup)%" } else { "$($speedup)%" }
    $p99RedText = if ($p99Red -ge 0) { "-$($p99Red)%" } else { "+$([math]::Abs($p99Red))%" }
    
    $md += "`n| **$testName** | **TitanHTTP** | **$($t.Rps)** | **$($t.P50)** | **$($t.P99)** | **$speedupText** | **$p99RedText** |"
    $md += "`n| | net/http | $($n.Rps) | $($n.P50) | $($n.P99) | | |"
}

$tChurn = $AllResults["ConnectionChurn"].TitanHTTP
$nChurn = $AllResults["ConnectionChurn"].NetHttp
$tTLS = $AllResults["TLS"].TitanHTTP
$nTLS = $AllResults["TLS"].NetHttp

$churnSpeedup = 0
if ($nChurn.Rps -gt 0) { $churnSpeedup = [math]::Round((($tChurn.Rps - $nChurn.Rps) / $nChurn.Rps) * 100, 1) }
$churnP99 = 0
if ($nChurn.P99 -gt 0) { $churnP99 = [math]::Round((($nChurn.P99 - $tChurn.P99) / $nChurn.P99) * 100, 1) }

$tlsSpeedup = 0
if ($nTLS.Rps -gt 0) { $tlsSpeedup = [math]::Round((($tTLS.Rps - $nTLS.Rps) / $nTLS.Rps) * 100, 1) }
$tlsP99 = 0
if ($nTLS.P99 -gt 0) { $tlsP99 = [math]::Round((($nTLS.P99 - $tTLS.P99) / $nTLS.P99) * 100, 1) }

$churnSpeedupText = if ($churnSpeedup -ge 0) { "+$($churnSpeedup)%" } else { "$($churnSpeedup)%" }
$churnP99Text = if ($churnP99 -ge 0) { "-$($churnP99)%" } else { "+$([math]::Abs($churnP99))%" }
$tlsSpeedupText = if ($tlsSpeedup -ge 0) { "+$($tlsSpeedup)%" } else { "$($tlsSpeedup)%" }
$tlsP99Text = if ($tlsP99 -ge 0) { "-$($tlsP99)%" } else { "+$([math]::Abs($tlsP99))%" }

$mdLines3 = @(
    ""
    ""
    "## 3. High-Churn and TLS Infrastructure"
    "Evaluating TCP handshake cycling and TLS session overhead."
    ""
    "| Profile | Server | Throughput (Req/s) | P50 Latency (ms) | P99 Latency (ms) | RPS Difference | P99 Difference |"
    "| :--- | :--- | :--- | :--- | :--- | :--- | :--- |"
    "| **Connection Churn** | **TitanHTTP** | **$($tChurn.Rps)** | **$($tChurn.P50)** | **$($tChurn.P99)** | **$churnSpeedupText** | **$churnP99Text** |"
    "| (Connection: close) | net/http | $($nChurn.Rps) | $($nChurn.P50) | $($nChurn.P99) | | |"
    "| **TLS Overhead** | **TitanHTTP** | **$($tTLS.Rps)** | **$($tTLS.P50)** | **$($tTLS.P99)** | **$tlsSpeedupText** | **$tlsP99Text** |"
    "| | net/http | $($nTLS.Rps) | $($nTLS.P50) | $($nTLS.P99) | | |"
    ""
    "## 4. Memory Profiling"
    "Analyzing memory leaks during a 500,000 request stress run."
    ""
    "- **TitanHTTP Leak:** $($AllResults['MemoryLeakMB'].TitanHTTP) MB"
    "- **net/http Leak:** $($AllResults['MemoryLeakMB'].NetHttp) MB"
)
$md += "`r`n" + ($mdLines3 -join "`r`n")

$md | Out-File -FilePath $reportPath -Encoding utf8

Write-Host "`n[+] Exporting JSON Baseline Stats for Showcase Platform..." -ForegroundColor Green
$AllResults | ConvertTo-Json -Depth 5 | Out-File -FilePath "benchmarks/data/unified_baseline.json" -Encoding utf8

Write-Host "`n  ============================================================" -ForegroundColor Cyan
Write-Host "  ::                                                        ::" -ForegroundColor Cyan
Write-Host "  ::            UNIFIED CINEMATIC SUITE COMPLETE            ::" -ForegroundColor Yellow
Write-Host "  ::                                                        ::" -ForegroundColor Cyan
Write-Host "  ============================================================`n" -ForegroundColor Cyan
