# Prototype Admin Dashboard coding_challenge
Prototype admin dashboard that collects and displays usage metrics submitted by a fakeIOT client.

### Installation:

```
go get github.com/gmo-personal/coding_challenge
cd ~/go/src/github.com/gmo-personal/coding_challenge
```

### Setup database: 
Assuming mysql is installed, 

Use the following command in mysql to create the schema.

`create schema teleport`

### Setup server:

Run the following command to spin up the docker container. The site can be accessed on https://localhost:443

replace \<username> and \<password> with local mysql database username and password respectively.

`db_user=<username> db_pass=<password> docker-compose up --build`

### Run tests:

Use the following command to run the tests. 

Note that this only works if the server setup was run exactly previously, otherwise manually use docker container ls to find the correct docker container ID

`docker exec -it $(docker container ls -q | head -1) bash`

`./test.sh`

### Site:

The default login for the base user `testacct-0000-0000-0000-000000000000`:

username: `t@gmail.com`

password: `t`

