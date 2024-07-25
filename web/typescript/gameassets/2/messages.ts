interface LevelMessage {
    blocks: {};
    players: {};

    config: LevelConfig;

    tps: number;

    playerRectId: number;
}

interface LevelUpdateMessage {
    movedPlayers: {};
}

interface ConnectMessage {
    rect: AbstractPlayer;
    rectId: number;
}

interface DisconnectMessage {
    rectId: number;
}