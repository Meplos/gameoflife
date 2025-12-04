const SCALE = 5;
const DEBUG = "#6DC3BB"
const ALIVE_COLOR = "#F4B342"
const DEAD_COLOR = "black"

class Drawer {

  constructor() {
    this.canvas = document.getElementById("canvas")
    this.ctx = this.canvas.getContext("2d")
    this.debug = true


  }

  setDimension() {
    const dpr = window.devicePixelRatio || 1;
    const rect = this.canvas.getBoundingClientRect();

    // grid size in cells based on CSS pixels
    this.w = Math.floor(rect.width / SCALE);
    this.h = Math.floor(rect.height / SCALE);

    // set internal canvas size in device pixels
    const width = Math.max(1, Math.round(rect.width * dpr));
    const height = Math.max(1, Math.round(rect.height * dpr));
    if (this.canvas.width !== width || this.canvas.height !== height) {
      this.canvas.width = width;
      this.canvas.height = height;
    }

    // ensure drawing uses CSS pixel coordinates
    this.ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
  }

  clean() {
    this.ctx.fillStyle = DEAD_COLOR
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

}

const CMD = {
  PLAY: "play",
  PAUSE: "pause",
  RESTART: "restart",
  INIT: "init"
}

class WindowState {

  constructor() {
    this.h = window.innerHeight / SCALE;
    this.w = window.innerWidth / SCALE;
    this.playPause = document.getElementById("playPause");
    this.setPlayPause(true)
    this.playPause.addEventListener("click", function () {
      if (this.pause) {
        this.send(CMD.PLAY);
        return
      }
      this.send(CMD.PAUSE);
    }.bind(this))

    this.restart = document.getElementById("restart");
    this.restart.addEventListener("click", function () {
      this.send(CMD.RESTART);
    }.bind(this))

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
    renderFrame();
  }

  setState({ latestBoard, dirty }) {
    this.latestBoard = latestBoard;
    this.dirty = dirty;
  }

  init() {
    this.send(CMD.INIT, { w: this.w, h: this.h })
  }


  send(cmd, option) {
    this.ws.send(cmd, option)
  }
}


class WsClient {
  constructor(windowState) {
    this.windowState = windowState
  }

  connect() {
    console.log("WsClient.Connect")
    return new Promise((resolve) => {

      this.ws = new WebSocket((location.protocol === 'https:' ? 'wss://' : 'ws://') + location.host + '/ws');

      this.ws.onmessage = (ev) => {
        const latestBoard = JSON.parse(ev.data);

        this.windowState.setState({
          latestBoard,
          dirty: true
        });


        this.windowState.setPlayPause(latestBoard.pause);
      };
      this.ws.addEventListener("open", () => resolve())
    })
  }

  send(cmd, options) {
    const payload = { cmd: cmd }
    if (options) {
      payload.options = options
    }
    this.ws.send(JSON.stringify(payload))
  }

}

let state;
let dw;
let ws;
async function init() {
  state = new WindowState()
  dw = new Drawer();

  // initialize canvas to full viewport and set grid size
  dw.setDimension();
  state.w = dw.w;
  state.h = dw.h;

  ws = new WsClient(state)
  state.setSocket(ws)
  await ws.connect()
  dw.clean()
  state.init()

  // handle viewport resize
  window.addEventListener('resize', () => {
    dw.setDimension();
    state.w = dw.w;
    state.h = dw.h;
    state.send(CMD.INIT, { w: state.w, h: state.h });
    dw.clean();
  });
}


function renderFrame() {
  function frame() {
    if (state.pause) return
    if (state.dirty && state.latestBoard) {
      const start = Date.now()
      console.log("TYPE", state.latestBoard.type)
      if (state.latestBoard.type === "init") {
        dw.render(state.latestBoard.state)
      } else {
        dw.renderChanges(state.latestBoard.changes);
      }
      state.dirty = false;
      const end = Date.now()
      console.log(`Frame took: ${end - start}ms`)
    }
    requestAnimationFrame(frame)
  }
  requestAnimationFrame(frame);
}
init()
