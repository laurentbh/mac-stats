CREATE DATABASE macstats

CREATE TABLE battery (
    host    TEXT,
    stamp   TIMESTAMPTZ,
    metrics JSON
);

CREATE TABLE ssd (
    host    TEXT,
    stamp   TIMESTAMPTZ,
    metrics JSON
);