"""
Job Model & In-Memory Storage.

Representasi struktur DB yang disepakati dengan divisi backend:

Table: video_jobs
  id         → UUID, primary key
  prompt     → string, prompt yang dikirim
  job_id     → string, sama dengan id (alias untuk kejelasan)
  status     → enum: pending | processing | done | failed
  video_url  → string nullable, diisi setelah done
  created_at → datetime UTC
  updated_at → datetime UTC

Untuk production, ganti InMemoryJobStore dengan PostgreSQL / Redis.
"""

from __future__ import annotations

import uuid
from dataclasses import dataclass, field
from datetime import datetime, timezone
from enum import Enum
from typing import Dict, Optional


class JobStatus(str, Enum):
    PENDING    = "pending"
    PROCESSING = "processing"
    DONE       = "done"
    FAILED     = "failed"


@dataclass
class VideoJob:
    """Satu record job — merepresentasikan satu baris di DB."""
    id: str
    prompt: str
    status: JobStatus = JobStatus.PENDING
    video_url: Optional[str] = None
    thumbnail_url: Optional[str] = None
    error: Optional[str] = None
    created_at: str = field(default_factory=lambda: datetime.now(timezone.utc).isoformat())
    updated_at: str = field(default_factory=lambda: datetime.now(timezone.utc).isoformat())

    # metadata tambahan (tidak masuk DB utama, untuk debugging)
    provider: Optional[str] = None
    model: Optional[str] = None
    duration: int = 10
    ratio: str = "16:9"

    @property
    def job_id(self) -> str:
        """Alias untuk id — sesuai naming convention backend."""
        return self.id

    def to_dict(self) -> Dict:
        """Serialisasi ke dict — format yang dikirim ke backend / response API."""
        return {
            "id": self.id,
            "job_id": self.job_id,
            "prompt": self.prompt,
            "status": self.status.value,
            "video_url": self.video_url,
            "thumbnail_url": self.thumbnail_url,
            "error": self.error,
            "created_at": self.created_at,
            "updated_at": self.updated_at,
            "meta": {
                "provider": self.provider,
                "model": self.model,
                "duration": self.duration,
                "ratio": self.ratio,
            }
        }

    def _touch(self):
        self.updated_at = datetime.now(timezone.utc).isoformat()

    def mark_processing(self, provider: str, model: str):
        self.status = JobStatus.PROCESSING
        self.provider = provider
        self.model = model
        self._touch()

    def mark_done(self, video_url: str, thumbnail_url: Optional[str] = None):
        self.status = JobStatus.DONE
        self.video_url = video_url
        self.thumbnail_url = thumbnail_url
        self._touch()

    def mark_failed(self, error: str):
        self.status = JobStatus.FAILED
        self.error = error
        self._touch()


def new_job_id() -> str:
    return str(uuid.uuid4())


# ---------------------------------------------------------------------------
# In-Memory Job Store (ganti dengan DB di production)
# ---------------------------------------------------------------------------

class InMemoryJobStore:
    """
    Simple in-memory store untuk development/testing.
    
    Interface ini SAMA dengan yang akan diimplementasi backend
    menggunakan PostgreSQL / Redis — sehingga swapnya mudah.
    """

    def __init__(self):
        self._store: Dict[str, VideoJob] = {}

    def save(self, job: VideoJob) -> VideoJob:
        self._store[job.id] = job
        return job

    def get(self, job_id: str) -> Optional[VideoJob]:
        return self._store.get(job_id)

    def update(self, job: VideoJob) -> VideoJob:
        self._store[job.id] = job
        return job

    def all(self) -> Dict[str, VideoJob]:
        return dict(self._store)


# Singleton store — di production diganti dengan DB session
job_store = InMemoryJobStore()