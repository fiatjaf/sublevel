# sublevel

Separate sections of the same LevelDB. Compatible (at least in the basics -- I didn't test it minuciously) with [the nodejs sublevel](https://github.com/dominictarr/level-sublevel).

[![Travis-CI build status](https://travis-ci.org/fiatjaf/sublevel.svg)](https://travis-ci.org/fiatjaf/sublevel)
[![API documentation on Godoc.org](https://img.shields.io/badge/godoc-reference-orange.svg)](https://godoc.org/github.com/fiatjaf/sublevel)

```go
import "github.com/fiatjaf/sublevel"

sub := sublevel.MustOpen("example.db").Sub("specific-stuff")

sub.Put([]byte("this"), []byte("2007-04-01"), nil)
dateOfThis := sub.Get([]byte("this"), nil)
sub.Delete([]byte("this"))
```

### Batch operations

```go
db, err := sublevel.Open("example.db")
if err != nil {
    panic(err)
}

// batch on a single sublevel:
sub := db.Sub("some-things")
batch := sub.NewBatch()
batch.Put([]byte("newthing"), []byte("true"))
batch.Delete([]byte("oldthing"))
_ := sub.Write(batch)
// committed.
// ~
// batch on different sublevels:
othersub := db.Sub("other-things")
otherbatch := othersub.NewBatch()
otherbatch.Put([]byte("new-other-thing"), []byte("true"))
otherbatch.Delete([]byte("old-other-thing"))

batchagain := sub.NewBatch()
batchagain.Delete([]byte("newthing"))
batchagain.Put([]byte("newestthing"), []byte("true"))

superbatch := db.MultiBatch(otherbatch, batchagain)
_ := db.Write(superbatch, nil)
// committed.
// ~
// or
b := db.NewBatch()
b.Put("some-things", []byte("newthing"), []byte("true"))
b.Delete("other-things", []byte("old-other-thing"))
b.Put("other-things", []byte("new-other-thing"), []byte("true"))
b.Delete("some-things", []byte("newthing"))
b.Put("some-things", []byte("newestthing"), []byte("true"))
db.Write(b, nil)
```

**sublevel** is built on top of [goleveldb](http://godoc.org/github.com/syndtr/goleveldb/leveldb) and supports most methods from there (not all, but in most cases everything you'll need).
