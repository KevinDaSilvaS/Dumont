# DUMONT

#### Please don't add it to your prod database yet :)

### About
Dumont is a binlog processor for MariaDb written in Go, similar to Maxwell

<img src="Alberto_Santos-Dumont_portrait.jpg" alt="Santos Dumont Brazilian Father of Aviation" width="200"/> <img src="Le_Petit_Journal_Santos_Dumont_25_Novembre_1906.jpg" alt="Santos Dumont Brazilian Father of Aviation" width="200"/>

** About the picture: **  [Santos Dumont Brazilian Father of Aviation.](https://en.wikipedia.org/wiki/Alberto_Santos-Dumont)

### What it does?
It reads the mariadb-binlog files parse them and publishes on rabbitmq concurrently from time to time

### Actions supported:
Currently it only processes INSERT UPDATE and DELETE queries

### Roadmap
For future releases the main goals would be:

- Add web server to provide statistics and prometheus metrics
- Add support to other remaining actions(~~DELETE~~, CREATE TABLE, ALTER TABLE, DROP TABLE, etc)
- ~~Add fetch from remote feature as it will be unlikely the db will be at the same level than dumont~~ [ :white_check_mark: ] Done :rocket:
- ~~Make Dockerfile image work on compose file(currently receiving a permission denied error)~~ [ :white_check_mark: ] Done :rocket:
- ~~Add logs instead of Println~~ [ :white_check_mark: ] Done :rocket:
- Add more producers
- ~~Add previous values to the Old field map on updates~~ [ :white_check_mark: ] Done :rocket:
- ~~Add time stamp to Ts field~~ [ :white_check_mark: ] Done :rocket:

### JSON output provided
It was planned to be as Maxwell compatible as possible

Insert: ` {"Database":"example-db","Table":"people","Type":"INSERT","Ts":1766432444,"Data":{"id_people":"1","last_name":"'Dumont'","name":"'Santos'"},"Old":{},"RawQuery":"INSERT INTO people (name, last_name) VALUES ('Santos', 'Dumont') ","DbTs":"251222 15:48:04"} `

Update: ` {"Database":"example-db","Table":"people","Type":"UPDATE","Ts":1766432444,"Data":{"id":"'2'","id_people":"2","last_name":"'da Silva'","name":"'Kevin'"},"Old":{"id":"2","last_name":"'da SIlva'","name":"'Kevin'"},"RawQuery":"UPDATE people SET id = '2', name = 'Kevin', last_name = 'da Silva' WHERE id = '2' ","DbTs":"251222 16:17:53"} ` 

Delete: ` {"Database":"example-db","Table":"people","Type":"DELETE","Ts":1766432444,"Data":{},"Old":{"id":"1","last_name":"'Dumont!'","name":"'Santos'"},"RawQuery":"DELETE FROM people WHERE ((id = '1')) ","DbTs":"251222 16:15:25"} `

- Database = db name
- Table    = table
- Type     = action(INSERT, UPDATE, etc)
- Ts       = timestamp of execution
- DbTs     = time when the query was executed on db engine YY/MM/DD hh:mm:ss
- Data     = map containing the modified/inserted data
- Old      = map containing the previous data
- RawQuery = new field added containing the raw query executed on db for debug

### Env set up

- DATABASE_NAME       = "database_name"
- DATABASE_PASSWORD   = "paswd"
- DATABASE_USER       = "root"
- DATABASE_HOST       = "127.0.0.1"
- DATABASE_PORT       = "3306"
- EXECUTE_INTERVAL    = "3" time in seconds, dumont will read binlog every X seconds set on EXECUTE_INTERVAL
- MAX_CONSUMERS       = "3" Number of concurrent reader consumers for bin log files, Ex: if I have 6 binlog files and three consumers ideally it will very likely distribute 2 files per consumer
- PRODUCER_HOST       = "amqp://admin:admin@localhost:5672/"
- PRODUCER_QUEUE_NAME = "dumont"

### Running it

  - ### Using podman( or your favorite container tool :) )
     - podman compose -f ./docker-compose.yml up --detach

  - ### Running manually
    - install mariadb-server and mariadb-client
    - install rabbitmq
    - cd dumont
    - set env
    - <span style="background:dimgrey;padding:0.5%;color:orange;border-radius:5%"> <span style="color: blue; font-weight: bold;">go</span> build </span>
    - ./dumont