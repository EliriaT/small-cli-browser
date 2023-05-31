# small-cli-browser
* Laboratory work nr 2 for Web Programming course.
* CLI application that allows making HTTP requests from the terminal and to the Google search engine.
* No HTTP client is used. 
* HTTP requests are made and received only through TCP connection only


## Run the application

To run the CLI app make sure you are in the project's root directory, where the `main.go` or `go2web` file is located. 

1. Run the help command:

```sh
go run main.go -h
```
 or 
 
 ```sh
./go2web -h
```

2. Run the search command:

```sh
go run main.go -s [a list of words describing what you want to search using the google search engine ]
```


```sh
go run main.go -s websockets history
```
 or 
 
 ```sh
./go2web -s what is the weather today
```


3. Run the url command:

```sh
go run main.go -u https://www.lipsum.com/
```

 or 
 
 ```sh
./go2web -u https://www.rabenhorst.de/en/science-of-juice/fruit-vegetable-encyclopedia/sea-buckthorn/
```
