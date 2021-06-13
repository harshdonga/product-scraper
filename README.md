To run the containers:

1. Clone the repo
2. change directory to main directory
3. docker-compose build ( Builds both the services - scrape & mongo interactor )
4. docker-compose up    ( Runs all the containers )

Container Name            Port            Description

scraper                   5000            Takes x-www-form-urlencoded amazon product URL as POST data

dbapi                     5001            Interacts with mongoDB to store & update documents

mongoDB                   27017           mongoDB service to store product documents

