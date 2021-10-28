require "kemal-form"

class RegisterForm < Kemal::Form
  field username : Kemal::Form::TextField,
    attrs: {"class" => "input"},
    label: Kemal::Form::Label.new(
      "username",
      "Username",
      {"class" => "label"}),
    validators: [Kemal::FormValidator::Required.new]
  field email : Kemal::Form::EmailField,
    attrs: {"class" => "input"},
    label: Kemal::Form::Label.new(
      "email",
      "Email",
      {"class" => "label"}),
    validators: [
      Kemal::FormValidator::Required.new,
      Kemal::FormValidator::Email.new,
    ]
  field password : Kemal::Form::PasswordField,
    attrs: {"class" => "input"},
    label: Kemal::Form::Label.new(
      "password",
      "Password",
      {"class" => "label"}),
    validators: [Kemal::FormValidator::Required.new]
  field confirm_password : Kemal::Form::PasswordField,
    attrs: {"class" => "input"},
    label: Kemal::Form::Label.new(
      "confirm_password",
      "Confirm Password",
      {"class" => "label"}),
    validators: [Kemal::FormValidator::Required.new]

  def valid?
    if @password.value != @confirm_password.value
      @password.add_error "Please confirm your password"
      @password.value = ""
      @confirm_password.value = ""
      return false
    end
    super
  end
end
