"""
    CMPT 383 Final Project
    FastApi Server
    Author: Yiming Zhang
"""
import uvicorn
from fastapi import FastAPI
from starlette.middleware.cors import CORSMiddleware
from ctypes import *
import json
import model

lib = cdll.LoadLibrary("lib/MonteCarloTreeSearch.so")

app = FastAPI()

# CORS Policy
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"])

lib.startGame.restype = c_char_p
lib.nextStep.restype = c_char_p
lib.waitForAI.restype = c_char_p
lib.passCurrentTurn.restype = c_char_p


@app.get("/start/{order}")
async def start_game(order: int):
    reversi_game_result = lib.startGame(order).decode()
    json_obj = json.loads(reversi_game_result)
    return json_obj


@app.get("/next/{positionX}/{positionY}")
async def play_next_step(positionX: int, positionY: int):
    reversi_game_result = lib.nextStep(positionX, positionY).decode()
    json_obj = json.loads(reversi_game_result)
    return json_obj


@app.get("/wait")
async def wait_for_AI():
    reversi_game_result = lib.waitForAI().decode()
    json_obj = json.loads(reversi_game_result)
    return json_obj


@app.get("/pass")
async def pass_current_turn():
    reversi_game_result = lib.passCurrentTurn.decode()
    json_obj = json.loads(reversi_game_result)
    return json_obj


if __name__ == '__main__':
    uvicorn.run(app='main:app', host="0.0.0.0", port=8000, reload=True, debug=True)
