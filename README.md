# jphotos

[![builds.sr.ht
status](https://builds.sr.ht/~mjorgensen/jphotos.svg)](https://builds.sr.ht/~mjorgensen/jphotos?)

jphotos is a simple HTTP server for sharing pictures.

## Features

* No "social media" features
* No JavaScript!
* Blazingly fast

## Getting Started

Read our [**Getting Started**][getting-started] guide.

[getting-started]:https://man.sr.ht/~mjorgensen/jphotos/getting_started.md

## Developing

### Live reloading for development

Live reloading can be accomplished with [codegangsta/gin][gin] by
running:

```
$ gin --build ./cmd/jphotos-server [OPTION...]
```

in the root of the repository.

Options:

* `--all` - to rebuild on template modifications
* `--excludeDir data` - so gin doesn't rebuild after a photo upload

**Note:** I've had issues attempting to `go get` gin from within
jphotos, so I'd recommend installing gin from outside the project
directory. Try `cd ..; go get github.com/codegangsta/gin; cd -` instead.

[gin]: https://github.com/codegangsta/gin

### Initializing the database

First, set up your database.

* [Debian PostgreSQL Help][debian-postgres]
* [macOS PostgreSQL Help][macos-postgres]

Then, read about [database migrations][db-migrations]

[debian-postgres]:https://man.sr.ht/~mjorgensen/jphotos/debian_postgresql.md
[macos-postgres]:https://man.sr.ht/~mjorgensen/jphotos/macos_postgresql.md
[db-migrations]:https://man.sr.ht/~mjorgensen/jphotos/database_migrations.md

## Resources

Comprehensive documentation [can be found here][man].

Discussion and patches are welcome and should be directed to my public
inbox for now: [~mjorgensen/public-inbox@lists.sr.ht][lists]. Please use
`--subject-prefix PATCH jphotos` for clarity when sending
patches.

Bugs, issues, planning, and tasks can all be found at the tracker: 
[~mjorgensen/jphotos][todo]

[man]: https://man.sr.ht/~mjorgensen/jphotos
[lists]: https://lists.sr.ht/~mjorgensen/public-inbox
[todo]: https://todo.sr.ht./~mjorgensen/jphotos
