interface AbstractBlock extends AbstractRect {

}

function isAbstractBlock(obj: any) : obj is AbstractBlock{
    return isAbstractRect(obj)
}

class Block extends Rect {
    constructor(abstract: AbstractBlock){
        super(abstract)
    }
}