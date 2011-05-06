include make.$(GOARCH)

all: compile bf

compile:
	$(c) ./src/bf2go.go
	$(c) ./src/app.go
	$(l) -o bf2go app.$(obj)

clean:
	rm -f *.$(obj)

include make.bf