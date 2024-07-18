function collideKinematicWithStatic(kinematicRect: KinematicRect, staticRect: PhysicalRect, dt: number){
    if(!staticRect.doApplyCollisions){
        return
    }
    kinematicRect.setCollisionDir(Direction.None)

    let futureKinematic = new Rect(kinematicRect);
    futureKinematic.setPosition(kinematicRect.targetPosition.x, kinematicRect.targetPosition.y)

    const velX = kinematicRect.velocity.x * dt
    const velY = kinematicRect.velocity.y * dt

    if(velY > 0) {
        // down
        futureKinematic.size.y += velY

        if(futureKinematic.intersects(staticRect)){
            kinematicRect.setTargetPos(kinematicRect.targetPosition.x, staticRect.position.y - kinematicRect.size.y, true)
            kinematicRect.velocity.y = 0

            kinematicRect.setCollisionDir(Direction.Down)
        }

    } else if(velY < 0) {
        // up
        futureKinematic.size.y += Math.abs(velY)
        futureKinematic.position.y -= Math.abs(velY)

        if(futureKinematic.intersects(staticRect)){
            kinematicRect.setTargetPos(kinematicRect.targetPosition.x, staticRect.position.y + staticRect.size.y, true)
            kinematicRect.velocity.y = 0

            kinematicRect.setCollisionDir(Direction.Up)
        }
    }

    futureKinematic = new Rect(kinematicRect);
    futureKinematic.setPosition(kinematicRect.targetPosition.x, kinematicRect.targetPosition.y)

    if(velX > 0) {
        // right
        futureKinematic.size.x += velX

        if(futureKinematic.intersects(staticRect)){
            kinematicRect.setTargetPos(staticRect.position.x - kinematicRect.size.x, kinematicRect.targetPosition.y, true)
            kinematicRect.velocity.x = 0

            kinematicRect.setCollisionDir(Direction.Right)
        }

    } else if(velX < 0) {
        // left
        futureKinematic.size.x += Math.abs(velX)
        futureKinematic.position.x -= Math.abs(velX)

        if(futureKinematic.intersects(staticRect)){
            kinematicRect.setTargetPos(staticRect.position.x + staticRect.size.x, kinematicRect.targetPosition.y, true)
            kinematicRect.velocity.x = 0

            kinematicRect.setCollisionDir(Direction.Left)
        }
    }
}