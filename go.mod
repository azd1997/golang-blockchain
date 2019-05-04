module github.com/azd1997/golang-blockchain

require (
	github.com/AndreasBriese/bbloom v0.0.0-20190306092124-e2d15f34fcf9 // indirect
	github.com/dgraph-io/badger v1.5.4
	github.com/dgryski/go-farm v0.0.0-20190423205320-6a90982ecee2 // indirect
	github.com/golang/protobuf v1.3.1 // indirect
	github.com/mr-tron/base58 v1.1.2
	github.com/pkg/errors v0.8.1 // indirect
	github.com/stretchr/testify v1.3.0 // indirect
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2
	golang.org/x/net v0.0.0-20190311183353-d8887717615a // indirect
	golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a
	golang.org/x/text v0.3.2
	golang.org/x/tools v0.0.0-20190328211700-ab21143f2384
)

replace (
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2 => github.com/golang/crypto v0.0.0-20190308221718-c2843e01d9a2
	golang.org/x/net v0.0.0-20190311183353-d8887717615a => github.com/golang/net v0.0.0-20190311183353-d8887717615a
	golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a => github.com/golang/sys v0.0.0-20190215142949-d0b11bdaac8a
	golang.org/x/tools v0.0.0-20190328211700-ab21143f2384 => github.com/golang/tools v0.0.0-20190328211700-ab21143f2384
	golang.org/x/text v0.3.2 => github.com/golang/text v0.3.2
)

go 1.12
