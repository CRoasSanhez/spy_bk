package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"reflect"
	"time"
)

// Transaction represents the struct used for a transaction
// Payload is the amount to be sent
// Signature is the sender public key
// PrevTx is the previous transaction
type Transaction struct {
	Header    TransactionHeader
	Signature []byte
	Payload   []byte
}

// TransactionHeader ...
// From is the user who is sending the transaction
// To is the user who is receiving the transaction
type TransactionHeader struct {
	From          []byte
	To            []byte
	Timestamp     uint32
	PayloadHash   []byte
	PayloadLength uint32
	Nonce         uint32
}

// NewTransaction eturns the transaction to be sent
func NewTransaction(from, to, payload []byte) *Transaction {

	t := Transaction{
		Header: TransactionHeader{
			From: from,
			To:   to,
		},
		Payload: payload,
	}

	t.Header.Timestamp = uint32(time.Now().Unix())
	hash := sha256.New()
	hash.Write(t.Payload)
	t.Header.PayloadHash = hash.Sum(nil)
	t.Header.PayloadLength = uint32(len(t.Payload))

	return &t
}

// Hash Creates transaction hash
func (t *Transaction) Hash() []byte {

	headerBytes, _ := t.Header.MarshalBinary()
	hash := sha256.New()
	hash.Write(headerBytes)
	return hash.Sum(nil)
}

// Sign returns the transaction signed
func (t *Transaction) Sign(keypair *Keypair) []byte {

	s, _ := keypair.Sign(t.Hash())

	return s
}

// VerifyTransaction verify if a transaction authenticity
func (t *Transaction) VerifyTransaction(pow []byte) bool {

	headerHash := t.Hash()
	tempHash := sha256.New()
	tempHash.Write(t.Payload)
	payloadHash := tempHash.Sum(nil)

	return reflect.DeepEqual(payloadHash, t.Header.PayloadHash) && CheckProofOfWork(pow, headerHash) && SignatureVerify(t.Header.From, t.Signature, headerHash)
}

// GenerateNonce generates the valid transaction nonce
// Returns an integer as valid nonce
func (t *Transaction) GenerateNonce(prefix []byte) uint32 {

	newT := t
	for {

		if CheckProofOfWork(prefix, newT.Hash()) {
			break
		}

		newT.Header.Nonce++
	}

	return newT.Header.Nonce
}

// TransactionSlice nedded for block body
// the slice of transactions
type TransactionSlice []Transaction

// Len gets the length of TransactionSlice
func (slice TransactionSlice) Len() int {

	return len(slice)
}

// Exists verify if the transaction exists within the transactionSlice
func (slice TransactionSlice) Exists(tr Transaction) bool {

	for _, t := range slice {
		if reflect.DeepEqual(t.Signature, tr.Signature) {
			return true
		}
	}
	return false
}

// AddTransaction append a transaction into TransactionSlice
func (slice TransactionSlice) AddTransaction(t Transaction) TransactionSlice {

	// Inserted sorted by timestamp
	for i, tr := range slice {
		if tr.Header.Timestamp >= t.Header.Timestamp {
			return append(append(slice[:i], t), slice[i:]...)
		}
	}

	return append(slice, t)
}

// PreviousTransaction get the previous transaction
func (slice TransactionSlice) PreviousTransaction() *Transaction {

	l := len(slice)
	if l == 0 {
		return nil
	}
	return &slice[l-1]
}

// MarshalBinary is Transaction Transformation into []bytes
func (t *Transaction) MarshalBinary() ([]byte, error) {
	buff := new(bytes.Buffer)
	buff.Write(FitBytesInto(t.Payload, NETWORK_KEY_SIZE))
	buff.Write(FitBytesInto(t.Signature, NETWORK_KEY_SIZE))

	thb, err := t.Header.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buff.Write(FitBytesInto(thb, TRANSACTION_HEADER_SIZE))

	return buff.Bytes(), nil
}

// MarshalBinary is TransactionHeader transformation into []byte
// little-endian is needed for some implementations and calculations
func (th *TransactionHeader) MarshalBinary() ([]byte, error) {
	buff := new(bytes.Buffer)

	buff.Write(FitBytesInto(th.From, NETWORK_KEY_SIZE))
	buff.Write(FitBytesInto(th.To, NETWORK_KEY_SIZE))
	binary.Write(buff, binary.LittleEndian, th.Timestamp)
	buff.Write(FitBytesInto(th.PayloadHash, 32))
	binary.Write(buff, binary.LittleEndian, th.PayloadLength)
	binary.Write(buff, binary.LittleEndian, th.Nonce)

	return buff.Bytes(), nil
}

// MarshalBinary is TransactionSlice transformation into []byte
func (slice *TransactionSlice) MarshalBinary() ([]byte, error) {
	buff := new(bytes.Buffer)

	// Iterate through slice toget the binary of the transaction
	// and white it to the buffer
	for _, t := range *slice {

		binaryTx, err := t.MarshalBinary()
		if err != nil {
			return nil, err
		}

		buff.Write(binaryTx)
	}

	return buff.Bytes(), nil
}
