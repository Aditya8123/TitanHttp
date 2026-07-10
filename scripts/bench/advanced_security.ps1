param (
    [string]$TargetHost = "localhost",
    [int]$TargetPort = 8080
)

$ErrorActionPreference = "Stop"

Write-Host "`n  ============================================================" -ForegroundColor Cyan
Write-Host "  ::                                                        ::" -ForegroundColor Cyan
Write-Host "  ::           TITANHTTP ADVANCED SECURITY SCENARIOS        ::" -ForegroundColor Yellow
Write-Host "  ::                                                        ::" -ForegroundColor Cyan
Write-Host "  ============================================================`n" -ForegroundColor Cyan
Write-Host "Target: ${TargetHost}:${TargetPort}"
Write-Host "----------------------------------------------------------" -ForegroundColor Gray

function Test-SlowPOST {
    Write-Host "`n[Scenario 1] Slow POST Body Trickle" -ForegroundColor Yellow
    Write-Host "Testing if server correctly times out on a slowly sent body..."
    
    try {
        $tcpClient = New-Object System.Net.Sockets.TcpClient($TargetHost, $TargetPort)
        $stream = $tcpClient.GetStream()
        
        $request = "POST / HTTP/1.1`r`nHost: $TargetHost`r`nContent-Length: 100`r`n`r`n"
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
        if ($timedOut) {
            Write-Host "SUCCESS: Server closed connection during Slow POST trickle (Duration: $([Math]::Round($duration.TotalSeconds, 2))s)" -ForegroundColor Green
        } else {
            Write-Host "FAIL: Server allowed Slow POST to continue indefinitely (Duration: $([Math]::Round($duration.TotalSeconds, 2))s)" -ForegroundColor Red
        }
        
        $tcpClient.Close()
    } catch {
        Write-Host "Error during Slow POST test: $_" -ForegroundColor Red
    }
}

function Test-HeaderExhaustion {
    Write-Host "`n[Scenario 2] Header Exhaustion/Bomb" -ForegroundColor Yellow
    Write-Host "Testing if server limits massive number of headers..."
    
    try {
        $tcpClient = New-Object System.Net.Sockets.TcpClient($TargetHost, $TargetPort)
        $stream = $tcpClient.GetStream()
        
        $requestLine = "GET / HTTP/1.1`r`nHost: $TargetHost`r`n"
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
        
        if (-not $connectionClosed) {
            # Try to complete the request
            try {
                $endBytes = [System.Text.Encoding]::ASCII.GetBytes("`r`n")
                $stream.Write($endBytes, 0, $endBytes.Length)
                
                # Read response
                $reader = New-Object System.IO.StreamReader($stream)
                $response = $reader.ReadLine()
                if ($response -match "431|400|500") {
                    Write-Host "SUCCESS: Server rejected massive headers with status: $response" -ForegroundColor Green
                } else {
                    Write-Host "WARNING: Server accepted all $headersSent headers (Response: $response)" -ForegroundColor Yellow
                }
            } catch {
                $connectionClosed = $true
            }
        }
        
        if ($connectionClosed) {
            Write-Host "SUCCESS: Server proactively closed connection after $headersSent headers." -ForegroundColor Green
        }
        
        $tcpClient.Close()
    } catch {
        Write-Host "Error during Header Exhaustion test: $_" -ForegroundColor Red
    }
}

function Test-IdleConnectionFlood {
    Write-Host "`n[Scenario 3] Idle Connection Flood" -ForegroundColor Yellow
    Write-Host "Opening 200 connections and doing nothing. Verifying they are closed after IdleTimeout."
    
    $sockets = @()
    for ($i = 0; $i -lt 200; $i++) {
        try {
            $tcpClient = New-Object System.Net.Sockets.TcpClient($TargetHost, $TargetPort)
            $sockets += $tcpClient
        } catch { }
    }
    
    $opened = $sockets.Count
    Write-Host "Opened $opened connections successfully."
    Write-Host "Waiting for server to timeout connections (approx 6 seconds)..."
    
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
    
    if ($active -eq 0) {
        Write-Host "SUCCESS: All $opened idle connections were successfully timed out and closed by server." -ForegroundColor Green
    } else {
        Write-Host "FAIL: $active connections are still active." -ForegroundColor Red
    }
}

function Test-MalformedHTTP {
    Write-Host "`n[Scenario 4] Malformed HTTP Request" -ForegroundColor Yellow
    Write-Host "Sending garbage bytes instead of HTTP."
    
    try {
        $tcpClient = New-Object System.Net.Sockets.TcpClient($TargetHost, $TargetPort)
        $stream = $tcpClient.GetStream()
        
        $garbage = [System.Text.Encoding]::ASCII.GetBytes("JUST_GARBAGE_NO_SPACES_OR_PROTOCOL`r`n`r`n")
        $stream.Write($garbage, 0, $garbage.Length)
        
        $reader = New-Object System.IO.StreamReader($stream)
        $response = $reader.ReadLine()
        
        if ($response -match "400") {
            Write-Host "SUCCESS: Server handled garbage request with status: $response" -ForegroundColor Green
        } elseif ([string]::IsNullOrEmpty($response)) {
            Write-Host "SUCCESS: Server immediately closed connection on malformed request." -ForegroundColor Green
        } else {
            Write-Host "FAIL: Unexpected response to garbage: $response" -ForegroundColor Red
        }
        
        $tcpClient.Close()
    } catch {
        Write-Host "SUCCESS: Server immediately closed connection on malformed request." -ForegroundColor Green
    }
}

Test-SlowPOST
Test-HeaderExhaustion
Test-IdleConnectionFlood
Test-MalformedHTTP

Write-Host "`nAdvanced Security Tests Complete." -ForegroundColor Cyan
