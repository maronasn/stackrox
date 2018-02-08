package boltdb

import (
	"fmt"

	"bitbucket.org/stack-rox/apollo/central/db"
	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	"bitbucket.org/stack-rox/apollo/pkg/uuid"
	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
)

const notifierBucket = "notifiers"

func (b *BoltDB) getNotifier(id string, bucket *bolt.Bucket) (notifier *v1.Notifier, exists bool, err error) {
	notifier = new(v1.Notifier)
	val := bucket.Get([]byte(id))
	if val == nil {
		return
	}
	exists = true
	err = proto.Unmarshal(val, notifier)
	return
}

// GetNotifier returns notifier with given id.
func (b *BoltDB) GetNotifier(id string) (notifier *v1.Notifier, exists bool, err error) {
	err = b.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(notifierBucket))
		notifier, exists, err = b.getNotifier(id, bucket)
		return err
	})
	return
}

// GetNotifiers retrieves notifiers matching the request from bolt
func (b *BoltDB) GetNotifiers(request *v1.GetNotifiersRequest) ([]*v1.Notifier, error) {
	var notifiers []*v1.Notifier
	err := b.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(notifierBucket))
		b.ForEach(func(k, v []byte) error {
			var notifier v1.Notifier
			if err := proto.Unmarshal(v, &notifier); err != nil {
				return err
			}
			notifiers = append(notifiers, &notifier)
			return nil
		})
		return nil
	})
	return notifiers, err
}

// AddNotifier adds a notifier to bolt
func (b *BoltDB) AddNotifier(notifier *v1.Notifier) (string, error) {
	notifier.Id = uuid.NewV4().String()
	err := b.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(notifierBucket))
		_, exists, err := b.getNotifier(notifier.GetId(), bucket)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("Notifier %v (%v) cannot be added because it already exists", notifier.GetName(), notifier.GetId())
		}
		bytes, err := proto.Marshal(notifier)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(notifier.GetId()), bytes)
	})
	return notifier.Id, err
}

// UpdateNotifier updates a notifier to bolt
func (b *BoltDB) UpdateNotifier(notifier *v1.Notifier) error {
	return b.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(notifierBucket))
		bytes, err := proto.Marshal(notifier)
		if err != nil {
			return err
		}
		return b.Put([]byte(notifier.GetId()), bytes)
	})
}

// RemoveNotifier removes a notifier.
func (b *BoltDB) RemoveNotifier(id string) error {
	return b.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(notifierBucket))
		key := []byte(id)
		if exists := b.Get(key) != nil; !exists {
			return db.ErrNotFound{Type: "Notifier", ID: string(key)}
		}
		return b.Delete(key)
	})
}
