from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware

app = FastAPI(title="AI Video Gen Service")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.get("/health")
def health():
    return {"success": True, "status": "ok", "message": "AI service is running"}


@app.get("/status")
def status():
    return {"success": True, "message": "Connected to AI service"}


@app.post("/generate")
async def generate(request: Request):
    body = await request.json()
    user_id = request.headers.get("X-User-ID", "unknown")
    user_email = request.headers.get("X-User-Email", "unknown")

    return {
        "success": True,
        "message": "Video generation request received",
        "user_id": user_id,
        "user_email": user_email,
        "received_data": body,
    }


@app.post("/generate-content-pillars")
async def generate_content_pillars(request: Request):
    body = await request.json()
    user_id = request.headers.get("X-User-ID", "unknown")

    return {
        "success": True,
        "message": "Content pillars generated",
        "user_id": user_id,
        "data": {
            "pillars": [
                {"title": "Brand Awareness", "description": "Content focused on increasing brand visibility"},
                {"title": "Product Education", "description": "Content that educates about product features"},
                {"title": "Social Proof", "description": "Content showcasing testimonials and success stories"},
            ]
        },
    }


@app.post("/generate-storyboard")
async def generate_storyboard(request: Request):
    body = await request.json()
    user_id = request.headers.get("X-User-ID", "unknown")

    return {
        "success": True,
        "message": "Storyboard generated",
        "user_id": user_id,
        "data": {
            "storyboards": [
                {
                    "title": "Dynamic Storyboard",
                    "scenes": [
                        {"scene_number": 1, "title": "Opening Hook", "duration": 5},
                        {"scene_number": 2, "title": "Problem Statement", "duration": 8},
                        {"scene_number": 3, "title": "Solution", "duration": 10},
                        {"scene_number": 4, "title": "Call to Action", "duration": 7},
                    ],
                }
            ]
        },
    }


@app.post("/generate-video")
async def generate_video(request: Request):
    body = await request.json()
    user_id = request.headers.get("X-User-ID", "unknown")

    return {
        "success": True,
        "message": "Video generation completed",
        "user_id": user_id,
        "data": {
            "status": "completed",
            "video_url": "/videos/sample-output.mp4",
            "duration": 30,
            "format": "mp4",
            "resolution": "1080p",
        },
    }


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
