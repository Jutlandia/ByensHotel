require "./spec_helper"
require "../src/forms/**"

describe "Forms" do
  describe "LoginForm" do
    form = LoginForm.new

    before_each do
      form.username.value = "alice"
      form.password.value = "pw123"
    end

    it "is valid if username and password is non-empty" do
      form.valid?.should be_true
    end

    {% begin %}
      {% fields = %w(username password) %}
      {% for field in fields %}
        it "is invalid if {{field.id}} is empty" do
          form.{{field.id}}.value = ""

          form.valid?.should be_false
          form.{{field.id}}.errors.includes?("This field is required").should be_true
        end
      {% end %}
    {% end %}
  end

  describe "RegisterForm" do
    form = RegisterForm.new

    before_each do
      form.username.value = "alice"
      form.email.value = "alice@mail.com"
      form.password.value = "pw123"
      form.confirm_password.value = "pw123"
    end

    it "is valid if all fields are non-empty and passwords are equal" do
      form.valid?.should be_true
    end

    {% begin %}
      {% fields = %w(username email password confirm_password) %}
      {% for field in fields %}
        it "is invalid if {{field.id}} is empty" do
          form.{{field.id}}.value = ""

          form.valid?.should be_false
        end
      {% end %}
    {% end %}

    it "is invalid if passwords are not equal" do
      form.password.value = "123456"
      form.confirm_password.value = "1234567"

      form.valid?.should be_false
      form.password.errors.includes?("Please confirm your password").should be_true
    end
  end
end
