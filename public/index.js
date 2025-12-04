const SCALE = 5;
const DEBUG = "#6DC3BB"
const ALIVE_COLOR = "#F4B342"
const DEAD_COLOR = "black"

class Drawer {

  constructor() {
    this.canvas = document.getElementById("canvas")
    this.ctx = this.canvas.getContext("2d")
    this.debug = true

    this.handleDebugButton()

  }

  setDimension(w, h) {
    if (this.w && this.h) return
    this.w = w;
    this.h = h;
    this.canvas.width = w * SCALE;
    this.canvas.height = h * SCALE;
  }

  clean() {
    this.ctx.fillStyle = "black"
    this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);
  }

  drawCell(x, y) {
    this.ctx.rect(x * SCALE, y * SCALE, SCALE, SCALE)
  }


  renderChanges(changes) {
    const alives = []
    const deads = []
    changes.forEach(c => {
      if (c.state) {
        alives.push(c)
        return
      }
      deads.push(c)
    });


    this.renderCells(deads, DEAD_COLOR)
    this.renderCells(alives, ALIVE_COLOR)
  }

  renderCells(array, color) {
    this.ctx.beginPath()
    this.ctx.fillStyle = color
    for (const c of array) {
      this.drawCell(c.x, c.y)

    }
    this.ctx.fill()
  }



  render(data) {
    console.log("RenderFull")
    this.clean()
    this.ctx.fillStyle = ALIVE_COLOR
    this.ctx.beginPath();
    for (let y = 0; y < this.h; y++) {
      const decodedState = data[y]
      for (let x = 0; x < this.w; x++) {
        const cell = decodedState[x]
        if (!cell) {
          continue;
        }
        this.drawCell(x, y)
      }
    }

    this.ctx.fill();
  }

  handleDebugButton() {
    document.getElementById("debugBtn").addEventListener("click", () => this.debug = !this.debug)
  }
}

const CMD = {
  PLAY: "play",
  PAUSE: "pause",
  RESTART: "restart"
}

class WindowState {

  constructor() {
    this.playPause = document.getElementById("playPause");
    this.playPause.addEventListener("click", function () {
      if (this.pause) {
        this.send(CMD.PLAY);
        return
      }
      this.send(CMD.PAUSE);
    }.bind(this))

    this.restart = document.getElementById("restart");
    this.restart.addEventListener("click", function () { this.send(CMD.RESTART) }.bind(this))

  }

  setSocket(socket) {
    this.ws = socket;
  }

  setPlayPause(value) {
    if (value === this.pause) return;
    this.pause = value;
    if (this.pause) {
      this.playPause.innerHTML = "PLAY";
      return
    }
    this.playPause.innerHTML = "PAUSE";
  }

  setState({ latestBoard, dirty }) {
    this.latestBoard = latestBoard;
    this.dirty = dirty;
  }

  send(cmd) {
    this.ws.send(cmd)
  }

}


class WsClient {
  constructor(windowState) {
    this.windowState = windowState
  }

  connect() {
    console.log("WsClient.Connect")
    this.ws = new WebSocket((location.protocol === 'https:' ? 'wss://' : 'ws://') + location.host + '/ws');

    this.ws.onmessage = (ev) => {
      const latestBoard = JSON.parse(ev.data);

      this.windowState.setState({
        latestBoard,
        dirty: true
      });

      this.windowState.setPlayPause(latestBoard.pause);
    };
  }

  send(cmd) {
    this.ws.send(JSON.stringify({ cmd: cmd }))
  }

}

function init() {
  const state = new WindowState()
  const ws = new WsClient(state)
  state.setSocket(ws)
  ws.connect()
  const dw = new Drawer();


  function frame() {
    if (state.dirty && state.latestBoard && !state.pause) {
      const start = Date.now()
      if (dw.w !== state.latestBoard.w || dw.h !== state.latestBoard.h) {
        dw.setDimension(state.latestBoard.w, state.latestBoard.h);
        dw.clean()
      }

      dw.renderChanges(state.latestBoard.changes);
      //dw.render(state.latestBoard.state)
      state.dirty = false;
      const end = Date.now()
      console.log(`Frame took: ${end - start}ms`)
    }
    requestAnimationFrame(frame);
  }
  requestAnimationFrame(frame);
}
init()
