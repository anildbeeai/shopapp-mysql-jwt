angular.module('shopApp')
  .controller('HomeCtrl', ['$scope', 'AuthService', 'ProductService', function($scope, AuthService, ProductService) {
    $scope.isLoggedIn = AuthService.isLoggedIn();
    $scope.featuredProducts = [];

    ProductService.getAll().then(function(res) {
      $scope.featuredProducts = res.data.data.slice(0, 3);
    });
  }]);
