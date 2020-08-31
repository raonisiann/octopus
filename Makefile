

all:
	go build -o ./bin/octopus octopus

clean:
	rm -f bin/octopus