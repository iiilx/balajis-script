package main

import (
	"fmt"
	"github.com/dgraph-io/badger/v3"
)

const (
	PubKeyBytesLenCompressed = 33
)


func PrintEdges(dir string) {
	db, _ := badger.Open(badger.DefaultOptions(dir))
	defer db.Close()
	prefixFollowerPKIDToFollowedPKID := byte(28)
	iterateFollowEntries(db, []byte{prefixFollowerPKIDToFollowedPKID})
}


func iterateFollowEntries(db *badger.DB, dbPrefix []byte) {
	db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		nodeIterator := txn.NewIterator(opts)
		defer nodeIterator.Close()
		prefix := dbPrefix

		for nodeIterator.Seek(prefix); nodeIterator.ValidForPrefix(prefix); nodeIterator.Next() {
			key := nodeIterator.Item().Key()
			followerPKIDBytes := key[1:PubKeyBytesLenCompressed]
			followerPKID := &PKID{}
			copy(followerPKID[:], followerPKIDBytes)

			followedPKIDBytes := key[1+PubKeyBytesLenCompressed:]
			followedPKID := &PKID{}
			copy(followedPKID[:], followedPKIDBytes)
			fmt.Println(string(followerPKID[:]), ' ', string(followedPKID[:]))
		}
		return nil
	})
}
type PKID [33]byte