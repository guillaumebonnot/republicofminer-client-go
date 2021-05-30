# Republic of Miner
Republic of Miner is a blockchain game where you need to mine resources with your CPU and craft stuff.\
More informations here : http://republicofminer.com

# The project
The code included in this repository was written for fun, feel free to contribute for more fun !

Basically, it includes a web server used a rest api to query the blockchain like an explorer,\
and a client that will continuously mine resources on your behalf.

You can switch between those 2 behaviour by modifying main.go

## web
We setup a web server where you can query blocks and transactions :

Get transaction by hash : http://localhost:3000/tx/zIJZB67U0gTUnGq649baM%2F5ylbUE1ydm5WpJ7xn2XfQ%3D \
Get block by height : http://localhost:3000/block/10 \
Get block by hash : http://localhost:3000/block/XZ6uMgCbVnDmj10QMjD8g7pquULgtfBuYwbm5mDAKzs= \

The server connects to the explorer to get the data and displays it in json format.

## explorer
The package to access the blockchain explorer.

## miner
The miner loads the wallet and continuously requests a mining task to the server then solves it.\
The miner connects both to the explorer and the game server.

## republicofminer
The package to access the game server.

## protocol
The procotocol folder contains code related the representation of the elements of the blockchain.\
More informations can be found here : https://github.com/caasiope/caasiope-blockchain

## wallet
The wallet loads or create a private key in the database.

## vault
The vault is a SQLite database where you can store data encrypted by the password associated with a key.




