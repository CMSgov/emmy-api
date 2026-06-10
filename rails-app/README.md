# README

This README would normally document whatever steps are necessary to get the
application up and running.

Things you may want to cover:

* Ruby version

* System dependencies

* Configuration

* Database creation

* Database initialization

* Services (job queues, cache servers, search engines, etc.)

* Deployment instructions

* ...

# How to run the test suite

There is a Docker compose setup for this Rails application. It should be
possible to start the application with

```bash
docker compose up app_rails
```

and then the test suite can be run with

```bash
docker compose exec app_rails bundle exec bin/rails spec
```
