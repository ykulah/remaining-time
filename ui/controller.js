angular.module('remainingTimeApp', ['ngMaterial'])
  .controller('remainingTimeController', function($scope, $http, $mdDialog) {
    var backendURL = "https://backend-dot-remaining-time-c9dd7.appspot.com"
    $scope.showRegister = true
    $scope.showViewInfo = true

    $scope.getInfo = function(){
      $scope.showProgressBar=true
      $scope.showViewInfo = false
      $http.get(backendURL + "/api/count/"+$scope.username)
      .then(function(response){
        var cnt = parseInt(response.data)
        if (cnt <= 0){
          $scope.existStatus = "User not found. Register first."
          $scope.userExistsError = true
          $scope.showProgressBar=false
          $scope.showViewInfo = true
        } else {
          $scope.showRegister = false
          $scope.showNewTripForm = true
          $scope.userExistsError = false
          updateTripList()
          calculateRemainingTime()
        }
      }); 
    };

    $scope.addTrip =function(){
      $scope.showProgressBar=true
      $http.get(backendURL + "/api/addTrip/"+$scope.username+"/"+humanReadableTime($scope.startDate)+"/"+humanReadableTime($scope.endDate))
      .then(function(response){
        $scope.showStatus = true
        $scope.status = "Trip added"
        updateTripList(delay=100)
        calculateRemainingTime(delay=100)
      }); 
    };

    $scope.addUser = function(){
      $scope.userExistsError = false
      $http.get(backendURL + "/api/count/"+$scope.username)
      .then(function(response){
        var cnt = parseInt(response.data)
        if (cnt == 0){
          $http.get(backendURL + "/api/addUser/"+$scope.username+"/"+humanReadableTime($scope.jobStartDate))
          .then(function(response){
            $scope.showRegisterResult = true
            $scope.registerStatus = "User registered."
            $scope.getInfo()
          });
        } else if (cnt >= 1){
          $scope.registerStatus = "User exitsts. Pick another username or select view info."
          $scope.showRegisterResult = true
        }
      })
    };

    $scope.deleteTrip = function(index){
      $scope.showProgressBar=true
      tripDetails = $scope.trips[index]
   
      var confirm = $mdDialog.confirm()
            .title('Would you like to delete your trip?')
            .textContent("Departure:"+ tripDetails.start + "  Arrival:" + tripDetails.end)
            .ariaLabel('Lucky day')
            .ok('Please do it!')
            .cancel('Back to list');
  
      $mdDialog.show(confirm).then(function() {
        $http.get(backendURL + "/api/removeTrip/"+$scope.username+"/"+tripDetails.start+"/"+tripDetails.end)
        .then(function(response){
          updateTripList(delay=100)
          calculateRemainingTime(delay=100)
        }); 


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

    function calculateRemainingTime(delay=0){
      setTimeout(function() { 
        $http.get(backendURL + "/api/getRemainingDays/"+$scope.username)
        .then(function(response){
          $scope.remainingDays = response.data.requiredDays
          $scope.exactDate = addDays(new Date(), $scope.remainingDays)
          $scope.showProgressBar=false
        });
      })
    }

    function updateTripList(delay=0){
      setTimeout(function() { 
        $http.get(backendURL + "/api/getTrips/"+$scope.username)
        .then(function(response){
          $scope.trips = [];
          var vals = response.data.reverse()
          for (i=0; i<vals.length; i++){
            vals[i].start = humanReadableTime(vals[i].start)
            vals[i].end = humanReadableTime(vals[i].end)
          }
          angular.forEach(vals, function(t) {
            $scope.trips.push(t)
          });
        }); 
      }, delay);
    }

    function humanReadableTime(inputT){
      var t = new Date(inputT)
      var ret = ""
      if (t.getDate() < 10){
        ret += "0"
      }
      ret += t.getDate() 
      ret += "-"
      if (t.getMonth()+1 < 10){
        ret += "0"
      }
      var tmp = t.getMonth() + 1
      ret += tmp + "-" + t.getFullYear()
      return ret
    }
  });