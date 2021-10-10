build-windows:
	go build -ldflags="-H windowsgui" -o bin/currents.exe

build-linux:
	go build -o bin/currents