"use strict";
var __extends = (this && this.__extends) || (function () {
    var extendStatics = function (d, b) {
        extendStatics = Object.setPrototypeOf ||
            ({ __proto__: [] } instanceof Array && function (d, b) { d.__proto__ = b; }) ||
            function (d, b) { for (var p in b) if (Object.prototype.hasOwnProperty.call(b, p)) d[p] = b[p]; };
        return extendStatics(d, b);
    };
    return function (d, b) {
        if (typeof b !== "function" && b !== null)
            throw new TypeError("Class extends value " + String(b) + " is not a constructor or null");
        extendStatics(d, b);
        function __() { this.constructor = d; }
        d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
    };
})();
var GameCanvas = /** @class */ (function () {
    function GameCanvas(canvasID, layersCount) {
        var _this = this;
        this.canvas = document.getElementById(canvasID);
        var ctx = this.canvas.getContext("2d");
        if (!ctx) {
            throw new Error("Failed to get context");
        }
        this.ctx = ctx;
        this.layers = new Array();
        for (var i = 0; i < layersCount; i++) {
            this.layers.push(new Map());
        }
        this.drawables = new Map();
        this.mousePos = new Vector2(0, 0);
        this.backgroundColor = RGB(0, 0, 0);
        this.resize();
        window.addEventListener('resize', function () { return _this.resize(); }, false);
        this.canvas.addEventListener("mousemove", function (e) {
            _this.updateMousePos(e);
        });
    }
    GameCanvas.prototype.stop = function () {
        this.canvas.remove();
    };
    GameCanvas.prototype.resize = function () {
        var canvas = this.canvas;
        canvas.width = window.innerWidth;
        canvas.height = window.innerHeight - canvas.getBoundingClientRect().top;
        this.draw();
    };
    GameCanvas.prototype.insertDrawable = function (drawable, layerNum, id) {
        this.drawables.set(id, drawable);
        var layer = this.layers[layerNum];
        if (layer) {
            layer.set(id, drawable);
        }
    };
    GameCanvas.prototype.deleteDrawable = function (id) {
        this.drawables.delete(id);
        this.layers.forEach(function (layer) {
            layer.delete(id);
        });
    };
    GameCanvas.prototype.draw = function () {
        var ctx = this.ctx;
        this.clear();
        this.layers.forEach(function (layer) {
            layer.forEach(function (drawable) {
                drawable.draw(ctx);
            });
        });
    };
    GameCanvas.prototype.clear = function () {
        var ctx = this.ctx;
        var canvas = this.canvas;
        ctx.clearRect(0, 0, canvas.width, canvas.height);
        ctx.fillStyle = this.backgroundColor;
        ctx.fillRect(0, 0, canvas.width, canvas.height);
    };
    GameCanvas.prototype.setBackgroundColor = function (color) {
        this.backgroundColor = color;
    };
    GameCanvas.prototype.updateMousePos = function (e) {
        var rect = this.canvas.getBoundingClientRect();
        var x = e.clientX - rect.left;
        var y = e.clientY - rect.top;
        this.mousePos.setPosition(x, y);
    };
    GameCanvas.prototype.getMousePos = function () {
        return this.mousePos;
    };
    GameCanvas.prototype.drawableExists = function (id) {
        return this.drawables.has(id);
    };
    GameCanvas.prototype.getDrawable = function (id) {
        return this.drawables.get(id);
    };
    return GameCanvas;
}());
var Controls = /** @class */ (function () {
    function Controls() {
        var _this = this;
        // using map instead of set here because golang doesn't have set implementation yet
        this.heldControls = new Map();
        this.bindings = new Map();
        // on key press
        document.addEventListener("keypress", function (e) {
            // only single press
            if (e.repeat) {
                return;
            }
            // if no key in bindings
            if (!_this.bindings.has(e.key)) {
                return;
            }
            // get control from binding
            var control = _this.bindings.get(e.key);
            if (control) {
                _this.heldControls.set(control, true);
            }
        });
        document.addEventListener("keyup", function (e) {
            if (!_this.bindings.has(e.key)) {
                return;
            }
            var control = _this.bindings.get(e.key);
            if (control) {
                _this.heldControls.delete(control);
            }
        });
    }
    Controls.prototype.bindControl = function (key, control) {
        this.bindings.set(key, control);
    };
    Controls.prototype.isHeld = function (control) {
        return this.heldControls.has(control);
    };
    Controls.prototype.getHeldControls = function () {
        return this.heldControls;
    };
    return Controls;
}());
var DeltaTimer = /** @class */ (function () {
    function DeltaTimer() {
        this.lastTick = performance.now();
    }
    DeltaTimer.prototype.getDeltaTime = function () {
        var now = performance.now();
        var dt = now - this.lastTick;
        this.lastTick = now;
        return dt;
    };
    return DeltaTimer;
}());
var Drawable = /** @class */ (function () {
    function Drawable() {
    }
    Drawable.prototype.draw = function (_ctx) {
    };
    return Drawable;
}());
var Gui = /** @class */ (function () {
    function Gui() {
        var _this = this;
        this.navVisible = true;
        this.onclick();
        $("#show-nav").on("click", function (_e) {
            _this.onclick();
        });
    }
    Gui.prototype.showMessage = function (text) {
        $("#message").text(text);
    };
    Gui.prototype.setNavBarVisibility = function (visibility) {
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
    };
    Gui.prototype.onclick = function () {
        this.navVisible = !this.navVisible;
        this.setNavBarVisibility(this.navVisible);
    };
    return Gui;
}());
function lerp(a, b, alpha) {
    return a + alpha * (b - a);
}
var Rect = /** @class */ (function () {
    function Rect() {
        this.position = new Vector2(0, 0);
        this.size = new Vector2(0, 0);
    }
    Rect.prototype.setPosition = function (x, y) {
        this.position.setPosition(x, y);
    };
    Rect.prototype.setSize = function (x, y) {
        this.size.setPosition(x, y);
    };
    Rect.prototype.containsPoint = function (x, y) {
        var pos = this.position;
        var size = this.size;
        if (x >= pos.x && x <= pos.x + size.x &&
            y >= pos.y && y <= pos.y + size.y) {
            return true;
        }
        return false;
    };
    Rect.prototype.getPosition = function () {
        return this.position;
    };
    Rect.prototype.getSize = function () {
        return this.size;
    };
    return Rect;
}());
var RectangleShape = /** @class */ (function (_super) {
    __extends(RectangleShape, _super);
    function RectangleShape(width, height) {
        var _this = _super.call(this) || this;
        _this.rect = new Rect();
        _this.rect.setSize(width, height);
        _this.rect.setPosition(0, 0);
        _this.color = RGB(255, 255, 255);
        return _this;
    }
    RectangleShape.prototype.setSize = function (width, height) {
        this.rect.setSize(width, height);
    };
    RectangleShape.prototype.setPosition = function (x, y) {
        this.rect.setPosition(x, y);
    };
    RectangleShape.prototype.setColor = function (color) {
        this.color = color;
    };
    RectangleShape.prototype.draw = function (ctx) {
        var pos = this.rect.position;
        var size = this.rect.size;
        ctx.fillStyle = this.color;
        ctx.fillRect(pos.x, pos.y, size.x, size.y);
    };
    return RectangleShape;
}(Drawable));
function RGB(r, g, b) {
    return "rgb(".concat(r, " ").concat(g, " ").concat(b, ")");
}
var DrawableText = /** @class */ (function (_super) {
    __extends(DrawableText, _super);
    function DrawableText(string, size) {
        var _this = _super.call(this) || this;
        _this.string = string;
        _this.color = RGB(255, 255, 255);
        _this.font = "serif";
        _this.size = size;
        _this.position = new Vector2(0, 0);
        return _this;
    }
    DrawableText.prototype.setPosition = function (x, y) {
        this.position.setPosition(x, y);
    };
    DrawableText.prototype.setString = function (string) {
        this.string = string;
    };
    DrawableText.prototype.setColor = function (color) {
        this.color = color;
    };
    DrawableText.prototype.setFont = function (font) {
        this.font = font;
    };
    DrawableText.prototype.setSize = function (size) {
        this.size = size;
    };
    DrawableText.prototype.draw = function (ctx) {
        ctx.fillStyle = this.color;
        ctx.font = "".concat(this.size, "px ").concat(this.font);
        // adding size to y because text's origin is located on the bottom
        ctx.fillText(this.string, this.position.x, this.position.y + this.size);
    };
    return DrawableText;
}(Drawable));
var Ticker = /** @class */ (function () {
    function Ticker() {
        this.timer = new DeltaTimer();
    }
    Ticker.prototype.tick = function (callback) {
        var _this = this;
        var dt = this.timer.getDeltaTime();
        callback(dt);
        requestAnimationFrame(function () { return _this.tick(callback); });
    };
    return Ticker;
}());
var Vector2 = /** @class */ (function () {
    function Vector2(x, y) {
        this.x = x;
        this.y = y;
    }
    Vector2.prototype.setPosition = function (x, y) {
        this.x = x;
        this.y = y;
    };
    return Vector2;
}());
var GameWebSocket = /** @class */ (function () {
    function GameWebSocket() {
        this.handleClose = function (_body) { };
        this.handleMessage = function (_type, _body) { };
        this.active = false;
        this.websocket = null;
    }
    GameWebSocket.prototype.openConnection = function (gameID, roomID) {
        var _this = this;
        var protocol = location.protocol == "https:" ? "wss:" : "ws:";
        this.websocket = new WebSocket("".concat(protocol, "//").concat(document.location.host, "/rt/ws/game/").concat(gameID, "/room/").concat(roomID));
        var ws = this.websocket;
        ws.onopen = function (_e) {
            _this.active = true;
        };
        ws.onclose = function (_e) {
            if (!_this.active) {
                return;
            }
            _this.handleClose("Connection closed");
            _this.active = false;
        };
        ws.onerror = function (_e) {
            if (!_this.active) {
                return;
            }
            _this.handleClose("Something went wrong");
            _this.active = false;
        };
        ws.onmessage = function (e) {
            if (!_this.active) {
                return;
            }
            var data = JSON.parse(e.data);
            if (data.type == "close") {
                _this.handleClose(data.body);
                _this.active = false;
            }
            else {
                _this.handleMessage(data.type, data.body);
            }
        };
    };
    GameWebSocket.prototype.sendMessage = function (type, body) {
        if (!this.active) {
            return;
        }
        var ws = this.websocket;
        if (!ws) {
            return;
        }
        ws.send(JSON.stringify({
            type: type,
            body: body,
        }));
    };
    return GameWebSocket;
}());
