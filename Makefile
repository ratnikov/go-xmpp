include $(GOROOT)/src/Make.inc

TARG=xmpp
GOFILES=auth.go connection.go

include $(GOROOT)/src/Make.pkg $<

main: package
	6g -I_obj/ -o main.6 main.go
	6l -L_obj/ -o main.out main.6
