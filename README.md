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
- Make Dockerfile and compose file work
- ~~Add logs instead of Println~~ [ :white_check_mark: ] Done :rocket:
- Add more producers
- ~~Add previous values to the Old field map on updates~~ [ :white_check_mark: ] Done :rocket:
- ~~Add time stamp to Ts field~~ [ :white_check_mark: ] Done :rocket:

### JSON output provided
It was planned to be as Maxwell compatible as possible

Insert: ` {"Database":"example-db","Table":"names","Type":"INSERT","Ts":1766405680,"Data":{"id_names":"1","name":"'oi'"},"Old":{"@1":"1","@2":"'oi'"},"RawQuery":"INSERT INTO names (name) VALUES ('oi') ","DbTs":"251221  9:27:34"} `

Update: ` {"Database":"example-db","Table":"names","Type":"UPDATE","Ts":1766405680,"Data":{"id":"'1'","id_names":"1","name":"'Santos Dumont'"},"Old":{"@1":"1","@2":"'Santos Dumont'"},"RawQuery":"UPDATE names SET id = '1', name = 'Santos Dumont' WHERE id = '1' ","DbTs":"251221  9:27:49"} ` 

Delete: ` {"Database":"example-db","Table":"table2","Type":"DELETE","Ts":1766406191,"Data":{},"Old":{"@1":"1","@2":"123","@3":"'oii'"},"RawQuery":"DELETE FROM table2 WHERE ((id = '1')) ","DbTs":"251222  8:43:42"} `

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
- MAX_CONSUMERS       = "3" Number of concurrent reader consumers for bin log files, Ex: if I have 6 binlog files and three consumers it will likely be distributed between the consumers so each will read 2 files
- PRODUCER_HOST       = "amqp://admin:admin@localhost:5672/"
- PRODUCER_QUEUE_NAME = "dumont"

### Running it
- cd dumont
- go build
- ./dumont