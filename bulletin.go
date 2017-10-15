package bulletin

import (
	"fmt"

	"github.com/couchbase/moss"
)

// KV stores one key-value pair.
type KV struct {
	Key []byte
	Val []byte
}

// Board is the bulletin board for posting and exchanging KV.
type Board struct {
	cfg        *Config
	store      *moss.Store
	collection moss.Collection
}

func NewBoard(cfg *Config) *Board {
	b := &Board{
		cfg: cfg,
	}
	return b
}

func (b *Board) openMoss() error {

	mossStore, c, err := moss.OpenStoreCollection(b.cfg.MossPath,
		moss.StoreOptions{}, moss.StorePersistOptions{})
	if err != nil {
		return err
	}
	b.store = mossStore
	b.collection = c
	return nil
}

func (b *Board) writeMossKV(kv *KV) error {
	batch, err := b.collection.NewBatch(0, 0)
	if err != nil {
		return err
	}
	defer batch.Close()

	err = batch.Set(kv.Key, kv.Val)
	if err != nil {
		return err
	}

	err = b.collection.ExecuteBatch(batch, moss.WriteOptions{})
	// On my mac book pro laptop, adding this comment will
	// make test 101 succeed, where it fails otherwise. So
	// it appears there is a background data race to complete
	// the write.
	//p("ExecuteBatch finished with err = %v", err)
	return err
}

func (b *Board) readMoss(key []byte) (val []byte, err error) {
	ss, err := b.collection.Snapshot()
	if err != nil {
		return val, err
	}
	if ss == nil {
		return val, fmt.Errorf("Unable to take moss collection snapshot")
	}
	defer ss.Close()

	ropts := moss.ReadOptions{}

	return ss.Get(key, ropts)
}

// listMoss gives back just the keys
func (b *Board) listMoss() (kv []*KV, err error) {

	ss, err := b.collection.Snapshot()
	if err != nil {
		return kv, err
	}
	if ss == nil {
		return kv, fmt.Errorf("Unable to take moss collection snapshot")
	}
	defer ss.Close()

	iter, err := ss.StartIterator(nil, nil, moss.IteratorOptions{})
	if err != nil {
		return kv, err
	}
	if iter == nil {
		return kv, fmt.Errorf("Unable to get moss collection iterator")
	}
	defer iter.Close()

	for {
		k, _, err := iter.Current()
		if err == moss.ErrIteratorDone {
			return kv, nil
		}
		key2 := make([]byte, len(k))
		copy(key2, k)

		kv = append(kv, &KV{Key: key2})

		err = iter.Next()
		if err == moss.ErrIteratorDone {
			break
		}
	}

	return kv, nil
}

func (b *Board) closeMoss() error {

	b.store.Persist(nil, moss.StorePersistOptions{})

	errs := ""
	err := b.collection.Close()
	if err != nil {
		errs = fmt.Sprintf("b.collection.Close error: '%s'", err.Error())
	}
	err2 := b.store.Close()
	if err2 != nil {
		errs += fmt.Sprintf(" and b.store.Close error: '%s'", err2.Error())
	}
	if errs == "" {
		return nil
	} else {
		return fmt.Errorf(errs)
	}
}
