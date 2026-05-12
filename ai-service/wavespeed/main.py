import os
import asyncio
import httpx
import uvicorn
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from fastapi.middleware.cors import CORSMiddleware

app = FastAPI()

# Enable CORS so your frontend can talk to this backend
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Configuration
WAVESPEED_API_KEY = ""
BASE_URL = "https://api.wavespeed.ai/api/v3"

# Request Model
class GenerateVideoRequest(BaseModel):
    prompt: str
    storyboard_id: str | None = None

@app.post("/submit")
async def submit_video(req: GenerateVideoRequest): # Used the Pydantic model for safety
    if not WAVESPEED_API_KEY:
        raise HTTPException(status_code=500, detail="API Key not found in environment variables")

    headers = {
        "Authorization": f"Bearer {WAVESPEED_API_KEY}",
        "Content-Type": "application/json"
    }

    async with httpx.AsyncClient() as client:
        res = await client.post(
            f"{BASE_URL}/wavespeed-ai/flux-dev",
            json={"prompt": req.prompt},
            headers=headers
        )

    data = res.json()
    if data.get("code") != 200:
        raise HTTPException(status_code=500, detail=data)

    return {
        "task_id": data["data"]["id"],
        "result_url": data["data"]["urls"]["get"]
    }

@app.get("/video-status/{task_id}")
async def get_status(task_id: str):
    headers = {"Authorization": f"Bearer {WAVESPEED_API_KEY}"}

    async with httpx.AsyncClient(timeout=30) as client:
        res = await client.get(
            f"{BASE_URL}/predictions/{task_id}/result",
            headers=headers
        )

    data = res.json()
    if data.get("code") != 200:
        raise HTTPException(status_code=500, detail=data)

    result = data["data"]
    return {
        "task_id": result["id"],
        "status": result["status"],
        "outputs": result.get("outputs", []),
        "error": result.get("error")
    }

@app.get("/video-wait/{task_id}")
async def wait_video(task_id: str):
    headers = {"Authorization": f"Bearer {WAVESPEED_API_KEY}"}

    async with httpx.AsyncClient(timeout=30) as client:
        for _ in range(60): 
            res = await client.get(
                f"{BASE_URL}/predictions/{task_id}/result",
                headers=headers
            )
            data = res.json()

            if data.get("code") != 200:
                raise HTTPException(status_code=500, detail=data)

            result = data["data"]
            status = result["status"]

            if status == "completed":
                return {
                    "status": status,
                    "video_urls": result.get("outputs", [])
                }
            if status == "failed":
                raise HTTPException(status_code=500, detail=result.get("error", "Task failed"))

            await asyncio.sleep(2) # Increased sleep slightly to be kind to the API

    raise HTTPException(status_code=408, detail="Timeout waiting for video")

# This is the part that was missing/commented out!
if __name__ == "__main__":
    uvicorn.run("main:app", host="localhost", port=8000, reload=True)