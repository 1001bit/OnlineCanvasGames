interface LevelMessage {
    blocks: Map<number, AbstractBlock>
    players: Map<number, AbstractPlayer>

    config: LevelConfig

    tps: number;

    playerRectId: number
}

interface ConnectMessage {
    rect: AbstractPlayer;
    rectId: number;
}

interface DisconnectMessage {
    rectId: number;
}