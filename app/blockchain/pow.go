package blockchain

import "reflect"

var (
	TRANSACTION_POW = SliceOfBytes(TRANSACTION_POW_COMPLEXITY, POW_PREFIX)
	BLOCK_POW       = SliceOfBytes(BLOCK_POW_COMPLEXITY, POW_PREFIX)
)

// CheckProofOfWork check if
// proofType (transaction or block)
func CheckProofOfWork(proofType []byte, hash []byte) bool {

	if len(proofType) > 0 {
		return reflect.DeepEqual(proofType, hash[:len(proofType)])
	}
	return true
}
