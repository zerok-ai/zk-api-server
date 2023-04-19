GOOS=linux GOARCH=amd64 go build -o main main.go
docker build -t zk-api-server:dev .
sh ./gcp-artifact-deploy-go.sh