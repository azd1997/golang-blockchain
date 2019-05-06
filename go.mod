module github.com/azd1997/golang-blockchain

require (
	github.com/AndreasBriese/bbloom v0.0.0-20190306092124-e2d15f34fcf9 // indirect
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575 // indirect
	github.com/dgraph-io/badger v1.5.4
	github.com/dgryski/go-farm v0.0.0-20190423205320-6a90982ecee2 // indirect
	github.com/golang/protobuf v1.3.1 // indirect
	github.com/mr-tron/base58 v1.1.2
	github.com/pkg/errors v0.8.1 // indirect
	github.com/smartystreets/goconvey v0.0.0-20190330032615-68dc04aab96a // indirect
	github.com/stretchr/testify v1.3.0 // indirect
	github.com/vrecan/death v3.0.1+incompatible
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2
)

replace (
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2 => github.com/golang/crypto v0.0.0-20190308221718-c2843e01d9a2
	golang.org/x/net v0.0.0-20190311183353-d8887717615a => github.com/golang/net v0.0.0-20190311183353-d8887717615a
	golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a => github.com/golang/sys v0.0.0-20190215142949-d0b11bdaac8a
	golang.org/x/text v0.3.0 => github.com/golang/text v0.3.0
	golang.org/x/tools v0.0.0-20180917221912-90fa682c2a6e => github.com/golang/tools v0.0.0-20180917221912-90fa682c2a6e
	golang.org/x/tools v0.0.0-20190328211700-ab21143f2384 => github.com/golang/tools v0.0.0-20190328211700-ab21143f2384
)

go 1.12
