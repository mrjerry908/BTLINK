package block

import (
	"bytes"
	"encoding/gob"
	"math/big"
	"time"
)

type Block struct {
	
	PreHash []byte

	Transactions []Transaction
	
	TimeStamp int64
	
	Height int

	Nonce int64
	
	Hash []byte
}


func mineBlock(transaction []Transaction, preHash []byte, height int) (*Block, error) {
	timeStamp := time.Now().Unix()
	
	block := Block{preHash, transaction, timeStamp, height, 0, nil}
	pow := NewProofOfWork(&block)
	nonce, hash, err := pow.run()
	if err != nil {
		return nil, err
	}
	block.Nonce = nonce
	block.Hash = hash[:]
	log.Info("pow verify : ", pow.Verify())
	log.Infof( block.Height)
	return &block, nil
}


func newGenesisBlock(transaction []Transaction) *Block {
	
	preHash := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	
	genesisBlock, err := mineBlock(transaction, preHash, 1)
	if err != nil {
		log.Error(err)
	}
	return genesisBlock
}


func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}

func (v *Block) Deserialize(d []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(v)
	if err != nil {
		log.Panic(err)
	}
}

func isGenesisBlock(block *Block) bool {
	var hashInt big.Int
	hashInt.SetBytes(block.PreHash)
	if big.NewInt(0).Cmp(&hashInt) == 0 {
		return true
	}
	return false
}
