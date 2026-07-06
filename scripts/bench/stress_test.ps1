$Timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
$OutputFile = "benchmarks/stress/stress_$Timestamp.txt"

Write-Host "========================================"
Write-Host " TitanHTTP Stress Test"
Write-Host "========================================"
Write-Host "Target: http://localhost:8080/ping"
Write-Host "Connections: 1000"
Write-Host "Duration: 30s"
Write-Host "Output File: $OutputFile"
Write-Host "----------------------------------------"

bombardier -c 1000 -d 30s -l http://localhost:8080/ping | Tee-Object -FilePath $OutputFile
