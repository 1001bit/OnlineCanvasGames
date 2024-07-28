"use strict";
class GameCanvas {
    constructor(canvasID, layersCount) {
        this.onMouseClick = (_e) => { };
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
        this.canvas.addEventListener("click", e => {
            this.onMouseClick(e);
        });
    }
    resize() {
        const canvas = this.canvas;
        canvas.width = window.innerWidth;
        canvas.height = window.innerHeight - canvas.getBoundingClientRect().top;
        this.draw();
    }
    clear() {
        const ctx = this.ctx;
        const canvas = this.canvas;
        ctx.clearRect(0, 0, canvas.width, canvas.height);
        ctx.fillStyle = this.backgroundColor;
        ctx.fillRect(0, 0, canvas.width, canvas.height);
    }
    updateMousePos(e) {
        let rect = this.canvas.getBoundingClientRect();
        let x = e.clientX - rect.left;
        let y = e.clientY - rect.top;
        this.mousePos.setPosition(x, y);
    }
    stop() {
        this.canvas.remove();
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
    draw(viewCenter) {
        const ctx = this.ctx;
        ctx.save();
        this.clear();
        if (viewCenter) {
            ctx.translate(this.canvas.width / 2 - viewCenter.x, this.canvas.height / 2 - viewCenter.y);
        }
        this.layers.forEach(layer => {
            layer.forEach(drawable => {
                drawable.draw(ctx);
            });
        });
        ctx.restore();
    }
    setBackgroundColor(color) {
        this.backgroundColor = color;
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
        this.heldControls = new Map();
        this.bindings = new Map();
        this.heldControlsTicks = new Map();
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
    addTick(control) {
        if (!this.heldControls.has(control)) {
            return;
        }
        const ticks = this.heldControlsTicks.get(control);
        if (!ticks) {
            this.heldControlsTicks.set(control, 1);
            return;
        }
        this.heldControlsTicks.set(control, ticks + 1);
    }
    resetHeldControlsTicks(serverTPS, clientTPS) {
        const maxTicks = Math.ceil(clientTPS / serverTPS);
        for (const [control, ticks] of this.heldControlsTicks) {
            if (ticks <= maxTicks) {
                // delete controls, that didn't bypass the limit
                this.heldControlsTicks.delete(control);
                continue;
            }
            // postpone ticks, that are beyond for the future, since can't send any more.
            this.heldControlsTicks.set(control, ticks - maxTicks);
        }
    }
    getHeldControlsTicks() {
        return this.heldControlsTicks;
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
        $("#show-nav").on("click", () => this.onclick());
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
function isAbstractRect(obj) {
    return "position" in obj && "size" in obj;
}
class Rect {
    constructor(abstractRect) {
        this.position = new Vector2(0, 0);
        this.size = new Vector2(0, 0);
        if (abstractRect) {
            this.setPosition(abstractRect.position.x, abstractRect.position.y);
            this.setSize(abstractRect.size.x, abstractRect.size.y);
        }
    }
    setPosition(x, y) {
        this.position.setPosition(x, y);
    }
    setSize(x, y) {
        this.size.setPosition(x, y);
    }
    extend(extX, extY) {
        this.size.x += Math.abs(extX);
        this.size.y += Math.abs(extY);
        if (extX < 0) {
            this.position.x -= Math.abs(extX);
        }
        if (extY < 0) {
            this.position.y -= Math.abs(extY);
        }
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
    intersects(rect) {
        if (this.position.x + this.size.x <= rect.position.x ||
            this.position.x >= rect.position.x + rect.size.x ||
            this.position.y + this.size.y <= rect.position.y ||
            this.position.y >= rect.position.y + rect.size.y) {
            return false;
        }
        return true;
    }
    getPosition() {
        return new Vector2(this.position.x, this.position.y);
    }
    getSize() {
        return new Vector2(this.size.x, this.size.y);
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
        const pos = this.rect.getPosition();
        const size = this.rect.getSize();
        ctx.fillStyle = this.color;
        ctx.fillRect(pos.x, pos.y, size.x, size.y);
    }
    getSize() {
        return this.rect.getSize();
    }
    getPosition() {
        return this.rect.getPosition();
    }
    getRect() {
        return this.rect;
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
        this.previousTime = 0;
    }
    tick(callback, time) {
        const dt = time - this.previousTime;
        this.previousTime = time;
        callback(dt);
        requestAnimationFrame((time) => {
            this.tick(callback, time);
        });
    }
    start(callback) {
        requestAnimationFrame((time) => {
            this.tick(callback, time);
        });
    }
}
class FixedTicker {
    constructor(tps) {
        this.tps = tps;
        this.accumulator = 0;
    }
    update(dt, callback) {
        this.accumulator += dt;
        const maxAccumulator = 1000 / this.tps;
        while (this.accumulator >= maxAccumulator) {
            callback(maxAccumulator);
            this.accumulator -= maxAccumulator;
        }
    }
    setTPS(tps) {
        this.tps = tps;
    }
    getAlpha() {
        return this.accumulator / (1000 / this.tps);
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
    interpolateBetween(v1, v2, a) {
        this.x = lerp(v1.x, v2.x, a);
        this.y = lerp(v1.y, v2.y, a);
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
class Clicker {
    constructor() {
        const layers = 2;
        this.clicks = 0;
        this.canvas = new GameCanvas("canvas", layers);
        this.canvas.setBackgroundColor(RGB(60, 70, 70));
        this.websocket = new GameWebSocket();
        const gameID = $("main").data("game-id");
        const roomID = $("main").data("room-id");
        this.initWebsocket(gameID, roomID);
        this.drawables = new Map();
        this.initDrawables();
        this.ticker = new Ticker();
        this.ticker.start(dt => this.tick(dt));
    }
    tick(_dt) {
        this.canvas.draw();
    }
    initWebsocket(gameID, roomID) {
        this.websocket.handleMessage = (type, body) => {
            switch (type) {
                case "clicks":
                    this.click(Number(body));
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
    stopWithText(text) {
        this.canvas.stop();
        roomGui.showMessage(text);
        roomGui.setNavBarVisibility(true);
    }
    initDrawables() {
        const button = new RectangleShape();
        this.drawables.set("button", button);
        button.setColor(RGB(150, 150, 40));
        button.setSize(300, 200);
        button.setPosition((window.innerWidth - button.getSize().x) / 2, (window.innerHeight - button.getSize().y) / 2);
        this.canvas.insertDrawable(button, 0, 0);
        const text = new DrawableText("0 clicks", 48);
        this.drawables.set("text", text);
        text.setPosition(button.getPosition().x + 10, button.getPosition().y + 10);
        this.canvas.insertDrawable(text, 1, 1);
        // button click
        this.canvas.onMouseClick = (_e) => {
            let mPos = this.canvas.getMousePos();
            if (button.getRect().containsPoint(mPos.x, mPos.y)) {
                this.click(this.clicks + 1);
                this.websocket.sendMessage("click", "");
            }
        };
    }
    click(clicks) {
        this.clicks = clicks;
        const text = this.drawables.get("text");
        text.setString(`${this.clicks} clicks`);
    }
}
new Clicker();
