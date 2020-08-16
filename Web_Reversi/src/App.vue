<template>
  <div>
    <header class="header" id="header">
      <div class="title">
        <h1>CMPT 383 Final Project</h1>
        <ul class="header-list">
          <li><a href="">About</a></li>
          <li><a href="hello_page.html">New Game</a></li>
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
    <div v-show="isStart" style="margin-top: 100px; margin-left: 300px; width: 1300px">
      <div>
        <canvas class="layer" height="600px" ref="availablePositions" width="600px"></canvas>
        <canvas class="layer" height="600px" ref="pieces" width="600px"></canvas>
        <canvas @click="canvasOnClick" class="board" height="600px" ref="reversi_board" width="600px"></canvas>
      </div>
      <div style="float: right; width: 230px; height: 560px; box-shadow: -2px -2px 2px #EFEFEF, 5px 5px 5px #B9B9B9; padding: 20px">
        <h2>Current Turn: {{ currTurn === -1 ? 'Black' : 'White' }}</h2>
        <h3 v-show="loading" style="margin-top: 20px; color:#2B486B;">AI needs 5 sec to think....</h3>
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
  }),
  methods: {
    canvasOnClick(e) {
      if (this.loading) {
        return;
      }
      let positionX = Math.round((e.offsetX - 35) / 70);
      let positionY = Math.round((e.offsetY - 35) / 70);
      for (let position of this.availablePosition) {
        if ((positionX === position.Y) && (positionY === position.X)) {
          this.loading = true;
          this.playNextStep(positionY, positionX);
          let context = this.$refs.availablePositions.getContext("2d");
          context.canvas.width = 600;
          break;
        }
      }

    },
    drawBoard() {
      let board_context = this.$refs.reversi_board.getContext("2d");
      board_context.strokeStyle = "#708090"
      for (let i = 0; i < 9; i++) {
        board_context.moveTo(20, 20 + i * 70)
        board_context.lineTo(580, 20 + i * 70);
        board_context.stroke();
        board_context.moveTo(20 + i * 70, 20);
        board_context.lineTo(20 + i * 70, 580);
        board_context.stroke();
      }
    },
    drawPiece(positionX, positionY, isBlack) {
      const context = this.$refs.pieces.getContext("2d");
      context.beginPath();
      context.arc(55 + positionX * 70, 55 + positionY * 70, 20, 0, 2 * Math.PI);
      context.closePath();
      if (isBlack) {
        context.fillStyle = "#242424"
      } else {
        context.fillStyle = "#DEDEDE"
      }
      context.fill();
    },
    playNextStep(positionX, positionY) {
      let that = this;
      this.$http.get('/next/' + positionX + '/' + positionY).then(function (response) {
        that.gameBoardState = response.data.State;
        that.drawAllPieces();
        that.currTurn = response.data.CurrTurn;
        that.waitForAI();
      })
    },
    waitForAI() {
      let that = this;

      this.$http.get('/wait').then(function (response) {
        that.gameBoardState = response.data.State;
        that.drawAllPieces();
        that.availablePosition = response.data.AvailablePos;
        for (let position of that.availablePosition) {
          that.drawAvailablePosition(position.Y, position.X);
        }
        that.currTurn = response.data.CurrTurn;
        that.loading = false;
      })
    },
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

    startGame(role) {
      let that = this;
      this.loading = true;
      that.isStart = true;
      that.gameBoardState = [[0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0],
                              [0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0],
                              [0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0],]
      that.gameBoardState[3][4] = 1
      that.gameBoardState[4][3] = 1
      that.gameBoardState[3][3] = -1
      that.gameBoardState[4][4] = -1
      that.drawAllPieces();
      this.$http.get('/start/' + role).then(function (response) {
        that.gameBoardState = response.data.State;
        that.availablePosition = response.data.AvailablePos;
        for (let position of that.availablePosition) {
          that.drawAvailablePosition(position.Y, position.X);
        }
        that.currTurn = response.data.CurrTurn;
        that.showAvailablePos = true;
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
    font-family: Avenir, Helvetica, Arial, sans-serif;
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

  body {
    width: 100%;
    min-width: 1200px;
    font-family: Arial, serif;
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

  input {
    height: 25px;
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
