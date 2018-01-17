# very simple NEO transaction catcher 
Catches if transaction was made in the range of predetermined and last block. 
Send information to console

Running rpc client with open wallet required. Use neo-cli
```dotnet neo-cli.dll /rpc
   open wallet <path_to_wallet>
```  
enter your pass

dependencies: 
**neo-go-sdk**
```
go get github.com/CityOfZion/neo-go-sdk
```
**gorilla/mux**
```
go get -u github.com/gorilla/mux
```

run:
```
go run main.go rpcClient.go
```
