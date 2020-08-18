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

# Foreign function interface
lib = cdll.LoadLibrary("lib/MonteCarloTreeSearch.so")

app = FastAPI()

# CORS Policy
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"])

# Define the return types from ctypes
lib.startGame.restype = c_void_p
lib.nextStep.restype = c_void_p
lib.waitForAI.restype = c_void_p
lib.passCurrentTurn.restype = c_void_p


@app.get("/start/{order}")
async def start_game(order: int):
    """
        API 1: start the Reversi game

        Args:
            order: the player order (1 or 2)

        Returns:
            a json dict containing game info
    """
    string = lib.startGame(order)  # call foreign function from Go
    reversi_game_result = cast(string, c_char_p).value.decode()
    lib.freeMemory(string)  # free memory
    json_obj = json.loads(reversi_game_result)  # convert string to json data
    return json_obj


@app.get("/next/{positionX}/{positionY}")
async def play_next_step(positionX: int, positionY: int):
    """
        API 2: receive the user's play

        Args:
            positionX: the user's next piece position X on game board
            positionY: the user's next piece position Y on game board

        Returns:
            a json dict containing game info
    """
    string = lib.nextStep(positionX, positionY)  # call foreign function from Go
    reversi_game_result = cast(string, c_char_p).value.decode()
    lib.freeMemory(string)  # free memory
    json_obj = json.loads(reversi_game_result)  # convert string to json data including the game states
    return json_obj


@app.get("/wait")
async def wait_for_AI():
    """
        API 3: wait for AI's play

        Returns:
            a json dict containing game info
    """
    string = lib.waitForAI()  # call foreign function from Go
    reversi_game_result = cast(string, c_char_p).value.decode()
    lib.freeMemory(string)  # free memory
    json_obj = json.loads(reversi_game_result)
    return json_obj


@app.get("/pass")
async def pass_current_turn():
    """
        API 4: pass current turn
        Called when user has no available positions to play

        Returns:
            a json dict containing game info
    """
    string = lib.passCurrentTurn  # call foreign function from Go
    reversi_game_result = cast(string, c_char_p).value.decode()
    lib.freeMemory(string)    # free memory
    json_obj = json.loads(reversi_game_result)
    return json_obj


if __name__ == '__main__':
    uvicorn.run(app='main:app', host="0.0.0.0", port=8000, reload=False, debug=False)
