# Add Credits dengan Admin Credentials
# Login sebagai admin, dapatkan token, lalu add credits

$BackendUrl = "http://localhost:3000"
$AdminEmail = "admin@example.com"
$AdminPassword = "admin123"
$TestUserId = "2c4a7410-ef75-4f02-86ff-d54445528c98"
$CreditsToAdd = 5000

Write-Host "[STEP 1] Backend Connection Check..." -ForegroundColor Cyan

try {
    $health = Invoke-RestMethod -Uri "$BackendUrl/health" -Method GET -TimeoutSec 5
    Write-Host "[OK] Backend is running" -ForegroundColor Green
} catch {
    Write-Host "[FAIL] Backend not running" -ForegroundColor Red
    exit 1
}

Write-Host "`n[STEP 2] Admin Login..." -ForegroundColor Cyan
Write-Host "Email: $AdminEmail" -ForegroundColor Gray

$headers = @{"Content-Type" = "application/json"}

$loginBody = @{
    email = $AdminEmail
    password = $AdminPassword
} | ConvertTo-Json

try {
    $loginResponse = Invoke-RestMethod -Uri "$BackendUrl/api/auth/login" `
        -Method POST `
        -Headers $headers `
        -Body $loginBody `
        -TimeoutSec 10
    
    $AdminToken = $loginResponse.data.access_token
    $AdminRole = $loginResponse.data.user.role
    $AdminUserId = $loginResponse.data.user.id
    
    Write-Host "[OK] Login successful!" -ForegroundColor Green
    Write-Host "User ID: $AdminUserId" -ForegroundColor Gray
    Write-Host "Role: $AdminRole" -ForegroundColor Gray
    
    if ($AdminRole -ne "admin") {
        Write-Host "[WARNING] Role is '$AdminRole' not 'admin' - Add credits may fail" -ForegroundColor Yellow
    }
    
} catch {
    $statusCode = $_.Exception.Response.StatusCode.Value__
    Write-Host "[FAIL] Login failed (Status: $statusCode)" -ForegroundColor Red
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

Write-Host "`n[STEP 3] Adding $CreditsToAdd credits to user..." -ForegroundColor Cyan
Write-Host "Endpoint: POST $BackendUrl/api/admin/credits" -ForegroundColor Gray
Write-Host "Target User ID: $TestUserId" -ForegroundColor Gray

$adminHeaders = @{
    "Authorization" = "Bearer $AdminToken"
    "Content-Type" = "application/json"
}

$creditBody = @{
    user_id = $TestUserId
    amount = $CreditsToAdd
} | ConvertTo-Json

try {
    $creditResponse = Invoke-RestMethod -Uri "$BackendUrl/api/admin/credits" `
        -Method POST `
        -Headers $adminHeaders `
        -Body $creditBody `
        -TimeoutSec 10
    
    Write-Host "[OK] Credits added successfully!" -ForegroundColor Green
    Write-Host "`nResponse:" -ForegroundColor Gray
    Write-Host ($creditResponse | ConvertTo-Json -Depth 10) -ForegroundColor Gray
    
    Write-Host "`n[COMPLETE] Ready to run full test suite!" -ForegroundColor Green
    Write-Host "User now has $CreditsToAdd credits" -ForegroundColor Green
    
} catch {
    $statusCode = $_.Exception.Response.StatusCode.Value__
    
    if ($statusCode -eq 401) {
        Write-Host "`n[FAIL] Unauthorized (401) - Admin token not valid or user is not admin" -ForegroundColor Red
    } elseif ($statusCode -eq 404) {
        Write-Host "`n[FAIL] Endpoint not found (404)" -ForegroundColor Red
    } else {
        Write-Host "`n[FAIL] Request failed (Status: $statusCode)" -ForegroundColor Red
    }
    
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}
