package blockchain

import (
	"crypto/sha256"
	"reflect"
)

// BlockSlice is a slice of Block structures
type BlockSlice []Block

// Exists Check if the block already exists
func (bs BlockSlice) Exists(b Block) bool {

	//Traverse array in reverse order because if a block exists is more likely to be on top.
	l := len(bs)
	for i := l - 1; i >= 0; i-- {

		bb := bs[i]
		if reflect.DeepEqual(b.Signature, bb.Signature) {
			return true
		}
	}

	return false
}

// PreviousBlock returns previous block
func (bs BlockSlice) PreviousBlock() *Block {
	l := len(bs)
	if l == 0 {
		return nil
	}
	return &bs[l-1]
}

// Block structure
type Block struct {
	*BlockHeader
	Signature []byte
	*TransactionSlice
}

// BlockHeader struct
// Signature is the identifier
// Origin
// PrevBlock is the previous block identifier
// Nonce the value needed t get the proper format for the signature hash
type BlockHeader struct {
	Origin     []byte
	PrevBlock  []byte
	MerkelRoot []byte
	Timestamp  uint32
	Nonce      uint32
}

// NewBlock cerates the new block with its prevblock hash as
// header of previousBlock incoming
func NewBlock(previousBlock []byte) Block {

	header := &BlockHeader{PrevBlock: previousBlock}
	return Block{header, nil, new(TransactionSlice)}
}

// AddTransaction addds new transaction to TransactionSlice in Block
func (b *Block) AddTransaction(t *Transaction) {
	newSlice := b.TransactionSlice.AddTransaction(*t)
	b.TransactionSlice = &newSlice
}

// Sign block
func (b *Block) Sign(keypair *Keypair) []byte {

	s, _ := keypair.Sign(b.Hash())
	return s
}

// VerifyBlock verify headers authenticity of the block
func (b *Block) VerifyBlock(prefix []byte) bool {

	headerHash := b.Hash()
	merkel := new([]byte) //b.GenerateMerkelRoot()

	return reflect.DeepEqual(merkel, b.BlockHeader.MerkelRoot) && CheckProofOfWork(prefix, headerHash) && SignatureVerify(b.BlockHeader.Origin, b.Signature, headerHash)
}

// Hash cerates the hash for the block header
func (b *Block) Hash() []byte {

	headerHash, _ := b.BlockHeader.MarshalBinary()
	hash := sha256.New()
	hash.Write(headerHash)
	return hash.Sum(nil)
}

// GenerateNonce generates the valid block nonce
// Returns an integer as valid nonce
func (b *Block) GenerateNonce(prefix []byte) uint32 {

	newB := b
	for {

		// ensures the block header bits which indicate the target difficulty
		// is in min/max range and that the block hash is less than the
		// target difficulty as claimed
		if CheckProofOfWork(prefix, newB.Hash()) {
			break
		}

		newB.BlockHeader.Nonce++
	}

	return newB.BlockHeader.Nonce
}

///////////////////////////////
///////////--------------_NOT USABLE FROM NOW

// GenerateMerkelRoot Generates the hash(hash()) of the block identifier
func (b *Block) GenerateMerkelRoot() []byte {

	//var merkell func(hashes [][]byte) []byte
	/*
			merkell = func(hashes [][]byte) []byte {

				//
				numNodes := len(hashes)
				select {
				case numNodes == 0:
					return nil
				case numNodes == 1:
					return hashes[0]
				case numNodes%2 == 1:
					return merkell([][]byte{
						merkell(hashes[:numNodes-1]), hashes[numNodes-1]
					})
				default:
					bs := make([][]byte, numNodes/2)
					for i := range bs {
						j, k := i*2, (i*2)+1
						hash := sha256.New()
						bs[i]  = hash.Write(append(hashes[j], hashes[k]...))
					}
					return merkell(bs)
				}
			}

		ts := functional.Map(
			func(t Transaction) []byte {
				return t.Hash()
			},
			[]Transaction(*b.TransactionSlice)).([][]byte)
	*/
	return nil //merkell(ts)

}

//////////////------------end test

// Transformation of BlockHeader in []bytes in common structure size
func (b *Block) MarshalBinary() ([]byte, error) {
	return nil, nil
}

func (b *Block) UnmarshallBinary(data []byte) error {
	return nil
}

func (bh *BlockHeader) MarshalBinary() ([]byte, error) {
	return nil, nil
}

func (bh *BlockHeader) UnmarshallBinary() ([]byte, error) {
	return nil, nil
}
