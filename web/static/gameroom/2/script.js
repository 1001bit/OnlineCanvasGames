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
function collideKinematicWithStatic(kinematicRect, staticRect, dt) {
    if (!staticRect.canCollide) {
        return;
    }
    kinematicRect.setCollisionDir(Direction.None);
    let futureKinematic = new Rect(kinematicRect);
    futureKinematic.setPosition(kinematicRect.targetPosition.x, kinematicRect.targetPosition.y);
    const velX = kinematicRect.velocity.x * dt;
    const velY = kinematicRect.velocity.y * dt;
    if (velY > 0) {
        // down
        futureKinematic.size.y += velY;
        if (futureKinematic.intersects(staticRect)) {
            kinematicRect.setTargetPos(kinematicRect.targetPosition.x, staticRect.position.y - kinematicRect.size.y, true);
            kinematicRect.velocity.y = 0;
            kinematicRect.setCollisionDir(Direction.Down);
        }
    }
    else if (velY < 0) {
        // up
        futureKinematic.size.y += Math.abs(velY);
        futureKinematic.position.y -= Math.abs(velY);
        if (futureKinematic.intersects(staticRect)) {
            kinematicRect.setTargetPos(kinematicRect.targetPosition.x, staticRect.position.y + staticRect.size.y, true);
            kinematicRect.velocity.y = 0;
            kinematicRect.setCollisionDir(Direction.Up);
        }
    }
    futureKinematic = new Rect(kinematicRect);
    futureKinematic.setPosition(kinematicRect.targetPosition.x, kinematicRect.targetPosition.y);
    if (velX > 0) {
        // right
        futureKinematic.size.x += velX;
        if (futureKinematic.intersects(staticRect)) {
            kinematicRect.setTargetPos(staticRect.position.x - kinematicRect.size.x, kinematicRect.targetPosition.y, true);
            kinematicRect.velocity.x = 0;
            kinematicRect.setCollisionDir(Direction.Right);
        }
    }
    else if (velX < 0) {
        // left
        futureKinematic.size.x += Math.abs(velX);
        futureKinematic.position.x -= Math.abs(velX);
        if (futureKinematic.intersects(staticRect)) {
            kinematicRect.setTargetPos(staticRect.position.x + staticRect.size.x, kinematicRect.targetPosition.y, true);
            kinematicRect.velocity.x = 0;
            kinematicRect.setCollisionDir(Direction.Left);
        }
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
class PhysicsEngine {
    constructor() {
        this.staticRects = new Map();
        this.kinematicRects = new Map();
        this.interpolatedRects = new Map();
    }
    insertStaticRect(id, rect) {
        this.staticRects.set(id, rect);
    }
    insertInterpolatedRect(id, rect) {
        this.interpolatedRects.set(id, rect);
    }
    insertKinematicRect(id, rect) {
        this.kinematicRects.set(id, rect);
    }
    deleteRect(id) {
        this.staticRects.delete(id);
        this.kinematicRects.delete(id);
        this.interpolatedRects.delete(id);
    }
    update(dt, constants) {
        for (const [_id, rect] of this.kinematicRects) {
            this.applyGravityToVel(rect, constants.gravity, dt);
            this.applyFrictionToVel(rect, constants.friction);
            this.applyCollisions(rect, dt);
            this.applyVelToPos(rect, dt);
        }
    }
    updateKinematicsInterpolation() {
        for (const [_id, rect] of this.kinematicRects) {
            rect.updateStartPos();
        }
    }
    updateInterpolatedInterpolation() {
        for (const [_id, rect] of this.interpolatedRects) {
            rect.updateStartPos();
        }
    }
    interpolate(interpolatedAlpha, kinematicAlpha) {
        for (const [_id, rect] of this.kinematicRects) {
            rect.interpolate(kinematicAlpha);
        }
        for (const [_id, rect] of this.interpolatedRects) {
            rect.interpolate(interpolatedAlpha);
        }
    }
    applyGravityToVel(rect, gravity, dt) {
        if (!rect.doApplyForce(ForceType.Gravity)) {
            rect;
        }
        rect.velocity.y += gravity * dt;
    }
    applyFrictionToVel(rect, friction) {
        if (!rect.doApplyForce(ForceType.Friction)) {
            return;
        }
        rect.velocity.x -= rect.velocity.x * friction;
        // also do friction on y axis if non gravitable
        if (!rect.doApplyForce(ForceType.Gravity)) {
            rect.velocity.y -= rect.velocity.y * friction;
        }
    }
    applyCollisions(rect, dt) {
        if (!rect.canCollide) {
            return;
        }
        for (const [_id, staticRect] of this.staticRects) {
            collideKinematicWithStatic(rect, staticRect, dt);
        }
    }
    applyVelToPos(rect, dt) {
        const posX = rect.targetPosition.x + rect.velocity.x * dt;
        const posY = rect.targetPosition.y + rect.velocity.y * dt;
        rect.setTargetPos(posX, posY);
    }
    setMultiplePositions(positions) {
        for (const [key, val] of Object.entries(positions)) {
            const id = Number(key);
            const position = val;
            const staticRect = this.staticRects.get(id);
            if (staticRect) {
                staticRect.setPosition(position.x, position.y);
                continue;
            }
            const kinematicRect = this.kinematicRects.get(id);
            if (kinematicRect) {
                // TODO: Correct sometimes
                const correct = false;
                if (correct) {
                    kinematicRect.setTargetPos(position.x, position.y);
                }
                continue;
            }
            const interpolatedRect = this.interpolatedRects.get(id);
            if (interpolatedRect) {
                interpolatedRect.setTargetPos(position.x, position.y);
            }
        }
    }
}
var ForceType;
(function (ForceType) {
    ForceType["Friction"] = "friction";
    ForceType["Gravity"] = "gravity";
})(ForceType || (ForceType = {}));
function isAbstractPhysicalRect(obj) {
    return isAbstractRect(obj) && "canCollide" in obj;
}
class PhysicalRect extends Rect {
    constructor(abstractRect) {
        super(abstractRect);
        this.canCollide = abstractRect.canCollide;
    }
}
/// <reference path="physicalrect.ts"/>
class InterpolatedRect extends PhysicalRect {
    constructor(abstractRect) {
        super(abstractRect);
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
/// <reference path="interpolatedrect.ts"/>
function isAbstractKinematicRect(obj) {
    return isAbstractPhysicalRect(obj) && "velocity" in obj;
}
class KinematicRect extends InterpolatedRect {
    constructor(abstractRect) {
        super(abstractRect);
        this.velocity = new Vector2(abstractRect.velocity.x, abstractRect.velocity.y);
        this.forcesToApply = new Set();
        for (const [key, _val] of Object.entries(abstractRect.forcesToApply)) {
            this.forcesToApply.add(key);
        }
        this.collisionHorizontal = Direction.None;
        this.collisionVertical = Direction.None;
    }
    setVelocity(x, y) {
        this.velocity.setPosition(x, y);
    }
    setCollisionDir(dir) {
        if (dir == Direction.Down || dir == Direction.Up) {
            this.collisionVertical = dir;
        }
        else if (dir == Direction.Left || dir == Direction.Right) {
            this.collisionHorizontal = dir;
        }
        else {
            this.collisionVertical = dir;
            this.collisionHorizontal = dir;
        }
    }
    doApplyForce(force) {
        return this.forcesToApply.has(force);
    }
}
class Platformer {
    constructor() {
        const layers = 2;
        this.playerRectID = 0;
        this.constants = {
            physics: {
                friction: 0,
                gravity: 0,
            },
            playerSpeed: 0,
            playerJump: 0,
        };
        this.physTps = 30;
        this.physTicker = new FixedTicker(this.physTps);
        this.serverAccumulator = 0;
        this.serverTPS = 0;
        this.canvas = new GameCanvas("canvas", layers);
        this.canvas.setBackgroundColor(RGB(30, 100, 100));
        this.controls = new Controls();
        this.bindControls();
        this.websocket = new GameWebSocket();
        const gameID = $("main").data("game-id");
        const roomID = $("main").data("room-id");
        this.initWebsocket(gameID, roomID);
        this.physicsEngine = new PhysicsEngine();
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
                case "update":
                    this.handleUpdateMessage(body);
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
        // physics
        this.physTicker.update(dt, (fixedDT) => {
            // update phys/server tps coeffs
            this.controls.updateCoeffs(this.serverTPS, this.physTps);
            this.physicsEngine.updateKinematicsInterpolation();
            this.handleControls();
            this.physicsEngine.update(fixedDT, this.constants.physics);
        });
        // interpolation
        this.serverAccumulator += dt;
        const interpolatedAlpha = Math.min(1, this.serverAccumulator / (1000 / this.serverTPS));
        const kinematicAlpha = this.physTicker.getAlpha();
        this.physicsEngine.interpolate(interpolatedAlpha, kinematicAlpha);
        // draw
        this.canvas.draw();
    }
    handleControls() {
        const playerRect = this.physicsEngine.kinematicRects.get(this.playerRectID);
        if (!playerRect) {
            return;
        }
        if (this.controls.isHeld("left")) {
            playerRect.velocity.x -= this.constants.playerSpeed;
        }
        if (this.controls.isHeld("right")) {
            playerRect.velocity.x += this.constants.playerSpeed;
        }
        if (this.controls.isHeld("jump") && playerRect.collisionVertical == Direction.Down) {
            playerRect.velocity.y -= this.constants.playerJump;
        }
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
        if (isAbstractKinematicRect(serverRect) && rectID == this.playerRectID) {
            // Doing physics prediction only for player rect
            const rect = new KinematicRect(serverRect);
            this.physicsEngine.insertKinematicRect(rectID, rect);
            rectangle = new RectangleShape(rect);
        }
        else if (isAbstractKinematicRect(serverRect)) {
            // Interpolated rect for other kinematic rects
            const rect = new InterpolatedRect(serverRect);
            this.physicsEngine.insertInterpolatedRect(rectID, rect);
            rectangle = new RectangleShape(rect);
        }
        else if (isAbstractPhysicalRect(serverRect)) {
            // Static rects
            const rect = new PhysicalRect(serverRect);
            this.physicsEngine.insertStaticRect(rectID, rect);
            rectangle = new RectangleShape(rect);
        }
        else {
            // non rects
            return;
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
    handleUpdateMessage(body) {
        // physics operations
        this.serverAccumulator = 0;
        this.physicsEngine.updateInterpolatedInterpolation();
        this.physicsEngine.setMultiplePositions(body.rectsMoved);
        // send controls to server
        const controlsCoeffs = this.controls.getCoeffs();
        if (controlsCoeffs.size > 0) {
            const json = JSON.stringify(Object.fromEntries(controlsCoeffs.entries()));
            this.controls.resetCoeffs();
            this.websocket.sendMessage("input", json);
        }
    }
    handleDeleteMessage(body) {
        this.canvas.deleteDrawable(body.id);
        this.physicsEngine.deleteRect(body.id);
    }
    handleCreateMessage(body) {
        let serverRect = body.rect;
        let rectID = body.id;
        this.createRectangleShape(serverRect, rectID);
    }
    handleGameInfoMessage(body) {
        this.playerRectID = body.rectID;
        this.serverTPS = body.tps;
        this.constants = body.constants;
    }
}
new Platformer();
