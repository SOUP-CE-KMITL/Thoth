angular.module('myApp').controller('loginController',
  ['$scope', '$location', 'AuthService',
  function ($scope, $location, AuthService) {

    console.log(AuthService.getUserStatus());

    $scope.login = function () {

      // initial values
      $scope.error = false;
      $scope.disabled = true;

      // call login from service
      AuthService.login($scope.loginForm.username, $scope.loginForm.password)
        // handle success
        .then(function () {
          $location.path('/');
          $scope.disabled = false;
          $scope.loginForm = {};
        })
        // handle error
        .catch(function () {
          $scope.error = true;
          $scope.errorMessage = "Invalid username and/or password";
          $scope.disabled = false;
          $scope.loginForm = {};
        });

    };

}]);

angular.module('myApp').controller('logoutController',
  ['$scope', '$location', 'AuthService',
  function ($scope, $location, AuthService) {

    $scope.logout = function () {

      console.log(AuthService.getUserStatus());

      // call logout from service
      AuthService.logout()
        .then(function () {
          $location.path('/login');
        });

    };

}]);

angular.module('myApp').controller('registerController',
  ['$scope', '$location', 'AuthService',
  function ($scope, $location, AuthService) {

    console.log(AuthService.getUserStatus());

    $scope.register = function () {

      // initial values
      $scope.error = false;
      $scope.disabled = true;

      // call register from service
      AuthService.register($scope.registerForm.username, $scope.registerForm.password)
        // handle success
        .then(function () {
          $location.path('/login');
          $scope.disabled = false;
          $scope.registerForm = {};
        })
        // handle error
        .catch(function () {
          $scope.error = true;
          $scope.errorMessage = "Something went wrong!";
          $scope.disabled = false;
          $scope.registerForm = {};
        });

    };

}]);

angular.module('myApp').controller('deployController', 
  ['$scope', '$location', 'AuthService',
  function ($scope, $location, AuthService) {
}]);

angular.module('myApp').controller('configureController', 
  ['$scope', '$location', 'AuthService',
  function ($scope, $location, AuthService) {
    console.log(AuthService.getUserStatus());
}]);

// chart configure
var chart_options = {
            chart: {
                type: 'lineChart',
                height: 250,
                margin : {
                    top: 20,
                    right: 20,
                    bottom: 40,
                    left: 55
                },
                x: function(d){ return d.x; },  
                y: function(d){ return d.y; },
                useInteractiveGuideline: true,
                dispatch: {
                    stateChange: function(e){ console.log("stateChange"); },
                    changeState: function(e){ console.log("changeState"); },
                    tooltipShow: function(e){ console.log("tooltipShow"); },
                    tooltipHide: function(e){ console.log("tooltipHide"); }
                },
                xAxis: {
                    axisLabel: 'Time (s)'
                },
                yAxis: {
                    axisLabel: 'Resource (%)',
                    tickFormat: function(d){
                        return d3.format('.02f')(d);
                    },
                    axisLabelDistance: -10
                },
                callback: function(chart){
                    console.log("!!! lineChart callback !!!");
                }
            },
            title: {
                enable: true,
                text: 'Application Resource Usage'
            },
            subtitle: {
                enable: true,
                text: 'realtime application resource usage',
                css: {
                    'text-align': 'center',
                    'margin': '10px 13px 0px 7px'
                }
            },
            caption: {
                enable: true,
                html: '<b>ResourceUsageGraph</b> : This resource usage graph show resource metrics',
                css: {
                    'text-align': 'justify',
                    'margin': '10px 13px 0px 7px'
                }
            }
        }; 

// others
// resource usage graph
angular.module('myApp').controller('AppResourceUsageController',
  ['$scope', '$http', '$q', 'AuthService',
    function ($scope, $http, $q, AuthService) {
      // get user details
      AuthService.getUser().then(function( user ) {
        $scope.user = user;
        console.log(" user : " + user.user);
        // array of Application
        var apps = [];
        // http get application lists.
        $http.get("http://localhost:8182/apps/"+user.user)
        .success(function(response) {

          console.log("length : " + response.items.length);
          for(var i = 0; i < response.items.length; i++){
            // get application profile
            apps[i] = {};
            apps[i].name = response.items[i].metadata.name;
            apps[i].namespace = response.items[i].metadata.namespace;
            apps[i].internal_port = response.items[i].spec.template.spec.containers[0].ports[0].containerPort;
            apps[i].replicas = response.items[i].spec.replicas
            // initail array
            apps[i].data = [];
            apps[i].data[0] = {};
            apps[i].data[0].values = [];
            apps[i].data[1] = {};
            apps[i].data[1].values = [];

            apps[i].data[0].key = 'cpu';
            apps[i].data[0].color = '#ff7f0e';
            apps[i].data[1].key = 'memory';
            apps[i].data[1].color = '#2ca02c';
            console.log("created "+i)
          }

          $scope.options = chart_options;
          //$scope.data = [{values: [], key: 'cpu', color: '#ff7f0e'},{values: [], key: 'memory', color: '#2ca02c'}];
          $scope.apps = apps;
          // pause/play btn
          $scope.run = true;
          var app_datas = [];

          var t = 0;

          setInterval(function(){
            if (!$scope.run) return;
            for(var c = 0; c < apps.length; c++){
              console.log(apps[c].name);
              app_datas[c] = $http.get("http://localhost:8182/app/"+apps[c].name+"/metrics/"+user.user)
            }
            // request resource usage from api
            $q.all(app_datas).then(function(response) {
              console.log(response);
              for(var r = 0; r < response.length; r++){
                console.log(response[r].data.cpu); 
                $scope.apps[r].data[0].values.push({ x:t, y:response[r].data.cpu })
                var percent_mem = response[r].data.memory[0].mem_usage_in_bytes/1000000000 * 100;
                $scope.apps[r].data[1].values.push({ x:t, y:percent_mem })
                console.log("cpu : "+$scope.apps[r].data[1].values.slice(-1)[0].x);
                if($scope.apps[r].data[0].values.length > 20) $scope.apps[r].data[0].values.shift();
                if($scope.apps[r].data[1].values.length > 20) $scope.apps[r].data[1].values.shift();
              }
                t++;
            });
          }, 1000);
        })
      })    
    }
]);

// navbar active control
angular.module('myApp').controller('HeaderController', 
  ['$scope', '$location', 'AuthService',
  function HeaderController($scope, $location, AuthService){
    // check navbar active
    console.log($location.path());
    $scope.isActive = function (viewLocation) {
      return $location.path().indexOf(viewLocation) == 0;
    };
    // get user details
    AuthService.getUser().then(function( user ) {
      $scope.user = user;
    });
    $scope.user_status = AuthService.getUserStatus();
  }
]);