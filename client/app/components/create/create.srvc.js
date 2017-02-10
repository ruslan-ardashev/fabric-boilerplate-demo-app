app.service('CreateService', ["$http", "$q", function($http, $q) {
    
    var createLogic = function(createDetails){
        var deferred = $q.defer();
        
        $http({
            method: 'POST',
            url: '/api/v1/create',
            data: { mto: createDetails.mto, firstName: createDetails.firstName, lastName: createDetails.lastName, accountNumber: createDetails.accountNumber, balance: createDetails.balance  }
        }).then(function success(response) {
            deferred.resolve(response.data);

        }, function error(error) {
            deferred.reject(error);
        });

        return deferred.promise;

    }

    var afterCreateLogic = function(){
        var deferred = $q.defer();

        $http({
            method: 'GET',
            url: '/api/v1/create/after',
            data: { }
        }).then(function success(response) {
            deferred.resolve(response.data);

        }, function error(error) {
            deferred.reject(error);
        });

        return deferred.promise;

    }

    return {
        create: createLogic,
        afterCreate: afterCreateLogic
    }

  }
]);

