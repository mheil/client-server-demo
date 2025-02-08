# go-client-server
Go demo project for simple client / server implementation

Hopefully someone might find it helpfull to look at this code to address
some of of the problems a developer has when beginning with golang.
Especially the client code may be usefull.

Please feel free to make suggestions on how to improve the code. 

## function
The function of the client is to read single lines from stdin and send them
to the server which displays the received messages. When the client sends
a message containing "exit" the server has to close the connection after sending
a "bye..." message. The client should afterwards gracefully shutdown.

## implementation / difficulties
This project allowed me to test how use tcp sockets networks in golang.
Especially the following aspects were interesting: 
 - detect disconnection
 - properly close connections when application stops

### client
The client was far more difficult to implement, cause as it points out,
it is not possible to detect a closed connection while writing to a connection.
Disconnection detection can only be done while reading from the connection.
After a disconnect has been detected the routine reading from stdin had to be
stopped from the (parallel executing) routine that is reading from the connection.

The first attempt was to simply close stdin, cause the docs say that closing
it would cause blocking reader to stop reading. This was not the case;
closing stdin it not possible and another way had to be found to stop
the sending routine.

### server
The server itself was relatively easy to implement. The only point to really
take care of was to close the connections in all cases.
 
## build
To build client and server, a simple `go install ./cmd/*` should be sufficiant.
After that, the client can be simply started by `client` wheras the server
can be simply started by `server`.

Per default, client and server use the port `7666` to communicate, which can
be changed by command line flag.

## open quirks / unhandled cases

### Packet loss, faulty connections
Detecting of disconnections due to cable fault or dropped tcp packets by 
firewall are not yet handled, as this has to use timeouts, which may be
handled in a future case. 
