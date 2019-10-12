# jphotos

jphotos is a simple HTTP server for sharing pictures.

jphotos will be written in Go and use Amazon S3 as the storage backend.
The project will fill a long-time desire to have a platform to share my
photos, that isn't Flickr, or 500px, or Format, etc.

photos.jrgnsn.net is eventually where this app will go live and will be
where I host pictures, photography, etc.

## Building

**Dependencies:**

* go (>-=1.12)

	cd $REPO/cmd/jbooks-server; go build

in the root of the repository.

## Developing

### Live reloading for development

Live reloading can be accomplished with [codegangsta/gin][gin] by
running:

	gin --build cmd/jbooks-server

in the root of the repository. Add `--all` to watch HTML templates,
also. 

### Initializing the database

First, make sure you have postgres installed.

Then, import the schema:

	psql < $REPO/db/sql/schema.sql


## Resources

Once documentation exists, it will be able to be [found here][man].

Discussion and patches are welcome and should be directed to my public
inbox for now: [~mjorgensen/public-inbox@lists.sr.ht][lists]. Please use
`--subject-prefix PATCH jphotos` for clarity when sending
patches.

Bugs, issues, planning, and tasks can all be found at the tracker: 
[~mjorgensen/photos.jrgnsn.net][todo].

[man]: https://man.sr.ht/~mjorgensen/jphotos
[lists]: https://lists.sr.ht/~mjorgensen/public-inbox
[todo]: https://todo.sr.ht./~mjorgensen/jphotos