import os
import logging
import time
from typing import Any, Dict, List, Optional
from .base import VideoProvider, VideoRequest, VideoResponse
from .http_client import build_http_session, safe_post, APIError

logger = logging.getLogger(__name__)

class WavespeedProvider(VideoProvider):
    """
    Wavespeed AI Video Aggregator Provider.
    
    Provider ini mengirimkan payload terpadu (Full Storyboard) ke API Wavespeed
    yang mendukung model seperti 'veo3'. Wavespeed menangani parsing dan stitching secara internal.
    """

    def __init__(
        self, 
        api_key: str = None, 
        api_url: str = None,
        model: str = "veo3"
    ):
        self._api_key = api_key or os.getenv("WAVESPEED_API_KEY", "")
        self._api_url = api_url or os.getenv("WAVESPEED_API_URL", "https://api.wavespeed.ai/v1/generate")
        self._model = model
        self._session = build_http_session()

    @property
    def name(self) -> str:
        return f"wavespeed:{self._model}"

    def generate(self, request: VideoRequest, output_path: str) -> VideoResponse:
        """
        Kirim payload ke Wavespeed API.
        """
        logger.info(f"[Wavespeed] Mengirim request ke {self._api_url} menggunakan model {self._model}")
        
        # Build payload sesuai format yang diminta user
        payload = {
            "model": self._model,
            "prompt": request.input,
            "reference_images": request.context.get("reference_images", [])
        }

        headers = {
            "Authorization": f"Bearer {self._api_key}",
            "Content-Type": "application/json"
        }

        try:
            # Kirim request
            resp = safe_post(
                session=self._session,
                url=self._api_url,
                headers=headers,
                payload=payload,
                timeout=600 # Video generation bisa lama
            )
            
            # Asumsi response berisi binary video atau JSON dengan URL
            # Jika Wavespeed mengembalikan binary:
            os.makedirs(os.path.dirname(output_path), exist_ok=True)
            with open(output_path, "wb") as f:
                f.write(resp.content)
            
            logger.info(f"[Wavespeed] Video berhasil diterima dan disimpan di {output_path}")
            return VideoResponse.success(self.name, {"video_path": output_path, "size": len(resp.content)})

        except APIError as e:
            logger.error(f"[Wavespeed] API Error: {e}")
            return VideoResponse.failure(self.name, f"Wavespeed API Error: {str(e)}")
        except Exception as e:
            logger.error(f"[Wavespeed] Unexpected Error: {e}")
            return VideoResponse.failure(self.name, f"Wavespeed Unexpected Error: {str(e)}")
