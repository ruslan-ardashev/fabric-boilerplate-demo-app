app.controller("AddRemoveController", ["things", "$localStorage", "$location", function (things, $localStorage, $location) {

    var vm = this;

    vm.things = things;
    
    vm.openThing = function(thingId){
        
        $localStorage.selectedThing = thingId;
        $location.path('/detail');
    }

    vm.accounts_create = function() {
    		console.log("selected accounts create!");
    		// $location.path("/");
    }

    vm.accounts_addremove = function() {
    		console.log("selected accounts addremove!");
    }

    vm.transfer = function() {
    		console.log("selected trasfer!");
    }
    
}]);
