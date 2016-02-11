## How to use Thoth api
We provided api for control thoth's core over kubernetes and HAproxy

#### Running api server
`# ./start.sh`

#### get infomation via GET
 + `/nodes` = list  all nodes
 + `/node/{name}` = list detials specific node by name
 + `/pods` = list all pods
 + `/pod/{Name}` = list detials specific pod by name

#### POST to api
 + `/pod/create` = send json from to create new pod 

#### pod json structure

``` 
	{
		"name":  "your_pod_name",
		"image": "your_pod_image",
		"port":   your_port_number(int),
		"memory": your_cpu_limit(number in bytes),
		"cpu":    your_memory_limit,
	}
```


# !!!! Thoth API is under development !!! 