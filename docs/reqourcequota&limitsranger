## Thoth :: How to configure ResourceQuota & LimitsRanger

#### ResourceQuota
สามารถใช้เมื่อต้องการกำหนดทรัพยากรของเครื่องที่ผู้ใช้ใน  Paas คนหนึ่งจะใช้ได้โดยจะถูกผูกเข้ากับ namespace ซึ่ง namespace คือผู้ใช้งานนันเอง แต่ ResourceQuota จะไม่สามารถทำงานได้หากไม่มีการกำหนด Limits

#### LimitsRanger

Limits เปรียบเสมือนกรอบขอบเขตของ Resource ที่ Pod จะนำไปใช้งานได้ ซึ่งแต่ละชนิดของ pod เราสามารกำหนดให้มีขอบเขตที่ต่างกันได้

เนื่องจากทั้ง 2 อย่างที่กล่าวไปข้างต้นเป็น plugin ของ apiserver ที่เราต้องเปิดขึ้นมาก่อนค่าต่างๆที่เซตถึงจะทำงาน

#### ขั้นตอนการ configure

1. เข้าไปแก้ image ที่จะใช้สร้าง environment จำลอง เพื่อกำหนด่าให้ plugin ของ ResourceQuota และ LimitsRanger ทำงาน ขั้นตอนแรกทำการ run container ของ image `gcr.io/google_containers/hyperkube:v1.0.1` ด้วยคำสั่ง
```
	$ docker run -it  gcr.io/google_containers/hyperkube:v1.0.1 bash
```
2. หลังจากเข้าสู่ container แล้วผ่าน bash ให้เข้าไปที่โฟลเดอร์ที่เก็บการเซตค่าเริ่มต้นของ kubernetes apiserver ดังนี้ 
```
	$ vi /etc/kubernetes/manifests/master.json
```

3. สร้าง environment จำลองขึ้นมาเพื่อให้ได้ containner ของ kubernetes ที่ีรันอยู่บน docker ตามลิงค์นี้ [Kubernetes locally intallation ](https://github.com/kubernetes/kubernetes/blob/master/docs/getting-started-guides/docker.md)

2. 
