# FAMPAY ASSIGNMENT

This repository contains the code for the assignment given by Fampay.

The assignment is [here](https://www.notion.so/fampay/Backend-Assignment-FamPay-32aa100dbd8a4479878f174ad8f9d990).

### About my submission:

I have used the following technologies:

- Golang
- Gorilla Mux
- MongoDB (database)
- Docker

The application is built as mentioned in the Assignment.

### Instructions to Run the App:

- First, a .env file needs to be put in the root directory. The .env.example files contain the format for the env file.
- To run the application one needs to run the docker command

      `docker build –rm -t <name_of_docker_container_> .`


      `docker run -p 8080:8080 <name_of_docker_container_>`

  That’s all. It will run the application on localhost:8080

- Another method to run the application is simply opening the terminal at the root and running the following commands:

  `go mod tidy` – it will install all the packages into your local machine

  `go run main.go` – it will start the server at localhost:8080

### API Routes:

- There are two APIs
  - First, is on the home route `/` which returns the list of videos from database according to pagination where you have to pass the page parameter. Example: [http://localhost:8080/?page=2](http://localhost:8080/?page=2)
  - Second is the search query route where we can search for data in database using a specific string. Example: [http://localhost:8080/search?searchQuery="your_query"&searchType="titleOrdescription"](http://localhost:8080/search?searchQuery=)
- The main.go file contains the function for executing the YouTube Data API and it runs at an interval of 20 seconds searching for the query: Programming.
