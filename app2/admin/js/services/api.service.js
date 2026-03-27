angular.module('adminApp')
  .service('ApiService', ['$http', 'API_URL', function($http, API_URL) {

    // Stats
    this.getStats = function() { return $http.get(API_URL + '/admin/stats'); };

    // Users
    this.getUsers    = function()   { return $http.get(API_URL + '/admin/users'); };
    this.deleteUser  = function(id) { return $http.delete(API_URL + '/admin/users/' + id); };

    // Products
    this.getProducts    = function()     { return $http.get(API_URL + '/products'); };
    this.createProduct  = function(data) { return $http.post(API_URL + '/admin/products', data); };
    this.updateProduct  = function(id, data) { return $http.put(API_URL + '/admin/products/' + id, data); };
    this.deleteProduct  = function(id)   { return $http.delete(API_URL + '/admin/products/' + id); };
  }]);
