angular.module('adminApp')
  .controller('LoginCtrl', ['$scope', '$location', 'AuthService', function($scope, $location, AuthService) {

    if (AuthService.isLoggedIn()) { $location.path('/dashboard'); return; }

    $scope.form = {};
    $scope.loading = false;
    $scope.error = '';

    $scope.login = function() {
      if (!$scope.form.email || !$scope.form.password) {
        $scope.error = 'Please enter both email and password.';
        return;
      }
      $scope.loading = true;
      $scope.error = '';
      AuthService.login($scope.form)
        .then(function(res) {
          var user = res.data.data.user;
          if (user.role !== 'admin') {
            $scope.error = 'Access denied. Admin privileges required.';
            $scope.loading = false;
            return;
          }
          AuthService.saveToken(res.data.data.token);
          AuthService.saveUser(user);
          $location.path('/dashboard');
        })
        .catch(function(err) {
          $scope.error = err.data ? err.data.message : 'Login failed.';
          $scope.loading = false;
        });
    };
  }]);
