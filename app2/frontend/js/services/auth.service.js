angular.module('shopApp')
  .service('AuthService', ['$http', '$window', 'API_URL', function($http, $window, API_URL) {

    this.register = function(data) {
      return $http.post(API_URL + '/auth/register', data);
    };

    this.login = function(data) {
      return $http.post(API_URL + '/auth/login', data);
    };

    this.getProfile = function() {
      return $http.get(API_URL + '/auth/profile');
    };

    this.saveToken = function(token) {
      $window.localStorage.setItem('jwt_token', token);
    };

    this.saveUser = function(user) {
      $window.localStorage.setItem('current_user', JSON.stringify(user));
    };

    this.getToken = function() {
      return $window.localStorage.getItem('jwt_token');
    };

    this.getUser = function() {
      var user = $window.localStorage.getItem('current_user');
      return user ? JSON.parse(user) : null;
    };

    this.isLoggedIn = function() {
      return !!this.getToken();
    };

    this.isAdmin = function() {
      var user = this.getUser();
      return user && user.role === 'admin';
    };

    this.logout = function() {
      $window.localStorage.removeItem('jwt_token');
      $window.localStorage.removeItem('current_user');
    };
  }]);
