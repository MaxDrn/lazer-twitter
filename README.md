# Lazer Twitter



## Database setup

To setup the Database,
you need to have it running on an individual Docker Container.

Type in the following to start you Container:

```
docker run --name your_name -e POSTGRES_PASSWORD=your_password -p your_port container_name 
```

After your Database is running on a container, you can access it in order to manage your data from the command line.

To do that, type in:

```
docker exec -it postgres psql -U username
``` 

## GO Server

In order to start your GO Server, move to the path of the project and type in:

``` 
go run cmd/lazer-twitter/main.go --rest-listen-port server_port --db-name database_name --db-user database_user --db-pw database_password --db-port database_port
```

Now you can go to your default Browser and access the server, example:

```
localhost:5432
```