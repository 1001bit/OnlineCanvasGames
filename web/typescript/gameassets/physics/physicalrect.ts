interface AbstractPhysicalRect extends AbstractRect {
    canCollide: boolean;
}

function isAbstractPhysicalRect(obj: any) : obj is PhysicalRect{
    return isAbstractRect(obj) && "canCollide" in obj
}

class PhysicalRect extends Rect {
    canCollide: boolean;

    constructor(abstractRect: AbstractPhysicalRect){
        super(abstractRect)

        this.canCollide = abstractRect.canCollide
    }
}