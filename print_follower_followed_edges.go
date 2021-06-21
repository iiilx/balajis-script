package main

import (
	"bytes"
	"fmt"
	"encoding/gob"
	"github.com/dgraph-io/badger/v3"
)
var _PrefixPKIDToProfileEntry = []byte{23}
const (
	PubKeyBytesLenCompressed = 33
	PrefixFollowerPKIDToFollowedPKID = byte(28)
)

func _dbKeyForPKIDToProfileEntry(pkid *PKID) []byte {
	prefixCopy := append([]byte{}, _PrefixPKIDToProfileEntry...)
	key := append(prefixCopy, pkid[:]...)
	return key
}

func PrintEdges(dir string) {
	db, _ := badger.Open(badger.DefaultOptions(dir))
	defer db.Close()
	iterateFollowEntries(db, []byte{PrefixFollowerPKIDToFollowedPKID})
}


func iterateFollowEntries(db *badger.DB, dbPrefix []byte) {
	db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		nodeIterator := txn.NewIterator(opts)
		defer nodeIterator.Close()
		prefix := dbPrefix

		for nodeIterator.Seek(prefix); nodeIterator.ValidForPrefix(prefix); nodeIterator.Next() {
			key := nodeIterator.Item().Key()
			followerPKIDBytes := key[1:PubKeyBytesLenCompressed+1]
			followerPKID := &PKID{}
			copy(followerPKID[:], followerPKIDBytes)

			followedPKIDBytes := key[1+PubKeyBytesLenCompressed:]
			followedPKID := &PKID{}
			copy(followedPKID[:], followedPKIDBytes)
			followedProfileKey := _dbKeyForPKIDToProfileEntry(followedPKID)
			followedProfileItem, err := txn.Get(followedProfileKey)
			if err != nil {
				fmt.Println(err)
				return err
			}
			followedProfileEntry := &ProfileEntry{}
			err = followedProfileItem.Value(func(val []byte) error {
				gob.NewDecoder(bytes.NewReader(val)).Decode(followedProfileEntry)
				return nil
			})
			if err != nil {
				fmt.Println(err)
				return err
			}

			followerProfileKey := _dbKeyForPKIDToProfileEntry(followerPKID)
			followerProfileItem, err := txn.Get(followerProfileKey)
			if err != nil {
				fmt.Println(err)
				return err
			}
			followerProfileEntry := &ProfileEntry{}
			err = followerProfileItem.Value(func(val []byte) error {
				gob.NewDecoder(bytes.NewReader(val)).Decode(followerProfileEntry)
				return nil
			})
			if err != nil {
				fmt.Println(err)
				return err
			}
			fmt.Println(string(followedProfileEntry.Username), string(followerProfileEntry.Username))
		}
		return nil
	})
}
