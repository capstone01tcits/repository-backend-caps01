$baseUrl = "http://localhost:3000"
$global:TOKEN = ""
$global:REFRESH = ""
$global:PROJECT_ID = ""
$global:BRIEF_ID = ""
$global:CREATIVE_ID = ""
$global:PILLAR_ID = ""
$global:THEME_ID = ""
$global:STORY_ID = ""
$global:JOB_ID = ""
$global:VAR_ID = ""
$global:SCENE_ID = ""

Write-Host "====== 50 ENDPOINT COMPREHENSIVE TEST ======" -ForegroundColor Cyan
Write-Host ""

function T {
    param([int]$n, [string]$m, [string]$p, [object]$b, [string]$t)
    try {
        $h = @{"Content-Type" = "application/json"}
        if ($t) { $h["Authorization"] = "Bearer $t" }
        $pa = @{Uri = "$baseUrl$p"; Method = $m; Headers = $h; UseBasicParsing = $true; ErrorAction = "Stop"}
        if ($b) { $pa["Body"] = ($b | ConvertTo-Json -Depth 10) }
        $r = Invoke-RestMethod @pa
        Write-Host "[$n] OK $m $p" -ForegroundColor Green
        return $r
    }
    catch {
        $c = "ERROR"
        if ($_.Exception.Response) { $c = $_.Exception.Response.StatusCode.Value }
        Write-Host "[$n] FAIL $m $p ($c)" -ForegroundColor Red
        return $null
    }
}

# Setup unique email
$emailUnique = "test$(Get-Random)@test.com"
$global:ADMIN_TOKEN = ""
$global:USER_ID = ""

# 1-2: Health
T 1 "GET" "/health" $null $null | Out-Null
T 2 "GET" "/api/ai/health" $null $null | Out-Null

# 3-8: Auth
$r = T 3 "POST" "/api/auth/register" @{name="Test User"; email=$emailUnique; password="secure123"} $null
$global:TOKEN = $r.data.access_token
$global:REFRESH = $r.data.refresh_token
$global:USER_ID = if ($r) { $r.data.user.id } else { $null }

# Auto-add 500 credits via admin
if ($global:USER_ID) {
    $adminLoginResp = Invoke-RestMethod -Uri "$baseUrl/api/auth/login" -Method POST -Headers @{"Content-Type" = "application/json"} -Body (@{email = "admin@example.com"; password = "admin123"} | ConvertTo-Json) -ErrorAction SilentlyContinue
    if ($adminLoginResp) {
        $global:ADMIN_TOKEN = $adminLoginResp.data.access_token
        Invoke-RestMethod -Uri "$baseUrl/api/admin/credits" -Method POST -Headers @{"Content-Type" = "application/json"; "Authorization" = "Bearer $global:ADMIN_TOKEN"} -Body (@{user_id = $global:USER_ID; amount = 500} | ConvertTo-Json) -ErrorAction SilentlyContinue
    }
}

$r = T 4 "POST" "/api/auth/login" @{email=$emailUnique; password="secure123"} $null
if ($r) { $global:TOKEN = $r.data.access_token }

T 5 "GET" "/api/auth/me" $null $global:TOKEN | Out-Null

$r = T 6 "POST" "/api/auth/refresh" @{refresh_token = $global:REFRESH} $null
if ($r) { $global:TOKEN = $r.data.access_token }

T 7 "POST" "/api/auth/change-password" @{old_password = "secure123"; new_password = "newSecure456"} $global:TOKEN | Out-Null

$r = T 8 "POST" "/api/auth/login" @{email = $emailUnique; password = "newSecure456"} $null
if ($r) { $global:TOKEN = $r.data.access_token }

# 9-15: Projects & Briefs
Write-Host ""
$r = T 9 "POST" "/api/projects" @{name = "Test Project Q2"; description = "Testing"; theme = "Corporate Branding"} $global:TOKEN
$global:PROJECT_ID = $r.data.id

T 10 "GET" "/api/projects" $null $global:TOKEN | Out-Null
T 11 "GET" "/api/projects/$global:PROJECT_ID" $null $global:TOKEN | Out-Null
T 12 "PUT" "/api/projects/$global:PROJECT_ID" @{name = "Updated Project"; theme = "Modern"} $global:TOKEN | Out-Null

$r = T 13 "POST" "/api/briefs/business" @{project_id = $global:PROJECT_ID; project_name = "Test Project Q2"; company_name = "PT Tech"; education = "Higher"; target_audience = "Tech Professionals"; deadline = "2026-12-31T00:00:00Z"} $global:TOKEN
$global:BRIEF_ID = if ($r) { $r.data.id } else { $null }

T 14 "GET" "/api/briefs/business" $null $global:TOKEN | Out-Null
if ($global:BRIEF_ID) { T 15 "GET" "/api/briefs/business/$global:BRIEF_ID" $null $global:TOKEN | Out-Null } else { Write-Host "[15] SKIP GET /api/briefs/business/{id} (no brief id)" -ForegroundColor Yellow }

# 16-20: Creative Briefs
Write-Host ""
if ($global:BRIEF_ID) { T 16 "PUT" "/api/briefs/business/$global:BRIEF_ID" @{company_name = "Updated Company"} $global:TOKEN | Out-Null } else { Write-Host "[16] SKIP PUT /api/briefs/business/{id} (no brief id)" -ForegroundColor Yellow }

if ($global:BRIEF_ID) {
    $r = T 17 "POST" "/api/briefs/creative" @{business_brief_id = $global:BRIEF_ID; title = "Product Video"; video_type = "promo"; duration = 60} $global:TOKEN
    $global:CREATIVE_ID = if ($r) { $r.data.id } else { $null }
} else {
    Write-Host "[17] SKIP POST /api/briefs/creative (no brief id)" -ForegroundColor Yellow
}

T 18 "GET" "/api/briefs/creative" $null $global:TOKEN | Out-Null
if ($global:CREATIVE_ID) { T 19 "GET" "/api/briefs/creative/$global:CREATIVE_ID" $null $global:TOKEN | Out-Null } else { Write-Host "[19] SKIP GET /api/briefs/creative/{id} (no creative id)" -ForegroundColor Yellow }
if ($global:CREATIVE_ID) { T 20 "PUT" "/api/briefs/creative/$global:CREATIVE_ID" @{title = "Updated"; duration = 90} $global:TOKEN | Out-Null } else { Write-Host "[20] SKIP PUT /api/briefs/creative/{id} (no creative id)" -ForegroundColor Yellow }

# 21-27: Content Pillars
Write-Host ""
$r = T 21 "POST" "/api/projects/$global:PROJECT_ID/content-pillars/generate" @{} $global:TOKEN
if ($r.data -and @($r.data).Count -gt 0) {
    $global:PILLAR_ID = @($r.data)[0].id
    $global:THEME_ID = @(@($r.data)[0].content_themes)[0].id
}

T 22 "GET" "/api/projects/$global:PROJECT_ID/content-pillars" $null $global:TOKEN | Out-Null
T 23 "GET" "/api/content-pillars/$global:PILLAR_ID" $null $global:TOKEN | Out-Null
T 24 "POST" "/api/content-pillars/$global:PILLAR_ID/select" @{} $global:TOKEN | Out-Null
T 25 "PUT" "/api/content-pillars/$global:PILLAR_ID" @{prompt = "Custom prompt"} $global:TOKEN | Out-Null
T 26 "GET" "/api/content-pillars/$global:PILLAR_ID/themes" $null $global:TOKEN | Out-Null
T 27 "POST" "/api/content-themes/$global:THEME_ID/select" @{} $global:TOKEN | Out-Null

# 28-33: Storyboards
Write-Host ""
$r = T 28 "POST" "/api/projects/$global:PROJECT_ID/storyboards/generate" @{content_theme_id = $global:THEME_ID; prompt = "Cinematic product demo"} $global:TOKEN
if ($r.data -and @($r.data).Count -gt 0) {
    $global:STORY_ID = @($r.data)[0].id
}

T 29 "GET" "/api/projects/$global:PROJECT_ID/storyboards" $null $global:TOKEN | Out-Null
T 30 "GET" "/api/storyboards/$global:STORY_ID" $null $global:TOKEN | Out-Null
T 31 "POST" "/api/storyboards/$global:STORY_ID/select" @{} $global:TOKEN | Out-Null
T 32 "PUT" "/api/storyboards/$global:STORY_ID" @{prompt = "Updated prompt"} $global:TOKEN | Out-Null

$r = T 33 "GET" "/api/storyboards/$global:STORY_ID/scenes" $null $global:TOKEN
if ($r.data -and @($r.data).Count -gt 0) {
    $global:SCENE_ID = @($r.data)[0].id
}

# 34-40: Videos
Write-Host ""
if ($global:STORY_ID) {
    $r = T 34 "POST" "/api/videos/generate" @{project_id = $global:PROJECT_ID; storyboard_id = $global:STORY_ID} $global:TOKEN
    $global:JOB_ID = if ($r) { $r.data.generation_job_id } else { $null }
} else {
    Write-Host "[34] SKIP POST /api/videos/generate (no storyboard id)" -ForegroundColor Yellow
}

if ($global:JOB_ID) { T 35 "GET" "/api/videos/generation/$global:JOB_ID" $null $global:TOKEN | Out-Null } else { Write-Host "[35] SKIP GET /api/videos/generation/{jobId} (no job id)" -ForegroundColor Yellow }

if ($global:STORY_ID) {
    $r = T 36 "GET" "/api/videos/storyboard/$global:STORY_ID" $null $global:TOKEN
    if ($r.data -and @($r.data).Count -gt 0) {
        $global:VAR_ID = @($r.data)[0].id
        # Extract scene ID from variant if available
        if ($r.data[0].scenes -and @($r.data[0].scenes).Count -gt 0) {
            $global:SCENE_ID = @($r.data[0].scenes)[0].id
        }
    }
} else {
    Write-Host "[36] SKIP GET /api/videos/storyboard/{id} (no storyboard id)" -ForegroundColor Yellow
}

if ($global:VAR_ID) { T 37 "GET" "/api/videos/$global:VAR_ID" $null $global:TOKEN | Out-Null } else { Write-Host "[37] SKIP GET /api/videos/{id} (no variant id)" -ForegroundColor Yellow }
if ($global:VAR_ID) { T 38 "GET" "/api/videos/$global:VAR_ID/download" $null $global:TOKEN | Out-Null } else { Write-Host "[38] SKIP GET /api/videos/{id}/download (no variant id)" -ForegroundColor Yellow }
if ($global:VAR_ID) { T 39 "POST" "/api/videos/$global:VAR_ID/regenerate" @{new_prompt = "Vibrant colors"} $global:TOKEN | Out-Null } else { Write-Host "[39] SKIP POST /api/videos/{id}/regenerate (no variant id)" -ForegroundColor Yellow }
if ($global:SCENE_ID) { T 40 "POST" "/api/videos/scene/$global:SCENE_ID/regenerate" @{new_prompt = "Close-up detail"} $global:TOKEN | Out-Null } else { Write-Host "[40] SKIP POST /api/videos/scene/{id}/regenerate (no scene id)" -ForegroundColor Yellow }

# 41-45: Credits & Admin
Write-Host ""
T 41 "GET" "/api/credits/" $null $global:TOKEN | Out-Null

T 42 "POST" "/api/ai/generate" @{prompt = "AI video"; duration = 30} $global:TOKEN | Out-Null
T 43 "GET" "/api/ai/status" $null $global:TOKEN | Out-Null

$r = T 44 "POST" "/api/auth/register" @{name = "Admin2"; email = "admin2_$(Get-Random)@test.com"; password = "admin123"} $null
$targetUserId = if ($r) { $r.data.user.id } else { $null }

if ($global:ADMIN_TOKEN -and $targetUserId) { T 45 "POST" "/api/admin/credits" @{user_id = $targetUserId; amount = 10} $global:ADMIN_TOKEN | Out-Null } else { Write-Host "[45] SKIP POST /api/admin/credits (no admin token)" -ForegroundColor Yellow }

# 46-50: Cleanup
Write-Host ""
if ($global:CREATIVE_ID) { T 46 "DELETE" "/api/briefs/creative/$global:CREATIVE_ID" $null $global:TOKEN | Out-Null } else { Write-Host "[46] SKIP DELETE /api/briefs/creative/{id} (no creative id)" -ForegroundColor Yellow }
if ($global:BRIEF_ID) { T 47 "DELETE" "/api/briefs/business/$global:BRIEF_ID" $null $global:TOKEN | Out-Null } else { Write-Host "[47] SKIP DELETE /api/briefs/business/{id} (no brief id)" -ForegroundColor Yellow }
if ($global:PROJECT_ID) { T 48 "DELETE" "/api/projects/$global:PROJECT_ID" $null $global:TOKEN | Out-Null } else { Write-Host "[48] SKIP DELETE /api/projects/{id} (no project id)" -ForegroundColor Yellow }
T 49 "DELETE" "/api/auth/account" $null $global:TOKEN | Out-Null
T 50 "POST" "/api/auth/restore" @{refresh_token = $global:REFRESH} $null | Out-Null

Write-Host ""
Write-Host "TEST COMPLETE" -ForegroundColor Cyan
