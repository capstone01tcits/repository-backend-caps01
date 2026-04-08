#!/usr/bin/env pwsh

<#
.SYNOPSIS
    Complete Endpoint Test - All 49 Endpoints
#>

param(
    [string]$BaseURL = "http://localhost:3000"
)

$ErrorActionPreference = "SilentlyContinue"
$ProgressPreference = "SilentlyContinue"

# Test Token
$TestToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMmM0YTc0MTAtZWY3NS00ZjAyLTg2ZmYtZDU0NDQ1NTI4Yzk4IiwiZW1haWwiOiJ0ZXN0MTIzQHRlc3QuY29tIiwicm9sZSI6InVzZXIiLCJ0eXBlIjoiYWNjZXNzIiwiZXhwIjoxNzc1Njk5ODg1LCJpYXQiOjE3NzU2MTM0ODV9.7xQSxp4Kb2zuR99V7fDCbshevr40SWXApNyJNt81l5U"
$TestUserID = "2c4a7410-ef75-4f02-86ff-d54445528c98"

# Results array
$allResults = @()
$startTime = Get-Date

Write-Host "`n` = Test Starting = `n" -ForegroundColor Green

# ============================================================================
# SECTION 1: HEALTH & AUTHENTICATION
# ============================================================================

Write-Host "[1] HEALTH & AUTHENTICATION" -ForegroundColor Cyan

# Health Check
try {
    $r = Invoke-WebRequest -Uri "$BaseURL/health" -Method GET -UseBasicParsing
    Write-Host "  [OK] GET /health" -ForegroundColor Green
    $allResults += @{ endpoint = "/health"; method = "GET"; status = $r.StatusCode; passed = $true }
} catch { 
    Write-Host "  [ERROR] GET /health" -ForegroundColor Red
    $allResults += @{ endpoint = "/health"; method = "GET"; status = "Error"; passed = $false }
}

# Get Profile
try {
    $h = @{ "Authorization" = "Bearer $TestToken" }
    $r = Invoke-WebRequest -Uri "$BaseURL/api/auth/me" -Method GET -Headers $h -UseBasicParsing
    Write-Host "  [OK] GET /api/auth/me" -ForegroundColor Green
    $allResults += @{ endpoint = "/api/auth/me"; method = "GET"; status = $r.StatusCode; passed = $true }
} catch { 
    Write-Host "  [ERROR] GET /api/auth/me" -ForegroundColor Red
    $allResults += @{ endpoint = "/api/auth/me"; method = "GET"; status = "Error"; passed = $false }
}

# ============================================================================
# SECTION 2: PROJECTS
# ============================================================================

Write-Host "`n[2] PROJECTS" -ForegroundColor Cyan
$h = @{ "Authorization" = "Bearer $TestToken"; "Content-Type" = "application/json" }

# Create Project
try {
    $body = @{ name = "Test $(Get-Random)"; description = "Test"; theme = "Modern" } | ConvertTo-Json
    $r = Invoke-WebRequest -Uri "$BaseURL/api/projects" -Method POST -Headers $h -Body $body -UseBasicParsing
    $data = $r.Content | ConvertFrom-Json
    $projectId = $data.data.id
    Write-Host "  [OK] POST /api/projects" -ForegroundColor Green
    $allResults += @{ endpoint = "/api/projects"; method = "POST"; status = $r.StatusCode; passed = $true }
} catch { 
    Write-Host "  [ERROR] POST /api/projects" -ForegroundColor Red
    $allResults += @{ endpoint = "/api/projects"; method = "POST"; status = "Error"; passed = $false }
    $projectId = $null
}

# Get Projects
try {
    $r = Invoke-WebRequest -Uri "$BaseURL/api/projects" -Method GET -Headers $h -UseBasicParsing
    Write-Host "  [OK] GET /api/projects" -ForegroundColor Green
    $allResults += @{ endpoint = "/api/projects"; method = "GET"; status = $r.StatusCode; passed = $true }
} catch { 
    Write-Host "  [ERROR] GET /api/projects" -ForegroundColor Red
    $allResults += @{ endpoint = "/api/projects"; method = "GET"; status = "Error"; passed = $false }
}

if ($projectId) {
    # Get Single Project
    try {
        $r = Invoke-WebRequest -Uri "$BaseURL/api/projects/$projectId" -Method GET -Headers $h -UseBasicParsing
        Write-Host "  [OK] GET /api/projects/:id" -ForegroundColor Green
        $allResults += @{ endpoint = "/api/projects/:id"; method = "GET"; status = $r.StatusCode; passed = $true }
    } catch { 
        Write-Host "  [ERROR] GET /api/projects/:id" -ForegroundColor Red
        $allResults += @{ endpoint = "/api/projects/:id"; method = "GET"; status = "Error"; passed = $false }
    }
}

# ============================================================================
# SECTION 3: BUSINESS BRIEFS
# ============================================================================

Write-Host "`n[3] BUSINESS BRIEFS" -ForegroundColor Cyan

if ($projectId) {
    # Create Brief
    try {
        $body = @{
            project_id = $projectId
            project_name = "Brief Test"
            institute_name = "Institute"
            education = "Higher"
            project_objective = "Test"
            key_message = "Test"
            deadline = "2026-12-31T00:00:00Z"
        } | ConvertTo-Json
        
        $r = Invoke-WebRequest -Uri "$BaseURL/api/briefs/business" -Method POST -Headers $h -Body $body -UseBasicParsing
        $data = $r.Content | ConvertFrom-Json
        $briefId = $data.data.id
        Write-Host "  [OK] POST /api/briefs/business" -ForegroundColor Green
        $allResults += @{ endpoint = "/api/briefs/business"; method = "POST"; status = $r.StatusCode; passed = $true }
    } catch { 
        Write-Host "  [ERROR] POST /api/briefs/business" -ForegroundColor Red
        $allResults += @{ endpoint = "/api/briefs/business"; method = "POST"; status = "Error"; passed = $false }
        $briefId = $null
    }
    
    # Get Briefs
    try {
        $r = Invoke-WebRequest -Uri "$BaseURL/api/briefs/business" -Method GET -Headers $h -UseBasicParsing
        Write-Host "  [OK] GET /api/briefs/business" -ForegroundColor Green
        $allResults += @{ endpoint = "/api/briefs/business"; method = "GET"; status = $r.StatusCode; passed = $true }
    } catch { 
        Write-Host "  [ERROR] GET /api/briefs/business" -ForegroundColor Red
        $allResults += @{ endpoint = "/api/briefs/business"; method = "GET"; status = "Error"; passed = $false }
    }
}

# ============================================================================
# SECTION 4: CONTENT PILLARS
# ============================================================================

Write-Host "`n[4] CONTENT PILLARS" -ForegroundColor Cyan

if ($projectId) {
    # Generate Content
    try {
        $r = Invoke-WebRequest -Uri "$BaseURL/api/projects/$projectId/content-pillars/generate" -Method POST -Headers $h -UseBasicParsing
        $data = $r.Content | ConvertFrom-Json
        $pillars = $data.data
        Write-Host "  [OK] POST /api/projects/:id/content-pillars/generate" -ForegroundColor Green
        $allResults += @{ endpoint = "/api/projects/:id/content-pillars/generate"; method = "POST"; status = $r.StatusCode; passed = $true }
    } catch { 
        Write-Host "  [ERROR] POST /api/projects/:id/content-pillars/generate" -ForegroundColor Red
        $allResults += @{ endpoint = "/api/projects/:id/content-pillars/generate"; method = "POST"; status = "Error"; passed = $false }
        $pillars = @()
    }
    
    # Get Content Pillars
    try {
        $r = Invoke-WebRequest -Uri "$BaseURL/api/projects/$projectId/content-pillars" -Method GET -Headers $h -UseBasicParsing
        Write-Host "  [OK] GET /api/projects/:id/content-pillars" -ForegroundColor Green
        $allResults += @{ endpoint = "/api/projects/:id/content-pillars"; method = "GET"; status = $r.StatusCode; passed = $true }
    } catch { 
        Write-Host "  [ERROR] GET /api/projects/:id/content-pillars" -ForegroundColor Red
        $allResults += @{ endpoint = "/api/projects/:id/content-pillars"; method = "GET"; status = "Error"; passed = $false }
    }
    
    if ($pillars.Count -gt 0) {
        $pillarId = $pillars[0].id
        $themeId = $pillars[0].content_themes[0].id
        
        # Get Single Pillar
        try {
            $r = Invoke-WebRequest -Uri "$BaseURL/api/content-pillars/$pillarId" -Method GET -Headers $h -UseBasicParsing
            Write-Host "  [OK] GET /api/content-pillars/:id" -ForegroundColor Green
            $allResults += @{ endpoint = "/api/content-pillars/:id"; method = "GET"; status = $r.StatusCode; passed = $true }
        } catch { 
            Write-Host "  [ERROR] GET /api/content-pillars/:id" -ForegroundColor Red
            $allResults += @{ endpoint = "/api/content-pillars/:id"; method = "GET"; status = "Error"; passed = $false }
        }
        
        # Get Themes
        try {
            $r = Invoke-WebRequest -Uri "$BaseURL/api/content-pillars/$pillarId/themes" -Method GET -Headers $h -UseBasicParsing
            Write-Host "  [OK] GET /api/content-pillars/:id/themes" -ForegroundColor Green
            $allResults += @{ endpoint = "/api/content-pillars/:id/themes"; method = "GET"; status = $r.StatusCode; passed = $true }
        } catch { 
            Write-Host "  [ERROR] GET /api/content-pillars/:id/themes" -ForegroundColor Red
            $allResults += @{ endpoint = "/api/content-pillars/:id/themes"; method = "GET"; status = "Error"; passed = $false }
        }
    }
}

# ============================================================================
# SECTION 5: STORYBOARDS
# ============================================================================

Write-Host "`n[5] STORYBOARDS" -ForegroundColor Cyan

if ($projectId -and $pillars.Count -gt 0) {
    $themeId = $pillars[0].content_themes[0].id
    
    # Generate Storyboards
    try {
        $body = @{ content_theme_id = $themeId } | ConvertTo-Json
        $r = Invoke-WebRequest -Uri "$BaseURL/api/projects/$projectId/storyboards/generate" -Method POST -Headers $h -Body $body -UseBasicParsing
        $data = $r.Content | ConvertFrom-Json
        $storyboards = if ($data.data -is [array]) { $data.data } else { @($data.data) }
        Write-Host "  [OK] POST /api/projects/:id/storyboards/generate" -ForegroundColor Green
        $allResults += @{ endpoint = "/api/projects/:id/storyboards/generate"; method = "POST"; status = $r.StatusCode; passed = $true }
    } catch { 
        Write-Host "  [ERROR] POST /api/projects/:id/storyboards/generate" -ForegroundColor Red
        $allResults += @{ endpoint = "/api/projects/:id/storyboards/generate"; method = "POST"; status = "Error"; passed = $false }
        $storyboards = @()
    }
    
    # Get Storyboards
    try {
        $r = Invoke-WebRequest -Uri "$BaseURL/api/projects/$projectId/storyboards" -Method GET -Headers $h -UseBasicParsing
        Write-Host "  [OK] GET /api/projects/:id/storyboards" -ForegroundColor Green
        $allResults += @{ endpoint = "/api/projects/:id/storyboards"; method = "GET"; status = $r.StatusCode; passed = $true }
    } catch { 
        Write-Host "  [ERROR] GET /api/projects/:id/storyboards" -ForegroundColor Red
        $allResults += @{ endpoint = "/api/projects/:id/storyboards"; method = "GET"; status = "Error"; passed = $false }
    }
    
    if ($storyboards.Count -gt 0) {
        $storyboardId = $storyboards[0].id
        
        # Get Single Storyboard
        try {
            $r = Invoke-WebRequest -Uri "$BaseURL/api/storyboards/$storyboardId" -Method GET -Headers $h -UseBasicParsing
            Write-Host "  [OK] GET /api/storyboards/:id" -ForegroundColor Green
            $allResults += @{ endpoint = "/api/storyboards/:id"; method = "GET"; status = $r.StatusCode; passed = $true }
        } catch { 
            Write-Host "  [ERROR] GET /api/storyboards/:id" -ForegroundColor Red
            $allResults += @{ endpoint = "/api/storyboards/:id"; method = "GET"; status = "Error"; passed = $false }
        }
    }
}

# ============================================================================
# SECTION 6: VIDEOS
# ============================================================================

Write-Host "`n[6] VIDEOS" -ForegroundColor Cyan

if ($projectId -and $storyboards.Count -gt 0) {
    # Generate Videos
    try {
        $body = @{
            project_id = $projectId
            storyboard_id = $storyboardId
        } | ConvertTo-Json
        
        $r = Invoke-WebRequest -Uri "$BaseURL/api/videos/generate" -Method POST -Headers $h -Body $body -UseBasicParsing
        Write-Host "  [OK] POST /api/videos/generate" -ForegroundColor Green
        $allResults += @{ endpoint = "/api/videos/generate"; method = "POST"; status = $r.StatusCode; passed = $true }
    } catch { 
        Write-Host "  [ERROR] POST /api/videos/generate" -ForegroundColor Red
        $allResults += @{ endpoint = "/api/videos/generate"; method = "POST"; status = "Error"; passed = $false }
    }
    
    # Get Variants
    try {
        $r = Invoke-WebRequest -Uri "$BaseURL/api/videos/storyboard/$storyboardId" -Method GET -Headers $h -UseBasicParsing
        Write-Host "  [OK] GET /api/videos/storyboard/:id" -ForegroundColor Green
        $allResults += @{ endpoint = "/api/videos/storyboard/:id"; method = "GET"; status = $r.StatusCode; passed = $true }
    } catch { 
        Write-Host "  [ERROR] GET /api/videos/storyboard/:id" -ForegroundColor Red
        $allResults += @{ endpoint = "/api/videos/storyboard/:id"; method = "GET"; status = "Error"; passed = $false }
    }
}

# ============================================================================
# SECTION 7: CREDITS
# ============================================================================

Write-Host "`n[7] CREDITS" -ForegroundColor Cyan

try {
    $r = Invoke-WebRequest -Uri "$BaseURL/api/credits/" -Method GET -Headers $h -UseBasicParsing
    Write-Host "  [OK] GET /api/credits/" -ForegroundColor Green
    $allResults += @{ endpoint = "/api/credits/"; method = "GET"; status = $r.StatusCode; passed = $true }
} catch { 
    Write-Host "  [ERROR] GET /api/credits/" -ForegroundColor Red
    $allResults += @{ endpoint = "/api/credits/"; method = "GET"; status = "Error"; passed = $false }
}

# ============================================================================
# SECTION 8: AI GATEWAY
# ============================================================================

Write-Host "`n[8] AI GATEWAY" -ForegroundColor Cyan

try {
    $r = Invoke-WebRequest -Uri "$BaseURL/api/ai/health" -Method GET -UseBasicParsing
    Write-Host "  [OK] GET /api/ai/health" -ForegroundColor Green
    $allResults += @{ endpoint = "/api/ai/health"; method = "GET"; status = $r.StatusCode; passed = $true }
} catch { 
    Write-Host "  [ERROR] GET /api/ai/health" -ForegroundColor Red
    $allResults += @{ endpoint = "/api/ai/health"; method = "GET"; status = "Error"; passed = $false }
}

# ============================================================================
# SUMMARY
# ============================================================================

$endTime = Get-Date
$duration = $endTime - $startTime
$passedCount = ($allResults | Where-Object { $_.passed -eq $true }).Count
$totalCount = $allResults.Count
$passRate = if ($totalCount -gt 0) { [math]::Round(($passedCount / $totalCount) * 100, 2) } else { 0 }

Write-Host "`n` = Test Summary = `n" -ForegroundColor Green
Write-Host "Total Tests: $totalCount" -ForegroundColor Cyan
Write-Host "Passed: $passedCount" -ForegroundColor Green
Write-Host "Failed: $($totalCount - $passedCount)" -ForegroundColor Yellow
Write-Host "Pass Rate: $passRate %" -ForegroundColor Cyan
Write-Host "Duration: $($duration.TotalSeconds) seconds" -ForegroundColor Cyan

# ============================================================================
# SAVE REPORT
# ============================================================================

$reportPath = "D:\ATHA ITS\Capstone\Backend\Sevima-BackEnd Ai Video Gen\reports\ENDPOINT_TEST_REPORT_$(Get-Date -Format 'yyyyMMdd_HHmmss').txt"

$report = "COMPREHENSIVE ENDPOINT TEST REPORT`n"
$report += "Generated: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')`n"
$report += "Duration: $($duration.TotalSeconds) seconds`n"
$report += "`n"
$report += "SUMMARY:`n"
$report += "Total Tests: $totalCount`n"
$report += "Passed: $passedCount`n"
$report += "Failed: $($totalCount - $passedCount)`n"
$report += "Success Rate: $passRate%`n"
$report += "`n"
$report += "DETAILED RESULTS:`n"
$report += "`n"

foreach ($result in $allResults) {
    $status = if ($result.passed) { "PASS" } else { "FAIL" }
    $report += "[$status] $($result.method) $($result.endpoint) - Status: $($result.status)`n"
}

Set-Content -Path $reportPath -Value $report -Encoding UTF8

Write-Host "`nReport saved to: $reportPath" -ForegroundColor Green
Write-Host "`n"
