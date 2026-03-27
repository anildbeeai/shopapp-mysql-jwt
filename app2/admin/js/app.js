angular.module('adminApp', ['ngRoute'])
  .config(['$routeProvider', '$httpProvider', function($routeProvider, $httpProvider) {

    $routeProvider
      .when('/login', {
        templateUrl: 'views/login.html',
        controller: 'LoginCtrl'
      })
      .when('/dashboard', {
        templateUrl: 'views/dashboard.html',
        controller: 'DashboardCtrl'
      })
      .when('/products', {
        templateUrl: 'views/products.html',
        controller: 'ProductsCtrl'
      })
      .when('/users', {
        templateUrl: 'views/users.html',
        controller: 'UsersCtrl'
      })
      .otherwise({ redirectTo: '/dashboard' });

    // JWT Interceptor
    $httpProvider.interceptors.push('AuthInterceptor');
  }])

  .factory('AuthInterceptor', ['$q', '$injector', function($q, $injector) {
    return {
      request: function(config) {
        var AuthService = $injector.get('AuthService');
        var token = AuthService.getToken();
        if (token) { config.headers['Authorization'] = 'Bearer ' + token; }
        return config;
      },
      responseError: function(rejection) {
        if (rejection.status === 401 || rejection.status === 403) {
          var AuthService = $injector.get('AuthService');
          AuthService.logout();
          var $location = $injector.get('$location');
          $location.path('/login');
        }
        return $q.reject(rejection);
      }
    };
  }])

  .constant('API_URL', 'http://localhost:8080/api');
