import os
import logging
import subprocess
import time
import re
from typing import List, Dict, Any
from .base import VideoProvider, VideoRequest, VideoResponse

logger = logging.getLogger(__name__)

class Veo3Provider(VideoProvider):
    """
    Veo 3 Provider — Automated Storyboard Processor & Video Stitcher.
    
    Provider ini menerima prompt lengkap (Hook, Value, CTA),
    memprosesnya per scene, dan menggabungkannya menggunakan FFmpeg.
    """

    def __init__(self, model: str = "veo3"):
        self._model = model
        self._output_dir = os.getenv("VIDEO_OUTPUT_DIR", "outputs/videos")
        os.makedirs(self._output_dir, exist_ok=True)

    @property
    def name(self) -> str:
        return f"veo3:{self._model}"

    def generate(self, request: VideoRequest, output_path: str) -> VideoResponse:
        """
        Main entry point untuk pemrosesan Veo 3.
        """
        prompt = request.input
        logger.info(f"[Veo3] Memulai proses generate untuk job: {request.context.get('job_id')}")
        
        # 1. Parse prompt ke dalam list of scenes
        scenes = self._parse_scenes(prompt)
        
        if not scenes:
            logger.error("[Veo3] Gagal memparsing prompt. Pastikan format SCENE 1 (Xs-Ys) terpenuhi.")
            # Fallback: jika gagal parse, anggap seluruh prompt sebagai satu scene
            scenes = [{"type": "FULL", "description": prompt, "duration": request.constraints.get("duration", 10)}]

        logger.info(f"[Veo3] Terdeteksi {len(scenes)} scene. Menjalankan generation per scene...")

        # 2. Generate video per scene
        scene_files = []
        try:
            for i, scene in enumerate(scenes):
                scene_filename = f"temp_{request.context.get('job_id')}_scene_{i+1}.mp4"
                scene_path = os.path.join(self._output_dir, scene_filename)
                
                logger.info(f"[Veo3] Generating Scene {i+1}: {scene['type']} - {scene['description'][:50]}...")
                
                # Di implementasi nyata, ini akan memanggil LTX/Runway
                # Untuk tugas ini, kita gunakan FFmpeg generator untuk mensimulasikan clip
                self._generate_clip_simulation(scene, scene_path)
                scene_files.append(scene_path)

            # 3. Gabungkan semua clip menggunakan FFmpeg
            logger.info(f"[Veo3] Menggabungkan {len(scene_files)} clip menjadi satu video final...")
            self._stitch_videos(scene_files, output_path)
            
            # Cleanup temp files
            for f in scene_files:
                if os.path.exists(f):
                    os.remove(f)

            logger.info(f"[Veo3] Sukses! Video final disimpan di: {output_path}")
            return VideoResponse.success(self.name, {"video_path": output_path, "scenes_processed": len(scenes)})

        except Exception as e:
            logger.error(f"[Veo3] Error dalam pipeline: {e}")
            # Cleanup temp files on error
            for f in scene_files:
                if os.path.exists(f):
                    os.remove(f)
            return VideoResponse.failure(self.name, f"Veo3 Pipeline Error: {str(e)}")

    def _parse_scenes(self, prompt: str) -> List[Dict[str, Any]]:
        """
        Parsing prompt format: SCENE 1 (0–5s): HOOK [Description]
        """
        # Regex yang fleksibel untuk menangkap format dari Go worker
        pattern = r"SCENE (\d+) \((\d+)–(\d+)s\): (\w+) \[(.*?)\]"
        matches = re.findall(pattern, prompt, re.DOTALL)
        
        scenes = []
        for match in matches:
            scenes.append({
                "index": int(match[0]),
                "start": int(match[1]),
                "end": int(match[2]),
                "duration": int(match[2]) - int(match[1]),
                "type": match[3],
                "description": match[4].strip()
            })
        
        # Sort by index
        scenes.sort(key=lambda x: x["index"])
        return scenes

    def _generate_clip_simulation(self, scene: Dict, output_path: str):
        """
        Simulasi pembuatan clip video menggunakan FFmpeg.
        Menambahkan teks overlay sesuai deskripsi scene.
        """
        duration = scene.get("duration", 5)
        description = scene.get("description", "No description")
        scene_type = scene.get("type", "SCENE")

        # Command FFmpeg untuk membuat video placeholder berkualitas
        # Menggunakan color source dan drawtext filter
        cmd = [
            "ffmpeg", "-y",
            "-f", "lavfi", "-i", f"color=c=black:s=1280x720:d={duration}",
            "-vf", f"drawtext=text='{scene_type}':fontcolor=white:fontsize=48:x=(w-text_w)/2:y=(h-text_h)/2-50,"
                   f"drawtext=text='{description[:40]}...':fontcolor=gray:fontsize=24:x=(w-text_w)/2:y=(h-text_h)/2+20",
            "-c:v", "libx264", "-t", str(duration), "-pix_fmt", "yuv420p",
            output_path
        ]
        
        result = subprocess.run(cmd, capture_output=True, text=True)
        if result.returncode != 0:
            raise RuntimeError(f"FFmpeg failed to generate clip: {result.stderr}")

    def _stitch_videos(self, video_paths: List[str], output_path: str):
        """
        Menggabungkan list video menjadi satu menggunakan FFmpeg concat demuxer.
        """
        if not video_paths:
            raise ValueError("Tidak ada video untuk digabungkan")

        # Buat file list untuk concat
        list_filename = f"list_{int(time.time())}.txt"
        list_path = os.path.join(self._output_dir, list_filename)
        
        with open(list_path, "w", encoding="utf-8") as f:
            for path in video_paths:
                # FFmpeg concat butuh path absolut dengan escape
                abs_path = os.path.abspath(path).replace('\\', '/')
                f.write(f"file '{abs_path}'\n")

        # Command FFmpeg concat
        cmd = [
            "ffmpeg", "-y",
            "-f", "concat", "-safe", "0",
            "-i", list_path,
            "-c", "copy", # Copy codec agar cepat (asumsi semua clip sama formatnya)
            output_path
        ]

        try:
            result = subprocess.run(cmd, capture_output=True, text=True)
            if result.returncode != 0:
                raise RuntimeError(f"FFmpeg failed to stitch videos: {result.stderr}")
        finally:
            # Pastikan file list dihapus
            if os.path.exists(list_path):
                os.remove(list_path)
