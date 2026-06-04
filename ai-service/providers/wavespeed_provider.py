import os
import time
import logging
import requests
from typing import Any, Dict, List, Optional

from .base import VideoProvider, VideoRequest, VideoResponse

logger = logging.getLogger(__name__)

WAVESPEED_BASE_URL = "https://api.wavespeed.ai/api/v3"

# Text-to-video model available on Wavespeed (free-tier friendly)
DEFAULT_VIDEO_MODEL = "google/veo3.1-lite/text-to-video"

# MODE=test  → skip Wavespeed, hanya log request (hemat quota)
# MODE=production → kirim ke Wavespeed seperti biasa
_SERVICE_MODE = os.getenv("MODE", "production").strip().lower()


class WavespeedProvider(VideoProvider):
    """
    Wavespeed AI Video Provider.

    Submits a text-to-video job to api.wavespeed.ai, polls for completion,
    and downloads the resulting video file.

    Env vars:
        WAVESPEED_API_KEY  – Bearer token for Wavespeed (required)
        WAVESPEED_MODEL    – override the model slug (optional)
    """

    def __init__(
        self,
        api_key: str = None,
        model: str = None,
        **kwargs,
    ):
        self._api_key = api_key or os.getenv("WAVESPEED_API_KEY", "")
        self._model = model or os.getenv("WAVESPEED_MODEL", DEFAULT_VIDEO_MODEL)

        if not self._api_key:
            raise EnvironmentError(
                "WAVESPEED_API_KEY tidak ditemukan di environment"
            )

    # ------------------------------------------------------------------
    # VideoProvider interface
    # ------------------------------------------------------------------

    @property
    def name(self) -> str:
        return f"wavespeed:{self._model}"

    def generate(self, request: VideoRequest, output_path: str) -> VideoResponse:
        """
        1. Submit job to Wavespeed
        2. Poll until completed or failed
        3. Download video bytes → save to output_path

        Jika MODE=test, request TIDAK dikirim ke Wavespeed.
        Hanya di-log agar bisa di-verifikasi tanpa membuang quota.
        """
        duration = request.constraints.get("duration", 5)
        ctx = request.context or {}

        # ── TEST MODE: skip Wavespeed ──────────────────────────────────
        if _SERVICE_MODE == "test":
            logger.info(
                "[Wavespeed][TEST MODE] Request diterima — TIDAK dikirim ke Wavespeed.\n"
                "  prompt      : %.120s\n"
                "  duration    : %s detik\n"
                "  video_mode  : %s\n"
                "  start_image : %s\n"
                "  end_image   : %s\n"
                "  resolution  : %s\n"
                "  negative    : %s",
                request.input,
                duration,
                ctx.get("video_mode", "text-to-video"),
                ctx.get("start_image", "—"),
                ctx.get("end_image",   "—"),
                ctx.get("resolution",  "—"),
                ctx.get("negative_prompt", "—"),
            )
            return VideoResponse.success(
                self.name,
                {"video_url": "", "video_path": output_path, "test_mode": True},
            )
        # ──────────────────────────────────────────────────────────────

        logger.info(
            "[Wavespeed] Submitting job duration=%s prompt_preview=%.80s",
            duration, request.input,
        )

        # ── 1. Submit ──────────────────────────────────────────────────
        task_id, get_url = self._submit(request, duration)
        if not task_id:
            return VideoResponse.failure(self.name, "Failed to submit job to Wavespeed")

        logger.info("[Wavespeed] Job submitted task_id=%s", task_id)

        # ── 2. Poll ────────────────────────────────────────────────────
        video_url = self._poll(task_id, get_url, max_wait=600)
        if not video_url:
            return VideoResponse.failure(self.name, "Wavespeed job timed-out or failed")

        logger.info("[Wavespeed] Job completed video_url=%s", video_url)

        # ── 3. Download ────────────────────────────────────────────────
        try:
            os.makedirs(os.path.dirname(output_path) or ".", exist_ok=True)
            video_bytes = self._download(video_url)
            with open(output_path, "wb") as f:
                f.write(video_bytes)
            logger.info("[Wavespeed] Video saved to %s (%d bytes)", output_path, len(video_bytes))
            return VideoResponse.success(
                self.name,
                {"video_path": output_path, "video_url": video_url, "size": len(video_bytes)},
            )
        except Exception as exc:
            logger.error("[Wavespeed] Download/save error: %s", exc)
            # Return success with remote URL so callers can still use it
            return VideoResponse.success(
                self.name,
                {"video_url": video_url, "video_path": output_path},
            )

    # ------------------------------------------------------------------
    # Internal helpers
    # ------------------------------------------------------------------

    def _headers(self) -> Dict[str, str]:
        return {
            "Authorization": f"Bearer {self._api_key}",
            "Content-Type": "application/json",
        }

    # Maps resolution string → Wavespeed size parameter
    _SIZE_MAP = {
        "1080p": "1920*1080",
        "720p":  "1280*720",
        "480p":  "832*480",
    }

    def _submit(self, request: VideoRequest, duration: int) -> tuple[Optional[str], Optional[str]]:
        """POST to Wavespeed; return (task_id, get_url)."""
        ctx = request.context or {}
        video_mode      = ctx.get("video_mode", "text-to-video")
        start_image     = ctx.get("start_image", "")
        end_image       = ctx.get("end_image", "")
        negative_prompt = ctx.get("negative_prompt", "")
        generate_audio  = ctx.get("generate_audio", False)
        seed            = ctx.get("seed", -1)
        resolution      = ctx.get("resolution", "480p")

        # Select model slug based on explicit video_mode (not image presence)
        if video_mode == "start-end-to-video":
            model = "google/veo3.1-lite/start-end-to-video"
        elif video_mode == "image-to-video":
            model = "google/veo3.1-lite/image-to-video"
        else:
            model = self._model  # text-to-video

        url = f"{WAVESPEED_BASE_URL}/{model}"

        # Ensure duration is one of [4, 6, 8] to satisfy Veo 3.1 Lite strict constraints
        safe_duration = duration if duration in [4, 6, 8] else 6

        size = self._SIZE_MAP.get(resolution, "832*480")

        payload: dict = {
            "prompt":   request.input,
            "duration": safe_duration,
            "size":     size,
        }

        # Only add image fields when the mode requires them
        if video_mode in ("image-to-video", "start-end-to-video") and start_image:
            payload["image_url"] = start_image
        if video_mode == "start-end-to-video" and end_image:
            payload["last_frame_image_url"] = end_image
        if negative_prompt:
            payload["negative_prompt"] = negative_prompt
        # generate_audio only supported on text-to-video
        if generate_audio and video_mode == "text-to-video":
            payload["generate_audio"] = True
        if seed != -1:
            payload["seed"] = seed

        logger.info("[Wavespeed] model=%s size=%s video_mode=%s",
                    model, size, video_mode)

        try:
            resp = requests.post(url, json=payload, headers=self._headers(), timeout=30)
            data = resp.json()
            logger.debug("[Wavespeed] Submit response: %s", data)

            if data.get("code") != 200:
                logger.error("[Wavespeed] Submit error: %s", data)
                return None, None

            task_data = data.get("data", {})
            task_id = task_data.get("id")
            get_url = task_data.get("urls", {}).get("get")
            return task_id, get_url

        except Exception as exc:
            logger.error("[Wavespeed] Submit exception: %s", exc)
            return None, None

    def _poll(self, task_id: str, get_url: str, max_wait: int = 600) -> Optional[str]:
        """Poll Wavespeed until completed; return video URL or None."""
        poll_url = get_url or f"{WAVESPEED_BASE_URL}/predictions/{task_id}"

        deadline = time.time() + max_wait
        interval = 5  # seconds between polls

        while time.time() < deadline:
            try:
                resp = requests.get(poll_url, headers=self._headers(), timeout=30)
                data = resp.json()

                if data.get("code") != 200:
                    logger.error("[Wavespeed] Poll error: %s", data)
                    return None

                result = data.get("data", {})
                status = result.get("status", "unknown")
                logger.info("[Wavespeed] Poll task_id=%s status=%s", task_id, status)

                if status == "completed":
                    outputs = result.get("outputs", [])
                    if outputs:
                        return outputs[0]
                    return None

                if status == "failed":
                    logger.error("[Wavespeed] Job failed: %s", result.get("error"))
                    return None

                # pending / processing → wait
                time.sleep(interval)

            except Exception as exc:
                logger.warning("[Wavespeed] Poll exception: %s — retrying in %ds", exc, interval)
                time.sleep(interval)

        logger.error("[Wavespeed] Polling timed out after %ds for task_id=%s", max_wait, task_id)
        return None

    def _download(self, video_url: str) -> bytes:
        resp = requests.get(video_url, timeout=120)
        resp.raise_for_status()
        return resp.content
