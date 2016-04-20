// user model
var mongoose = require('mongoose'),
    Schema = mongoose.Schema,
    passportLocalMongoose = require('passport-local-mongoose');

var User = new Schema({
  username: String,
  password: String,
  app:[{
	app_name: { Type: String, index: { unique: true }}
	image_name: String,
	runtime_env: String,
	internal_port: Number,
	external_port: Number,
	vamp_port: Number,
	max_instance: Number,
	min_instance: Number
	}],
});

User.plugin(passportLocalMongoose);

module.exports = mongoose.model('users', User);
