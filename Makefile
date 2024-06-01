game:=capyclick
desktopBin:=bin/desktop
webBin:=bin/web
wasmHtml:=build/capyclick.html
wasmExec:=$(shell go env GOROOT)

current: clean environment
	cd src && go build && mv $(game)* ../$(desktopBin)

web: environment
	cd src && env GOOS=js GOARCH=wasm go build -o $(game).wasm . && mv $(game).wasm ../$(webBin)
	cp $(wasmHtml) $(webBin)
	cp $(wasmExec)/misc/wasm/wasm_exec.js $(webBin)

desktop: clean environment
	cd src && GOOS=windows GOARCH=amd64 go build && mv $(game)* ../$(desktopBin)
	cd src && GOOS=linux GOARCH=amd64 go build && mv $(game)* ../$(desktopBin)

cross: clean environment web desktop

environment:
	mkdir -p $(desktopBin) $(webBin)

clean:
	rm -rf bin