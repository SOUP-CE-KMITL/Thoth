angular.module('thothEyes', ['nvd3'])
.controller('DashboardCtrl', function($scope){

  function getPieChartData() {
    return  [
      { 
        "label": "Extreme",
        "value" : 10,
        "color" : '#e74c3c'


      } , 
      { 
        "label": "Dangerous",
        "value" : Math.floor((Math.random() * 20) + 1),
        "color" : '#e67e22'

      } , 
      { 
        "label": "Becareful",
        "value" : Math.floor((Math.random() * 20) + 1),
        "color" : '#f1c40f'

      } , 
      { 
        "label": "Normal",
        "value" : Math.floor((Math.random() * 20) + 1),
        "color" : '#2ecc71'

      } 

    ];
  }
    
    $scope.computeOptions = {
      chart: {
        type: 'pieChart',
        height: 300,
        x: function(d){return d.label;},
        y: function(d){return d.value;},
        showLabels: true,
        labelThreshold: 0.03,
        labelType: 'percent',
        donut: true,
        donutRatio: 0.3

      }
    }

    //set initial mock data, going to remove after connect to API.
    $scope.computeData = getPieChartData();

    $scope.dataOptions = angular.copy($scope.computeOptions);
    $scope.dataData = angular.copy($scope.computeData);

    $scope.fileOptions = angular.copy($scope.computeOptions);
    $scope.fileData = angular.copy($scope.computeData);

    setInterval(function(){
      // update data from API here.
      $scope.computeData = getPieChartData();
      $scope.dataData = getPieChartData();
      $scope.fileData = getPieChartData();
      //$scope.api.updateWithData($scope.data);
      $scope.$apply();
    }, 5000);

  })
.controller('NodesCtrl', ['$scope', '$http', '$q', function($scope, $http, $q){
  //for contain all nodes.
  nodes = ['thotheyes.cloudapp.net'];
  node_datas = [];
  $scope.nodes = [];

  for(var i = 0; i < nodes.length; i++){
    // get application profile
    // initail array
    $scope.nodes[i] = {}
    $scope.nodes[i].data = [];
    $scope.nodes[i].data[0] = {};
    $scope.nodes[i].data[0].values = [];
    $scope.nodes[i].data[1] = {};
    $scope.nodes[i].data[1].values = [];

    $scope.nodes[i].data[0].key = 'cpu';
    $scope.nodes[i].data[0].color = '#e74c3c';
    $scope.nodes[i].data[1].key = 'memory';
    $scope.nodes[i].data[1].color = '#f39c12';
  }
  function getNodesData() {
    return  [
    ];
  }
    
    $scope.resourceOptions = {
      chart: {
        type: 'lineChart',
        height: 190,
        margin: {
          left: 100
        },
        duration: 0,
        yDomain: [0,100],
        x: function(d){ return d.x },
        y: function(d){ return d.y },
        useInteractiveGuideline: true,
        xAxis: {
          axisLabel: 'Time (s)'
        },
        yAxis: {
          axisLabel: 'Resource (%)',
          tickFormat: function(d){
            return d3.format('.02f')(d);
          },
          axisLabelDistance: -10
        }
      }
    }

    //set initial mock data, going to remove after connect to API.


    var t = 0;

    setInterval(function(){
      for(var c = 0; c < nodes.length; c++){
       node_datas[c] = $http.get("http://"+nodes[c]+":8182/metrics");
      }

      $q.all(node_datas).then(function(response){
        for(var c=0; c < response.length; c++){
         $scope.nodes[c].name = nodes[c] ;
         $scope.nodes[c].data[0].values.push({ x:t, y:response[c].data.cpu });
         $scope.nodes[c].data[1].values.push({ x:t, y:response[c].data.memory.used_percent});
         //console.log("cpu : "+$scope.apps[r].data[1].values.slice(-1)[0].x);
         if($scope.nodes[c].data[0].values.length > 20) $scope.nodes[c].data[0].values.shift();
         if($scope.nodes[c].data[1].values.length > 20) $scope.nodes[c].data[1].values.shift();
         t++;
        }
       console.log($scope.nodes[0]);
      });;
      
    }, 1000);

  }]);
