$Timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
$OutputFile = "benchmarks/load/load_$Timestamp.txt"

Write-Host "========================================"
Write-Host " TitanHTTP Standard Load Test"
Write-Host "========================================"
Write-Host "Target: http://localhost:8080/ping"
Write-Host "Connections: 125"
Write-Host "Duration: 10s"
Write-Host "Output File: $OutputFile"
Write-Host "----------------------------------------"

bombardier -c 125 -d 10s -l http://localhost:8080/ping | Tee-Object -FilePath $OutputFile
