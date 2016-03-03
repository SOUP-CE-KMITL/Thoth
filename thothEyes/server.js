var app 	 = require('express')(),
    express  = require('express'),
	Sequence = require('sequence').Sequence,
	http 	 = require('http'),
	request	 = require('request'),
	async    = require('async'),
	swig 	 = require('swig'),
	people;

// setup swig for render file before serve to user 
app.engine('html', swig.renderFile);

app.use(express.static(__dirname + '/public'));
app.set('view engine', 'html');
app.set('views', __dirname + '/views');

app.set('view cache', false);

swig.setDefaults({ cache: false });

//app.use('/bower_components/nvd3', express.static( __dirname + '/bower_components/nvd3'));

app.get('/', function (reg, res) {
	res.render('Thoth');
});

app.get('/node', function (reg, res) {
	res.render('node');
});

app.get('/nodes', function (reg, res) {
	res.render('nodes');
});

var api_server_ip   = 'localhost'
var api_server_port = '8182'

app.get('/node/:nodeName', function(reg, res){

	var name = reg.params.nodeName
	var info = {}

	// connect to api server 
	var options = {
		host: api_server_ip,
		port: api_server_port,
		path: '/node/'+name,
		method: 'GET'
	};

	//  waiting until http request is already finish.  
	sequence = Sequence.create();
	sequence
		.then(function(next) {
			var req = http.request(options, function(res){
			  console.log('STATUS: ' + res.statusCode);
			  console.log('HEADERS: ' + JSON.stringify(res.headers));
			  res.setEncoding('utf8');
			  res.on('data', function (data) {
			    data = JSON.parse(data);	
			    node = data;

			    // create node information
				info = {
					nodeName: name,
					create_at: node.metadata.creationTimestamp,
					limit_cpu: node.status.capacity.cpu,
					limit_memory: node.status.capacity.memory,
					limit_pods: node.status.capacity.pods,
					status: node.status.conditions[0].type,
					address: node.status.addresses[0].address
				}
				// pass value to next function
				next("", info)
			  });
			});
			// handle error request
			req.on('error', function(e) {
			  console.log('problem with request: ' + e.message);
			});
			req.end();

		}).then(function(next, err, info){
			console.log(info);
			console.log("it's already serve")
			// render view with variable by swig
			if (err) return console.log(err);
			res.render('nodes', info)
		});
});

app.get('/error/node', function(reg, res) {
	// This number should configuration by user.
	err5xx_threshold = 0;
	
	namespace = "thoth";
	apps_name = [];
	apps_error = [];
	sequence = Sequence.create();

	var get_app_name = {
		host: api_server_ip,
		port: api_server_port,
		path: '/apps/'+namespace,
		method:  'GET'
	}

	sequence
	.then(function(next) {
		var req = http.request(get_app_name, function(res) {
			console.log('STATUS: ' + res.statusCode);
			console.log('HEADERS: ' + JSON.stringify(res.headers));
			res.setEncoding('utf8');
			res.on('data', function (data) {
				data = JSON.parse(data);
				var app = {};
				// get all application name
				for(var i = 0;i < data.items.length;i++){
					app.name = data.items[i].metadata.name;
					app.replicas  = data.items[i].spec.replicas;
					app.containers  = data.items[i].spec.template.spec.containers; 
					apps_name.push(app);
				}
				// not have any error and pass app name to next
				next("", apps_name)
			});
		});
		req.on('error', function(e) {
		  console.log('problem with request: ' + e.message);
		});
		req.end();
	}).then(function(next, err, apps){
		console.log(apps);
		if (err) return console.log(err);
		// check 5xx status
		var get_app_error;
		apps_error = [];

		async.forEach(apps, function(item, callback){
			get_app_error = {
				host: api_server_ip,
				port: api_server_port,
				path: '/app/'+item.name+'/metrics/'+namespace
			}
			var req = http.request(get_app_error, function(res){
				console.log(res.statusCode)
				res.on('data', function (data) {
					data = JSON.parse(data);
					if(data.Response5xx >= err5xx_threshold){
						apps_error.push(item);
					}
					callback();
				});
			});
			req.on('error', function(e) {
			  console.log('problem with request: ' + e.message);
			  res.status(500).json("error");
			});
			req.end();
		}, function() {
			res.status(200).json(apps_error);
		});
	});
});


app.listen(1337);
console.log('Application started at port 1337');
