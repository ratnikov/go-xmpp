include $(GOROOT)/src/Make.inc

TARG=xmpp
GOFILES=auth.go connection.go util.go test_util.go

_xmpp_:
	6g -o _xmpp_.6 -Itls/_obj/ $(GOFILES)

xmpp.a: _obj _xmpp_
	gopack grc _obj/xmpp.a _xmpp_.6

main: xmpp.a
	6g -I_obj/ -o main.6 main.go
	6l -L_obj/ -o main.out main.6

include $(GOROOT)/src/Make.pkg
