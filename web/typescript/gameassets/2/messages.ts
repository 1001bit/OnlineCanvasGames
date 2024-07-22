interface GameInfoMessage {
    tps: number;
}

interface LevelMessage {
    blocks: Map<number, AbstractBlock>
    players: Map<number, AbstractPlayer>

    config: LevelConfig

    playerRectId: number
}

interface PlayerMovementMessage {
    playersMoved: Map<number, {x: number, y: number}>
}

interface ConnectMessage {
    rect: AbstractPlayer;
    rectId: number;
}

interface DisconnectMessage {
    rectId: number;
}