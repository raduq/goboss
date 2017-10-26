( function() {
    angular.module( 'goboss', ['ngResource'] )
        .controller( 'gobossController',[  '$scope','$resource', function($scope, $resource){
   
        function initialize(){
            setScopeMethods();
            getStatus();
        }

        function setScopeMethods(){
            $scope.start = doStart;
            $scope.deploy = doDeploy;
            $scope.undeploy = doUndeploy;
            $scope.kill = doKill;
            $scope.refresh = doRefresh;
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

        function doUndeploy(){
            $resource('/goboss/unbuild').save( function(data){
                $scope.deployed = false;
            });
        }

        function doKill(){
            $resource('/goboss/kill').save( function(data){
                $scope.started = false;                
            });
        }

        function doRefresh(){
            getStatus();
        }

        function getStatus(){
            $resource('/goboss/status').get( function(data){
                console.log('refresh ' + data);
                $scope.started = data.started;
                $scope.deployed = data.deployed;
            } );
        }

        initialize();
    }]);
} )();