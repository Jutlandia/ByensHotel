require "./spec_helper"
require "../src/jutlandia"

describe "views" do
  describe "Home view" do
    it "renders /" do
      get "/"

      response.status_code.should eq 200
      response.headers["Content-Type"].should eq "text/html"
      response.body.should contain %(<h1 class="title has-text-centered">Jutlandia</h1>)
    end
  end

  describe "Log in view" do
    it "renders /login" do
      get "/login"

      response.status_code.should eq 200
      response.headers["Content-Type"].should eq "text/html"
      response.body.should contain %(<h1 class="title has-text-centered my-6">Log in</h1>)
      response.body.should contain %(<input type="text" id="username" name="username" class="input" value="" required/>)
      response.body.should contain %(<input type="password" id="password" name="password" class="input" value="" required/>)
    end
  end

  describe "Register view" do
    it "renders /register" do
      get "/register"

      response.status_code.should eq 200
      response.headers["Content-Type"].should eq "text/html"
      response.body.should contain %(<h1 class="title has-text-centered my-6">Register</h1>)
      response.body.should contain %(<input type="text" id="username" name="username" class="input" value="" required/>)
      response.body.should contain %(<input type="email" id="email" name="email" class="input" value="" required/>)
      response.body.should contain %(<input type="password" id="password" name="password" class="input" value="" required/>)
      response.body.should contain %(<input type="password" id="confirm_password" name="confirm_password" class="input" value="" required/>)
    end
  end
end
