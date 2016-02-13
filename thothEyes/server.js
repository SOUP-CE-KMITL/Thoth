var app 	 = require('express')(),
        express = require('express'),
	Sequence = require('sequence').Sequence,
	http 	 = require('http'),
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


app.listen(1337);
console.log('Application started at port 1337');
