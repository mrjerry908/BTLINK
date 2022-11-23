package database

import (
	"errors"
	"github.com/boltdb/"
	log "github.com/corgi-kx/logcustom"
)

//goland:noinspection ALL
func (bd *BlockchainDB) Put(k, v []byte, bt BucketType) {
	var DBFileName = "blockchain_" + ListenPort + ".db"
	db, err := bolt.Open(DBFileName, 0600, nil)
	defer db.Close()
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bt))
		if bucket == nil {
			var err error
			bucket, err = tx.CreateBucket([]byte(bt))
			if err != nil {
				log.Panic(err)
			}
		}
		err := bucket.Put(k, v)
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}


func (bd *BlockchainDB) View(k []byte, bt BucketType) []byte {
	var DBFileName = "blockchain_" + ListenPort + ".db"
	db, err := bolt.Open(DBFileName, 0600, nil)
	defer db.Close()
	if err != nil {
		log.Panic(err)
	}

	var result []byte
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bt))
		if bucket == nil {

			return errors.New(msg)
		}
		result = bucket.Get(k)
		return nil
	})
	if err != nil {
		//log.Warn(err)
		return nil
	}
	realResult := make([]byte, len(result))
	copy(realResult, result)
	return realResult
}


//goland:noinspection ALL
func (bd *BlockchainDB) Delete(k []byte, bt BucketType) bool {
	var DBFileName = "blockchain_" + ListenPort + ".db"
	db, err := bolt.Open(DBFileName, 0600, nil)
	defer db.Close()
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bt))
		if bucket == nil {
			msg := "datebase delete warnning:" + string(bt)
			return errors.New(msg)
		}
		err := bucket.Delete(k)
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return true
}


//goland:noinspection ALL
func (bd *BlockchainDB) DeleteBucket(bt BucketType) bool {
	var DBFileName = "blockchain_" + ListenPort + ".db"
	db, err := bolt.Open(DBFileName, 0600, nil)
	defer db.Close()
	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(bt))
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return true
}
