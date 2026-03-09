from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware

app = FastAPI(title="AI Video Gen Service")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.get("/health")
def health():
    return {"success": True, "status": "ok", "message": "AI service is running"}


@app.get("/status")
def status():
    return {"success": True, "message": "Connected to AI service"}


@app.post("/generate")
async def generate(request: Request):
    body = await request.json()
    user_id = request.headers.get("X-User-ID", "unknown")
    user_email = request.headers.get("X-User-Email", "unknown")

    return {
        "success": True,
        "message": "Connected - generate endpoint reached",
        "user_id": user_id,
        "user_email": user_email,
        "received_data": body,
    }


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
