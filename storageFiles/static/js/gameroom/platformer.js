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
    openConnection(gameTitle, roomID) {
        const protocol = location.protocol == "https:" ? "wss:" : "ws:";
        this.websocket = new WebSocket(`${protocol}//${document.location.host}/rt/ws/game/${gameTitle}/room/${roomID}`);
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
class SmoothCamera {
    constructor() {
        this.position = new Vector2(0, 0);
        this.target = undefined;
        this.strength = 0.96;
    }
    setTarget(target) {
        this.target = target;
    }
    update(dt) {
        if (this.target == undefined) {
            return;
        }
        const centerPos = new Vector2(this.target.getPosition().x + this.target.getSize().x / 2, this.target.getPosition().y + this.target.getSize().y / 2);
        this.position.interpolateBetween(this.position, centerPos, Math.pow(this.strength, dt));
    }
    getPosition() {
        return this.position;
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
        this.camera = new SmoothCamera();
        this.config = {
            playerSpeed: 0,
            playerJump: 0,
            playerGravity: 0,
            playerFriction: 0,
        },
            this.playerRectID = 0;
        this.fixedTicker = new FixedTicker(10);
        this.serverTPS = 0;
        this.serverAccumulator = 0;
    }
    setConfig(config) {
        this.config = config;
    }
    setPlayerRectID(id) {
        this.playerRectID = id;
    }
    setTPS(client, server) {
        this.serverTPS = server;
        this.fixedTicker.setTPS(client);
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
            // camera
            this.camera.setTarget(rect);
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
    tick(dt, controls) {
        // interpolate interpolated players
        this.serverAccumulator += dt;
        const interpolatedAlpha = Math.min(this.serverAccumulator / (1000 / this.serverTPS), 1);
        for (const [_, player] of this.interpolatedPlayers) {
            player.interpolate(interpolatedAlpha);
        }
        // update kinematic players
        this.fixedTicker.update(dt, fixedDT => {
            for (const [rectID, player] of this.kinematicPlayers) {
                // update interpolation
                player.updateStartPos();
                // Control
                if (rectID == this.playerRectID) {
                    player.control(this.config.playerSpeed, this.config.playerJump, controls);
                }
                // Forces
                player.applyGravity(this.config.playerGravity, fixedDT);
                player.applyFriction(this.config.playerFriction);
                // Collisions and movement
                player.setCollisionDir(Direction.None);
                // Horizontal
                for (const [_, block] of this.blocks) {
                    const dir = player.detectHorizontalCollision(block, fixedDT);
                    if (dir != Direction.None) {
                        player.resolveCollision(block, dir);
                        break;
                    }
                }
                player.applyHorizontalVelToPos(fixedDT);
                // Vertical
                for (const [_, block] of this.blocks) {
                    const dir = player.detectVerticalCollision(block, fixedDT);
                    if (dir != Direction.None) {
                        player.resolveCollision(block, dir);
                        break;
                    }
                }
                player.applyVerticalVelToPos(fixedDT);
            }
        });
        // interpolate kinematic players
        const kinematicAlpha = this.fixedTicker.getAlpha();
        for (const [_, player] of this.kinematicPlayers) {
            player.interpolate(kinematicAlpha);
        }
        // camera follow
        this.camera.update(dt);
    }
    handlePlayerMovement(moved, correct) {
        // update interpolated rects interpolation
        this.serverAccumulator = 0;
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
                if (correct) {
                    kinematic.correctDivergence(pos.x, pos.y);
                }
            }
        }
    }
    getCameraPosition() {
        return this.camera.getPosition();
    }
}
class Platformer {
    constructor() {
        const layers = 2;
        this.level = new Level();
        this.canvas = new GameCanvas("canvas", layers);
        this.canvas.setBackgroundColor(RGB(30, 100, 100));
        this.controls = new Controls();
        this.bindControls();
        this.websocket = new GameWebSocket();
        const gameTitle = $("main").data("game-title");
        const roomID = $("main").data("room-id");
        this.initWebsocket(gameTitle, roomID);
        this.serverTPS = 0;
        this.clientTPS = 0;
        this.ticker = new Ticker();
        this.ticker.start(dt => this.tick(dt));
    }
    bindControls() {
        const controls = this.controls;
        controls.bindControl("d", "right");
        controls.bindControl("a", "left");
        controls.bindControl("w", "jump");
        controls.bindControl(" ", "jump");
    }
    initWebsocket(gameTitle, roomID) {
        this.websocket.handleMessage = (type, body) => {
            switch (type) {
                case "level":
                    this.handleLevelMessage(body);
                    break;
                case "levelUpdate":
                    this.handleLevelUpdateMessage(body);
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
        this.websocket.openConnection(gameTitle, roomID);
    }
    tick(dt) {
        // draw
        this.canvas.draw(this.level.getCameraPosition());
        // level
        this.level.tick(dt, this.controls);
    }
    stopWithText(text) {
        this.canvas.stop();
        roomGui.showMessage(text);
        roomGui.setNavBarVisibility(true);
    }
    handleLevelMessage(body) {
        this.level.setConfig(body.config);
        this.level.setPlayerRectID(body.playerRectId);
        this.level.setTPS(body.clientTps, body.serverTps);
        this.serverTPS = body.serverTps;
        this.clientTPS = body.clientTps;
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
    handleLevelUpdateMessage(body) {
        this.level.handlePlayerMovement(body.players, body.correct);
        // send controls right after level message, because server allows sending messages right after sending level message
        const heldControlsTicks = this.controls.getHeldControlsTicks();
        if (heldControlsTicks.size != 0) {
            // no need of cutting ticks in map, that is being sent to server, since ticks are being limited there
            const json = JSON.stringify(Object.fromEntries(heldControlsTicks));
            this.websocket.sendMessage("input", json);
            // cutting ticks after sending
            this.controls.resetHeldControlsTicks(this.serverTPS, this.clientTPS);
        }
    }
    handleConnectMessage(body) {
        let serverRect = body.rect;
        let rectID = body.rectId;
        const rectangle = this.level.createPlayerRectangle(serverRect, rectID);
        if (rectangle) {
            this.canvas.insertDrawable(rectangle, 0, rectID);
        }
    }
    handleDisconnectMessage(body) {
        this.canvas.deleteDrawable(body.rectId);
        this.level.disconnectPlayer(body.rectId);
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
    moveTargetPos(x, y) {
        this.targetPosition.x += x;
        this.targetPosition.y += y;
    }
    updateStartPos() {
        this.startPosition.setPosition(this.targetPosition.x, this.targetPosition.y);
    }
    interpolate(alpha) {
        this.position.interpolateBetween(this.startPosition, this.targetPosition, alpha);
    }
}
class KinematicPlayer extends InterpolatedPlayer {
    constructor(abstract) {
        super(abstract);
        this.velocity = new Vector2(0, 0);
        this.collisionHorizontal = Direction.None;
        this.collisionVertical = Direction.None;
        this.futurePath = new Rect();
    }
    control(speed, jump, controls) {
        if (controls.isHeld("left")) {
            this.velocity.x -= speed;
            controls.addTick("left");
        }
        if (controls.isHeld("right")) {
            this.velocity.x += speed;
            controls.addTick("right");
        }
        if (controls.isHeld("jump") && this.isCollisionInDirection(Direction.Down)) {
            this.velocity.y -= jump;
            controls.addTick("jump");
        }
    }
    applyGravity(force, dt) {
        this.velocity.y += force * dt;
    }
    applyFriction(force) {
        this.velocity.x *= force;
        // this.velocity.y *= force
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
    applyHorizontalVelToPos(dt) {
        this.targetPosition.x += this.velocity.x * dt;
    }
    applyVerticalVelToPos(dt) {
        this.targetPosition.y += this.velocity.y * dt;
    }
    detectHorizontalCollision(block, dtMs) {
        if (this.velocity.x == 0) {
            return Direction.None;
        }
        this.futurePath.setPosition(this.targetPosition.x, this.targetPosition.y);
        this.futurePath.setSize(this.size.x, this.size.y);
        this.futurePath.extend(this.velocity.x * dtMs, 0);
        if (!this.futurePath.intersects(block)) {
            return Direction.None;
        }
        if (this.velocity.x > 0) {
            return Direction.Right;
        }
        else {
            return Direction.Left;
        }
    }
    detectVerticalCollision(block, dtMs) {
        if (this.velocity.y == 0) {
            return Direction.None;
        }
        this.futurePath.setPosition(this.targetPosition.x, this.targetPosition.y);
        this.futurePath.setSize(this.size.x, this.size.y);
        this.futurePath.extend(0, this.velocity.y * dtMs);
        if (!this.futurePath.intersects(block)) {
            return Direction.None;
        }
        if (this.velocity.y > 0) {
            return Direction.Down;
        }
        else {
            return Direction.Up;
        }
    }
    resolveCollision(block, dir) {
        if (dir == Direction.None) {
            return;
        }
        this.setCollisionDir(dir);
        switch (dir) {
            case Direction.Up:
                this.velocity.y = 0;
                this.targetPosition.y = block.getPosition().y + block.getSize().y;
                break;
            case Direction.Down:
                this.velocity.y = 0;
                this.targetPosition.y = block.getPosition().y - this.getSize().y;
                break;
            case Direction.Left:
                this.velocity.x = 0;
                this.targetPosition.x = block.getPosition().x + block.getSize().x;
                break;
            case Direction.Right:
                this.velocity.x = 0;
                this.targetPosition.x = block.getPosition().x - this.getSize().x;
                break;
        }
    }
    correctDivergence(posX, posY) {
        const divergenceTolerance = 30;
        const distX = Math.abs(posX - this.targetPosition.x);
        if (distX >= divergenceTolerance && Math.abs(this.velocity.x) < 0.1) {
            this.targetPosition.x = posX;
        }
        const distY = Math.abs(posY - this.targetPosition.y);
        if (distY >= divergenceTolerance && Math.abs(this.velocity.y) < 0.1) {
            this.targetPosition.y = posY;
        }
    }
    isCollisionInDirection(dir) {
        if (dir == Direction.None) {
            return false;
        }
        return this.collisionHorizontal == dir || this.collisionVertical == dir;
    }
}
