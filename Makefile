build:
	go build -ldflags="-H windowsgui" -o bin/currents.exe

debug:
	go build -o bin/currents.exe