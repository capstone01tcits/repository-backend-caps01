"""
FastAPI App — AI Video Generator Service.

Endpoint yang disepakati dengan divisi backend:

  POST /generate          → submit job, return job_id
  GET  /status/{job_id}   → cek status: pending/processing/done/failed
  GET  /result/{job_id}   → ambil video_url setelah done
  GET  /health            → health check (untuk load balancer / monitoring)

Backend (Go/FastAPI lain) tinggal hit endpoint ini.
Tidak perlu tahu apapun tentang routing.
"""

import sys
import os
import logging
from typing import List




# Add current directory to Python path for imports
sys.path.insert(0, os.path.dirname(__file__))

# FIX 1: load_dotenv() dipanggil SEBELUM import local
# agar env vars sudah ter-load saat module diinisialisasi
from dotenv import load_dotenv
load_dotenv(override=True)

# FIX 2: Import local SETELAH load_dotenv() — urutan ini penting
from logging_config import setup_logging
from services.ai_video_service import AIVideoService

from fastapi import FastAPI, HTTPException
from fastapi.responses import FileResponse
from pydantic import BaseModel, Field
from fastapi.middleware.cors import CORSMiddleware

app = FastAPI(
    title="AI Video Generator — ITS Marketing",
    description="Service untuk generate video promosi kampus.",
    version="1.0.0",
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=[
        "http://localhost:3000",
    ],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

setup_logging(
    log_dir=os.getenv("LOG_DIR", "outputs/logs"),
    log_level=os.getenv("LOG_LEVEL", "INFO"),
)

logger = logging.getLogger(__name__)

service = AIVideoService()


# ---------------------------------------------------------------------------
# Request / Response Schema
# ---------------------------------------------------------------------------

class GenerateRequest(BaseModel):
    """Request body untuk POST /generate."""
    prompt: str = Field(
        ...,
        min_length=10,
        description="Deskripsi video yang ingin digenerate",
    )
    duration: int = Field(
        default=10,
        ge=1,
        le=30,
        description="Durasi video dalam detik (1-30)",
    )
    ratio: str = Field(
        default="16:9",
        description="Aspect ratio: 16:9 | 9:16 | 1:1 | 4:3",
    )
    task_type: str = Field(
        default="text_to_video",
        description="Jenis task: text_to_video | text_to_video_hq | image_to_video",
    )

    # Video mode (menentukan model WaveSpeed yang dipakai)
    video_mode: str = Field(default="text-to-video", description="text-to-video | image-to-video | start-end-to-video")

    # Image inputs (CDN URLs)
    start_image: str = Field(default="", description="URL gambar start frame (image-to-video / start-end-to-video)")
    end_image: str = Field(default="", description="URL gambar end frame (start-end-to-video)")

    # Generation controls
    negative_prompt: str = Field(default="", description="Hal yang ingin dihindari dalam video")
    generate_audio: bool = Field(default=False, description="Generate audio sinkron (text-to-video only)")
    seed: int = Field(default=-1, description="-1 = random")
    resolution: str = Field(default="480p", description="Resolusi output: 480p | 720p | 1080p")

    model_config = {
        "json_schema_extra": {
            "example": {
                "prompt": "Cinematic video of Institut Teknologi Sepuluh Nopember Surabaya campus at golden hour, drone shot",
                "duration": 6,
                "ratio": "16:9",
                "task_type": "veo3",
                "resolution": "1080p",
                "generate_audio": True,
            }
        }
    }


class GenerateResponse(BaseModel):
    """Response dari POST /generate."""
    job_id: str
    status: str
    message: str


# ---------------------------------------------------------------------------
# Endpoints
# ---------------------------------------------------------------------------

@app.get("/health")
def health_check():
    """Health check endpoint untuk monitoring."""
    return {"status": "ok", "service": "ai-video-generator"}


class Veo3Payload(BaseModel):
    """Request body untuk POST /api/veo3/generate."""
    model: str = Field(default="veo-3.1-lite")
    prompt: str = Field(..., description="Prompt lengkap dengan format SCENE")
    reference_images: List[str] = Field(default_factory=list)
    video_mode: str = Field(default="text-to-video")
    start_image: str = Field(default="")
    end_image: str = Field(default="")
    negative_prompt: str = Field(default="")
    generate_audio: bool = Field(default=False)
    seed: int = Field(default=-1)
    resolution: str = Field(default="480p")

@app.post("/generate", response_model=GenerateResponse, status_code=202)
def generate_video(body: GenerateRequest):
    """
    Submit job generate video.
    """
    logger.info("[API] POST /generate prompt_preview=%.60s", body.prompt)
    try:
        job = service.submit(
            prompt=body.prompt,
            duration=body.duration,
            ratio=body.ratio,
            task_type=body.task_type,
            video_mode=body.video_mode,
            start_image=body.start_image,
            end_image=body.end_image,
            negative_prompt=body.negative_prompt,
            generate_audio=body.generate_audio,
            seed=body.seed,
            resolution=body.resolution,
        )
    except ValueError as exc:
        raise HTTPException(status_code=422, detail=str(exc)) from exc

    return GenerateResponse(
        job_id=job.job_id,
        status=job.status.value,
        message="Job diterima. Gunakan GET /status/{job_id} untuk cek progress.",
    )


@app.post("/api/veo3/generate", response_model=GenerateResponse, status_code=202)
def generate_veo3(body: Veo3Payload):
    """
    Endpoint khusus Veo 3 untuk pemrosesan prompt otomatis dan stitching.
    """
    logger.info("[API] POST /api/veo3/generate model=%s", body.model)
    try:
        # Kita gunakan task_type 'veo3' agar router memilih Veo3Provider
        job = service.submit(
            prompt=body.prompt,
            duration=15,
            ratio="16:9",
            task_type="veo3",
            reference_images=body.reference_images,
            video_mode=body.video_mode,
            start_image=body.start_image,
            end_image=body.end_image,
            negative_prompt=body.negative_prompt,
            generate_audio=body.generate_audio,
            seed=body.seed,
            resolution=body.resolution,
        )
    except ValueError as exc:
        raise HTTPException(status_code=422, detail=str(exc)) from exc

    return GenerateResponse(
        job_id=job.job_id,
        status=job.status.value,
        message="Veo 3 job diterima. Video akan digenerate dan digabung secara otomatis.",
    )


@app.get("/status/{job_id}")
def get_status(job_id: str):
    """
    Cek status job.

    Response:
    {
        "job_id": "...",
        "status": "pending" | "processing" | "done" | "failed",
        "created_at": "...",
        "updated_at": "...",
        "error": null
    }
    """
    logger.debug("[API] GET /status/%s", job_id)
    try:
        return service.get_status(job_id)
    except KeyError as exc:
        raise HTTPException(
            status_code=404,
            detail=f"Job '{job_id}' tidak ditemukan",
        ) from exc


@app.get("/result/{job_id}")
def get_result(job_id: str):
    """
    Ambil hasil video setelah status done.

    Response saat done:
    {
        "job_id": "...",
        "status": "done",
        "video_url": "http://127.0.0.1:8000/video/xxx.mp4",
        "prompt": "...",
        "meta": { "provider": ..., "model": ..., "duration": ..., "ratio": ... }
    }
    """
    logger.debug("[API] GET /result/%s", job_id)
    try:
        return service.get_result(job_id)
    except KeyError as exc:
        raise HTTPException(
            status_code=404,
            detail=f"Job '{job_id}' tidak ditemukan",
        ) from exc


@app.get("/video/{filename}")
def serve_video(filename: str):
    """
    Serve file video langsung (untuk development).
    Di production: serve via Nginx / CDN.
    """
    video_dir = os.getenv("VIDEO_OUTPUT_DIR", "outputs/videos")
    filepath = os.path.join(video_dir, filename)
    if not os.path.exists(filepath):
        raise HTTPException(status_code=404, detail="File tidak ditemukan")
    media_type = "image/jpeg" if filename.endswith(".jpg") else "video/mp4"
    return FileResponse(filepath, media_type=media_type)


# ===== RUN SERVER =====
if __name__ == "__main__":
    import uvicorn
    
    host = os.getenv("AI_SERVICE_HOST", "127.0.0.1")
    port = int(os.getenv("AI_SERVICE_PORT", "8000"))
    
    logger.info(f"Starting AI Video Service on {host}:{port}...")
    uvicorn.run(
        app,
        host=host,
        port=port,
        reload=False,  # Disable reload for direct script execution
        log_level="info",
    )