
[![Build Status](https://github.com/laurentbh/mac-stats/workflows/Go/badge.svg)](https://github.com/laurentbh/mac-stats/actions)

# mac-stats

Fetches battery and SSD info and stores in a Postgres DB.

Tested under Catalina (10.15.7)

## install
- build `go` binary
- create the database, [see schema](https://github.com/laurentbh/mac-stats/schema.sql)
- change DB credentials, as there is no config yet [see](https://github.com/laurentbh/mac-stats/blob/main/postgres.go#L15-L21)
## requirements:
- [Postgres DB](https://www.postgresql.org/)
- [smartcl](http://www.smartmontools.org)
