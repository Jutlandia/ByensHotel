require "dotenv"
require "kemal"
require "kemal-session"
require "kemal-csrf"
require "./**"

if Kemal.config.env != "production"
  Dotenv.load
end

# TODO: configure session correctly
Kemal::Session.config do |config|
  config.secret = ENV["SESSION_SECRET"]
  config.secure = Kemal.config.env == "production"
end

# TODO: configure csrf correctly
add_handler CSRF.new

get "/" do
  page_title = "Home"
  render "src/views/index.ecr", "src/views/layouts/layout.ecr"
end

get "/login" do |env|
  page_title = "Log in"
  form = LoginForm.new
  render "src/views/auth/login.ecr", "src/views/layouts/layout.ecr"
end

post "/login" do |env|
  page_title = "Log in"
  form = LoginForm.new env
  if form.valid?
    # TODO: verify that the credentials are correct
    env.redirect "/"
    next
  end
  render "src/views/auth/login.ecr", "src/views/layouts/layout.ecr"
end

get "/register" do |env|
  page_title = "Register"
  form = RegisterForm.new
  render "src/views/auth/register.ecr", "src/views/layouts/layout.ecr"
end

post "/register" do |env|
  page_title = "Register"
  form = RegisterForm.new env
  if form.valid?
    # TODO: do stuff
    env.redirect "/login"
    next
  end
  render "src/views/auth/register.ecr", "src/views/layouts/layout.ecr"
end

Kemal.run
