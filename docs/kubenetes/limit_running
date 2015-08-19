### How to configure runtime resource constraint
----
1. เราสามารถทำการ limit resource ที่ user ใช้ได้โดยการกำหนดเป็น parameter ใน file json หรือไฟล์ YAML และทำการสร้าง pods, replication controller ดังนี้

```
{
  "kind": "ReplicationController",
  "apiVersion": "v1",
  "metadata": {
    "name": "goweb-controller",
    "labels": {
      "state": "serving"
    }
  },
  "spec": {
    "replicas": 2,
    "selector": {
      "app": "goweb-app"
    },
    "template": {
      "metadata": {
        "labels": {
          "app": "goweb-app"
        }
      },
      "spec": {
        "volumes": null,
        "containers": [
          {
            "name": "goweb-server",
            "image": "goweb",
            "resources": {
              "limits": {
                "cpu": "100m",
                "memory": "50Mi"
               }
            },
            "ports": [
              {
                "containerPort": 80,
                "protocol": "TCP"
              }
            ],
            "imagePullPolicy": "IfNotPresent"
          }
        ],
        "restartPolicy": "Always",
        "dnsPolicy": "ClusterFirst"
      }
    }
  }
}
```
หลังจากนั้นให้ทำการสร้าง pod และ replication controller จาก file 
```
  $ kubectl create -f filename.json
```
* สามารถสร้างเป็น json ก็ได้
