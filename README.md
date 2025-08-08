# bootdev_blogagg_go
bootdev project for an RSS feed aggregator in Go

build an RSS feed aggregator in Go! We'll call it "Gator"


Use [Go language](https://sqlc.dev/), [goose](https://github.com/pressly/goose), and [sqlc](https://sqlc.dev/).

`dnf install postgresql postgresql-contrib postgresql-server golang golang-github-pressly-goose`

[Fedora Postgresql quick](https://docs.fedoraproject.org/en-US/quick-docs/postgresql/)


`mkdir -p ~/.config/bootdev`
and place gatorconfig.json in it. 
`{"Db_url":"postgres://bootdev:@/gator?sslmode=disable","User_name":"kahya"}

1. Configure Postgrsql
2. cd sql/schema and do a gator up with the proper connection url
3. run it with help command to see list of commands. Requires to regester some users to use the features

