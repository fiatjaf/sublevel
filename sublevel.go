package sublevel

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"log"
)

/* basic db management */

func OpenFile(dbfile string, o *opt.Options) AbstractLevel {
	db, err := leveldb.OpenFile(dbfile, nil)
	return AbstractLevel{
		leveldb: db,
		err:     err,
	}
}

type AbstractLevel struct {
	leveldb *leveldb.DB
	err     error
}

func (a AbstractLevel) Close() error {
	return a.leveldb.Close()
}

func (a AbstractLevel) Sub(store string) (*Sublevel, error) {
	if a.err != nil {
		return &Sublevel{}, a.err
	}
	return &Sublevel{
		namespace: []byte("!" + store + "!"),
		db:        a.leveldb,
	}, nil
}

func (a AbstractLevel) MustSub(store string) *Sublevel {
	sub, err := a.Sub(store)
	if err != nil {
		log.Fatal("couldn't open database file. ", err)
	}
	return sub
}

type Sublevel struct {
	namespace []byte
	db        *leveldb.DB
}

func (s Sublevel) Close() error {
	return s.db.Close()
}

/* methods */

func (s Sublevel) Delete(key []byte, wo *opt.WriteOptions) error {
	key = append(append([]byte(nil), s.namespace...), key...)
	return s.db.Delete(key, wo)
}

func (s Sublevel) Get(key []byte, ro *opt.ReadOptions) (value []byte, err error) {
	key = append(append([]byte(nil), s.namespace...), key...)
	return s.db.Get(key, ro)
}

func (s Sublevel) Put(key []byte, value []byte, wo *opt.WriteOptions) error {
	key = append(append([]byte(nil), s.namespace...), key...)
	return s.db.Put(key, value, wo)
}

func (s Sublevel) Has(key []byte, ro *opt.ReadOptions) (ret bool, err error) {
	key = append(append([]byte(nil), s.namespace...), key...)
	return s.db.Has(key, ro)
}

/* iterator */
func (s Sublevel) NewIterator(slice *util.Range, ro *opt.ReadOptions) SubIterator {
	slice = &util.Range{
		Start: append(append([]byte(nil), s.namespace...), slice.Start...),
		Limit: append(append([]byte(nil), s.namespace...), slice.Limit...),
	}

	return SubIterator{
		namespace: s.namespace,
		iterator:  s.db.NewIterator(slice, ro),
	}
}

type SubIterator struct {
	namespace []byte
	iterator  iterator.Iterator
}

func (si SubIterator) Key() []byte {
	key := si.iterator.Key()
	return key[len(si.namespace):]
}
func (si SubIterator) Value() []byte {
	return si.iterator.Value()
}
func (si SubIterator) Next() bool {
	return si.iterator.Next()
}
func (si SubIterator) Prev() bool {
	return si.iterator.Prev()
}
func (si SubIterator) Last() bool {
	return si.iterator.Last()
}
func (si SubIterator) First() bool {
	return si.iterator.First()
}
func (si SubIterator) Seek(key []byte) bool {
	key = append(append([]byte(nil), si.namespace...), key...)
	return si.iterator.Seek(key)
}
func (si SubIterator) Release() {
	si.iterator.Release()
}
func (si SubIterator) Error() error {
	return si.iterator.Error()
}

/* transactions */
func (s Sublevel) NewBatch() *Batch {
	return &Batch{
		namespace: s.namespace,
		batch:     new(leveldb.Batch),
	}
}

type Batch struct {
	namespace []byte
	batch     *leveldb.Batch
}

func (b *Batch) Delete(key []byte) {
	key = append(append([]byte(nil), b.namespace...), key...)
	b.batch.Delete(key)
}

func (b *Batch) Put(key []byte, value []byte) {
	key = append(append([]byte(nil), b.namespace...), key...)
	b.batch.Put(key, value)
}

func (b *Batch) Len() int {
	return b.batch.Len()
}

func (b *Batch) Reset() {
	b.batch.Reset()
}

func (s Sublevel) Write(b *Batch, wo *opt.WriteOptions) (err error) {
	return s.db.Write(b.batch, wo)
}
