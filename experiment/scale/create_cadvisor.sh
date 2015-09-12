sudo docker run \
  --volume=/:/rootfs:ro \
  --volume=/var/run:/var/run:rw \
  --volume=/sys:/sys:ro \
  --volume=/var/lib/docker/:/var/lib/docker:ro \
  --publish=4194:8080 \
  --allow_dynamic_housekeeping=true \
  --housekeeping_interval=1s \
  --detach=true \
  --name=cadvisor \
  google/cadvisor:latest
