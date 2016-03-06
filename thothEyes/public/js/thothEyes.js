angular.module('thothEyes', ['nvd3'])
.controller('DashboardCtrl', ['$scope', '$http', function($scope, $http){
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
  var errorAppsNum = 0;
  function checkErrorApp() {
    $http.get('/monitor/apps')
     .success(function(data, status) {
        if(data.length != 0){
          console.log(data);
          $scope.errorApps = data;
        }
      });
  }
  checkErrorApp();

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
      donutRatio: 0.3,
      pie: {
        dispatch: {
          elementClick: function (element){ 
            document.location.href = document.location.href + "nodes#" + element.data.label.toLowerCase();
          }
        }
      }
    }
  }

    //set initial mock data, going to remove after connect to API.
    $scope.computeData = getPieChartData();

    $scope.dataOptions = angular.copy($scope.computeOptions);
    $scope.dataData = angular.copy($scope.computeData);

    $scope.fileOptions = angular.copy($scope.computeOptions);
    $scope.fileData = angular.copy($scope.computeData);

    setInterval(function(){
      checkErrorApp();
      // update data from API here.
      $scope.computeData = getPieChartData();
      $scope.dataData = getPieChartData();
      $scope.fileData = getPieChartData();
      //$scope.api.updateWithData($scope.data);
      $scope.$apply();
    }, 5000);
  }])
.controller('NodesCtrl', ['$scope', '$http', '$q', function($scope, $http, $q){
  $scope.testing = "testing";
  //for contain all nodes.
  nodes = ['161.246.6.235'];
  node_datas = [];
  $scope.nodes = [];

  for(var i = 0; i < nodes.length; i++){
    // get application profile
    // initail array
    $scope.nodes[i] = {}
    $scope.nodes[i].name = ""
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
       node_datas[c] = $http.get("https://"+nodes[c]+"/metrics");
      }

      $q.all(node_datas).then(function(response){
        for(var c=0; c < response.length; c++){
         $scope.nodes[c].name = nodes[c] ;
         $scope.nodes[c].data[0].values.push({ x:t, y:response[c].data.cpu });
         $scope.nodes[c].data[1].values.push({ x:t, y:response[c].data.memory.used_percent});
         //console.log("cpu : "+$scope.apps[r].data[1].values.slice(-1)[0].x);
         if($scope.nodes[c].data[0].values.length > 20) $scope.nodes[c].data[0].values.shift();
         if($scope.nodes[c].data[1].values.length > 20) $scope.nodes[c].data[1].values.shift();
         t = t+5;
        }
       console.log($scope.nodes[0]);
      });;
      
    }, 5000);

  }])
.controller('NodeCtrl', ['$scope', '$http', '$q', function($scope, $http, $q){
  $scope.node = {}
  $scope.node.name = "161.246.6.235"
  $scope.node.cpu = [{
    values: [],
    key: 'cpu',
    color: '#ff7f0e'
  }];
  $scope.node.memory = [{
    values: [],
    key: 'memory',
    color: '#ff7f0e'
  }]
  $scope.apps = [];
  appResourcePromises =  []

  $scope.cpuOptions = {
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

  $scope.memoryOptions = angular.copy($scope.cpuOptions);
  $scope.appOptions = {
    chart: {
      type: 'multiBarHorizontalChart',
      x: function(d){return d.label;},
      y: function(d){return d.value;},
      margin: {
        top: 30,
        rigth: 20,
        bottom: 50,
        left: 60
      },
      showValues: true,
      height: 200,
      groupSpacing: 0.5,
      tooltips: false,
      showLegend: false,
      showControls: false,
      forceY: [0, 100],
      yAxis: {
        tickFormat: function(d){
          return d3.format(',.2f')(d);
        }
      }
    }
  }

  var t = 1;
  function getNodeResource(){
    $http.get("https://"+ $scope.node.name+"/metrics").then(
      function(response){
        $scope.node.cpu[0].values.push({ x:t, y: response.data.cpu[0]});
        $scope.node.memory[0].values.push({ x:t, y:response.data.memory.used_percent});
        if($scope.node.cpu[0].values.length > 20) $scope.node.cpu[0].values.shift();
        if($scope.node.memory[0].values.length > 20) $scope.node.memory[0].values.shift();
      },
      function(err){
        console.log(err)
      });
  }

  function getAppsResource(){
      for(var c = 0; c < $scope.apps.length; c++){
        appResourcePromises[c] = $http.get("https://" + $scope.node.name + "/app/" + $scope.apps[c].key + "/metrics/thoth");
      }

      $q.all(appResourcePromises).then(function(response){
        for(var c=0; c < response.length; c++){
          $scope.apps[c].values[0] = {label: "CPU",  value: response[c].data.cpu};
          $scope.apps[c].values[1] = {label: "memory",  value: response[c].data.memory};
        }
        console.log("apps after resolve");
        console.log($scope.apps); });;
  }

  function getApps(){
    $http.get("https://"+ $scope.node.name+"/apps/thoth").then(
      function(response){
        for(var i=0; i < response.data.items.length; i++) {
          $scope.apps.push({key: response.data.items[i].spec.template.metadata.name, color: "#d62728", values: []});
        }
        getAppsResource();
      },
      function(err){
        console.log(err)
      });
  }


  getNodeResource();
  getApps();

  setInterval(function(){
    getNodeResource();
    getAppsResource();
    t = t+5;
  }, 5000);
  }]);
