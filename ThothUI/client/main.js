var myApp = angular.module('myApp', ['ngRoute', 'nvd3']);

myApp.config(function ($routeProvider) {
  $routeProvider
    .when('/', {templateUrl: 'partials/home.html'})
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
      access: {restricted: false}
    })
    .when('/configure', {
      templateUrl: 'partials/configure.html',
      controller: 'configureController',
      access: {restricted: false}
    })
    .otherwise({redirectTo: '/'});
});


myApp.run(function ($rootScope, $location, $route, AuthService) {
  $rootScope.$on('$routeChangeStart', function (event, next, current) {
    if (next.access.restricted && AuthService.isLoggedIn() === false) {
      $location.path('/login');
    }
  });
});
