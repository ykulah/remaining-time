angular.module('remainingTimeApp', [])
  .controller('remainingTimeController', function($scope, $http) {
    $scope.getInfo = function(){
      $scope.showNewTripForm = true
      var user = $scope.username
      $http.get("https://backend-dot-remaining-time-c9dd7.appspot.com/api/getTrips/"+user)
      .then(function(response){
        console.log(response.data)
        $scope.trips = [];
        angular.forEach(response.data, function(t) {
          $scope.trips.push(t)
        });
      });

      $http.get("https://backend-dot-remaining-time-c9dd7.appspot.com/api/getRemainingDays/"+user)
      .then(function(response){
        console.log(response)
        $scope.remainingDays = response.data.requiredDays
      });
    };

    $scope.addTrip =function(){
      var st = new Date($scope.startDate)
      console.log(st.getDate())
      var start = ""
      if (st.getDate() < 10){
        start += "0"
      }
      start = st.getDate() + "-" 
      if (st.getDay() < 10){
        start += "0"
      }
      var tmp = st.getMonth() + 1
      start += tmp + "-" + st.getFullYear()


      var e = new Date($scope.endDate)
      end = ""
      if (e.getDate() < 10){
        end += "0"
      }
      end += e.getDate() + "-"
      if (e.getMonth() < 10){
        end += "0"
      }
      var tmp = e.getMonth() + 1
      end += tmp + "-" + e.getFullYear()

      $http.get("https://backend-dot-remaining-time-c9dd7.appspot.com/api/addTrip/"+$scope.username+"/"+start+"/"+end)
      .then(function(response){
        $scope.showStatus = true
        $scope.status = "Trip added"
        $http.get("https://backend-dot-remaining-time-c9dd7.appspot.com/api/getTrips/"+$scope.username)
        .then(function(response){
          $scope.trips = [];
          angular.forEach(response.data, function(t) {
            $scope.trips.push(t)
          });
        });
        $http.get("https://backend-dot-remaining-time-c9dd7.appspot.com/api/getRemainingDays/"+$scope.username)
        .then(function(response){
          $scope.remainingDays = response.data.requiredDays
        });
      }); 
    }
  });