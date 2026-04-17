import os

# Set dummy API key agar build_provider tidak error environment
os.environ["LTX_API_KEY"] = "dummy_key"
os.environ["RUNWAY_API_KEY"] = "dummy_key"

from providers.base import VideoRequest
from providers.router import enrich_routing, build_provider

def check_model_routing():
    print("=== UJI COBA ROUTING TANPA MENGGUNAKAN API CREDIT ===")

    # Simulasi request seperti di main.py atau run_trial.py
    req = VideoRequest(
        instruction="text_to_video",
        input="Test prompt",
        routing={"task_type": "text_to_video"}
    )

    # 1. Cek Hasil Routing
    req = enrich_routing(req)
    provider_name = req.routing.get("resolved_provider")
    model_name = req.routing.get("resolved_model")

    print("\n[1] Hasil Routing:")
    print(f"    - Provider Target : {provider_name}")
    print(f"    - Model Target    : {model_name}")

    # 2. Cek Instansiasi Provider
    provider = build_provider(provider_name, model_name)
    print("\n[2] Hasil Build Provider:")
    print(f"    - Instance Name   : {provider.name}")

    # 3. Cek Payload yang akan dikirim (tanpa mengirimnya ke API)
    payload_builder = getattr(provider, "_build_payload", None)
    if callable(payload_builder):
        try:
            payload = payload_builder(req)
            print(f"\n[3] Payload akhir yang siap dikirim ke {provider_name.upper()}:")
            print(f"    - Model dalam payload : {payload.get('model')}")
            print(f"    - Generate Audio flag : {payload.get('generate_audio')}")
            print(f"    - Resolusi            : {payload.get('resolution')}")
        except (ValueError, TypeError, AttributeError) as exc:
            print(f"\n[3] Error saat build payload: {exc}")
    else:
        print(f"\n[3] Provider {provider_name} tidak memiliki method _build_payload untuk dicek")

if __name__ == "__main__":
    check_model_routing()