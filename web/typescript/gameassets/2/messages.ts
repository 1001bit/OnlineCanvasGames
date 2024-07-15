interface LevelMessage {
    kinematic: Map<number, Rect>
    static: Map<number, Rect>
}

interface DeleteMessage {
    ID: number;
}

interface CreateMessage {
    rect: Rect;
    id: number;
}

interface DeltasMessage {
    kinematic: Map<number, Rect>;
}

interface GameInfoMessage {
    rectID: number;
    constants: PlatformerConstants;
}