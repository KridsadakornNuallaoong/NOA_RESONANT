# filename: server.py
import json
import os

import joblib
import numpy as np
import tensorflow as tf
import uvicorn
from dotenv import load_dotenv
from fastapi import FastAPI, Request
from fastapi.responses import HTMLResponse
from fastapi.websockets import WebSocket

# Load environment variables from .env file
load_dotenv()

app = FastAPI()

# Load once at startup
model = tf.keras.models.load_model("model_latest.h5")
scaler = joblib.load("scaler.pkl")

@app.post("/predict")
async def predict(request: Request):
    data = await request.json()
    inputs = np.array(data["inputs"])
    inputs = scaler.transform(inputs)
    prediction = model.predict(inputs)
    return {"prediction": prediction.tolist()}


# create websocket to predict data
@app.websocket("/ws/predict")
async def websocket_predict(websocket: WebSocket):
    await websocket.accept()
    while True:
        data = await websocket.receive_text()
        inputs = np.array(json.loads(data)["inputs"])
        inputs = scaler.transform(inputs)
        prediction = model.predict(inputs)
        await websocket.send_text(json.dumps({"prediction": prediction.tolist()}))

if __name__ == "__main__":
    port = int(os.getenv("PORT", 8000))  # Default to 8000 if PORT is not set
    uvicorn.run(app, host="0.0.0.0", port=port)