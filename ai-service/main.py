"""
FastAPI App — AI Video Generator Service.

Endpoint utama:
  POST /generate          → submit job, proses di background, return job_id
  POST /api/veo3/generate → submit job khusus Veo3
  GET  /status/{job_id}   → cek status: pending/processing/done/failed
  GET  /result/{job_id}   → ambil video_url setelah done
  GET  /health            → health check
"""

from __future__ import annotations

import sys
import os
import logging
from pathlib import Path
from typing import List

# 1. Load environment variables
from dotenv import load_dotenv
load_dotenv(override=True)

# 2. Setup system path agar modul lokal (services, logging_config) terbaca
current_dir = Path(__file__).parent.absolute()
if str(current_dir) not in sys.path:
    sys.path.insert(0, str(current_dir))

from logging_config import setup_logging
from services.ai_video_service import AIVideoService

from fastapi import FastAPI, HTTPException
from fastapi.responses import FileResponse
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel, Field

# 3. Setup Logging
setup_logging(
    log_dir=os.getenv("LOG_DIR", "outputs/logs"),
    log_level=os.getenv("LOG_LEVEL", "INFO"),
)

logger = logging.getLogger(__name__)

# 4. Inisialisasi FastAPI
app = FastAPI(
    title="AI Video Generator — ITS Marketing",
    description="Service untuk generate video promosi kampus menggunakan LTX / Veo 3 / Wavespeed.",
    version="1.0.0",
)

# 5. Konfigurasi CORS (Penting untuk koneksi Frontend React/Next.js)
app.add_middleware(
    CORSMiddleware,
    allow_origins=[
        "http://localhost:3000",
        "http://127.0.0.1:3000",
    ],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

service = AIVideoService()


# ---------------------------------------------------------------------------
# Request / Response Schema
# ---------------------------------------------------------------------------

class GenerateRequest(BaseModel):
    """Request body untuk POST /generate."""
    prompt: str = Field(
        ..., min_length=5, description="Deskripsi video yang akan digenerate"
    )
    duration: int = Field(
        default=10, ge=1, le=30, description="Durasi video dalam detik"
    )
    ratio: str = Field(
        default="16:9", description="Aspect ratio: 16:9 | 9:16 | 1:1 | 4:3"
    )
    task_type: str = Field(
        default="text_to_video", description="Jenis task: text_to_video | veo3"
    )

    model_config = {
        "json_schema_extra": {
            "example": {
                "prompt": "Cinematic video of Institut Teknologi Sepuluh Nopember Surabaya campus at golden hour, drone shot",
                "duration": 10,
                "ratio": "16:9",
                "task_type": "text_to_video",
            }
        }
    }


class Veo3Payload(BaseModel):
    """Request body untuk POST /api/veo3/generate."""
    prompt: str = Field(
        ..., min_length=5, description="Prompt teks untuk generate Veo 3"
    )
    model: str = Field(
        default="google/veo3.1-lite/text-to-video",
        description="Model Wavespeed/Veo yang akan digunakan"
    )
    reference_images: List[str] = Field(
        default_factory=list, description="URL gambar referensi untuk Veo 3"
    )
    duration: int = Field(
        default=15, ge=1, le=30, description="Durasi video Veo 3 dalam detik"
    )
    ratio: str = Field(
        default="16:9", description="Aspect ratio untuk Veo 3"
    )


class GenerateResponse(BaseModel):
    """Response dari endpoint POST."""
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


@app.post("/generate", response_model=GenerateResponse, status_code=202)
def generate_video(body: GenerateRequest):
    """
    Submit job generate video. Proses inference berjalan asinkron di background.
    """
    logger.info("[API] POST /generate prompt_preview=%.60s", body.prompt)

    try:
        job = service.submit(
            prompt=body.prompt,
            duration=body.duration,
            ratio=body.ratio,
            task_type=body.task_type,
        )
    except ValueError as exc:
        raise HTTPException(status_code=422, detail=str(exc)) from exc

    return GenerateResponse(
        job_id=job.job_id,
        status=job.status.value,
        message="Job diterima dan diproses di background. Gunakan GET /status/{job_id} untuk cek progress.",
    )


@app.post("/api/veo3/generate", response_model=GenerateResponse, status_code=202)
def generate_veo3(body: Veo3Payload):
    """
    Endpoint khusus Veo 3 untuk pemrosesan prompt otomatis dan reference images.
    """
    logger.info("[API] POST /api/veo3/generate model=%s prompt_preview=%.60s", body.model, body.prompt)

    try:
        job = service.submit(
            prompt=body.prompt,
            duration=body.duration,
            ratio=body.ratio,
            task_type="veo3",
            reference_images=body.reference_images,
        )
    except ValueError as exc:
        raise HTTPException(status_code=422, detail=str(exc)) from exc

    return GenerateResponse(
        job_id=job.job_id,
        status=job.status.value,
        message="Veo 3 job diterima dan diproses di background.",
    )


@app.get("/status/{job_id}")
def get_status(job_id: str):
    """Cek status job saat ini."""
    logger.debug("[API] GET /status/%s", job_id)
    try:
        return service.get_status(job_id)
    except KeyError as exc:
        raise HTTPException(status_code=404, detail=f"Job '{job_id}' tidak ditemukan") from exc


@app.get("/result/{job_id}")
def get_result(job_id: str):
    """Ambil hasil video_url dan metadata setelah status done."""
    logger.debug("[API] GET /result/%s", job_id)
    try:
        return service.get_result(job_id)
    except KeyError as exc:
        raise HTTPException(status_code=404, detail=f"Job '{job_id}' tidak ditemukan") from exc


@app.get("/video/{filename}")
def serve_video(filename: str):
    """Serve file video langsung (untuk development/testing)."""
    video_dir = os.getenv("VIDEO_OUTPUT_DIR", "outputs/videos")
    filepath = os.path.join(video_dir, filename)
    if not os.path.exists(filepath):
        raise HTTPException(status_code=404, detail="File video tidak ditemukan")
    return FileResponse(filepath, media_type="video/mp4")


# ---------------------------------------------------------------------------
# Server Runner
# ---------------------------------------------------------------------------
if __name__ == "__main__":
    import uvicorn
    
    host = os.getenv("AI_SERVICE_HOST", "127.0.0.1")
    port = int(os.getenv("AI_SERVICE_PORT", "8000"))
    
    logger.info(f"Starting AI Video Service on {host}:{port}...")
    uvicorn.run(
        "main:app",
        host=host,
        port=port,
        reload=True,  # Hot-reload aktif untuk memudahkan development
        log_level="info",
    )