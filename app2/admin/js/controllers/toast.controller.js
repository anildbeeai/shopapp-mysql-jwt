angular.module('adminApp')
  .controller('ToastCtrl', ['$scope', '$rootScope', function($scope, $rootScope) {
    $scope.toast = { show: false };
    $rootScope.$on('toast', function(event, data) {
      $scope.toast = { show: true, message: data.message, type: data.type || 'info' };
      setTimeout(function() { $scope.$apply(function() { $scope.toast.show = false; }); }, 3500);
    });
  }])
  .factory('Toast', ['$rootScope', function($rootScope) {
    return {
      show: function(msg, type) { $rootScope.$emit('toast', { message: msg, type: type || 'info' }); },
      success: function(msg) { this.show(msg, 'success'); },
      error:   function(msg) { this.show(msg, 'error'); },
      info:    function(msg) { this.show(msg, 'info'); }
    };
  }]);
