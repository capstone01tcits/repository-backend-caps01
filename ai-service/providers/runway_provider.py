"""
Runway Video Provider — refactored sesuai ticket PM.
Support: Gen-4.5 text-to-video via async polling.
"""

from __future__ import annotations

import logging
import os
from typing import Any, Dict

from .base import VideoProvider, VideoRequest, VideoResponse
from .http_client import APIError, build_http_session, poll_until_done, safe_get, safe_post

logger = logging.getLogger(__name__)

# ---------------------------------------------------------------------------
# Konstanta
# ---------------------------------------------------------------------------

RUNWAY_BASE_URL = "https://api.dev.runwayml.com/v1"
RUNWAY_API_VERSION = "2024-11-06"

SUPPORTED_MODELS = {"gen4.5", "gen3a_turbo"}

RATIO_TO_RESOLUTION = {
    "16:9": "1280:720",
    "9:16": "720:1280",
    "1:1":  "960:960",
    "4:3":  "1024:768",
}


# ---------------------------------------------------------------------------
# Provider
# ---------------------------------------------------------------------------

class RunwayProvider(VideoProvider):
    """
    Runway ML video generation provider.

    Ticket PM — LTX Client Setup (berlaku sama untuk Runway):
    - HTTP session dengan retry digunakan
    - API key divalidasi saat init
    - Payload ditransformasi dari VideoRequest ke Runway format
    - Response distandarisasi ke VideoResponse
    """

    def __init__(
        self,
        api_key: str,
        model: str = "gen4.5",
        api_version: str = RUNWAY_API_VERSION,
        poll_max_attempts: int = 180,
        poll_interval: int = 5,
    ):
        if not api_key or not api_key.strip():
            raise ValueError("RUNWAY_API_KEY tidak boleh kosong")
        if model not in SUPPORTED_MODELS:
            logger.warning("Model '%s' tidak dikenali. Supported: %s", model, SUPPORTED_MODELS)

        self._api_key = api_key
        self._model = model
        self._api_version = api_version
        self._poll_max_attempts = poll_max_attempts
        self._poll_interval = poll_interval
        self._session = build_http_session()

    # ------------------------------------------------------------------
    # Contract implementation
    # ------------------------------------------------------------------

    @property
    def name(self) -> str:
        return f"runway:{self._model}"

    def generate(self, request: VideoRequest, output_path: str) -> VideoResponse:
        """
        Kirim task ke Runway, polling sampai selesai, download video.

        Ticket PM — Payload Mapping:
          input       → promptText
          constraints → ratio/duration
          routing     → menentukan endpoint yang dipakai
        """
        logger.info("[RUNWAY] START route=%s output=%s", self.name, output_path)

        headers = self._build_headers()
        payload = self._build_payload(request)

        # Tentukan endpoint berdasarkan model
        endpoint = self._resolve_endpoint(request)
        create_url = f"{RUNWAY_BASE_URL}/{endpoint}"

        logger.info("[RUNWAY] POST %s model=%s", create_url, self._model)

        try:
            create_resp = safe_post(
                session=self._session,
                url=create_url,
                headers=headers,
                payload=payload,
                timeout=60,
            )
        except APIError as e:
            logger.error("[RUNWAY] Create task failed: %s", e)
            return VideoResponse.failure(route_used=self.name, error=str(e))

        task = create_resp.json()
        task_id = task.get("id")
        if not task_id:
            return VideoResponse.failure(
                route_used=self.name,
                error=f"Runway tidak return task ID. Response: {task}"
            )

        logger.info("[RUNWAY] Task created id=%s", task_id)

        # Polling
        status_url = f"{RUNWAY_BASE_URL}/tasks/{task_id}"
        try:
            data = poll_until_done(
                session=self._session,
                status_url=status_url,
                headers=headers,
                success_status="SUCCEEDED",
                fail_statuses=("FAILED", "CANCELLED"),
                max_attempts=self._poll_max_attempts,
                poll_interval=self._poll_interval,
                provider_tag="RUNWAY",
            )
        except (APIError, TimeoutError) as e:
            logger.error("[RUNWAY] Polling failed: %s", e)
            return VideoResponse.failure(route_used=self.name, error=str(e))

        # Download
        output_urls = data.get("output", [])
        if not output_urls:
            return VideoResponse.failure(
                route_used=self.name,
                error="Task SUCCEEDED tapi tidak ada output URL."
            )

        video_url = output_urls[0]
        logger.info("[RUNWAY] Downloading from url=%s", video_url[:80])

        try:
            video_resp = safe_get(
                session=self._session,
                url=video_url,
                headers={},    # Download URL biasanya tidak perlu auth
                timeout=300,
            )
            os.makedirs(os.path.dirname(output_path), exist_ok=True)
            with open(output_path, "wb") as f:
                f.write(video_resp.content)
            logger.info("[RUNWAY] Saved path=%s size=%d bytes", output_path, len(video_resp.content))
        except (APIError, OSError) as e:
            return VideoResponse.failure(route_used=self.name, error=f"Download/save gagal: {e}")

        return VideoResponse.success(
            route_used=self.name,
            result={
                "output_path": output_path,
                "model": self._model,
                "task_id": task_id,
                "duration": request.constraints.get("duration"),
                "ratio": request.constraints.get("ratio"),
                "audio": False,  # Runway gen4.5 base juga tidak include audio
                "file_size_bytes": len(video_resp.content),
            },
        )

    # ------------------------------------------------------------------
    # Private helpers
    # ------------------------------------------------------------------

    def _build_headers(self) -> Dict[str, str]:
        return {
            "Authorization": f"Bearer {self._api_key}",
            "Content-Type": "application/json",
            "X-Runway-Version": self._api_version,
        }

    def _build_payload(self, request: VideoRequest) -> Dict[str, Any]:
        """Transform VideoRequest constraints ke Runway payload format."""
        constraints = request.constraints
        ratio = constraints.get("ratio", "16:9")
        resolution = RATIO_TO_RESOLUTION.get(ratio, "1280:720")
        duration = int(constraints.get("duration", 10))

        return {
            "model": self._model,
            "promptText": request.input,
            "ratio": resolution,
            "duration": duration,
        }

    def _resolve_endpoint(self, request: VideoRequest) -> str:
        """
        Ticket PM — Routing Logic Integration:
        Tentukan endpoint Runway berdasarkan task_type dan model.
        """
        task_type = request.routing.get("task_type", "text_to_video")

        if task_type == "image_to_video":
            return "image_to_video"
        if self._model == "gen4.5":
            # Gen-4.5 pure text-to-video lebih stabil via image_to_video tanpa promptImage
            return "image_to_video"
        return "text_to_video"