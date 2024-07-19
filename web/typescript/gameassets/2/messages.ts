interface LevelMessage {
    kinematic: Map<number, KinematicRect>
    static: Map<number, PhysicalRect>
}

interface DeleteMessage {
    id: number;
}

interface CreateMessage {
    rect: PhysicalRect | KinematicRect;
    id: number;
}

interface UpdateMessage {
    rectsMoved: Map<number, Vector2>
}

interface GameInfoMessage {
    rectID: number;
    tps: number;
    constants: PlatformerConstants;
}