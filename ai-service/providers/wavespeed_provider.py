from __future__ import annotations

import logging
import os
import time
from typing import Any, Dict, Optional, Tuple

import requests

from .base import VideoProvider, VideoRequest, VideoResponse
from .http_client import APIError, build_http_session, safe_get, safe_post

logger = logging.getLogger(__name__)

WAVESPEED_BASE_URL = "https://api.wavespeed.ai/api/v3"
DEFAULT_VIDEO_MODEL = "google/veo3.1-lite/text-to-video"
VALID_WAVESPEED_DURATIONS = [4, 6, 8]
DEFAULT_WAVESPEED_SIZE = "832*480"


class WavespeedProvider(VideoProvider):
    def __init__(
        self,
        api_key: str,
        model: str = DEFAULT_VIDEO_MODEL,
        request_timeout: int = 300,
        poll_interval: int = 5,
        poll_max_attempts: int = 120,
        size: str = DEFAULT_WAVESPEED_SIZE,
    ):
        if not api_key or not api_key.strip():
            raise ValueError("WAVESPEED_API_KEY tidak boleh kosong")

        self._api_key = api_key.strip()
        self._model = model.strip()
        self._timeout = request_timeout
        self._poll_interval = poll_interval
        self._poll_max_attempts = poll_max_attempts
        self._size = size
        self._session = build_http_session()

    @property
    def name(self) -> str:
        return f"wavespeed:{self._model}"

    def generate(self, request: VideoRequest, output_path: str) -> VideoResponse:
        prompt = request.input
        duration = int(request.constraints.get("duration", 5) or 5)

        logger.info(
            "[Wavespeed] Submitting job model=%s duration=%s prompt_preview=%.80s",
            self._model,
            duration,
            prompt,
        )

        task_id, status_url = self._submit(prompt, duration)
        if not task_id or not status_url:
            return VideoResponse.failure(self.name, "Failed to submit Wavespeed job")

        logger.info("[Wavespeed] Job submitted task_id=%s status_url=%s", task_id, status_url)

        video_url = self._poll(task_id, status_url)
        if not video_url:
            return VideoResponse.failure(self.name, "Wavespeed job timed-out or failed")

        logger.info("[Wavespeed] Job completed video_url=%s", video_url)

        try:
            os.makedirs(os.path.dirname(output_path) or ".", exist_ok=True)
            video_bytes = self._download(video_url)
            with open(output_path, "wb") as f:
                f.write(video_bytes)

            logger.info("[Wavespeed] Video saved to %s (%d bytes)", output_path, len(video_bytes))
            return VideoResponse.success(
                self.name,
                {
                    "output_path": output_path,
                    "video_url": video_url,
                    "duration": duration,
                    "size_bytes": len(video_bytes),
                },
            )
        except Exception as exc:
            logger.error("[Wavespeed] Download/save error: %s", exc)
            return VideoResponse.failure(self.name, f"Failed to save Wavespeed video: {exc}")

    def _headers(self) -> Dict[str, str]:
        return {
            "Authorization": f"Bearer {self._api_key}",
            "Content-Type": "application/json",
        }

    def _submit(self, prompt: str, duration: int) -> Tuple[Optional[str], Optional[str]]:
        url = f"{WAVESPEED_BASE_URL}/{self._model}"

        if duration <= 0:
            duration = 5
        
        # Logika pembulatan otomatis ke durasi yang diizinkan (4, 6, 8)
        safe_duration = min(VALID_WAVESPEED_DURATIONS, key=lambda x: abs(x - duration))

        payload = {
            "prompt": prompt,
            "duration": safe_duration,
            "size": self._size,
        }

        try:
            resp = safe_post(
                session=self._session,
                url=url,
                headers=self._headers(),
                payload=payload,
                timeout=self._timeout,
            )
            data = resp.json()

            task_data = data.get("data", {})
            task_id = task_data.get("id")
            status_url = task_data.get("urls", {}).get("get")

            if not task_id or not status_url:
                logger.error("[Wavespeed] Submit response incomplete: %s", data)
                return None, None

            return task_id, status_url
        except APIError as exc:
            logger.error("[Wavespeed] Submit API error: %s", exc)
            return None, None
        except ValueError as exc:
            logger.error("[Wavespeed] Submit parse error: %s", exc)
            return None, None

    def _poll(self, task_id: str, status_url: str) -> Optional[str]:
        deadline = time.time() + (self._poll_max_attempts * self._poll_interval)

        while time.time() < deadline:
            try:
                resp = self._session.get(status_url, headers=self._headers(), timeout=30)
            except requests.RequestException as exc:
                logger.warning("[Wavespeed] Poll request failed: %s", exc)
                time.sleep(self._poll_interval)
                continue

            if resp.status_code != 200:
                logger.warning("[Wavespeed] Poll HTTP %d for url=%s", resp.status_code, status_url)
                time.sleep(self._poll_interval)
                continue

            try:
                payload = resp.json()
            except ValueError as exc:
                logger.error("[Wavespeed] Poll response not JSON: %s", exc)
                time.sleep(self._poll_interval)
                continue

            status = self._extract_status(payload)
            logger.info("[Wavespeed] Poll status=%s task_id=%s", status, task_id)

            if status in ("success", "succeeded", "completed", "finished"):
                return self._extract_video_url(payload)

            if status in ("failed", "error", "cancelled"):
                logger.error("[Wavespeed] Task failed: %s", payload)
                return None

            time.sleep(self._poll_interval)

        logger.error("[Wavespeed] Poll timed out after %d seconds", self._poll_max_attempts * self._poll_interval)
        return None

    def _download(self, url: str) -> bytes:
        resp = safe_get(session=self._session, url=url, headers={}, timeout=300)
        return resp.content

    def _extract_status(self, payload: Dict[str, Any]) -> str:
        if not isinstance(payload, dict):
            return ""
        data = payload.get("data", payload)
        if isinstance(data, dict):
            return str(data.get("status", payload.get("status", ""))).lower()
        return str(payload.get("status", "")).lower()

    def _extract_video_url(self, payload: Dict[str, Any]) -> Optional[str]:
        if not isinstance(payload, dict):
            return None
        data = payload.get("data", payload)
        if not isinstance(data, dict):
            return None

        urls = data.get("urls", {}) or {}
        return (
            urls.get("get")
            or urls.get("video")
            or data.get("video_url")
            or data.get("url")
        )