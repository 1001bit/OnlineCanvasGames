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
    movedRects: Map<number, PhysicalRect>
}

interface GameInfoMessage {
    rectID: number;
    tps: number;
    constants: PlatformerConstants;
}