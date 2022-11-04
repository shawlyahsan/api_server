# api_server  
A simple Rest API using go, [gorilla mux](https://github.com/gorilla/mux), basic authentication and JWT token 
## Installation  
Run the server using the following commands

    $ git clone https://github.com/shawlyahsan/api_server.git  
    $ cd api_server
    $ bin/api_server  

## Implementation  
**api.go**  
- Handles all the endpoint requests using mux as router  
- Every endpoint request first undergoes basic authentication or jwt authorization  


**auth.go**  
- Contains the basic authentication and jwt authorization middleware
- Generates bearer token and verifies token

## Avaialble Endpoints  
| Method | API Endpoint | Authentication Type | Response |
| --- | --- | --- | --- |
| POST | /login | Basic | Generates a bearer token |
| GET | /books | Bearer token | Returns the description of all the books |
| GET | /books/{id} | Bearer token | Returns the description of the book with the specified id |
| POST | /books | Bearer token | Creates a new book description |
| PUT | /books/{id} | Bearer token | Updates the decription of the specified book |
| DELETE | /books/{id} | Bearer token | Deletes the book specified by id |  
## Environment Variables  
    $ export api_user=admin
    $ export api_pass=pass  

