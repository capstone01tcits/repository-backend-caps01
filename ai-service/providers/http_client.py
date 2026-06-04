"""
Shared HTTP client dengan retry mechanism dan timeout.
Ticket PM: Client Setup — timeout dan retry diimplementasikan.
Untuk memastikan backend stabil jika AI API provider delay atau fail.
"""

from __future__ import annotations

import logging
import time
from typing import Any, Dict, NoReturn, Optional, Tuple, Union

import requests
from requests.adapters import HTTPAdapter
from urllib3.util.retry import Retry

logger = logging.getLogger(__name__)


def build_http_session(
    total_retries: int = 3,
    backoff_factor: float = 1.5,
    # FIX 1: tuple → Tuple[int, ...] — type hint lebih spesifik
    status_forcelist: Tuple[int, ...] = (429, 500, 502, 503, 504),
) -> requests.Session:
    """
    Buat session requests dengan retry logic otomatis.

    - Retry 3x untuk status code yang ditetapkan
    - Exponential backoff: 1.5s, 3s, 6s
    - Tidak retry untuk 4xx client error kecuali 429
    """
    session = requests.Session()
    retry_strategy = Retry(
        total=total_retries,
        backoff_factor=backoff_factor,
        status_forcelist=status_forcelist,
        allowed_methods=["POST", "GET"],
        raise_on_status=False,
    )
    adapter = HTTPAdapter(max_retries=retry_strategy)
    session.mount("https://", adapter)
    session.mount("http://", adapter)
    return session


class APIError(Exception):
    """Custom exception untuk API error dengan detail response."""

    def __init__(
        self,
        message: str,
        status_code: Optional[int] = None,
        body: Optional[str] = None,
    ):
        super().__init__(message)
        self.status_code = status_code
        self.body = body

    def __str__(self) -> str:
        base = super().__str__()
        if self.status_code:
            return f"{base} [HTTP {self.status_code}]"
        return base


def safe_post(
    session: requests.Session,
    url: str,
    headers: Dict[str, str],
    payload: Dict[str, Any],
    timeout: int = 300,
) -> requests.Response:
    """POST request dengan error handling terpusat."""
    logger.debug("POST %s payload_keys=%s", url, list(payload.keys()))
    
    try:
        resp = session.post(url, json=payload, headers=headers, timeout=timeout)
        
        # --- PERBAIKAN DI SINI: Cek status code ---
        if not resp.ok:
            body_preview = resp.text[:500] if resp.text else "(empty)"
            logger.error("API error url=%s status=%d body=%s", url, resp.status_code, body_preview)
            # Kita raise APIError supaya ditangkap oleh AIVideoService
            raise APIError(
                f"API returned error {resp.status_code} untuk {url}",
                status_code=resp.status_code,
                body=body_preview,
            )
        # --- ---------------------------------- ---

        return resp # Pastikan return berada di luar blok pengecekan error

    except requests.exceptions.Timeout as exc:
        raise APIError(f"Request timeout setelah {timeout}s ke {url}") from exc
    except requests.exceptions.ConnectionError as exc:
        raise APIError(f"Connection error ke {url}: {exc}") from exc


def safe_get(
    session: requests.Session,
    url: str,
    headers: Dict[str, str],
    timeout: int = 30,
) -> requests.Response:
    """GET request dengan error handling terpusat."""
    try:
        # 1. Jalankan request
        resp = session.get(url, headers=headers, timeout=timeout)
        
        # 2. Cek status (Ini di dalam try)
        if not resp.ok:
            raise APIError(
                f"GET error {resp.status_code} untuk {url}",
                status_code=resp.status_code,
                body=resp.text[:500],
            )
        return resp

    # 3. Except HARUS sejajar dengan 'try' di atas
    except requests.exceptions.Timeout as exc:
        raise APIError(f"GET timeout setelah {timeout}s ke {url}") from exc
    except requests.exceptions.ConnectionError as exc:
        raise APIError(f"Connection error ke {url}: {exc}") from exc


def poll_until_done(
    session: requests.Session,
    status_url: str,
    headers: Dict[str, str],
    success_status: str = "SUCCEEDED",
    # FIX 2: tuple → Tuple[str, ...] — type hint lebih spesifik
    fail_statuses: Tuple[str, ...] = ("FAILED", "CANCELLED"),
    max_attempts: int = 180,
    poll_interval: int = 5,
    provider_tag: str = "PROVIDER",
    # FIX 3: return type Union[Dict, NoReturn] — karena fungsi bisa raise TimeoutError
) -> Union[Dict[str, Any], NoReturn]:
    """
    Generic polling loop untuk async task.
    Digunakan oleh provider async.

    Returns:
        Dict response dari API jika success.
    Raises:
        APIError: jika task failed/cancelled.
        TimeoutError: jika max_attempts tercapai.
    """
    for attempt in range(1, max_attempts + 1):
        logger.debug("[%s] POLL attempt=%d url=%s", provider_tag, attempt, status_url)
        try:
            resp = safe_get(session, status_url, headers)
            data: Dict[str, Any] = resp.json()
        except APIError as e:
            logger.warning(
                "[%s] POLL network error attempt=%d: %s",
                provider_tag, attempt, e,
            )
            time.sleep(poll_interval * 2)
            continue

        status = data.get("status", "")
        logger.info("[%s] POLL attempt=%d status=%s", provider_tag, attempt, status)

        if status == success_status:
            return data

        if status in fail_statuses:
            raise APIError(
                f"[{provider_tag}] Task ended with status={status}: {data}"
            )

        time.sleep(poll_interval)

    # FIX 3 lanjutan: eksplisit raise dengan pesan yang jelas
    raise TimeoutError(
        f"[{provider_tag}] Task timed out setelah {max_attempts * poll_interval}s "
        f"({max_attempts} attempts x {poll_interval}s interval)"
    )