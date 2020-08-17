"""
    CMPT 383 Final Project
    FastAPI Server
    Author: Yiming Zhang
"""
import uvicorn
from fastapi import FastAPI
from starlette.middleware.cors import CORSMiddleware
from ctypes import *
import json

lib = cdll.LoadLibrary("lib/MonteCarloTreeSearch.so")

app = FastAPI()

# CORS Policy
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"])

lib.startGame.restype = c_void_p
lib.nextStep.restype = c_void_p
lib.waitForAI.restype = c_void_p
lib.passCurrentTurn.restype = c_void_p


@app.get("/start/{order}")
async def start_game(order: int):
    string = lib.startGame(order)
    reversi_game_result = cast(string, c_char_p).value.decode()
    # free memory
    lib.freeMemory(string)
    json_obj = json.loads(reversi_game_result)
    return json_obj


@app.get("/next/{positionX}/{positionY}")
async def play_next_step(positionX: int, positionY: int):
    string = lib.nextStep(positionX, positionY)
    reversi_game_result = cast(string, c_char_p).value.decode()
    # free memory
    lib.freeMemory(string)
    json_obj = json.loads(reversi_game_result)
    return json_obj


@app.get("/wait")
async def wait_for_AI():
    string = lib.waitForAI()
    reversi_game_result = cast(string, c_char_p).value.decode()
    # free memory
    lib.freeMemory(string)
    json_obj = json.loads(reversi_game_result)
    return json_obj


@app.get("/pass")
async def pass_current_turn():
    string = lib.passCurrentTurn
    reversi_game_result = cast(string, c_char_p).value.decode()
    # free memory
    lib.freeMemory(string)
    json_obj = json.loads(reversi_game_result)
    return json_obj


if __name__ == '__main__':
    uvicorn.run(app='main:app', host="0.0.0.0", port=8000, reload=False, debug=False)
