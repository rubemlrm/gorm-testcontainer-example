# Motivation

One of the main motivations for this little project was to find a way of how test database integration with Golang and test containers.

# Used Technologies

For this project was used the following technologies:
* [Golang](https://go.dev/) 
* [TestContainers](https://golang.testcontainers.org/)
* [Gorm](https://gorm.io/) as ORM
* [Goose](https://pressly.github.io/goose/) as migration tool
* [go-faker](https://github.com/go-faker/faker) as faker data generator


# Implementation

For this implementation is possible to separate into several parts:
- How to handle migrations
- Preparing the environment for testing

## How to handle migrations

For the migrations we can create and apply migrations using on the following options:

### Gorm Way

Using the native way of Gorm we can apply migrations using the golang struct type and applying this changes with the [Auto Migration](https://gorm.io/docs/migration.html#Auto-Migration), but after a run 
test was possible to find that some field types were not created with the same settings, for instance the following struct will create the createdAt fields as *timestampz* while the sql counterpart will create as *timestamp*:

```go
type Book struct {
	ID        uint      `gorm:"primaryKey"`
	Title     string    `gorm:"not null;type:varchar(12)"`
	Author    string    `gorm:"not null;type:varchar(12)"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP()"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP()"`
}
```

```
create table books
(
  created_at timestamp without time zone default (now() at time zone 'utc'),
  updated_at timestamp without time zone default (now() at time zone 'utc')
);
```

Despite the miss configuration between sql and go struct approach, based on their documentation we will find another caveat, relative to the drop of 
unused fields. By implementation the Gorm auto migration util will not remove old unused fields from the tables, this is implemented as safe measure to avoid
loose information, but will require the user to manually delete that.Based on Gorm documentations looks that it's not possible to version the migrations changes between struct changes.

Another requirement to use this functionality, is to inject all structs manually on the `db.AutoMigrate` call, while this can useful for testing or small applications, when the applications
start to grow, could be troublesome to take track of every struct that was used and needs to migrate.

### Goose way

Since Goose is an agnostic migration tool, it's possible to be used with several project stacks. It's also possible to create new migrations in plain sql or using go code. The team that devoleped Goose
offer a way to run migrations inside the application, like the AutoMigrate but for this case are called [Embedded sql migrations](https://github.com/pressly/goose#embedded-sql-migrations).

Based on their example from documentation, it's possible to use this with the `go:embed` feature of golang that allows us to include additional files in compilation time, but this feature have 
a condition, is not possible to use paths that points to a directory above or a relative path to a directory outside the go file scope. For instance:

```
cmd
 -- app
    -- main.go
internal
    -- internal code
migrations
    -- migration files
```

More and more applications are following the above pattern to structure their files, with the this pattern our *migrations* path should be included next to the file that will use the `go:embed` directive.
This will mean that developer will include sql structure files in a directory that was not created for that porpuse and can create some entropy to the developers that came from other languages.
Since we want to use this for integration tests will mean that inside the **app** directory we will need to include an integration test directory. 
To surpass this limitation is possible to create an `fs.fs` variable pointing to a directory that we present and still use the embed migrations but in a custom way.

```go
func RunMigrations(dsn string) error {
	var sqlMigrations *sql.DB
	sqlMigrations, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../migrations")

	files := os.DirFS(dir)
	goose.SetBaseFS(files)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(sqlMigrations, "."); err != nil {
		panic(err)
	}
	return nil
}
```

This way it is possible to organize the project migrations structure in several ways, without losing the possibility to run the embed migrations. 

**Implementation**

For this implementation was used the *Goose way*, since it's possible to integrate this way with other ORM's, it's possible to have version for our migrations,
goose itself manages the migrations what already run, this is more important for application in prod than integration testing, since the container will be recreated and deleted by request.

## Preparing the environment for testing

To prepare the testing environment will be necessary the following aspects:
* Testcontainer that will run a database
* Testsuite implementation using [testify](https://pkg.go.dev/github.com/stretchr/testify/)
* Setup two database connection objects
  * sql connection that will be used to run migrations in tests
  * gorm connection that will handle all the operations against database

## Testcontainer that will run our database

The container that will run the testing database, will be initatiate at the `setup` step of this testing and be removed on the tear down step.
Would be possible to run per single test, but this would be more resource and time intensive for the tests. 

## Testsuite implementation

The test suite will run only the general setup and tear down steps for this proof of concept. It's possible to setup before/after tests too 
and for some situations can make sense to have that for dropping all migrations and run all migrations again to ensure the cleanest playground for 
each test.

## Setup two database connection objects

Since this project uses two tools to handle database operations, one for migrations and another one for data handling. It's necessary to 
set up two different connections, since both tools are waiting for a different kind of object to initiate the connection. For instance:

```golang

// goose
var sqlMigrations *sql.DB
sqlMigrations, err := sql.Open("postgres", dsn)
// gorm
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Error)})


```


# References

* [Integration-test With golang](https://github.com/underbek/integration-test-go)
* [Embend Migrations With goose](https://github.com/pressly/goose#embedded-sql-migrations)
* [Go Embed RFC relative to relative paths](https://go.googlesource.com/proposal/+/master/design/draft-embed.md#go_embed-directives)