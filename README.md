
[![Build Status](https://github.com/laurentbh/mac-stats/workflows/Go/badge.svg)](https://github.com/laurentbh/mac-stats/actions)

# mac-stats

Fetches battery and SSD info and stores in a Postgres DB.

Tested under Catalina (10.15.7)

## install
- build `go` binary
- create the database, [see schema](https://github.com/laurentbh/mac-stats/blob/schema.sql)
- config database, the config file needs to be either in the current dir or in ~/.config

 change the [config file](https://github.com/laurentbh/mac-stats/blob/main/mac-stats.yaml)

## display example
![](doc/unit_read_grafana.png)
With [grafana](https://grafana.com/) and query

```
SELECT
  stamp AS "time",
  ((metrics->>'UnitRead')::integer) AS "Unit Read"
FROM
  ssd
WHERE
  $__timeFilter(stamp)

```
## requirements:
- [Postgres DB](https://www.postgresql.org/)
- [smartcl](http://www.smartmontools.org)
