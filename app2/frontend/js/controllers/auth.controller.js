angular.module('shopApp')
  .controller('AuthCtrl', ['$scope', '$location', 'AuthService', 'Toast', function($scope, $location, AuthService, Toast) {

    if (AuthService.isLoggedIn()) {
      $location.path('/profile');
      return;
    }

    $scope.form = {};
    $scope.loading = false;
    $scope.error = '';

    $scope.login = function() {
      if (!$scope.form.email || !$scope.form.password) {
        $scope.error = 'Please fill in all fields.';
        return;
      }
      $scope.loading = true;
      $scope.error = '';
      AuthService.login($scope.form)
        .then(function(res) {
          AuthService.saveToken(res.data.data.token);
          AuthService.saveUser(res.data.data.user);
          Toast.success('Welcome back, ' + res.data.data.user.name + '!');
          $location.path('/profile');
        })
        .catch(function(err) {
          $scope.error = err.data ? err.data.message : 'Login failed. Please try again.';
        })
        .finally(function() { $scope.loading = false; });
    };

    $scope.register = function() {
      if (!$scope.form.name || !$scope.form.email || !$scope.form.password) {
        $scope.error = 'Please fill in all fields.';
        return;
      }
      if ($scope.form.password !== $scope.form.confirmPassword) {
        $scope.error = 'Passwords do not match.';
        return;
      }
      $scope.loading = true;
      $scope.error = '';
      AuthService.register($scope.form)
        .then(function(res) {
          AuthService.saveToken(res.data.data.token);
          AuthService.saveUser(res.data.data.user);
          Toast.success('Account created successfully!');
          $location.path('/profile');
        })
        .catch(function(err) {
          $scope.error = err.data ? err.data.message : 'Registration failed.';
        })
        .finally(function() { $scope.loading = false; });
    };
  }]);
