require "kemal-form"

class LoginForm < Kemal::Form
  field username : Kemal::Form::TextField,
    attrs: {"class" => "input"},
    label: Kemal::Form::Label.new(
      "username",
      "Username",
      {"class" => "label"}),
    validators: [Kemal::FormValidator::Required.new]
  field password : Kemal::Form::PasswordField,
    attrs: {"class" => "input"},
    label: Kemal::Form::Label.new(
      "password",
      "Password",
      {"class" => "label"}),
    validators: [Kemal::FormValidator::Required.new]
end
