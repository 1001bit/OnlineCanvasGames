- [X] Database models
- [X] Get rid of prepared statements
- [X] Switch to chi mux
- [X] Icon and thumbnail for games
- [X] Open access to some pages for unauthorized users
- [X] Header
- [X] Favicon
- [X] Refresh token
- [X] Use typescript instead of javascript
    - [X] Compiler
    - [X] Static scripts
    - [X] Game scripts
    
- [X] Use custom types for:
    - [X] RectID

- [X] Home page
    - [X] Use BaseNode's gamesJSON instead of database table for games list

- [X] Game Page
    - [X] Rooms list
    - [X] Room creater api
    - [X] Room connect (random/direct)

- [X] Profile page
    - [X] Basic info
    - [X] Logout button

- [X] Room Page
    - [X] WebSocket Room connection
    - [X] Make RT nodes independent from each other
    - [X] Split Nodes' run select statement into several goroutines (flows)
    - [X] Split WS/SSE handlers and basenode package
    - [X] Make WS connection safer by not allowing a user joining twice
    - [X] Different rooms for different games
    - [X] Basic multiplayer clicker game
    - [X] UI elements
        - [X] Server Message
        - [X] "Show nav bar" button
    - [X] Limit players amount

    - [X] Fix bug: concurrent map iteration and map write: roomnode/public.go:28
        - [X] Concurrent map
        - [X] Concurrent set

    - [X] Multiplayer platformer game
        - [X] Server
            - [X] Game loop
            - [X] Player and blocks
            - [X] Forces and collisions
            - [X] Client controls handling
                - [X] Controls receive
                - [X] Player control
                - [X] Limit input ticks
            - [X] Data send
                - [X] Level
                - [X] Info
                    - [X] Player rectID
                    - [X] Constants
                    - [X] TPS
                - [X] Rect Delete/Create
                - [X] Level update
                - [X] Level correction state
            - [ ] Fixed Timestep

        - [X] Client
            - [X] Level Draw
            - [X] Game loop
            - [X] Controls
                - [X] Frontend
                - [X] Send to Backend
                - [X] Postpone input ticks to next iteration, if the limit was bypassed
            - [X] Messages handling
            - [X] Smooth rect movements
                - [X] Client State -> Updated State interpolation
                - [X] Player physics replication
                    - [X] Forces
                    - [X] Control
                    - [X] Collisions
                - [X] Kinematic Players position correction
            - [X] Smooth following camera
            - [ ] Animated sprites

- [X] Microservices:
    - [X] Storage
    - [X] Users
    - [X] Games
    - [X] API gateway

- [X] use gRPC for gateway->service communications
    - [X] user service

- [X] Switch from html/templates -> a-h/templ
- [ ] Switch to htmx

- [ ] Admin page
    - [ ] Front end
    - [ ] SSE
    - [ ] RT nodes control
    - [ ] Users control
    - [ ] Games control