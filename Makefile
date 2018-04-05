run:
	- gin run main.go

build:
	- go build .

proto:
	- protoc -I . --go_out=plugins=grpc:. buildJob.proto
