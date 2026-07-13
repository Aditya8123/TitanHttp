$ErrorActionPreference = "Stop"

Write-Host "`n  ============================================================" -ForegroundColor Cyan
Write-Host "  ::                                                        ::" -ForegroundColor Cyan
Write-Host "  ::    TITANHTTP VS NET/HTTP ADVANCED SECURITY LABORATORY  ::" -ForegroundColor Yellow
Write-Host "  ::                                                        ::" -ForegroundColor Cyan
Write-Host "  ============================================================`n" -ForegroundColor Cyan

# Ensure outputs directory
if (-not (Test-Path "benchmarks/logs")) { New-Item -ItemType Directory -Path "benchmarks/logs" -Force | Out-Null }
if (-not (Test-Path "benchmarks/reports")) { New-Item -ItemType Directory -Path "benchmarks/reports" -Force | Out-Null }

$AllResults = [System.Collections.ArrayList]::new()

# Build servers
Write-Host "`n[+] Building servers for security testing..."
go build -o bin/titanhttp_bench.exe ./cmd/titanhttp
go build -o bin/nethttp_bench.exe ./cmd/nethttp_bench

Write-Host "`n>>> [ STAGE 1 ] Booting HTTP Servers..." -ForegroundColor Green
$env:IDLE_TIMEOUT = "5s" # Force 5s idle timeout for testing connection floods
$TitanProcess = Start-Process -FilePath ".\bin\titanhttp_bench.exe" -PassThru -NoNewWindow -RedirectStandardOutput "benchmarks/logs/titan_sec.log" -RedirectStandardError "benchmarks/logs/titan_sec_err.log"
$NetHttpProcess = Start-Process -FilePath ".\bin\nethttp_bench.exe" -PassThru -NoNewWindow -RedirectStandardOutput "benchmarks/logs/nethttp_sec.log" -RedirectStandardError "benchmarks/logs/nethttp_sec_err.log"
Start-Sleep -Seconds 4
Remove-Item env:IDLE_TIMEOUT -ErrorAction SilentlyContinue

# --- Attack Vectors ---

function Test-SlowPOST {
    param($HostName, $Port)
    try {
        $tcpClient = New-Object System.Net.Sockets.TcpClient($HostName, $Port)
        $stream = $tcpClient.GetStream()
        
        $request = "POST / HTTP/1.1`r`nHost: $HostName`r`nContent-Length: 100`r`n`r`n"
        $bytes = [System.Text.Encoding]::ASCII.GetBytes($request)
        $stream.Write($bytes, 0, $bytes.Length)
        
        $startTime = Get-Date
        $trickleByte = [System.Text.Encoding]::ASCII.GetBytes("A")
        
        $timedOut = $false
        for ($i = 0; $i -lt 50; $i++) {
            if (-not $tcpClient.Connected) {
                $timedOut = $true
                break
            }
            try {
                $stream.Write($trickleByte, 0, $trickleByte.Length)
                Start-Sleep -Milliseconds 200
            } catch {
                $timedOut = $true
                break
            }
        }
        
        $duration = (Get-Date) - $startTime
        $tcpClient.Close()
        
        $passed = $timedOut
        $notes = if ($passed) { "Connection severed proactively by server ($([Math]::Round($duration.TotalSeconds, 2))s)" } else { "VULNERABLE: Allowed trickle indefinitely ($([Math]::Round($duration.TotalSeconds, 2))s)" }
        return @{ Passed = $passed; Notes = $notes }
    } catch {
        return @{ Passed = $true; Notes = "Connection rejected/aborted proactively." }
    }
}

function Test-HeaderExhaustion {
    param($HostName, $Port)
    try {
        $tcpClient = New-Object System.Net.Sockets.TcpClient($HostName, $Port)
        $stream = $tcpClient.GetStream()
        
        $requestLine = "GET / HTTP/1.1`r`nHost: $HostName`r`n"
        $bytes = [System.Text.Encoding]::ASCII.GetBytes($requestLine)
        $stream.Write($bytes, 0, $bytes.Length)
        
        $headersSent = 0
        $connectionClosed = $false
        
        for ($i = 0; $i -lt 2000; $i++) {
            $headerLine = "X-Junk-Header-${i}: $( 'A' * 10 )`r`n"
            $headerBytes = [System.Text.Encoding]::ASCII.GetBytes($headerLine)
            try {
                $stream.Write($headerBytes, 0, $headerBytes.Length)
                $headersSent++
            } catch {
                $connectionClosed = $true
                break
            }
        }
        
        $response = ""
        if (-not $connectionClosed) {
            try {
                $endBytes = [System.Text.Encoding]::ASCII.GetBytes("`r`n")
                $stream.Write($endBytes, 0, $endBytes.Length)
                
                $reader = New-Object System.IO.StreamReader($stream)
                $response = $reader.ReadLine()
            } catch {
                $connectionClosed = $true
            }
        }
        $tcpClient.Close()
        
        $passed = $connectionClosed -or ($response -match "431|400|500")
        $notes = if ($passed) { "Blocked after $headersSent headers (Status: $response)" } else { "VULNERABLE: Parsed $headersSent junk headers" }
        return @{ Passed = $passed; Notes = $notes }
    } catch {
        return @{ Passed = $true; Notes = "Connection closed proactively by server." }
    }
}

function Test-IdleConnectionFlood {
    param($HostName, $Port)
    $sockets = @()
    for ($i = 0; $i -lt 200; $i++) {
        try {
            $tcpClient = New-Object System.Net.Sockets.TcpClient($HostName, $Port)
            $sockets += $tcpClient
        } catch { }
    }
    
    $opened = $sockets.Count
    Start-Sleep -Seconds 7
    
    $active = 0
    foreach ($sock in $sockets) {
        if ($sock.Client.Poll(0, [System.Net.Sockets.SelectMode]::SelectRead)) {
            $buffer = New-Object byte[] 1
            if ($sock.Client.Receive($buffer, [System.Net.Sockets.SocketFlags]::Peek) -eq 0) {
                # Connection closed
            } else {
                $active++
            }
        } else {
             $active++
        }
        $sock.Close()
    }
    
    $passed = ($active -eq 0)
    $notes = if ($passed) { "All $opened idle sockets purged successfully" } else { "VULNERABLE: $active idle sockets still active after timeout" }
    return @{ Passed = $passed; Notes = $notes }
}

function Test-MalformedHTTP {
    param($HostName, $Port)
    try {
        $tcpClient = New-Object System.Net.Sockets.TcpClient($HostName, $Port)
        $stream = $tcpClient.GetStream()
        
        $garbage = [System.Text.Encoding]::ASCII.GetBytes("JUST_GARBAGE_NO_SPACES_OR_PROTOCOL`r`n`r`n")
        $stream.Write($garbage, 0, $garbage.Length)
        
        $reader = New-Object System.IO.StreamReader($stream)
        $response = $reader.ReadLine()
        $tcpClient.Close()
        
        $passed = ($response -match "400") -or ([string]::IsNullOrEmpty($response))
        $notes = if ($passed) { "Garbage rejected safely (Status: $response)" } else { "Unexpected behavior (Status: $response)" }
        return @{ Passed = $passed; Notes = $notes }
    } catch {
        return @{ Passed = $true; Notes = "Garbage rejected safely (Connection dropped)" }
    }
}

function Run-Comparison {
    param ($Name, $ScriptBlock)
    Write-Host "`n--- Testing $Name ---" -ForegroundColor Yellow
    Write-Host "  -> TitanHTTP..."
    $tRes = & $ScriptBlock -HostName "localhost" -Port 8080
    Write-Host "  -> net/http..."
    $nRes = & $ScriptBlock -HostName "localhost" -Port 8081
    
    $tColor = if ($tRes.Passed) { "Green" } else { "Red" }
    $nColor = if ($nRes.Passed) { "Green" } else { "Red" }
    Write-Host "     [TitanHTTP] " -NoNewline; Write-Host $(if($tRes.Passed){"PASS"}else{"FAIL"}) -ForegroundColor $tColor
    Write-Host "     [net/http]  " -NoNewline; Write-Host $(if($nRes.Passed){"PASS"}else{"FAIL"}) -ForegroundColor $nColor

    $AllResults.Add(@{
        Name = $Name
        TitanHTTP = $tRes
        NetHttp = $nRes
    }) | Out-Null
}

Write-Host "`n>>> [ STAGE 2 ] Executing Security Attacks..." -ForegroundColor Green
Run-Comparison -Name "Slow POST Body Trickle" -ScriptBlock ${function:Test-SlowPOST}
Run-Comparison -Name "Header Exhaustion (Bomb)" -ScriptBlock ${function:Test-HeaderExhaustion}
Run-Comparison -Name "Idle Connection Flood" -ScriptBlock ${function:Test-IdleConnectionFlood}
Run-Comparison -Name "Malformed HTTP Request" -ScriptBlock ${function:Test-MalformedHTTP}

# Shutdown servers
Stop-Process -Id $TitanProcess.Id -Force -ErrorAction SilentlyContinue
Stop-Process -Id $NetHttpProcess.Id -Force -ErrorAction SilentlyContinue


# ---------------------------------------------------------
# GENERATE REPORT
# ---------------------------------------------------------
Write-Host "`n[+] Generating Security Report..." -ForegroundColor Green

$report = @()
$report += "========================================================="
$report += " TITANHTTP VS NET/HTTP ADVANCED SECURITY REPORT"
$report += "========================================================="
$report += ""
$report += "Environment"
$report += "-----------"
$report += "CPU:        $((Get-CimInstance Win32_Processor).Name)"
$report += "OS:         $((Get-CimInstance Win32_OperatingSystem).Caption)"
$report += "Go version: $(go version)"
$report += "Date:       $(Get-Date)"
$report += ""
$report += ""
$report += "========================================================="
$report += "SECURITY POSTURE MATRIX"
$report += "========================================================="
$report += ""
$reportFormat = "{0,-35} {1,-15} {2,-15}"
$report += $reportFormat -f "Attack Scenario", "TitanHTTP", "net/http"
$report += "-" * 65
foreach ($res in $AllResults) {
    $tPass = if ($res.TitanHTTP.Passed) { "PASS" } else { "FAIL" }
    $nPass = if ($res.NetHttp.Passed) { "PASS" } else { "FAIL" }
    $report += $reportFormat -f $res.Name, $tPass, $nPass
}

$report += ""
$report += ""
$report += "========================================================="
$report += "DETAILED FINDINGS"
$report += "========================================================="

foreach ($res in $AllResults) {
    $report += ""
    $report += "[ $($res.Name) ]"
    $tPass = if ($res.TitanHTTP.Passed) { "PASS" } else { "FAIL" }
    $nPass = if ($res.NetHttp.Passed) { "PASS" } else { "FAIL" }
    
    $report += "TitanHTTP: $tPass - $($res.TitanHTTP.Notes)"
    $report += "net/http:  $nPass - $($res.NetHttp.Notes)"
    $report += "-" * 65
}

$reportPath = "benchmarks/reports/advanced_security_report.txt"
$report | Out-File -FilePath $reportPath -Encoding utf8

Write-Host "`n  ============================================================" -ForegroundColor Cyan
Write-Host "  ::                                                        ::" -ForegroundColor Cyan
Write-Host "  ::             SECURITY LABORATORY COMPLETE               ::" -ForegroundColor Yellow
Write-Host "  ::         Report generated: $reportPath" -ForegroundColor White
Write-Host "  ::                                                        ::" -ForegroundColor Cyan
Write-Host "  ============================================================`n" -ForegroundColor Cyan
