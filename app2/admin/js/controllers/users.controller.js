angular.module('adminApp')
  .controller('UsersCtrl', ['$scope', 'ApiService', 'Toast', function($scope, ApiService, Toast) {
    $scope.users = [];
    $scope.loading = true;
    $scope.search = '';

    function load() {
      $scope.loading = true;
      ApiService.getUsers().then(function(res) {
        $scope.users = res.data.data;
      }).finally(function() { $scope.loading = false; });
    }

    load();

    $scope.filteredUsers = function() {
      if (!$scope.search) return $scope.users;
      var s = $scope.search.toLowerCase();
      return $scope.users.filter(function(u) {
        return u.name.toLowerCase().indexOf(s) !== -1 || u.email.toLowerCase().indexOf(s) !== -1;
      });
    };

    $scope.delete = function(user) {
      if (user.role === 'admin') {
        Toast.error('Cannot delete admin user.');
        return;
      }
      if (!confirm('Delete user "' + user.name + '"?')) return;
      ApiService.deleteUser(user.id).then(function() {
        Toast.success('User deleted.');
        load();
      }).catch(function(err) {
        Toast.error(err.data ? err.data.message : 'Failed to delete.');
      });
    };
  }]);
