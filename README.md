# jphotos

jphotos is a simple HTTP server for sharing pictures.

jphotos will be written in Go and use Amazon S3 as the storage backend.
The project will fill a long-time desire to have a platform to share my
photos, that isn't Flickr, or 500px, or Format, Imgur, etc.

photos.jrgnsn.net is eventually where this app will go live and will be
where I host pictures, photography, etc.

## Features

* No "social media" features
* No JavaScript!
* Blazingly fast

## Building and Running

**Dependencies:**

* go (>-=1.12)

		go build cmd/jphotos-server/
		./jphotos-server

in the root of the repository.

## Developing

### Live reloading for development

Live reloading can be accomplished with [codegangsta/gin][gin] by
running:

	gin --build cmd/jbooks-server

in the root of the repository.

* add `--all` to watch HTML templates
* add `--excludeDir data` so gin doesn't rebuild after a photo upload

***Note:*** *I've had issues attempting to `go get` gin from within
jphotos, so I'd recommend installing gin from outside the project
directory.*

[gin]: https://github.com/codegangsta/gin

### Initializing the database

First, make sure you have postgres installed.

Then, import the schema:

	psql < $REPO/db/sql/schema.sql

* [Debian PostgreSQL Help][debian-postgres]
* [macOS PostgreSQL Help][macos-postgres]

[debian-postgres]:https://man.sr.ht/~mjorgensen/jphotos/debian_postgresql.md
[macos-postgres]:https://man.sr.ht/~mjorgensen/jphotos/macos_postgresql.md

## Resources

Once documentation exists, it will be able to be [found here][man]

Discussion and patches are welcome and should be directed to my public
inbox for now: [~mjorgensen/public-inbox@lists.sr.ht][lists]. Please use
`--subject-prefix PATCH jphotos` for clarity when sending
patches.

Bugs, issues, planning, and tasks can all be found at the tracker: 
[~mjorgensen/jphotos][todo]

[man]: https://man.sr.ht/~mjorgensen/jphotos
[lists]: https://lists.sr.ht/~mjorgensen/public-inbox
[todo]: https://todo.sr.ht./~mjorgensen/jphotos