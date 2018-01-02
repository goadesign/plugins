build:
	cd calc/cmd/calcsvc && go build
	cd adder/cmd/addersvc && go build

adder: build
	cd adder/cmd/addersvc && ./addersvc

calc: build
	cd calc/cmd/calcsvc && ./calcsvc
