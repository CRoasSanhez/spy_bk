package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"

	"github.com/tv42/base58"
)

// Keypair struct
// Public and private key
type Keypair struct {
	Public  []byte `json:"public"`  // base58 (x y)
	Private []byte `json:"private"` // d (base58 encoded)
}

// GenerateNewKeypair generate the private and public
// base58 keys
func GenerateNewKeypair() *Keypair {

	pk, _ := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)

	b := bigJoin(KEY_SIZE, pk.PublicKey.X, pk.PublicKey.Y)

	public := base58.EncodeBig([]byte{}, b)
	private := base58.EncodeBig([]byte{}, pk.D)

	kp := Keypair{Public: public, Private: private}

	return &kp
}

// Sign the keypair using base58
func (k *Keypair) Sign(hash []byte) ([]byte, error) {

	prKeyDecoded, err := base58.DecodeToBig(k.Private)
	if err != nil {
		return nil, err
	}

	puKeyDecoded, _ := base58.DecodeToBig(k.Public)

	pub := splitBig(puKeyDecoded, 2)
	x, y := pub[0], pub[1]

	key := ecdsa.PrivateKey{
		ecdsa.PublicKey{
			elliptic.P224(),
			x,
			y,
		},
		prKeyDecoded,
	}

	r, s, _ := ecdsa.Sign(rand.Reader, &key, hash)

	return base58.EncodeBig([]byte{}, bigJoin(KEY_SIZE, r, s)), nil
}

// SignatureVerify validate the hash according to
// the given publicKey and the signature
func SignatureVerify(publicKey, sig, hash []byte) bool {

	bytesDecded, _ := base58.DecodeToBig(publicKey)
	publ := splitBig(bytesDecded, 2)
	x, y := publ[0], publ[1]

	bytesDecded, _ = base58.DecodeToBig(sig)
	sigg := splitBig(bytesDecded, 2)
	r, s := sigg[0], sigg[1]

	pub := ecdsa.PublicKey{elliptic.P224(), x, y}

	return ecdsa.Verify(&pub, hash, r, s)
}

// bigJoin appends the
func bigJoin(expectedLen int, bigs ...*big.Int) *big.Int {

	bs := []byte{}
	//newString := []byte{}

	for i, b := range bigs {

		bigIntBytes := b.Bytes()
		dif := expectedLen - len(bigIntBytes)
		if dif > 0 && i != 0 {

			bigIntBytes = append(SliceOfBytes(dif, 0), bigIntBytes...)
		}

		bs = append(bs, bigIntBytes...)
	}

	b := new(big.Int).SetBytes(bs)

	return b
}

// splitBig splits
func splitBig(bytesDecded *big.Int, parts int) []*big.Int {

	// returns the absolute value of bytesDecded
	// as a big-endian byte slice
	byteSlice := bytesDecded.Bytes()

	// if byteSlice length is none
	// a new  byte is append to byteSlice
	if len(byteSlice)%2 != 0 {
		byteSlice = append([]byte{0}, byteSlice...)
	}

	// clacultate the length of the middle byteSlice
	l := len(byteSlice) / parts

	// create a new slice with length equal
	// to number of parts
	as := make([]*big.Int, parts)

	for i := range as {

		// Generates a new bigInt pointer and set the bytes
		// of the incoming BigInt parameter
		as[i] = new(big.Int).SetBytes(byteSlice[i*l : (i+1)*l])
	}

	return as

}
