angular.module('myApp').factory('AuthService',
  ['$q', '$timeout', '$http',
  function ($q, $timeout, $http) {

    // create user variable
    var user = null;

    var user_obj;


    // return available functions for use in controllers
    return ({
      isLoggedIn: isLoggedIn,
      getUserStatus: getUserStatus,
      login: login,
      logout: logout,
      register: register,
      getUser: getUser,
      createApp: createApp
    });

    function isLoggedIn() {
        if(user) {
          return true;
        } else {
          return false;
        }
    }

    function getUserStatus() {
      return user;
    }

    function login(username, password) {

      // create a new instance of deferred
      var deferred = $q.defer();

      // send a post request to the server
      $http.post('/user/login', {username: username, password: password})
        // handle success
        .success(function (data, status) {
          if(status === 200 && data.status){
            user = true;
            console.log("user status : "+ isLoggedIn());
            deferred.resolve();
          } else {
            user = false;
            deferred.reject();
          }
        })
        // handle error
        .error(function (data) {
          user = false;
          deferred.reject();
        });
      // return promise object
      return deferred.promise;
    }

    function logout() {

      // create a new instance of deferred
      var deferred = $q.defer();

      // send a get request to the server
      $http.get('/user/logout')
        // handle success
        .success(function (data) {
          user = false;
          deferred.resolve();
        })
        // handle error
        .error(function (data) {
          user = false;
          deferred.reject();
        });

      // return promise object
      return deferred.promise;

    }

    function register(username, password) {

      // create a new instance of deferred
      var deferred = $q.defer();

      // send a post request to the server
      $http.post('/user/register', {username: username, password: password})
        // handle success
        .success(function (data, status) {
          if(status === 200 && data.status){
            deferred.resolve();
          } else {
            deferred.reject();
          }
        })
        // handle error
        .error(function (data) {
          deferred.reject();
        });

      // return promise object
      return deferred.promise;
    }

    function getUser() {
      var deferred = $q.defer();

      $http.get('/user/profile')
      .success(function (data) {
        console.log("user from backend (service) :"+data.user)
        user_obj = data.user; 
        deferred.resolve(data);
      })
      .error(function (data) {
        deferred.reject("error");
      });

      return deferred.promise;
    }

    function createApp(user_app) {
      console.log(user_obj);
      var deferred = $q.defer();
      $http.post('/user/create/app/'+user_obj, user_app)
      .success(function (data) {
        console.log("response from create user :"+data)
        // post to pull user's docker image to node  

        console.log(user_app);
        //create rc object
        // send essential infor for create replicaion controller
        var rc_obj = {
          name: user_app.image_hub,
          replicas: 1,
          namespace: user_obj,
          image: user_app.image_hub,
          port: user_app.internal_port
        };

        console.log(rc_obj);

        $http.post('http://localhost:8182/rc/create/', rc_obj)
        .success(function() {
          console.log("success to created RC")
          deferred.resolve(data);
        })
        .error(function(data) {
          deferred.reject("error to create replication.");  
        });                      
      })
      .error(function (data) {
        deferred.reject("error to create store application to DB.");
      });

      return deferred.promise;
  
    }
}]);