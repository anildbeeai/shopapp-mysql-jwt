angular.module('shopApp')
  .controller('NavCtrl', ['$scope', '$location', 'AuthService', function($scope, $location, AuthService) {

    $scope.isLoggedIn = function() { return AuthService.isLoggedIn(); };

    $scope.userName = function() {
      var user = AuthService.getUser();
      return user ? user.name : '';
    };

    $scope.isActive = function(route) {
      return $location.path().indexOf(route) !== -1;
    };

    $scope.logout = function() {
      AuthService.logout();
      $location.path('/home');
    };
  }]);
