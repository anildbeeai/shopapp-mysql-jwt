angular.module('shopApp', ['ngRoute'])
  .config(['$routeProvider', '$httpProvider', function($routeProvider, $httpProvider) {

    $routeProvider
      .when('/home', {
        templateUrl: 'views/home.html',
        controller: 'HomeCtrl'
      })
      .when('/products', {
        templateUrl: 'views/products.html',
        controller: 'ProductCtrl'
      })
      .when('/products/:id', {
        templateUrl: 'views/product-detail.html',
        controller: 'ProductDetailCtrl'
      })
      .when('/login', {
        templateUrl: 'views/login.html',
        controller: 'AuthCtrl'
      })
      .when('/register', {
        templateUrl: 'views/register.html',
        controller: 'AuthCtrl'
      })
      .when('/profile', {
        templateUrl: 'views/profile.html',
        controller: 'ProfileCtrl'
      })
      .otherwise({ redirectTo: '/home' });

    // JWT Interceptor
    $httpProvider.interceptors.push('AuthInterceptor');
  }])

  // JWT Interceptor Factory
  .factory('AuthInterceptor', ['$q', '$injector', function($q, $injector) {
    return {
      request: function(config) {
        var AuthService = $injector.get('AuthService');
        var token = AuthService.getToken();
        if (token) {
          config.headers['Authorization'] = 'Bearer ' + token;
        }
        return config;
      },
      responseError: function(rejection) {
        if (rejection.status === 401) {
          var AuthService = $injector.get('AuthService');
          AuthService.logout();
          var $location = $injector.get('$location');
          $location.path('/login');
        }
        return $q.reject(rejection);
      }
    };
  }])

  // Constants
  .constant('API_URL', 'http://localhost:8080/api');
