"use strict";
class GameCanvas {
    constructor(canvasID, layersCount) {
        this.canvas = document.getElementById(canvasID);
        const ctx = this.canvas.getContext("2d");
        if (!ctx) {
            throw new Error("Failed to get context");
        }
        this.ctx = ctx;
        this.layers = new Array();
        for (let i = 0; i < layersCount; i++) {
            this.layers.push(new Map());
        }
        this.drawables = new Map();
        this.mousePos = new Vector2(0, 0);
        this.backgroundColor = RGB(0, 0, 0);
        this.resize();
        window.addEventListener('resize', () => this.resize(), false);
        this.canvas.addEventListener("mousemove", e => {
            this.updateMousePos(e);
        });
    }
    stop() {
        this.canvas.remove();
    }
    resize() {
        const canvas = this.canvas;
        canvas.width = window.innerWidth;
        canvas.height = window.innerHeight - canvas.getBoundingClientRect().top;
        this.draw();
    }
    insertDrawable(drawable, layerNum, id) {
        if (this.drawables.has(id)) {
            return;
        }
        this.drawables.set(id, drawable);
        const layer = this.layers[layerNum];
        if (layer) {
            layer.set(id, drawable);
        }
    }
    deleteDrawable(id) {
        this.drawables.delete(id);
        this.layers.forEach(layer => {
            layer.delete(id);
        });
    }
    draw() {
        const ctx = this.ctx;
        this.clear();
        this.layers.forEach(layer => {
            layer.forEach(drawable => {
                drawable.draw(ctx);
            });
        });
    }
    clear() {
        const ctx = this.ctx;
        const canvas = this.canvas;
        ctx.clearRect(0, 0, canvas.width, canvas.height);
        ctx.fillStyle = this.backgroundColor;
        ctx.fillRect(0, 0, canvas.width, canvas.height);
    }
    setBackgroundColor(color) {
        this.backgroundColor = color;
    }
    updateMousePos(e) {
        let rect = this.canvas.getBoundingClientRect();
        let x = e.clientX - rect.left;
        let y = e.clientY - rect.top;
        this.mousePos.setPosition(x, y);
    }
    getMousePos() {
        return this.mousePos;
    }
    getDrawable(id) {
        return this.drawables.get(id);
    }
}
class Controls {
    constructor() {
        // using map instead of set here because golang doesn't have set implementation yet
        this.heldControls = new Map();
        this.bindings = new Map();
        // on key press
        document.addEventListener("keypress", (e) => {
            // only single press
            if (e.repeat) {
                return;
            }
            // if no key in bindings
            if (!this.bindings.has(e.key)) {
                return;
            }
            // get control from binding
            const control = this.bindings.get(e.key);
            if (control) {
                this.heldControls.set(control, true);
            }
        });
        document.addEventListener("keyup", (e) => {
            if (!this.bindings.has(e.key)) {
                return;
            }
            const control = this.bindings.get(e.key);
            if (control) {
                this.heldControls.delete(control);
            }
        });
    }
    bindControl(key, control) {
        this.bindings.set(key, control);
    }
    isHeld(control) {
        return this.heldControls.has(control);
    }
    getHeldControls() {
        return this.heldControls;
    }
}
class DeltaTimer {
    constructor() {
        this.lastTick = performance.now();
    }
    getDeltaTime() {
        let now = performance.now();
        let dt = now - this.lastTick;
        this.lastTick = now;
        return dt;
    }
}
function lerp(a, b, alpha) {
    return a + alpha * (b - a);
}
class Rect {
    constructor() {
        this.position = new Vector2(0, 0);
        this.size = new Vector2(0, 0);
    }
    setPosition(x, y) {
        this.position.setPosition(x, y);
    }
    setSize(x, y) {
        this.size.setPosition(x, y);
    }
    containsPoint(x, y) {
        let pos = this.position;
        let size = this.size;
        if (x >= pos.x && x <= pos.x + size.x &&
            y >= pos.y && y <= pos.y + size.y) {
            return true;
        }
        return false;
    }
    getPosition() {
        return this.position;
    }
    getSize() {
        return this.size;
    }
}
/// <reference path="rect.ts"/>
class Drawable extends Rect {
    constructor() {
        super();
    }
    draw(_ctx) { }
}
class RoomGui {
    constructor() {
        this.navVisible = true;
        this.onclick();
        $("#show-nav").click(() => {
            this.onclick();
        });
    }
    showMessage(text) {
        $("#message").text(text);
    }
    setNavBarVisibility(visibility) {
        this.navVisible = visibility;
        if (visibility) {
            $("#navigation").show();
            $("#gui").removeAttr("style");
            $("#show-nav").text("↑");
            return;
        }
        $("#navigation").hide();
        $("#gui").css("top", 0);
        $("#show-nav").text("↓");
    }
    onclick() {
        this.navVisible = !this.navVisible;
        this.setNavBarVisibility(this.navVisible);
    }
}
// using global variable, so other scripts can use it
const roomGui = new RoomGui();
class RectangleShape extends Drawable {
    constructor(width, height) {
        super();
        this.setSize(width, height);
        this.color = RGB(255, 255, 255);
    }
    setColor(color) {
        this.color = color;
    }
    draw(ctx) {
        let pos = this.position;
        let size = this.size;
        ctx.fillStyle = this.color;
        ctx.fillRect(pos.x, pos.y, size.x, size.y);
    }
}
function RGB(r, g, b) {
    return `rgb(${r} ${g} ${b})`;
}
class DrawableText extends Drawable {
    constructor(string, fontSize) {
        super();
        this.string = string;
        this.color = RGB(255, 255, 255);
        this.fontSize = fontSize;
        this.font = "serif";
    }
    setString(string) {
        this.string = string;
    }
    setColor(color) {
        this.color = color;
    }
    setFont(font) {
        this.font = font;
    }
    setFontSize(fontSize) {
        this.fontSize = fontSize;
    }
    draw(ctx) {
        ctx.fillStyle = this.color;
        ctx.font = `${this.fontSize}px ${this.font}`;
        // adding height to y because text's origin is located on the bottom
        const metrics = ctx.measureText(this.string);
        const height = metrics.actualBoundingBoxAscent + metrics.actualBoundingBoxDescent;
        ctx.fillText(this.string, this.position.x, this.position.y + height);
    }
}
class Ticker {
    constructor() {
        this.timer = new DeltaTimer();
    }
    tick(callback) {
        let dt = this.timer.getDeltaTime();
        callback(dt);
        requestAnimationFrame(() => this.tick(callback));
    }
}
class Vector2 {
    constructor(x, y) {
        this.x = x;
        this.y = y;
    }
    setPosition(x, y) {
        this.x = x;
        this.y = y;
    }
}
class GameWebSocket {
    constructor() {
        this.handleClose = (_body) => { };
        this.handleMessage = (_type, _body) => { };
        this.active = false;
        this.websocket = null;
    }
    openConnection(gameID, roomID) {
        const protocol = location.protocol == "https:" ? "wss:" : "ws:";
        this.websocket = new WebSocket(`${protocol}//${document.location.host}/rt/ws/game/${gameID}/room/${roomID}`);
        const ws = this.websocket;
        ws.onopen = _e => {
            this.active = true;
        };
        ws.onclose = _e => {
            if (!this.active) {
                return;
            }
            this.handleClose("Connection closed");
            this.active = false;
        };
        ws.onerror = _e => {
            if (!this.active) {
                return;
            }
            this.handleClose("Something went wrong");
            this.active = false;
        };
        ws.onmessage = (e) => {
            if (!this.active) {
                return;
            }
            const data = JSON.parse(e.data);
            if (data.type == "close") {
                this.handleClose(data.body);
                this.active = false;
            }
            else {
                this.handleMessage(data.type, data.body);
            }
        };
    }
    sendMessage(type, body) {
        if (!this.active) {
            return;
        }
        const ws = this.websocket;
        if (!ws) {
            return;
        }
        ws.send(JSON.stringify({
            type: type,
            body: body,
        }));
    }
}
class Platformer {
    constructor() {
        const layers = 2;
        this.playerRectID = 0;
        this.canSendControls = true;
        this.canvas = new GameCanvas("canvas", layers);
        this.canvas.setBackgroundColor(RGB(30, 100, 100));
        this.controls = new Controls();
        this.bindControls();
        this.websocket = new GameWebSocket();
        const gameID = $("main").data("game-id");
        const roomID = $("main").data("room-id");
        this.initWebsocket(gameID, roomID);
        this.ticker = new Ticker();
        this.ticker.tick(dt => this.tick(dt));
    }
    bindControls() {
        const controls = this.controls;
        controls.bindControl("d", "right");
        controls.bindControl("a", "left");
        controls.bindControl("w", "jump");
        controls.bindControl(" ", "jump");
    }
    initWebsocket(gameID, roomID) {
        this.websocket.handleMessage = (type, body) => {
            switch (type) {
                case "gameinfo":
                    this.handleGameInfoMessage(body);
                    break;
                case "level":
                    this.handleLevelMessage(body);
                    break;
                case "deltas":
                    this.handleDeltasMessage(body);
                    break;
                case "delete":
                    this.handleDeleteMessage(body);
                    break;
                case "create":
                    this.handleCreateMessage(body);
                    break;
                default:
                    break;
            }
        };
        this.websocket.handleClose = (body) => {
            this.stopWithText(body);
        };
        this.websocket.openConnection(gameID, roomID);
    }
    tick(_dt) {
        if (this.canSendControls) {
            let heldControls = this.controls.getHeldControls();
            if (heldControls.size > 0) {
                let json = JSON.stringify(Object.fromEntries(heldControls.entries()));
                this.websocket.sendMessage("input", json);
                this.canSendControls = false;
            }
        }
        this.canvas.draw();
    }
    stopWithText(text) {
        this.canvas.stop();
        roomGui.showMessage(text);
        roomGui.setNavBarVisibility(true);
    }
    createRectangleShape(serverRect, rectID) {
        let rectangle = new RectangleShape(serverRect.size.x, serverRect.size.y);
        rectangle.setPosition(serverRect.position.x, serverRect.position.y);
        this.canvas.insertDrawable(rectangle, 0, rectID);
    }
    handleLevelMessage(body) {
        for (let [key, val] of Object.entries(body.kinematic)) {
            const id = Number(key);
            const serverRect = val;
            this.createRectangleShape(serverRect, id);
        }
        for (let [key, val] of Object.entries(body.static)) {
            const id = Number(key);
            const serverRect = val;
            this.createRectangleShape(serverRect, id);
        }
    }
    handleDeleteMessage(body) {
        this.canvas.deleteDrawable(body.ID);
    }
    handleCreateMessage(body) {
        let serverRect = body.rect;
        let rectID = body.id;
        this.createRectangleShape(serverRect, rectID);
    }
    handleDeltasMessage(body) {
        this.canSendControls = true;
        for (let [key, val] of Object.entries(body.kinematic)) {
            const id = Number(key);
            const serverRect = val;
            const clientRect = this.canvas.getDrawable(id);
            if (clientRect) {
                clientRect.setPosition(serverRect.position.x, serverRect.position.y);
            }
        }
    }
    handleGameInfoMessage(body) {
        this.playerRectID = body.rectID;
        // TODO: Handle constants and make physical calculations based on them
        console.log(body.constants);
    }
}
new Platformer();