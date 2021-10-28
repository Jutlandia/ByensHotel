# jutlandia

TODO: Write a description here

## Requirements

+ [Crystal](https://crystal-lang.org/)
+ [npm](https://www.npmjs.com/)

## Setup

Copy `.env.sample` to `.env` and update accordingly.
```
cp .env.sample .env
```
To generate `SESSION_SECRET` run
```
crystal eval 'require "random/secure"; puts Random::Secure.hex(64)'
```

#### npm
+ Run `npm install` to install node modules.
+ Run `npm run css-build` to build required css styles.

#### Crystal
+ Run `shards install` to install required Crystal shards.

## Usage

Run `Crystal run src/jutlandia.cr` to start the server.

## Contributing

1. Fork it (<https://github.com/your-github-user/jutlandia/fork>)
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request

## Contributors

- [Henrik Christensen](https://github.com/your-github-user) - creator and maintainer
