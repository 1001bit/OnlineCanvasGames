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
class Drawable {
    constructor() { }
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
function lerp(a, b, alpha) {
    return a + alpha * (b - a);
}
class Rect {
    constructor(rect) {
        if (rect) {
            this.position = rect.position;
            this.size = rect.size;
            return;
        }
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
class RectangleShape extends Drawable {
    constructor(rect) {
        super();
        if (rect) {
            this.rect = rect;
        }
        else {
            this.rect = new Rect();
        }
        this.color = RGB(255, 255, 255);
    }
    setSize(x, y) {
        this.rect.setSize(x, y);
    }
    setPosition(x, y) {
        this.rect.setPosition(x, y);
    }
    setColor(color) {
        this.color = color;
    }
    draw(ctx) {
        let pos = this.rect.position;
        let size = this.rect.size;
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
        this.position = new Vector2(0, 0);
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
    setPosition(x, y) {
        this.position.setPosition(x, y);
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
class Physics {
    constructor() {
        this.constants = {
            friction: 0,
            gravity: 0,
            playerSpeed: 0,
            playerJump: 0
        };
        this.staticRects = new Map();
        this.kinematicRects = new Map();
    }
    setConstants(constants) {
        this.constants = constants;
    }
    insertKinematicRect(id, rect) {
        this.kinematicRects.set(id, rect);
    }
    insertStaticRect(id, rect) {
        this.staticRects.set(id, rect);
    }
    deleteRect(id) {
        this.staticRects.delete(id);
        this.kinematicRects.delete(id);
    }
    getRect(id) {
        const rect = this.kinematicRects.get(id);
        if (rect) {
            return rect;
        }
        return this.staticRects.get(id);
    }
    tick(dt) {
        for (const [_id, rect] of this.kinematicRects) {
            this.applyVelToPos(rect, dt);
        }
    }
    applyVelToPos(rect, dt) {
        rect.position.x += rect.velocity.x * dt;
        rect.position.y += rect.velocity.y * dt;
    }
}
/// <reference path="physicsEngine.ts"/>
class Platformer {
    constructor() {
        const layers = 2;
        this.playerRectID = 0;
        this.controlsAccumulator = 0;
        this.serverTPS = 0;
        this.canvas = new GameCanvas("canvas", layers);
        this.canvas.setBackgroundColor(RGB(30, 100, 100));
        this.controls = new Controls();
        this.bindControls();
        this.websocket = new GameWebSocket();
        const gameID = $("main").data("game-id");
        const roomID = $("main").data("room-id");
        this.initWebsocket(gameID, roomID);
        this.physicsEngine = new Physics();
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
    tick(dt) {
        this.physicsEngine.tick(dt);
        this.controlsAccumulator += dt;
        const maxControlsAccumulator = 1000 / (this.serverTPS * 4);
        while (this.controlsAccumulator > maxControlsAccumulator) {
            let heldControls = this.controls.getHeldControls();
            if (heldControls.size > 0) {
                let json = JSON.stringify(Object.fromEntries(heldControls.entries()));
                this.websocket.sendMessage("input", json);
            }
            this.controlsAccumulator -= maxControlsAccumulator;
        }
        this.canvas.draw();
    }
    stopWithText(text) {
        this.canvas.stop();
        roomGui.showMessage(text);
        roomGui.setNavBarVisibility(true);
    }
    createRectangleShape(serverRect, rectID) {
        if (this.canvas.getDrawable(rectID)) {
            return;
        }
        let rectangle;
        if (isKinematicRect(serverRect)) {
            const rect = new KinematicRect(serverRect);
            rect.setVelocity(serverRect.velocity.x, serverRect.velocity.y);
            this.physicsEngine.insertKinematicRect(rectID, rect);
            rectangle = new RectangleShape(rect);
        }
        else {
            const rect = new Rect(serverRect);
            this.physicsEngine.insertStaticRect(rectID, rect);
            rectangle = new RectangleShape(rect);
        }
        this.canvas.insertDrawable(rectangle, 0, rectID);
    }
    handleLevelMessage(body) {
        for (const [key, val] of Object.entries(body.kinematic)) {
            const id = Number(key);
            const serverRect = val;
            this.createRectangleShape(serverRect, id);
        }
        for (const [key, val] of Object.entries(body.static)) {
            const id = Number(key);
            const serverRect = val;
            this.createRectangleShape(serverRect, id);
        }
    }
    handleDeleteMessage(body) {
        this.canvas.deleteDrawable(body.ID);
        this.physicsEngine.deleteRect(body.ID);
    }
    handleCreateMessage(body) {
        let serverRect = body.rect;
        let rectID = body.id;
        this.createRectangleShape(serverRect, rectID);
    }
    handleGameInfoMessage(body) {
        this.playerRectID = body.rectID;
        this.serverTPS = body.tps;
        this.physicsEngine.setConstants(body.constants);
    }
}
new Platformer();
function isRect(obj) {
    return "position" in obj && "size" in obj;
}
function isKinematicRect(obj) {
    return isRect(obj) && "velocity" in obj;
}
class KinematicRect extends Rect {
    constructor(rect) {
        super(rect);
        this.velocity = new Vector2(0, 0);
    }
    setVelocity(x, y) {
        this.velocity.setPosition(x, y);
    }
}
