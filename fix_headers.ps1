$files = Get-ChildItem -Path internal,cmd -Recurse -Filter *.go
foreach ($file in $files) {
    $content = Get-Content $file.FullName -Raw
    
    $content = [regex]::Replace($content, '(?m)(\w+)\.Headers\["([^"]+)"\]\s*=\s*(.*?)$', '$1.Headers.Set("$2", $3)')
    $content = [regex]::Replace($content, '(?m)(\w+)\.Headers\[([a-zA-Z0-9_]+)\]\s*=\s*(.*?)$', '$1.Headers.Set($2, $3)')
    $content = [regex]::Replace($content, 'if ([a-zA-Z0-9_]+),\s*ok\s*:=\s*(\w+)\.Headers\["([^"]+)"\];\s*ok\s*\{', 'if $1 := $2.Headers.Get("$3"); $1 != "" {')
    $content = [regex]::Replace($content, 'if _,\s*ok\s*:=\s*(\w+)\.Headers\["([^"]+)"\];\s*!ok\s*\{', 'if $1.Headers.Get("$2") == "" {')
    $content = [regex]::Replace($content, 'if _,\s*has[A-Za-z0-9_]*\s*:=\s*(\w+)\.Headers\["([^"]+)"\];\s*has[A-Za-z0-9_]*\s*\{', 'if $1.Headers.Get("$2") != "" {')
    $content = [regex]::Replace($content, 'if _,\s*ok\s*:=\s*(\w+)\.Headers\["([^"]+)"\];\s*ok\s*\{', 'if $1.Headers.Get("$2") != "" {')
    
    $content = [regex]::Replace($content, '(\w+)\.Headers\["([^"]+)"\]', '$1.Headers.Get("$2")')
    $content = [regex]::Replace($content, '(\w+)\.Headers\[([a-zA-Z0-9_]+)\]', '$1.Headers.Get($2)')
    
    Set-Content -Path $file.FullName -Value $content -NoNewline
}
