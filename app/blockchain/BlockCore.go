package blockchain

import "fmt"

var BlockCore = struct {
	*Keypair
	*Blockchain
}{}

func Start(address string) {

	// Setup keys
	//keypair, _ := OpenConfiguration(HOME_DIRECTORY_CONFIG)

	fmt.Println("Generating keypair...")
	keypair := GenerateNewKeypair()
	//WriteConfiguration(HOME_DIRECTORY_CONFIG, keypair)

	BlockCore.Keypair = keypair

	// Setup blockchain
	BlockCore.Blockchain = SetupBlockchan()
	go BlockCore.Blockchain.Run()

	/*
		go func() {
			for {
				select {
				case msg := <-BlockCore.Network.IncomingMessages:
					HandleIncomingMessage(msg)
				}
			}
		}()
	*/
}

func CreateTransaction(txt string) *Transaction {

	t := NewTransaction(BlockCore.Keypair.Public, nil, []byte(txt))
	t.Header.Nonce = t.GenerateNonce(TRANSACTION_POW)
	t.Signature = t.Sign(BlockCore.Keypair)

	return t
}

func HandleIncomingMessage(msg Message) {
	// TODO handle incoming messsage
}
