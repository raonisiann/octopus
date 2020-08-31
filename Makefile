

all:
	go build -o ./bin/octopus **/*.go

clean:
	rm -f bin/octopus