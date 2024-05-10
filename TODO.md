- [X] Database models
- [X] Get rid of prepared statements
- [X] Switch to chi mux
- [X] Icon and thumbnail for games
- [X] Open access to some pages for unauthorized users
- [X] Header
- [X] Favicon
- [X] Refresh token

- [X] Home page
    - [X] Use from BaseNode gamesJSON instead of database table for games list

- [X] Game Page
    - [X] Rooms list
    - [X] Room creater api
    - [X] Room connect (random/direct)

- [X] Profile page
    - [X] Basic info
    - [X] Logout button

- [ ] Room Page
    - [X] WS Room connection
    - [X] Make RT nodes independent from each other
    - [X] Split Nodes' run select statement into several goroutines (flows)
    - [X] Split WS/SSE handlers and basenode package
    - [X] Make WS connection safer by not allowing single user joining twice
    - [X] Different rooms for different games
    - [X] Basic multiplayer clicker game
    - [ ] UI
    - [ ] Game engine
        - [ ] Backend
        - [ ] Frontend

- [ ] Admin page
    - [ ] UI
    - [ ] SSE
    - [ ] RT nodes control
    - [ ] Users control
    - [ ] Games control