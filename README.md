BlobStash
=========

**BlobStash** is both a content-addressed blob store and a key value store accessible via an HTTP API.

Key value pairs are stored as "meta" blobs, this mean you can build application on top of BlobStash without the need for another database.

Initially created to power [BlobSnap](https://github.com/tsileo/blobsnap) and [Blobs](http://blobs.co).

## Features

- [BLAKE2b](https://blake2.net) as hashing algorithm for the content-addressed blob store
- Backend routing, you can define rules to specify where blobs should be stored ("if-meta"...)
- Optional encryption (using [go.crypto/nacl secretbox](http://godoc.org/code.google.com/p/go.crypto/nacl))
- Possibility to incrementally archive blobs to AWS Glacier (with a recovery command-line tool)
- A full featured Go [client](http://godoc.org/github.com/tsileo/blobstash/client) using the HTTP API
- Can be embedded in your go app ([embedded client](http://godoc.org/github.com/tsileo/blobstash/embed))

## Getting started

```console
$ go get github.com/tsileo/blobstash/cmd/blobstash
$ $GOPATH/bin/blobstash
2015/08/13 21:32:27 Starting blobstash version 0.0.0; go1.4 (linux/amd64)
2015/08/13 21:32:27 BlobsFileBackend: starting, opening index
2015/08/13 21:32:27 BlobsFileBackend: scanning BlobsFiles...
2015/08/13 21:32:27 BlobsFileBackend: /data/blobs/blobs-00000 loaded
2015/08/13 21:32:27 BlobsFileBackend: opening /data/blobs/blobs-00000 for writing
2015/08/13 21:32:27 BlobsFileBackend: snappyCompression = false
2015/08/13 21:32:27 BlobsFileBackend: backend id => blobsfile-/data/blobs
2015/08/13 21:32:27 server: HTTP API listening on 0.0.0.0:8050
2015/08/13 21:32:38 Scan: done, 10596 blobs scanned in 11.114966366s, 0 blobs applied
```

## Blob store

You can deal directly with blobs when needed using the HTTP API, full docs [here](docs/blobstore.md).

```console
$ curl -F "c0f1480a26c2fd4deb8e738a52b7530ed111b9bcd17bbb09259ce03f129988c5=ok" http://0.0.0.0:8050/api/v1/blobstore/upload
```

## Key value store

Updates on keys are store in blobs, and automatically handled by BlobStash.

Perfect to keep a mutable pointer.

```console
$ curl -XPUT http://127.0.0.1:8050/api/v1/vkv/key/k1 -d value=v1
{"key":"k1","value":"v1","version":1421705651367957723}
```

```console
$ curl http://127.0.0.1:8050/api/v1/vkv/key/k1            
{"key":"k1","value":"v1","version":1421705651367957723}
```

## Extensions

BlobStash comes with few bundled extensions making it easier to build your own app on top of it.

Extensions only uses the blob store and the key value store, nothing else.

### (WIP) Files

A multipart file upload handler and a downalod handler.

### Document Store

A JSON document store running on top of an HTTP API. Support a subset of the MongoDB Query language.

JSON documents are stored as blobs and the key-value store handle the indexing.

Perfect for building app desined to only store your own data.

### Supported MongoDB query operators

Refers to MongpDB documentation: [query documents](https://docs.mongodb.org/manual/tutorial/query-documents/) and [query operators](https://docs.mongodb.org/manual/reference/operator/query/#query-selectors).

#### Features

- [ ] dot-notation support
- [x] `{}` - Select all documents
- [x] `{ <field>: <value> }` - equality, AND conditions

#### Query operators

##### Comparison

- [ ] `$eq`
- [ ] `$gt`
- [ ] `$gte`
- [ ] `$lt`
- [ ] `$lte`
- [ ] `$ne`
- [ ] `$in`
- [ ] `$nin`

##### Logical

- [x] `$or`
- [x] `$and`
- [ ] `$not`
- [ ] `$nor`

##### Element

- [ ] `$exists`
- [ ] `$type`

##### Evalutation

- [ ] `$regex`

## Backend

Blobs are stored in a backend.

The backend handle operations:

- Put
- Exists
- Get
- Delete
- Enumerate

### Available backends

- [BlobsFile](docs/blobsfile.md) (local disk)
- AWS S3
- Mirror
- AWS Glacier (only as a backup)
- A remote BlobStash instance
- Fallback backend (store failed upload locally and try to reupload them periodically)

- Submit a pull request!

You can combine backend as you wish, e.g. Mirror( Encrypt( S3() ), BlobsFile() ).

## Routing

You can define rules to specify where blobs should be stored, depending on whether it's a meta blob or not, or depending on the namespace it come from.

**Blobs are routed to the first matching rule backend, rules order is important.**

```json
[
    [["if-ns-myhost", "if-meta"], "customHandler2"],
    ["if-ns-myhost", "customHandler"],
    ["if-meta", "metaHandler"],
    ["default", "blobHandler"]
]
```

The minimal router config is:

```json
[["default", "blobHandler"]]
```

## Embedded mode

```go
package main

import (
	"github.com/tsileo/blobstash/server"
)

func main() {
	blobstash := server.New(nil)
	blobstash.SetUp()
	// wait till all meta blobs get scanned
	blobstash.TillReady()
	bs := blobstash.BlobStore()
	kvs := blobstash.KvStore()
	blobstash.TillShutdown()
}
```

## Projects built on top of BlobStash

 - [BlobSnap](https://github.com/tsileo/blobsnap)
 - [BlobFS](https://github.com/tsileo/blobfs)
 - [BlobFS-web](https://github.com/tsileo/blobfs-web)

Make a pull request if your project uses BlobStash as data store.

## Roadmap / Ideas

- A better documentation
- A web interface
- An S3-like HTTP API to store archive
- Fill an issue!

## Contribution

Pull requests are welcome but open an issue to start a discussion before starting something consequent.

Feel free to open an issue if you have any ideas/suggestions!

## Donation

[![Flattr this git repo](http://api.flattr.com/button/flattr-badge-large.png)](https://flattr.com/submit/auto?user_id=tsileo&url=https%3A%2F%2Fgithub.com%2Ftsileo%2Fblobstash)

BTC 12XKk3jEG9KZdZu2Jpr4DHgKVRqwctitvj

## License

Copyright (c) 2014-2015 Thomas Sileo and contributors. Released under the MIT license.
