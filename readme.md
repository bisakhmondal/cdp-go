# Chromium Starter work

### A fully functional dynamic parser and analyser for chromium git repo using chrome devtools protocol (CDP)

## How to run

### In Local Environment

```shell
go run main.go
#or
# go build main.go
#./main
```

**Note**: This program access different flags which can be passed during execution of the program.

#### flags

Flag Name | Default Flag Value | Usage |
---- | --- | --- |
-timeout  | 20 | context timeout in seconds|
-repo | https://chromium.googlesource.com/chromiumos/platform/tast |Repository URL |
-branch  | main | branch name where the parser should run |
-dir | ./commits |  folder where parsed commit messages is going to be stored |
-csvpath |stats.csv  | csv file location where the details statistics is going to be stored |
-commits | 20 | Number of commits to be parsed|

So to run with flags

```shell
go run main.go -FlagName1 FlagValue1 -FlagName2 FlagValue2 
```

### Run Using Docker

#### Steps

1. Build the image first

    ```shell
    docker build -t <Image Name> .
    ```

2. Run a container.

    ```shell
    docker run --name <Container Name> -v <Mount Volume Path>:/go-cdp/commits -d <Same Image Name>
   ```

3. Exec the following commands inside the container & the output will be stored into "mount volume path"

    ```shell
     docker exec -it <Same Container Name> /bin/bash
     ./main -csvpath ./commits/statistics.csv
    ```

A full example can be summarized as

```shell
 docker build -t hchrome .
 #The result will be stored in ~/Desktop/test in local machine.
 docker run --name scraper -p 9222:9222 -v ~/Desktop/test:/go-cdp/commits -d hchrome
 docker exec -it scraper /bin/bash
 #inside the container
 ./main -csvpath ./commits/statistics.csv
```
