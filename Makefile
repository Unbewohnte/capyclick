game:=capyclick
bin:=bin

all: clean
	mkdir -p bin
	cd src && go build && mv $(game)* ../bin

cross: clean
	mkdir -p bin
	cd src && GOOS=windows GOARCH=amd64 go build && mv $(game)* ../bin
	

clean:
	rm -rf bin