angular.module('shopApp')
  .controller('ProductCtrl', ['$scope', 'ProductService', 'AuthService', 'Toast', function($scope, ProductService, AuthService, Toast) {
    $scope.products = [];
    $scope.loading = true;
    $scope.isLoggedIn = AuthService.isLoggedIn;
    $scope.search = '';
    $scope.emojis = ['💻','🖱️','⌨️','🖥️','🔌','📱','🎮','🎧','📷','🔋'];

    $scope.getEmoji = function(index) { return $scope.emojis[index % $scope.emojis.length]; };

    function loadProducts() {
      $scope.loading = true;
      ProductService.getAll().then(function(res) {
        $scope.products = res.data.data;
      }).finally(function() { $scope.loading = false; });
    }

    loadProducts();

    $scope.filteredProducts = function() {
      if (!$scope.search) return $scope.products;
      var s = $scope.search.toLowerCase();
      return $scope.products.filter(function(p) {
        return p.name.toLowerCase().indexOf(s) !== -1 || p.description.toLowerCase().indexOf(s) !== -1;
      });
    };
  }])

  .controller('ProductDetailCtrl', ['$scope', '$routeParams', '$location', 'ProductService', 'AuthService', 'Toast', function($scope, $routeParams, $location, ProductService, AuthService, Toast) {
    $scope.product = null;
    $scope.loading = true;
    $scope.isLoggedIn = AuthService.isLoggedIn();

    ProductService.getOne($routeParams.id).then(function(res) {
      $scope.product = res.data.data;
    }).catch(function() {
      $location.path('/products');
    }).finally(function() { $scope.loading = false; });
  }]);
