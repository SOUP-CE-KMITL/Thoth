var express = require('express'),
    router = express.Router(),
    passport = require('passport'),
    util = require("util");
    User = require('../models/user.js');


router.post('/register', function(req, res) {
  User.register(new User({ username: req.body.username }), req.body.password, function(err, account) {
    if (err) {
      return res.status(500).json({err: err});
    } 
    passport.authenticate('local')(req, res, function () {
      return res.status(200).json({status: 'Registration successful!'});
    });
  });
});

router.post('/login', function(req, res, next) {
  passport.authenticate('local', function(err, user, info) {
    if (err) {
      return res.status(500).json({err: err});
    }
    if (!user) {
      return res.status(401).json({err: info});
    }
    req.logIn(user, function(err) {
      if (err) {
        return res.status(500).json({err: 'Could not log in user'});
      }
      res.status(200).json({status: 'Login successful!'});
    });
  })(req, res, next);
});

router.get('/logout', function(req, res) {
  req.logout();
  res.status(200).json({status: 'Bye!'});
});

router.get('/profile', function(req, res){
  if(req.user == undefined) {
    return "none";
  }
  console.log("backend user : "+req.user.username);
  res.status(200).json({user: req.user.username});
});
 
router.get('/get/apps/:username', function(req, res) {
    var username = req.params.username;
    console.log("username : " + username);
    User.findOne({username: username}, function(err, userDoc){
      console.log(userDoc);
      res.status(200).json(userDoc);
    });
});

router.post('/create/app/:username', function(req, res) {
    var username = req.params.username;
    console.log("username at : "+username);
    // save information to database
    User.findOne({username: username}, function(err, userDoc){
      console.log(userDoc.username);
      if(err){
        console.log('err to find user');
        res.json({err: 'err to find user'});
      }
      if(!userDoc){
        console.log('user is not returned');
        res.json({err: 'user is not returned'});
      }else{
        console.log("userDoc existed");
        console.log(userDoc.app);
        // find maximum external port
        /*
        this.findOne({username: username})
            .sort('-external_port')
            .exec(function(err, user) {
              var max_port = user.external_port;
              external_port = max_port + 1;
            });
        */

        // calculate vamp port from external_port
        var user_app = {
          app_name: req.body.app_name,
          image_name: req.body.image_name,
          runtime_env: req.body.runtime_env,
          internal_port: req.body.internal_port,
          external_port: req.body.external_port,
          vamp_port: req.body.external_port,
          max_instance: req.body.max_instance,
          min_instance: req.body.min_instance
        }
        // check app existed or not.
        if (typeof userDoc.app == 'undefined'){
          console.log("initial array inside userDoc");
          userDoc.app = [];
        }
        console.log("push app object");
        userDoc.app.push(user_app);
        userDoc.save(function (err) {
          console.log("come inside save")
          if (err){
            console.log("cannot save application to user"+err)
            res.json({err: 'cannot save application to user'});
          }
          res.status(200).json({status: "work"});
        });
      }
    });


});


module.exports = router;
