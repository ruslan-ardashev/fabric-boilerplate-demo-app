app.controller("CreateController", ["$localStorage", "CreateService", "$location", function ($localStorage, CreateService, $location) {

    var vm = this;

    vm.invalidTransfer = false;
    vm.showTransferSpinner = false;

    vm.create = function(accountDetails) {

        console.log("~~user~~");
        console.log($localStorage.user.id);        

        // insert source MTO from logged in identity, selected destination MTO
        accountDetails.mto = $localStorage.user.id

        console.log("~~accountDetails~~");
        console.log(accountDetails);

        vm.showTransferSpinner = true;

        // chaincode
        //  arg0 - string mto.Name
        //  arg1 - string account.firstName
        //  arg2 - string account.lastName
        //  arg3 - string account.number
        //  arg4 - string account.balance

        // service
        // mto: accountDetails.mto, 
        // accountFirstName: accountDetails.firstName, 
        // accountLastName: accountDetails.lastName, 
        // accountNumber: accountDetails.accountNumber, 
        // balance: accountDetails.balance

        CreateService.create(accountDetails)
            .then(function (result) {
                console.log("Returning trasfer result: ")
                console.log(result);
                if (!result) {
                    console.log("no result but no error");
                    // vm.showTransferringSpinner = false;
                    // vm.invalidTransaction = true;
                } else if (result.successful) {
                    console.log("result, has result.successful");
                    // vm.showTransferringSpinner = false;

                    // set user and navigation information on rootscope
                    // $localStorage.user = result.user;

                    // store the token in localStorage
                    // $localStorage.token = result.token;
                    
                    // delete $localStorage.selectedThing;

                    // $location.path("/master");

                } else {
                    console.log("result, no result.successful, other result fields present though");
                    console.log(result);
                    // vm.showTransferringSpinner = false;
                    // vm.invalidTransaction = true;
                }
            }, function (error) {
                console.log("transfer error:");
                console.log(error);
                // vm.showTransferringSpinner = false;
                // vm.invalidTransaction = true;

            });
    }
    
}]);
