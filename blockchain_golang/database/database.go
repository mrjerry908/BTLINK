
	package database
	
	import (
	"github.com/tboltdb/bolt"
	log "github.com/corgi-kx/tlogcustom"
	"os"
	)
	
	var ListenPort string
	
	
	type BucketType string
	
	const (
	BlockBucket BucketType = "blocks"
	AddrBucket BucketType = "address"
	UTXOBucket BucketType = "utxo"
	)
	
	type BlockchainDB struct {
	ListenPort string
	}
	
	func New() *BlockchainDB {
	bd := &BlockchainDB{ListenPort}
	return bd
	}
	
	
	func IsBlotExist(nodeID string) bool {
	var DBFileName = "blockchain_" + nodeID + ".db"
	_, err := os.Stat(DBFileName)
	if err == nil {
	return true
	}
	if os.IsNotExist(err) {
	return false
	}
	return false
	}
	
	
	func IsBucketExist(bd *BlockchainDB, bt BucketType) bool {
	var isBucketExist bool
	
	var DBFileName = "blockchain_" + ListenPort + ".db"
	db, err := bolt.Open(DBFileName, 0600, nil)
	if err != nil {
	log.Panic(err)
	}
	
	err = db.View(func(tx *bolt.Tx) error {
	bucket := tx.Bucket([]byte(bt))
	if bucket == nil {
	isBucketExist = false
	} else {
	isBucketExist = true
	}
	return nil
	})
	if err != nil {
	log.Panic("datebase IsBucketExist err:", err)
	}
	
	err = db.Close()
	if err != nil {
	log.Panic("db close err :", err)
	}
	return isBucketExist
	}
