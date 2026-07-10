param (
    [Parameter(Mandatory=$true)]
    [string]$OldBench,
    
    [Parameter(Mandatory=$true)]
    [string]$NewBench
)

Write-Host "========================================"
Write-Host " TitanHTTP Benchstat Comparison"
Write-Host "========================================"
Write-Host "Old: $OldBench"
Write-Host "New: $NewBench"
Write-Host "----------------------------------------"

benchstat $OldBench $NewBench
