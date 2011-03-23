include $(GOROOT)/src/Make.inc

TARG=xmpp
GOFILES=auth.go connection.go util.go

net.a:
	cd net && gomake

tls.a: net.a
	cd tls && gomake

tls-clean:
	cd tls && gomake clean

_xmpp_: tls.a
	6g -o _xmpp_.6 -Itls/_obj/ $(GOFILES)

_obj:
	mkdir _obj

xmpp.a: _obj _xmpp_
	gopack grc _obj/xmpp.a _xmpp_.6

main: xmpp.a
	6g -I_obj/ -Itls/_obj/ -o main.6 main.go
	6l -L_obj/ -Ltls/_obj/ -o main.out main.6

clean: tls-clean
	rm -rf *out _obj *6
