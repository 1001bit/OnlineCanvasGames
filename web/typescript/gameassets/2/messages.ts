interface LevelMessage {
    blocks: {};
    players: {};

    config: LevelConfig;

    serverTps: number;
    clientTps: number;

    playerRectId: number;
}

interface LevelUpdateMessage {
    players: {};
    correct: boolean;
}

interface ConnectMessage {
    rect: AbstractPlayer;
    rectId: number;
}

interface DisconnectMessage {
    rectId: number;
}