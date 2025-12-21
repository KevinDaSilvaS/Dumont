# DUMONT

#### Please don't add it to your prod database yet :)

### About
Dumont is a binlog processor for MariaDb written in Go, similar to Maxwell

<img src="Alberto_Santos-Dumont_portrait.jpg" alt="Santos Dumont Brazilian Father of Aviation" width="200"/> <img src="Le_Petit_Journal_Santos_Dumont_25_Novembre_1906.jpg" alt="Santos Dumont Brazilian Father of Aviation" width="200"/>

** About the picture: **  [Santos Dumont Brazilian Father of Aviation.](https://en.wikipedia.org/wiki/Alberto_Santos-Dumont)

### What it does?
It reads the mariadb-binlog files parse them and publishes on rabbitmq concurrently from time to time

### Actions supported:
Currently it only processes INSERT and UPDATE queries

### Roadmap
For future releases the main goals would be:

- Add web server to provide statistics and prometheus metrics
- Add support to other remaining actions(DELETE, CREATE TABLE, ALTER TABLE, DROP TABLE, etc)
- ~~Add fetch from remote feature as it will be unlikely the db will be at the same level than dumont~~ [ :white_check_mark: ] Done :rocket:
- Make Dockerfile and compose file work
- Add logs instead of Println
- Add more producers
- Add previous values to the Old field map on updates
- Add time stamp to Ts field

### JSON output provided
It was planned to be as Maxwell compatible as possible

Insert: `  {"Database":"database_name","Table":"users","Type":"INSERT","Ts":0,"Data":{"id":"2","id_users":"2","name":"\"kevin\""},"Old":{},"RawQuery":"INSERT INTO users (id, name) VALUES (2, \"kevin\")\n"} `

Update: ` {"Database":"database_name","Table":"users","Type":"UPDATE","Ts":0,"Data":{"id_users":"1","name":"'Stanislaw Lem'"},"Old":{},"RawQuery":"update users set name = 'Stanislaw Lem' WHERE id = 1\n"} ` 

- Database = db name
- Table    = table
- Type     = action(INSERT, UPDATE, etc)
- Ts       = timestamp of execution, currently not being filled
- Data     = map containing the modified/inserted data
- Old      = map containing the previous data before an updates, also currently not being filled
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
- READ_FROM_REMOTE    = "TRUE" set to TRUE if your db is on a remote server

### Running it
- cd dumont
- go build
- ./dumont