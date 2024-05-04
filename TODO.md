- [X] Database models
- [X] Get rid of prepared statements
- [X] Switch to chi mux
- [X] Icon and thumbnail for games
- [X] Open access to some pages for unauthorized users
- [X] Header
- [X] Favicon

- [X] Home page

- [X] Profile page

- [X] Game Page
    - [X] Room divs 
    - [X] New room creater api/ws
        - [X] Base
        - [X] Security
    - [X] Dynamic rooms list
    - [X] Room connect
    - [X] Random room returner api

- [ ] Room Page
    - [X] WS Room connection
    - [X] Make RT nodes independent from each other
    - [X] Split WS/SSE handlers and basenode package
    - [X] Make WS connection safer by not allowing single user joining twice
    - [X] Different rooms for different games
    - [X] Basic multiplayer clicker game
    - [ ] UI

- [ ] Admin page
    - [ ] UI
    - [ ] SSE
    - [ ] RT nodes control
    - [ ] Users control
    - [ ] Games control

- [X] Receiving games from BaseNode cache instead of querying database each home page load
- [X] Split Nodes' run select statement into several goroutines (flows)
- [ ] Put context where needed
- [ ] Logout functionality
- [ ] Profile editing functionality
- [ ] Longer JWT expiration date, JWT date for checking user existance in database

- [ ] Game engine
    - [ ] Backend
    - [ ] Frontend