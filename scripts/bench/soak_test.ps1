$Timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
$OutputFile = "benchmarks/soak/soak_$Timestamp.txt"

Write-Host "========================================"
Write-Host " TitanHTTP Soak Test"
Write-Host "========================================"
Write-Host "Target: http://localhost:8080/ping"
Write-Host "Connections: 100"
Write-Host "Duration: 1h"
Write-Host "Output File: $OutputFile"
Write-Host "----------------------------------------"
Write-Host "Note: Keep an eye on http://localhost:6060/debug/pprof/heap to monitor memory!"

bombardier -c 100 -d 1h -l http://localhost:8080/ping | Tee-Object -FilePath $OutputFile
