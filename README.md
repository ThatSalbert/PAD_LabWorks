### Steps to run the project:

1. Clone this Git repository
2. Download and import the Postman collection from [here](https://github.com/ThatSalbert/PAD_LabWorks/blob/main/Resources/Postman/Weather%20API.postman_collection.json) which has all the needed requests and the request bodies are already set up (you just need to change the values of the parameters)
3. Build and run the project with `docker-compose`
4. Run the requests from the Postman collection which can be found in `Weather API >> Gateway Module`

 **NOTE 1**: When building and running the project, it will take some time for the database to initialise, to create the database and to import the data from the SQL file.

 **NOTE 2**: List of variables and what values they can have can be found in the description of the Postman collection.

 **NOTE 3**: The only table that has data is the `location_table` in both databases. The only country that is available is `Moldova`.

 **NOTE 4**: Unit Test is added to `disaster-service-1`. Open up a terminal of `disaster-service-1` and run `go test main_test.go` to run the test.

 More information can be found [here](https://github.com/ThatSalbert/PAD_LabWorks/blob/main/Resources/ProjectInfo.md)
