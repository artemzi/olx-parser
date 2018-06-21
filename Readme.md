# TODO

     * olx client
     * select / change config from checkboxes
     * storage (mongo?)
     * cross-posting / mail

## usage

Set variables for running it localy with:
> export SERVICE_PORT="8000"

```Bash
go test -v && env CGO_ENABLED=0 GOOS=linux go build -o olx-parser .
docker build -t olx-parser -f ./Dockerfile .
docker run -p 8000:8000 olx-parser
```

request (will be changed) for debug:
```Bash
curl --header "Content-Type: application/json" \                                   52 ↵   dev ●
        --request POST \
        http://0.0.0.0:8000/adverts
```