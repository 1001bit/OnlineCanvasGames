interface LevelMessage {
    kinematic: Map<number, KinematicRect>
    static: Map<number, Rect>
}

interface DeleteMessage {
    ID: number;
}

interface CreateMessage {
    rect: Rect | KinematicRect;
    id: number;
}

interface GameInfoMessage {
    rectID: number;
    tps: number;
    constants: PlatformerConstants;
}