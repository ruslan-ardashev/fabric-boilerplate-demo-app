app.service('TransferService', ["$http", "$q", function($http, $q) {
  
    var transferLogic = function(transferDetails){
        var deferred = $q.defer();
        
        // req.sourceMTO, req.sourceAccountNumber, req.amount, req.destinationMTO, req.destinationAccountNumber

        $http({
            method: 'POST',
            url: '/api/v1/transfer',
            data: { sourceMTO: transferDetails.sourceMTO, sourceAccountNumber: transferDetails.sourceAccountNumber, amount: transferDetails.amount, destinationMTO: transferDetails.destinationMTO, destinationAccountNumber: transferDetails.destinationAccountNumber }
        }).then(function success(response) {
            deferred.resolve(response.data);

        }, function error(error) {
            deferred.reject(error);
        });

        return deferred.promise;

    }

    var afterTransferLogic = function(){
        var deferred = $q.defer();

        $http({
            method: 'GET',
            url: '/api/v1/transfer/after',
            data: { }
        }).then(function success(response) {
            deferred.resolve(response.data);

        }, function error(error) {
            deferred.reject(error);
        });

        return deferred.promise;

    }

  return {
    transfer: transferLogic,
    afterTransfer: afterTransferLogic
  }

}]);
