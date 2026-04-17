"""
Main runner — refactored sesuai ticket PM.

Perubahan dari versi sebelumnya:
- Payload lama { user_query, context } → VideoRequest structured format
- Routing logic ada di backend (router.py), bukan di .env saja
- Response distandarisasi via VideoResponse
- Logging proper (bukan print)
- Error handling per-item tanpa crash keseluruhan run
"""

from __future__ import annotations

import json
import logging
import os
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Dict, List

from dotenv import load_dotenv

from logging_config import setup_logging
from providers import VideoRequest, VideoResponse, build_provider, enrich_routing

logger = logging.getLogger(__name__)


# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------

def load_prompts(path: str) -> List[Dict[str, Any]]:
    """Load prompt list dari JSON file."""
    with open(path, "r", encoding="utf-8") as f:
        data = json.load(f)
    logger.info("Loaded %d prompts from %s", len(data), path)
    return data


def build_video_request(item: Dict[str, Any]) -> VideoRequest:
    """
    Ticket PM — Payload Mapping:
    Konversi item dari prompts.json ke VideoRequest (structured format).

    Input (format lama / prompts.json):
        { "id": ..., "prompt": ..., "duration": ..., "ratio": ... }

    Output (format baru):
        VideoRequest(
            instruction="text_to_video",
            input=prompt,
            context={ "prompt_id": ... },
            constraints={ "duration": ..., "ratio": ..., "resolution": ..., "fps": ... },
            routing={ "task_type": "text_to_video", "fallback": "ltx:ltx-2-fast" }
        )
    """
    task_type = item.get("task_type", "text_to_video")

    return VideoRequest(
        instruction=task_type,
        input=item["prompt"],
        context={
            "prompt_id": item.get("id", "unknown"),
            "source": "prompts.json",
        },
        constraints={
            "duration": int(item.get("duration", os.getenv("ACTIVE_DURATION", "10"))),
            "ratio": item.get("ratio", os.getenv("ACTIVE_RATIO", "16:9")),
            "resolution": item.get("resolution", os.getenv("LTX_RESOLUTION", "1920x1080")),
            "fps": int(item.get("fps", os.getenv("LTX_FPS", "25"))),
        },
        routing={
            "task_type": task_type,
            "fallback": "ltx:ltx-2-3-fast",
        },
    )


def make_output_path(video_dir: str, prompt_id: str, provider_name: str) -> str:
    """Buat path output yang aman (menghindari karakter spesial)."""
    safe_provider = provider_name.replace(":", "_")
    filename = f"{prompt_id}__{safe_provider}.mp4"
    return os.path.join(video_dir, filename)


def save_report(report_dir: str, results: List[Dict[str, Any]]) -> str:
    """Simpan report JSON dengan timestamp."""
    Path(report_dir).mkdir(parents=True, exist_ok=True)
    timestamp = datetime.now(timezone.utc).strftime("%Y%m%d_%H%M%S")
    report_path = os.path.join(report_dir, f"trial_report_{timestamp}.json")
    with open(report_path, "w", encoding="utf-8") as f:
        json.dump(results, f, indent=2, ensure_ascii=False)
    return report_path


# ---------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------

def main() -> None:
    load_dotenv()

    # Setup directories
    video_dir = os.getenv("VIDEO_OUTPUT_DIR", "outputs/videos")
    report_dir = os.getenv("REPORT_DIR", "outputs/reports")
    log_dir = os.getenv("LOG_DIR", "outputs/logs")
    log_level = os.getenv("LOG_LEVEL", "INFO")

    Path(video_dir).mkdir(parents=True, exist_ok=True)
    Path(report_dir).mkdir(parents=True, exist_ok=True)

    setup_logging(log_dir=log_dir, log_level=log_level)

    logger.info("=" * 60)
    logger.info("Video Generator — AI Marketing ITS")
    logger.info("=" * 60)

    # Load prompts
    prompts = load_prompts("prompts.json")

    results: List[Dict[str, Any]] = []
    success_count = 0
    fail_count = 0

    for item in prompts:
        prompt_id = item.get("id", "unknown")
        logger.info("-" * 60)
        logger.info("Processing prompt_id=%s", prompt_id)

        # Transform ke structured VideoRequest
        request = build_video_request(item)

        # Enrich routing di backend layer (ticket PM)
        request = enrich_routing(request)

        logger.info(
            "Routing: task_type=%s → provider=%s model=%s",
            request.routing.get("task_type"),
            request.routing.get("resolved_provider"),
            request.routing.get("resolved_model"),
        )

        provider_name = request.routing.get("resolved_provider")
        model_name = request.routing.get("resolved_model")

        try:
            provider = build_provider(provider_name, model_name)
        except (EnvironmentError, ValueError) as e:
            logger.error("Gagal inisialisasi provider untuk prompt_id=%s: %s", prompt_id, e)
            fail_count += 1
            continue

        output_path = make_output_path(video_dir, prompt_id, provider.name)

        # Generate
        response: VideoResponse = provider.generate(request=request, output_path=output_path)

        # Ticket PM — Standardized response logging
        logger.info(
            "Result: status=%s route_used=%s",
            response.status, response.route_used
        )
        if response.status == "error":
            logger.error("FAILED prompt_id=%s error=%s", prompt_id, response.error)
            fail_count += 1
        else:
            logger.info("SUCCESS prompt_id=%s output=%s", prompt_id, response.result.get("output_path"))
            if response.result.get("audio") is False:
                logger.warning(
                    "[AUDIO] prompt_id=%s — Video tidak memiliki audio. %s",
                    prompt_id,
                    response.result.get("audio_note", "")
                )
            success_count += 1

        # Build report entry sesuai standardized response
        report_entry: Dict[str, Any] = {
            **response.to_dict(),
            "prompt_id": prompt_id,
            "prompt_preview": request.input[:100] + "...",
            "request_payload": {
                "instruction": request.instruction,
                "constraints": request.constraints,
                "routing": request.routing,
            },
            "created_at": datetime.now(timezone.utc).isoformat(),
        }
        results.append(report_entry)

    # Simpan report
    report_path = save_report(report_dir, results)

    logger.info("=" * 60)
    logger.info("DONE — success=%d failed=%d total=%d", success_count, fail_count, len(prompts))
    logger.info("Report: %s", report_path)
    logger.info("=" * 60)


if __name__ == "__main__":
    main()