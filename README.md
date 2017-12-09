# geocom

Geocom is a P2P(Peer to Peer) chat room where a user can start a server and other users can connect to that server. The entire app is written in Go, including the TUI (Terminal User Interface).

![alt text](https://raw.githubusercontent.com/apache8080/geocom/master/geocom_example.png)

## technology

Geocom uses the Go networking library to connect two users via TCP. For the TUI we are using the tui-go library.

## building geocom

```
git clone https://github.com/apache8080/geocom.git
cd geocom
go get github.com/marcusolsson/tui-go
go build
```

## running geocom as the host
```
./geocom [nickname]
```

## connecting to a geocom host
```
./geocom [nickname] [IP address of host]
```

## contributors
James Wang (james9909) and Rishi Desai (apahce8080)
