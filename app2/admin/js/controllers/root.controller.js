angular.module('adminApp')
  .controller('RootCtrl', ['$scope', '$location', '$rootScope', 'AuthService', function($scope, $location, $rootScope, AuthService) {

    $scope.isLoggedIn = function() { return AuthService.isLoggedIn(); };
    $scope.getUserName = function() { var u = AuthService.getUser(); return u ? u.name : ''; };
    $scope.getInitial  = function() { var u = AuthService.getUser(); return u ? u.name[0].toUpperCase() : 'A'; };

    $scope.isActive = function(route) { return $location.path().indexOf(route) !== -1; };

    $rootScope.$on('$routeChangeSuccess', function() {
      if (!AuthService.isLoggedIn() && $location.path() !== '/login') {
        $location.path('/login');
      }
      var titles = { '/dashboard': '📊 Dashboard', '/products': '📦 Product Management', '/users': '👥 User Management' };
      $scope.pageTitle = titles[$location.path()] || 'Admin Panel';
    });

    $scope.logout = function() {
      AuthService.logout();
      $location.path('/login');
    };
  }]);
