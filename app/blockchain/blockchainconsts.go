package blockchain

const (
	BLOCKCHAIN_PORT            = "9119"
	MAX_NODE_CONNECTIONS       = 400
	NETWORK_KEY_SIZE           = 80
	TRANSACTION_HEADER_SIZE    = NETWORK_KEY_SIZE /* from key */ + NETWORK_KEY_SIZE /* to key */ + 4 /* int32 timestamp */ + 32 /* sha256 payload hash */ + 4 /* int32 payload length */ + 4 /* int32 nonce */
	BLOCK_HEADER_SIZE          = NETWORK_KEY_SIZE /* origin key */ + 4 /* int32 timestamp */ + 32 /* prev block hash */ + 32 /* merkel tree hash */ + 4                                      /* int32 nonce */
	KEY_POW_COMPLEXITY         = 0
	TRANSACTION_POW_COMPLEXITY = 1
	BLOCK_POW_COMPLEXITY       = 2
	KEY_SIZE                   = 28
	POW_PREFIX                 = 0
	MESSAGE_TYPE_SIZE          = 1
	// Length of payload in number of bytes
	MESSAGE_OPTIONS_SIZE = 4
)

const (
	MESSAGE_GET_NODES = iota + 20
	MESSAGE_SEND_NODES

	MESSAGE_GET_TRANSACTION
	MESSAGE_SEND_TRANSACTION

	MESSAGE_GET_BLOCK
	MESSAGE_SEND_BLOCK
)

// Needed for several nodes
func SEED_NODES() []string {
	nodes := []string{"10.0.5.33"}

	/*for i := 0; i < 100; i++ {
		nodes = append(nodes, fmt.Sprintf("172.17.0.%d", i))
	}*/

	return nodes
}
