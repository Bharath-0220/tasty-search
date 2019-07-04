# tasty-search

Steps -
1. Clone the repository
2. Copy the finefoods.txt file into the cloned directory
3. docker build -t test1 .
4. docker run -p 8080:8080 -it test1:latest


Ping URL-
http://localhost:8080/ping

Search URL-
http://localhost:8080/search?s={"tokens":["good"]}
