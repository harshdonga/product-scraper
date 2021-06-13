To run the containers:

1. Clone the repo
2. change directory to main directory
3. docker-compose build ( Builds both the services - scrape & mongo interactor )
4. docker-compose up    ( Runs all the containers )

Container Name            Port            Description

scraper                   5000            Takes x-www-form-urlencoded amazon product URL as POST data

dbapi                     5001            Interacts with mongoDB to store & update documents

mongoDB                   27017           mongoDB service to store product documents


References :

1. https://hub.docker.com/_/golang
2. https://www.youtube.com/watch?v=JNr5noDp6EM (playlist)
3. https://golangtutorial.dev/tips/http-post-json-go/#:~:text=Follow%20the%20below%20steps%20to,NewRequest%20method.&text=Second%20parameter%20is%20URL%20of%20the%20post%20request.
4. https://stackoverflow.com/questions/13582519/how-to-generate-hash-number-of-a-string-in-go
5. https://www.google.com/search?q=Get+%22%22%3A+unsupported+protocol+scheme+%22%22&oq=Get+%22%22%3A+unsupported+protocol+scheme+%22%22&aqs=chrome..69i57.9153j0j1&sourceid=chrome&ie=UTF-8

