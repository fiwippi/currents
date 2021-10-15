build:
	go build -ldflags="-H windowsgui" -o bin/currents.exe

debug:
	go build -o bin/currents.exe

icon:
	go-winres simply --icon icon.png