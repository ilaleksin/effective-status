-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE public.tag(id serial primary key, name text unique, type text);
CREATE TABLE public.outage(id serial primary key, ticket text unique, summary text, details text, start timestamp, outage_end timestamp, duration interval);
CREATE TABLE public.system(id serial primary key, name text unique, status integer);
CREATE TABLE public.health_check(id serial primary key, name text unique, monitor_name text, status integer, priority integer, monitor_system integer);
CREATE TABLE public.monitoring_system(id serial primary key, fqdn text, system_type text);
CREATE TABLE public.dependency(id serial primary key, parentID integer, childID integer);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE public.tag;
DROP TABLE public.outage;
DROP TABLE public.system;
DROP TABLE public.health_check;
DROP TABLE public.monitoring_system;
DROP TABLE public.dependency;
-- +goose StatementEnd
