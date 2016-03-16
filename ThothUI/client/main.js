var myApp = angular.module('myApp', ['ngRoute', 'nvd3']);

myApp.config(function ($routeProvider) {
  $routeProvider
    .when('/', { templateUrl: 'partials/home.html', access: {restricted: true}})
    .when('/login', {
      templateUrl: 'partials/login.html',
      controller: 'loginController',
      access: {restricted: false}
    })
    .when('/logout', {
      controller: 'logoutController',
      access: {restricted: true}
    })
    .when('/register', {
      templateUrl: 'partials/register.html',
      controller: 'registerController',
      access: {restricted: false}
    })
    .when('/deploy', {
      templateUrl: 'partials/deploy.html',
      controller: 'deployController',
      access: {restricted: true}
    })
    .when('/configure', {
      templateUrl: 'partials/configure.html',
      controller: 'configureController',
      access: {restricted: true}
    })
    .otherwise({redirectTo: '/'});
});


myApp.run(function ($rootScope, $location, $route, AuthService) {
  $rootScope.$on('$routeChangeStart', function (event, next, current) {
    if (next.access.restricted && AuthService.isLoggedIn() === false) {
      $location.path('/login');
    }
    console.log("restricted = "+next.access.restricted);
    console.log("AuthService = "+AuthService.isLoggedIn());
  });
});
