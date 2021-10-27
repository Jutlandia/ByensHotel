require "dotenv"
require "kemal"
require "kemal-session"
require "kemal-csrf"

if Kemal.config.env == "development"
  Dotenv.load
end

Kemal::Session.config do |config|
  config.secret = ENV["SESSION_SECRET"]
  config.secure = Kemal.config.env == "production"
end

add_handler CSRF.new

get "/" do
  page_title = "Home"
  render "src/views/index.ecr", "src/views/layouts/layout.ecr"
end

Kemal.run
