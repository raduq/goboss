( function() {
    angular.module( 'goboss', ['ngResource'] )
        .controller( 'gobossController',[  '$scope','$resource', function($scope, $resource){

        function initialize(){
            setScopeMethods();
        }

        function setScopeMethods(){
            $scope.start = doStart;
            $scope.deploy = doDeploy;
            $scope.clean = doClean;
            $scope.kill = doKill;
        }

        function doStart(){
            $resource('/goboss/start').save( function(data){
                $scope.started = true;
            });
        }

        function doDeploy(){
            $resource('/goboss/build').save( function(data){
                $scope.deployed = true;
            });
        }

        function doClean(){
          $resource('/goboss/clean').save( function(data){
              $scope.deployed = false;
          });
        }

        function doKill(){
            $resource('/goboss/kill').save( function(data){
                $scope.started = false;
            });
        }

        initialize();
    }]);
} )();
