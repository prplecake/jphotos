module github.com/prplecake/jphotos

go 1.18

require (
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/golang-migrate/migrate/v4 v4.17.1
	github.com/gorilla/mux v1.8.1
	github.com/lib/pq v1.10.9
	github.com/metal3d/go-slugify v0.0.0-20160607203414-7ac2014b2f23
	github.com/prplecake/go-thumbnail v0.1.6
	github.com/qustavo/dotsql v1.2.0
	golang.org/x/crypto v0.25.0
	golang.org/x/term v0.22.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	golang.org/x/image v0.18.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
)

//replace github.com/prplecake/go-thumbnail => /Users/mjorgensen/Projects/go-thumbnail
