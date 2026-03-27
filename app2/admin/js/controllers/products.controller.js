angular.module('adminApp')
  .controller('ProductsCtrl', ['$scope', 'ApiService', 'Toast', function($scope, ApiService, Toast) {
    $scope.products = [];
    $scope.loading = true;
    $scope.showModal = false;
    $scope.editMode = false;
    $scope.form = {};
    $scope.search = '';

    function load() {
      $scope.loading = true;
      ApiService.getProducts().then(function(res) {
        $scope.products = res.data.data;
      }).finally(function() { $scope.loading = false; });
    }

    load();

    $scope.filteredProducts = function() {
      if (!$scope.search) return $scope.products;
      var s = $scope.search.toLowerCase();
      return $scope.products.filter(function(p) {
        return p.name.toLowerCase().indexOf(s) !== -1;
      });
    };

    $scope.openCreate = function() {
      $scope.editMode = false;
      $scope.form = {};
      $scope.showModal = true;
    };

    $scope.openEdit = function(product) {
      $scope.editMode = true;
      $scope.form = angular.copy(product);
      $scope.showModal = true;
    };

    $scope.closeModal = function() {
      $scope.showModal = false;
      $scope.form = {};
    };

    $scope.save = function() {
      if (!$scope.form.name || !$scope.form.price) {
        Toast.error('Name and price are required.');
        return;
      }
      $scope.form.price = parseFloat($scope.form.price);
      $scope.form.stock = parseInt($scope.form.stock) || 0;
      $scope.saving = true;

      var promise = $scope.editMode
        ? ApiService.updateProduct($scope.form.id, $scope.form)
        : ApiService.createProduct($scope.form);

      promise.then(function() {
        Toast.success($scope.editMode ? 'Product updated!' : 'Product created!');
        $scope.closeModal();
        load();
      }).catch(function() {
        Toast.error('Operation failed.');
      }).finally(function() { $scope.saving = false; });
    };

    $scope.delete = function(product) {
      if (!confirm('Delete "' + product.name + '"? This cannot be undone.')) return;
      ApiService.deleteProduct(product.id).then(function() {
        Toast.success('Product deleted.');
        load();
      }).catch(function() { Toast.error('Failed to delete.'); });
    };
  }]);
