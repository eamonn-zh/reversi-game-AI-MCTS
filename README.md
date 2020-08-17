## Reversi Game AI Based on Concurrent Monte Carlo Tree Search (MCTS)

CMPT 383 Final Project

Yiming Zhang (301354482)



---------------------------------------------------------------------------------------------------------------

### 1. Introduction

This project implements an AI for Reversi game, it is based on Monte Carlo tree search (MCTS) algorithm, which is used in AlphaGo AI. Algorithm runs on the Python backend server, and also there is a web frontend, where users can play the Reversi games against the AI. Frontend and backend interact through REST APIs.

The secret to the success of the MCTS algorithm is to simulate as many games as possible in a limited time. So the algorithm is implemented in Go and is called by Python through the foreign function interface. In addition, I made this this algorithm concurrent.

---------------------------------------------------------------------------------------------------------------

### 2. Languages & Communication

- Languages:
   1. Go
   2. Python
   3. JavaScript
   4. Bash	

- Communication methods:
   1. Foreign function interface (between Python and Go through ctypes)
   2. REST APIs (between JavaScript and Python)

---------------------------------------------------------------------------------------------------------------


### 3. Project Structure

- **cookbooks/**  ---  Chef configurations
- **FastAPI_Server/**  ---  high-performance Python API backend server (similar to Django)
- **Go_MCTS/**  ---  Golang program (implemented concurrent Monte Carlo tree search algorithm)
- **Web_Reversi/**  ---  frontend web page (implemented via Vue.js, JavaScript code is in Web_Reversi/src/App.vue)
- **Vagrantfile**  ---  vagrant configuration

---------------------------------------------------------------------------------------------------------------

### 4. How to Run

```bash
# start the VM
vagrant up

# run the provisioner
vagrant provision
```
Then open your web browser and access http://localhost:8080/



*Notes: Due to unknown reasons, Chef cannot run the Node.js web server in background with command `npm run serve &`, so I use `npm run serve` instead, which will stall the bash when you execute `vagrant provision`, it does not affect the operation of the program.*