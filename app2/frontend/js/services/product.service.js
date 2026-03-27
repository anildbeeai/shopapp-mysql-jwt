angular.module('shopApp')
  .service('ProductService', ['$http', 'API_URL', function($http, API_URL) {

    this.getAll = function() {
      return $http.get(API_URL + '/products');
    };

    this.getOne = function(id) {
      return $http.get(API_URL + '/products/' + id);
    };

    this.create = function(data) {
      return $http.post(API_URL + '/products', data);
    };

    this.update = function(id, data) {
      return $http.put(API_URL + '/products/' + id, data);
    };

    this.delete = function(id) {
      return $http.delete(API_URL + '/products/' + id);
    };
  }]);
