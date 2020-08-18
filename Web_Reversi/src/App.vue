<template>
  <div>
    <header class="header" id="header">
      <div class="title">
        <h1>CMPT 383 Final Project</h1>
        <ul class="header-list">
          <li><a href="">New Game</a></li>
        </ul>
      </div>
    </header>
    <div style="text-align: center" v-if="!isStart">
      <div style="margin-top:160px;">
        <h2 style="margin-bottom: 30px">Reversi Game (Human vs AI)</h2>
        <div
            style="text-align: center; margin: 0 auto; height: 300px; width: 400px; box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1)">
          <h3 style="padding: 50px 0">Please select your role:</h3>
          <div class="btn">
            <a @click="startGame(1)">Black Piece</a>
          </div>
          <div class="btn">
            <a @click="startGame(2)">White Piece</a>
          </div>
        </div>
      </div>
    </div>
    <div v-show="isStart" style="margin-top: 50px; width: 100%; text-align: center">
      <div style="width: 1000px; margin: 0 auto">
        <div style="float: left; width: 700px">
          <div style="text-align: left">
            <h2 style="display: inline; color: #2c3e50">AI Player</h2>
            <h3 style="display: inline; color: green;">{{currTurn === humanPlayerPiece ? " (Current Turn) (AI needs 5 sec to think, please wait....)" : ""}}</h3>
            <h3 style="margin-bottom: 20px; opacity: 0.5;">{{humanPlayerPiece === 1 ? "White Piece" : "Black Piece"}}</h3>
          </div>
          <div style="position: relative; height: 600px; width: 600px" >
            <canvas class="layer" height="600px" ref="availablePositions" width="600px"></canvas>
            <canvas class="layer" height="600px" ref="pieces" width="600px"></canvas>
            <canvas class="layer" height="600px" ref="last_piece" width="600px"></canvas>
            <canvas @click="canvasOnClick" class="board" height="600px" ref="reversi_board" width="600px"></canvas>
          </div>
          <div style="margin-top: 20px; text-align: left; margin-bottom: 50px;">
            <h2 style="display: inline; color: #2c3e50">Human Player</h2>
            <h3 style="display: inline; color: green;">{{currTurn === humanPlayerPiece ? "" : " (Current Turn) (Please select a green circle.)"}}</h3>
            <h3 style="opacity: 0.5;">{{humanPlayerPiece === 1 ? "Black Piece" : "White Piece"}}</h3>
          </div>
        </div>
        <div style="float: left; margin-top: 70px; margin-left: 10px; box-shadow: -2px -2px 2px #EFEFEF, 5px 5px 5px #B9B9B9; padding: 20px; width: 230px">
          <h2 style="color:#2c3e50;">Scores:</h2>
          <hr style="margin: 20px 0"/>
          <div style="float: left; text-align: center">
            <h2>White</h2>
            <h2>{{whiteCount}}</h2>
          </div>
          <h2 style="float: left; margin: 0 30px">VS</h2>
          <div style="float: left; text-align: center;">
            <h2>Black</h2>
            <h2> {{blackCount}}</h2>
          </div>
          <div>
            <hr style="margin-top: 90px; margin-bottom: 20px"/>
            <h2 style="color: green">{{message}}</h2>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  export default {
    name: 'App',
    data: () => ({
      isStart: false,
      gameBoardState: null,
      availablePosition: null,
      currTurn: -1,
      loading: false,
      count: 1,
      humanPlayerPiece: 1,
      blackCount: 2,
      whiteCount: 2,
      message: "You have no available position to play, pass....",
    }),
    methods: {
      // listen on click event
      canvasOnClick(e) {
        if (this.loading) {
          return;
        }
        // calculate click position
        let positionX = Math.floor((e.offsetX - 35) / 70);
        let positionY = Math.floor((e.offsetY - 35) / 70);
        for (let position of this.availablePosition) {
          // check if the clicked position is an available position
          if ((positionX === position.Y) && (positionY === position.X)) {
            this.loading = true;
            // call api
            this.playNextStep(positionY, positionX);
            // clear canvas
            let context = this.$refs.availablePositions.getContext("2d");
            context.canvas.width = 600;
            break;
          }
        }

      },
      // draw game board using canvas
      drawBoard() {
        let board_context = this.$refs.reversi_board.getContext("2d");
        board_context.strokeStyle = "#708090"
        // vertical and horizontal lines
        for (let i = 0; i < 9; i++) {
          board_context.moveTo(20, 20 + i * 70)
          board_context.lineTo(580, 20 + i * 70);
          board_context.stroke();
          board_context.moveTo(20 + i * 70, 20);
          board_context.lineTo(20 + i * 70, 580);
          board_context.stroke();
        }
      },

      // draw the last piece position (blue circle)
      drawLastPiece(positionX, positionY) {
        const context = this.$refs.last_piece.getContext("2d");
        context.canvas.height = 600;
        context.fillStyle = 'rgba(255, 255, 255, 0)';
        context.beginPath();
        context.arc(55 + positionX * 70, 55 + positionY * 70, 20, 0, 2 * Math.PI);
        context.closePath();
        context.lineWidth = 3;
        context.strokeStyle = '#009BFF'
        context.stroke()
      },

      // draw black / white piece
      drawPiece(positionX, positionY, isBlack) {
        const context = this.$refs.pieces.getContext("2d");
        context.beginPath();
        context.arc(55 + positionX * 70, 55 + positionY * 70, 20, 0, 2 * Math.PI);
        context.closePath();
        if (isBlack) {
          context.fillStyle = "#242424"
        } else {
          context.fillStyle = "#D1D1D1"
        }
        context.fill();
      },
      // call api
      playNextStep(positionX, positionY) {
        this.message = "";
        let that = this;
        // send get request and get response from server
        this.$http.get('/next/' + positionX + '/' + positionY).then(function (response) {
          that.gameBoardState = response.data.State;
          that.drawAllPieces();
          that.drawLastPiece(response.data.LastPiece.Y, response.data.LastPiece.X);
          that.currTurn = response.data.CurrTurn;
          that.blackCount = response.data.BlackCount;
          that.whiteCount = response.data.WhiteCount;
          // check if game is over
          if (response.data.GameStatus !== 3) {
            if (response.data.GameStatus === 1){
              that.message = "Game Over! Black wins!"
            } else if (response.data.GameStatus === -1){
              that.message = "Game Over! White wins!"
            } else{
              that.message = "Game Over! It's a draw!"
            }
            that.loading = false;
            return;
          }
          that.waitForAI();
        })
      },
      // call api
      waitForAI() {
        let that = this;
        that.message = "";
        // send get request and get response from server
        this.$http.get('/wait').then(function (response) {
          that.gameBoardState = response.data.State;
          that.drawAllPieces();
          that.drawLastPiece(response.data.LastPiece.Y, response.data.LastPiece.X);
          that.availablePosition = response.data.AvailablePos;
          that.currTurn = response.data.CurrTurn;
          that.blackCount = response.data.BlackCount;
          that.whiteCount = response.data.WhiteCount;
          // check if game is over
          if (response.data.GameStatus !== 3) {
            if (response.data.GameStatus === 1){
              that.message = "Game Over! Black wins!"
            } else if (response.data.GameStatus === -1){
              that.message = "Game Over! White wins!"
            } else{
              that.message = "Game Over! It's a draw!"
            }
            that.loading = false;
            return;
          }
          // if the player has no available position, give up the current turn
          if (that.availablePosition === null || that.availablePosition.empty) {
            that.message = "You have no available position to play, pass...."
            that.pass();
          }else {
            for (let position of that.availablePosition) {
              that.drawAvailablePosition(position.Y, position.X);
            }
          }
          that.loading = false;
        })
      },
      // draw all available positions using green circle
      drawAvailablePosition(positionX, positionY) {
        const context = this.$refs.availablePositions.getContext("2d");
        context.fillStyle = 'rgba(255, 255, 255, 0)';
        context.beginPath();
        context.arc(55 + positionX * 70, 55 + positionY * 70, 20, 0, 2 * Math.PI);
        context.closePath();
        context.lineWidth = 3;
        context.strokeStyle = 'green'
        context.stroke()
      },
      // draw all black and white pieces
      drawAllPieces() {
        const context = this.$refs.pieces.getContext("2d");
        context.canvas.height = 600;
        for (let i = 0; i < 8; i++) {
          for (let j = 0; j < 8; j++) {
            if (this.gameBoardState[i][j] === 1) {
              this.drawPiece(i, j, true)
            } else if (this.gameBoardState[i][j] === -1) {
              this.drawPiece(i, j, false)
            }
          }
        }
      },
      // pass the current turn
      pass(){
        let that = this;
        // send get request and get response from server
        this.$http.get('/pass').then(function (response) {
          that.message = "";
          that.gameBoardState = response.data.State;
          that.drawAllPieces();
          that.drawLastPiece(response.data.LastPiece.Y, response.data.LastPiece.X);
          that.availablePosition = response.data.AvailablePos;
          that.currTurn = response.data.CurrTurn;
          that.blackCount = response.data.BlackCount;
          that.whiteCount = response.data.WhiteCount;
          if (response.data.GameStatus !== 3) {
            if (response.data.GameStatus === 1){
              that.message = "Game Over! Black wins!"
            } else if (response.data.GameStatus === -1){
              that.message = "Game Over! White wins!"
            } else{
              that.message = "Game Over! It's a draw!"
            }
            that.loading = false;
            return;
          }
          if (that.availablePosition === null || that.availablePosition.empty) {
            that.message = "You have no available position to play, pass...."
            that.pass();
          } else {
            for (let position of that.availablePosition) {
              that.drawAvailablePosition(position.Y, position.X);
            }
          }
          that.loading = false;
        })
      },

      startGame(role) {
        if (role === 2)
          this.humanPlayerPiece = -1
        this.message = "";
        let that = this;
        this.loading = true;
        this.isStart = true;
        this.gameBoardState = [[0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0],
                                [0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0],
                                [0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0],]
        this.gameBoardState[3][4] = 1
        this.gameBoardState[4][3] = 1
        this.gameBoardState[3][3] = -1
        this.gameBoardState[4][4] = -1
        this.drawAllPieces();

        // send get request and get response from server
        this.$http.get('/start/' + role).then(function (response) {
          that.gameBoardState = response.data.State;
          that.availablePosition = response.data.AvailablePos;
          for (let position of that.availablePosition) {
            that.drawAvailablePosition(position.Y, position.X);
          }
          that.drawAllPieces();
          that.drawLastPiece(response.data.LastPiece.Y, response.data.LastPiece.X);
          that.currTurn = response.data.CurrTurn;
          that.showAvailablePos = true;
          that.blackCount = response.data.BlackCount;
          that.whiteCount = response.data.WhiteCount;
          that.loading = false;
        })
      }
    },
    mounted() {
      this.drawBoard();
    }
  }
</script>

<style>
  #app {
    font-family: Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    text-align: center;
    color: #2c3e50;
    margin-top: 60px;
  }

  * {
    margin: 0;
    padding: 0;
  }

  .header-list {
    list-style: none;
    width: 60%;
    overflow: hidden;
    float: right;
  }

  .header-list li {
    float: right;
    line-height: 55px;
    margin: 0 20px;
  }

  .header-list li:hover {
    background: #fff;
  }

  .header-list li a:hover {
    color: #000;
  }

  .header-list li a {
    text-decoration: none;
    color: #fff;
    padding: 0 10px;
    display: block;
  }

  header {
    box-shadow: 0 2px 4px rgba(0, 0, 0, .12), 0 0 6px rgba(0, 0, 0, .04);
    background-color: #2B486B;
    width: 100%;
    color: white;
    height: 55px;
    min-width: 1200px;
  }

  .title {
    width: 1200px;
    height: 100%;
    padding: 0 100px;
    margin: 0 auto;
    box-sizing: border-box;
    position: relative;
  }

  a {
    font-size: 16px;
  }

  h1 {
    margin-top: 8px;
    float: left;
  }

  .board {
    display: block;
    position: absolute;
    box-shadow: -2px -2px 2px #EFEFEF, 5px 5px 5px #B9B9B9;
    top:50%;
    left:50%;
    transform: translate(-50%, -50%);
  }

  .layer {
    display: block;
    position: absolute;
    top:50%;
    left:50%;
    transform: translate(-50%, -50%);
  }

  .btn {
    margin-top: 30px;
    height: 50px;
  }

  .btn a {
    text-decoration: none;
    background: #2f435e;
    color: #f2f2f2;
    padding: 15px 30px 15px 30px;
    font-size: 16px;
    font-family: Arial, Helvetica, Verdana, sans-serif;
    font-weight: bold;
    border-radius: 4px;
    -webkit-transition: all linear 0.30s;
    -moz-transition: all linear 0.30s;
    transition: all linear 0.30s;

  }

  .btn a:hover {
    background: #385f9e;
  }
</style>
