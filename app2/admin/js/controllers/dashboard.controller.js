angular.module('adminApp')
  .controller('DashboardCtrl', ['$scope', 'ApiService', function($scope, ApiService) {
    $scope.stats = {};
    $scope.recentUsers = [];
    $scope.loading = true;

    ApiService.getStats().then(function(res) {
      $scope.stats = res.data.data;
    });

    ApiService.getUsers().then(function(res) {
      $scope.recentUsers = res.data.data.slice(-5).reverse();
    });

    ApiService.getProducts().then(function(res) {
      $scope.recentProducts = res.data.data.slice(-5).reverse();
    }).finally(function() { $scope.loading = false; });
  }]);
