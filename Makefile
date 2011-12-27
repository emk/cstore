include $(GOROOT)/src/Make.inc

TARG=cstore
GOFILES=\
	digest.go\
	registry.go\
	cstore.go\
	main.go\

include $(GOROOT)/src/Make.cmd
