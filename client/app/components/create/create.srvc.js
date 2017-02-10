app.service('CreateService', ["$http", "$q", function($http, $q) {
  return {
    create: function(accountDetails){
        var deferred = $q.defer();
        
        $http({
            method: 'POST',
            url: '/api/v1/create',
            data: { mto: accountDetails.mto, accountFirstName: accountDetails.firstName, accountLastName: accountDetails.lastName, accountNumber: accountDetails.accountNumber, balance: accountDetails.balance  }
        }).then(function success(response) {
            deferred.resolve(response.data);

        }, function error(error) {
            deferred.reject(error);
        });

        return deferred.promise;

        }
    }
  }
]);

