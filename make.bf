bf: cleanbf compilebf

compilebf:
	./bf2go -f bf/helloworld.bf -o helloworld.go
	$(c) ./helloworld.go
	$(l) -o helloworld helloworld.$(obj)

cleanbf:
	rm -f helloworld