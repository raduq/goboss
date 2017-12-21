( function() {
    angular.module( 'goboss', ['ngResource'] )
        .controller( 'gobossController',[  '$scope','$resource', function($scope, $resource){

        function initialize(){
            setScopeMethods();
            setScopeVars();
        }

        function setScopeMethods(){
            $scope.start = doStart;
            $scope.deploy = doDeploy;
            $scope.clean = doClean;
            $scope.kill = doKill;
        }

        function setScopeVars(){
          $scope.status = {
            started : false,
            deployed : false
          };
        }

        function doStart(){
          $scope.status.started = true;
          $resource('/goboss/start').save( function(data){
            console.log('Started!');
          });
        }

        function doDeploy(){
          $scope.status.deployed = true;
          $resource('/goboss/build').save( function(data){
            console.log('Deployed!');
          });
        }

        function doClean(){
          $scope.status.deployed = false;
          $resource('/goboss/clean').save( function(data){
            console.log('Cleaned!');
          });
        }

        function doKill(){
          $scope.status.started = false;
          $resource('/goboss/kill').save( function(data){
            console.log('Killed!');
          });
        }

        initialize();
    }]);
} )();
