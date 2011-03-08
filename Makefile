include $(GOROOT)/src/Make.inc

TARG=xmpp
GOFILES=auth.go connection.go

%.6: %.go
	6g -o $@ $<

include $(GOROOT)/src/Make.pkg $<

%.out: %.6
	6l -o $@ $<
