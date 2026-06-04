"""
Base abstract class for all video generation providers.
Defines the contract that every provider must fulfill.
"""

from __future__ import annotations

from abc import ABC, abstractmethod
from dataclasses import dataclass, field
from typing import Any, Dict, Optional


# ---------------------------------------------------------------------------
# Structured payload (ticket PM: payload mapping)
# ---------------------------------------------------------------------------

@dataclass
class VideoRequest:
    """
    Structured payload sesuai spesifikasi ticket PM.
    Menggantikan pendekatan prompt string tunggal.
    """
    instruction: str           # Jenis task: text_to_video, image_to_video, dll
    input: str                 # Raw user_query / prompt
    context: Dict[str, Any] = field(default_factory=dict)
    constraints: Dict[str, Any] = field(default_factory=dict)
    routing: Dict[str, Any] = field(default_factory=dict)

    @classmethod
    def from_legacy(cls, user_query: str, context: Optional[Dict] = None) -> "VideoRequest":
        """
        Backward-compatible factory: konversi format lama ke format baru.
        
        Format lama: { "user_query": "string", "context": {} }
        Format baru: { "instruction": ..., "input": ..., "context": ..., ... }
        """
        ctx = context or {}
        task_type = ctx.pop("task_type", "text_to_video")

        return cls(
            instruction=task_type,
            input=user_query,
            context=ctx,
            constraints={
                "duration": ctx.pop("duration", 10),
                "ratio": ctx.pop("ratio", "16:9"),
                "resolution": ctx.pop("resolution", "1920x1080"),
                "fps": ctx.pop("fps", 25),
            },
            routing={
                "task_type": task_type,
                "fallback": "wavespeed:wavespeed-ai/wan-2.1/t2v-480p",
            },
        )


# ---------------------------------------------------------------------------
# Standardized response (ticket PM: response handling)
# ---------------------------------------------------------------------------

@dataclass
class VideoResponse:
    """
    Standardized response format sesuai spesifikasi ticket PM.
    Semua provider wajib return format ini.

    {
        "status": "success" | "error",
        "route_used": "provider:model",
        "result": { ... }
    }
    """
    status: str                        # "success" | "error"
    route_used: str                    # e.g. "wavespeed:wavespeed-ai/wan-2.1/t2v-480p"
    result: Dict[str, Any] = field(default_factory=dict)
    error: Optional[str] = None

    def to_dict(self) -> Dict[str, Any]:
        base: Dict[str, Any] = {
            "status": self.status,
            "route_used": self.route_used,
            "result": self.result,
        }
        if self.error:
            base["error"] = self.error
        return base

    @classmethod
    def success(cls, route_used: str, result: Dict[str, Any]) -> "VideoResponse":
        return cls(status="success", route_used=route_used, result=result)

    @classmethod
    def failure(cls, route_used: str, error: str) -> "VideoResponse":
        return cls(status="error", route_used=route_used, result={}, error=error)


# ---------------------------------------------------------------------------
# Abstract provider
# ---------------------------------------------------------------------------

class VideoProvider(ABC):
    """Abstract base class for all video generation providers."""

    @property
    @abstractmethod
    def name(self) -> str:
        """Return provider identifier, e.g. 'wavespeed:wavespeed-ai/wan-2.1/t2v-480p'."""

    @abstractmethod
    def generate(self, request: VideoRequest, output_path: str) -> VideoResponse:
        """
        Main generation method.
        Semua provider wajib implement ini dan return VideoResponse.
        """