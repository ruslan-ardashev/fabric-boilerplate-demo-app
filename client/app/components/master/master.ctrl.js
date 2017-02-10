app.controller("MasterController", ["things", "$localStorage", "$location", function (things, $localStorage, $location) {

    var vm = this;

    vm.things = things;
    
    vm.openThing = function(thingId){
        
        $localStorage.selectedThing = thingId;
        $location.path('/detail');
    }

    vm.create = function() {
    		console.log("selected accounts create!");
    		$location.path("/create");
    }

    vm.accounts_addremove = function() {
    		console.log("selected accounts addremove!");
    }

    vm.transfer = function() {
    		console.log("selected trasfer!");
            $location.path("/transfer");
    }
    
}]);
