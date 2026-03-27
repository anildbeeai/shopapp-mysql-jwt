angular.module('adminApp')
  .service('AuthService', ['$http', '$window', 'API_URL', function($http, $window, API_URL) {

    this.login = function(data) { return $http.post(API_URL + '/auth/login', data); };

    this.saveToken = function(t) { $window.localStorage.setItem('admin_jwt', t); };
    this.saveUser  = function(u) { $window.localStorage.setItem('admin_user', JSON.stringify(u)); };
    this.getToken  = function()  { return $window.localStorage.getItem('admin_jwt'); };
    this.getUser   = function()  {
      var u = $window.localStorage.getItem('admin_user');
      return u ? JSON.parse(u) : null;
    };
    this.isLoggedIn = function() { return !!this.getToken() && this.getUser() && this.getUser().role === 'admin'; };
    this.logout = function() {
      $window.localStorage.removeItem('admin_jwt');
      $window.localStorage.removeItem('admin_user');
    };
  }]);
