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
        this.controlsCoeffs = new Map();
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
    resetCoeffs() {
        this.controlsCoeffs.clear();
    }
    updateCoeffs(serverTPS, clientTPS) {
        for (const [control, _] of this.heldControls) {
            let coeff = this.controlsCoeffs.get(control);
            if (coeff == undefined) {
                coeff = 0;
            }
            this.controlsCoeffs.set(control, coeff + serverTPS / clientTPS);
        }
    }
    getCoeffs() {
        return this.controlsCoeffs;
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
function lerpVector2(v1, v2, a) {
    return new Vector2(v1.x + a * (v2.x - v1.x), v1.y + a * (v2.y - v1.y));
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
function isAbstractBlock(obj) {
    return isAbstractRect(obj);
}
class Block extends Rect {
    constructor(abstract) {
        super(abstract);
    }
}
var Direction;
(function (Direction) {
    Direction[Direction["None"] = 0] = "None";
    Direction[Direction["Up"] = 1] = "Up";
    Direction[Direction["Down"] = 2] = "Down";
    Direction[Direction["Left"] = 3] = "Left";
    Direction[Direction["Right"] = 4] = "Right";
})(Direction || (Direction = {}));
class Level {
    constructor() {
        this.blocks = new Map();
        this.interpolatedPlayers = new Map();
        this.kinematicPlayers = new Map();
        this.config = {
            playerSpeed: 0,
            playerJump: 0,
            playerGravity: 0,
            playerFriction: 0,
        },
            this.playerRectID = 0;
        this.fixedTicker = new FixedTicker(20);
        this.serverAccumulator = 0;
    }
    setConfig(config) {
        this.config = config;
    }
    setPlayerRectID(id) {
        this.playerRectID = id;
    }
    createPlayerRectangle(serverRect, rectID) {
        if (this.interpolatedPlayers.has(rectID) || this.kinematicPlayers.has(rectID)) {
            return;
        }
        let rectangle;
        if (rectID == this.playerRectID) {
            const rect = new KinematicPlayer(serverRect);
            this.kinematicPlayers.set(rectID, rect);
            rectangle = new RectangleShape(rect);
        }
        else {
            const rect = new InterpolatedPlayer(serverRect);
            this.interpolatedPlayers.set(rectID, rect);
            rectangle = new RectangleShape(rect);
        }
        return rectangle;
    }
    disconnectPlayer(rectId) {
        this.interpolatedPlayers.delete(rectId);
        this.kinematicPlayers.delete(rectId);
    }
    createBlockRectangle(serverRect, rectID) {
        if (this.blocks.has(rectID)) {
            return;
        }
        const rect = new Block(serverRect);
        this.blocks.set(rectID, rect);
        const rectangle = new RectangleShape(rect);
        return rectangle;
    }
    tick(dt, serverTPS, controls) {
        this.serverAccumulator += dt;
        // interpolate kinematic players
        const kinematicAlpha = this.fixedTicker.getAlpha();
        for (const [_, player] of this.kinematicPlayers) {
            player.interpolate(kinematicAlpha);
        }
        // interpolate interpolated players
        const interpolatedAlpha = this.serverAccumulator / (1000 / serverTPS);
        for (const [_, player] of this.interpolatedPlayers) {
            player.interpolate(interpolatedAlpha);
        }
        // update kinematic players
        this.fixedTicker.update(dt, fixedDT => {
            controls.updateCoeffs(serverTPS, 1000 / fixedDT);
            for (const [_, player] of this.kinematicPlayers) {
                // update interpolation
                player.updateStartPos();
                // forces
                player.applyGravity(this.config.playerGravity, fixedDT);
                player.applyFriction(this.config.playerFriction);
                // control
                player.control(this.config.playerSpeed, this.config.playerJump, controls);
                // TODO: Collisions
                for (const [_, _block] of this.blocks) {
                    // Here
                }
                player.setCollisionDir(Direction.None);
                // Move rect
                player.applyVelToPos(fixedDT);
            }
        });
    }
    handlePlayerMovement(moved) {
        // update interpolated rects interpolation
        for (const [_, player] of this.interpolatedPlayers) {
            player.updateStartPos();
        }
        // set position
        for (const [key, val] of Object.entries(moved)) {
            const rectID = Number(key);
            const pos = val;
            const interpolated = this.interpolatedPlayers.get(rectID);
            if (interpolated) {
                interpolated.setTargetPos(pos.x, pos.y);
                continue;
            }
            const kinematic = this.kinematicPlayers.get(rectID);
            if (kinematic) {
                // TODO: Correction
                const correct = false;
                if (correct) {
                    kinematic.setTargetPos(pos.x, pos.y);
                }
            }
        }
    }
}
class Platformer {
    constructor() {
        const layers = 2;
        this.serverTPS = 0;
        this.level = new Level();
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
                case "playerMovement":
                    this.handlePlayerMovementMessage(body);
                    break;
                case "connect":
                    this.handleConnectMessage(body);
                    break;
                case "disconnect":
                    this.handleDisconnectMessage(body);
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
        // controls
        this.controls.updateCoeffs(this.serverTPS, 1000 / dt);
        // level
        this.level.tick(dt, this.serverTPS, this.controls);
        // draw
        this.canvas.draw();
    }
    stopWithText(text) {
        this.canvas.stop();
        roomGui.showMessage(text);
        roomGui.setNavBarVisibility(true);
    }
    handleLevelMessage(body) {
        this.level.setConfig(body.config);
        this.level.setPlayerRectID(body.playerRectId);
        for (const [key, val] of Object.entries(body.players)) {
            const id = Number(key);
            const serverRect = val;
            const rectangle = this.level.createPlayerRectangle(serverRect, id);
            if (rectangle) {
                this.canvas.insertDrawable(rectangle, 0, id);
            }
        }
        for (const [key, val] of Object.entries(body.blocks)) {
            const id = Number(key);
            const serverRect = val;
            const rectangle = this.level.createBlockRectangle(serverRect, id);
            if (rectangle) {
                this.canvas.insertDrawable(rectangle, 0, id);
            }
        }
    }
    handlePlayerMovementMessage(body) {
        // level
        this.level.handlePlayerMovement(body.playersMoved);
        // send controls to server
        const controlsCoeffs = this.controls.getCoeffs();
        if (controlsCoeffs.size > 0) {
            const json = JSON.stringify(Object.fromEntries(controlsCoeffs.entries()));
            this.controls.resetCoeffs();
            this.websocket.sendMessage("input", json);
        }
    }
    handleDisconnectMessage(body) {
        this.canvas.deleteDrawable(body.rectId);
        this.level.disconnectPlayer(body.rectId);
    }
    handleConnectMessage(body) {
        let serverRect = body.rect;
        let rectID = body.rectId;
        const rectangle = this.level.createPlayerRectangle(serverRect, rectID);
        if (rectangle) {
            this.canvas.insertDrawable(rectangle, 0, rectID);
        }
    }
    handleGameInfoMessage(body) {
        this.serverTPS = body.tps;
    }
}
new Platformer();
function isAbstractPlayer(obj) {
    return isAbstractRect(obj);
}
class InterpolatedPlayer extends Rect {
    constructor(abstract) {
        super(abstract);
        this.startPosition = this.getPosition();
        this.targetPosition = this.getPosition();
    }
    setTargetPos(x, y, teleport) {
        this.targetPosition.setPosition(x, y);
        if (teleport) {
            this.setPosition(x, y);
            this.startPosition.setPosition(x, y);
        }
    }
    updateStartPos() {
        this.startPosition.setPosition(this.targetPosition.x, this.targetPosition.y);
    }
    interpolate(alpha) {
        const newPos = lerpVector2(this.startPosition, this.targetPosition, alpha);
        this.setPosition(newPos.x, newPos.y);
    }
}
class KinematicPlayer extends InterpolatedPlayer {
    constructor(abstract) {
        super(abstract);
        this.velocity = new Vector2(0, 0);
        this.collisionHorizontal = Direction.None;
        this.collisionVertical = Direction.None;
    }
    control(speed, jump, controls) {
        if (controls.isHeld("left")) {
            this.velocity.x -= speed;
        }
        if (controls.isHeld("right")) {
            this.velocity.x += speed;
        }
        if (controls.isHeld("jump") && this.collisionVertical == Direction.Down) {
            this.velocity.y -= jump;
        }
    }
    applyGravity(force, dt) {
        this.velocity.y += force * dt;
    }
    applyFriction(force) {
        this.velocity.x *= force;
        // this.velocity.y *= force
    }
    applyVelToPos(dt) {
        this.targetPosition.x += this.velocity.x * dt;
        this.targetPosition.y += this.velocity.y * dt;
    }
    setCollisionDir(dir) {
        if (dir == Direction.Down || dir == Direction.Up) {
            this.collisionVertical = dir;
        }
        else if (dir == Direction.Left || dir == Direction.Right) {
            this.collisionHorizontal = dir;
        }
        else {
            this.collisionHorizontal = dir;
            this.collisionVertical = dir;
        }
    }
}
