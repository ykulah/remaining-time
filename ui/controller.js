angular.module('remainingTimeApp', ['ngMaterial'])
  .controller('remainingTimeController', function($scope, $http, $mdDialog) {
    var backendURL = "https://backend-dot-remaining-time-c9dd7.appspot.com"
    $scope.showRegister = true
    $scope.getInfo = function(){
      $scope.showRegister = false
      $scope.showNewTripForm = true
      var user = $scope.username
      $http.get(backendURL + "/api/getTrips/"+user)
      .then(function(response){
        $scope.trips = [];
        angular.forEach(response.data, function(t) {
          $scope.trips.push(t)
        });
      });

      $http.get(backendURL + "/api/getRemainingDays/"+user)
      .then(function(response){
        $scope.remainingDays = response.data.requiredDays
        $scope.exactDate = addDays(new Date(), $scope.remainingDays)
      });
    };

    $scope.addTrip =function(){
      var st = new Date($scope.startDate)
      var start = ""
      if (st.getDate() < 10){
        start += "0"
      }
      start += st.getDate() 
      start += "-" 
      if (st.getDay() < 10){
        start += "0"
      }
      var tmp = st.getMonth() + 1
      start = start + tmp + "-" + st.getFullYear()


      var e = new Date($scope.endDate)
      end = ""
      if (e.getDate() < 10){
        end += "0"
      }
      end = end + e.getDate() + "-"
      if (e.getMonth() < 10){
        end += "0"
      }
      var tmp = e.getMonth() + 1
      end += tmp + "-" + e.getFullYear()

      $http.get(backendURL + "/api/addTrip/"+$scope.username+"/"+start+"/"+end)
      .then(function(response){
        $scope.showStatus = true
        $scope.status = "Trip added"
        $http.get(backendURL + "/api/getTrips/"+$scope.username)
        .then(function(response){
          $scope.trips = [];
          angular.forEach(response.data, function(t) {
            $scope.trips.push(t)
          });
        });
        $http.get(backendURL + "/api/getRemainingDays/"+$scope.username)
        .then(function(response){
          $scope.remainingDays = response.data.requiredDays
        });
      }); 
    };

    $scope.addUser = function(){
      var st = new Date($scope.jobStartDate)
  
      var start = ""
      if (st.getDate() < 10){
        start += "0"
      }
      start += st.getDate() 
      start += "-"
      if (st.getMonth() < 10){
        start += "0"
      }
      var tmp = st.getMonth() + 1
      start += tmp + "-" + st.getFullYear()

      $http.get(backendURL + "/api/addUser/"+$scope.username+"/"+start)
      .then(function(response){
        $scope.showRegisterResult = true
        $scope.registerStatus = "User registered."
      });

    };

    $scope.deleteTrip = function(index){
      tripDetails = $scope.trips[index]
   
      var confirm = $mdDialog.confirm()
            .title('Would you like to delete your trip?')
            .textContent("Departure:"+ tripDetails.start + "\n Arrival:" + tripDetails.end)
            .ariaLabel('Lucky day')
            .ok('Please do it!')
            .cancel('Back to list');
  
      $mdDialog.show(confirm).then(function() {
        console.log("TODO: DELETE TRIP")
      }, function() {
        
      });

    };

    function addDays(date, days) {
      var result = new Date(date);
      result.setDate(result.getDate() + days);

      day = result.getDate()
      month = result.getMonth() + 1
      year = result.getFullYear()
      return day+"-"+month+"-"+year;
    }
  });