package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/dgraph-io/badger/v3"
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
			val, _ := nodeIterator.Item().ValueCopy(nil)

			followEntry := &FollowEntry{}
			gob.NewDecoder(bytes.NewReader(val)).Decode(followEntry)
    		fmt.Println(string(followEntry.FollowedPKID[:]), ' ', string(followEntry.FollowedPKID[:]))
		}
		return nil
	})
}

type FollowEntry struct {
	FollowerPKID *PKID
	FollowedPKID *PKID
	isDeleted bool
}

type PKID [33]byte