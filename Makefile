include $(GOROOT)/src/Make.inc

TARG=cstore
GOFILES=\
	digest.go\
	cstore.go\

include $(GOROOT)/src/Make.pkg
