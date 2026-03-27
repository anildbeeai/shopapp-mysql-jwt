angular.module('shopApp')
  .controller('ProfileCtrl', ['$scope', '$location', 'AuthService', 'Toast', function($scope, $location, AuthService, Toast) {

    if (!AuthService.isLoggedIn()) {
      $location.path('/login');
      return;
    }

    $scope.user = AuthService.getUser();

    $scope.logout = function() {
      AuthService.logout();
      Toast.info('You have been logged out.');
      $location.path('/home');
    };
  }]);
