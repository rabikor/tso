# Treatment Scheme Organizer

This project adheres to [Semantic Versioning](http://semver.org/)
and [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0-beta.2/).

## DB Models

### Drug

Fields:

```
- title (100 | varchar)
```

Actions:

```
1. Creating new items
2. Getting list
```

## Database

Command for getting status of database
```shell
sql-migrate status -config=dbconfig.yml -env {env_name}
```

Command for running migrations to database
```shell
sql-migrate up -config=dbconfig.yml -env {env_name}
```

Command for creating new migration in database
```shell
sql-migrate new -config=dbconfig.yml -env {env_name} {migration_name}
```
