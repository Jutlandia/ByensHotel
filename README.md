# Byens Hotel

TODO: Write a description here

## Requirements

+ [Go](https://golang.org/)
+ [Docker](https://www.docker.com/)
+ [npm](https://www.npmjs.com/)

## Setup

Copy `.env.sample` to `.env` and update accordingly.
```
cp .env.sample .env
```

ByensHotel will run in development mode if `HOTEL_ENV` is not set.

#### npm
+ Run `npm install` to install node modules.
+ Run `npm run css-build` to build required css styles.

#### LDAP test server

+ Run `docker pull rroemhild/test-openldap` to pull the ldap test server image.

## Usage


+ Run `docker run --rm -p 10389:10389 -p 10636:10636 rroemhild/test-openldap` to start the ldap test server.
+ Run `go run main.go` to start the server.

## Tests

Run `go test ./...` to run all specs.

*Note:* Make sure that the `rroemhild/test-openldap` container is running before you run `go test`.

## Contributing

1. Fork it (<https://github.com/Jutlandia/ByensHotel/fork>)
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request
