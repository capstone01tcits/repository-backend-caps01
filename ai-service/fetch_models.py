import os
import requests
import json
from dotenv import load_dotenv

load_dotenv()

api_key = os.getenv("WAVESPEED_API_KEY")
url = "https://api.wavespeed.ai/api/v3/models"

headers = {
    "Authorization": f"Bearer {api_key}"
}

try:
    print(f"Fetching models from {url}...")
    response = requests.get(url, headers=headers)
    data = response.json()
    
    if data.get("code") == 200:
        models = data.get("data", [])
        print(f"Found {len(models)} models:")
        for m in models:
            # Look for Veo or WAN in the slugs
            slug = m.get("slug")
            name = m.get("name")
            print(f" - {name}: {slug}")
    else:
        print(f"Error from API: {data}")
except Exception as e:
    print(f"Request failed: {e}")
