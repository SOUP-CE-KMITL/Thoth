Installing
---
This method https://docs.openshift.org/latest/admin_guide/install/advanced_install.html

1. Install ansible in your machine
1. Config Inventory file in /etc/ansible/host (https://docs.openshift.org/latest/admin_guide/install/advanced_install.html#single-master-multi-node)
1. Make sure that you can ssh to target machine without enter any password ( Use ssh-add )
1. Run `ansible-playbook ~/openshift-ansible/playbooks/byo/config.yml`

After install finished
- ssh to master node

### Grant Access to the Privileged SCC
- `oc edit scc privileged`
- Add your username to the end of content
```
#eg.
users:
- system:serviceaccount:openshift-infra:build-controller
- your-user
```




### Deploy docker registry
```
oadm registry --config=admin.kubeconfig --credentials=openshift-registry.kubeconfig 
```
- Get docker registry ip by `oc get all`
- For now can't connect to docker registry https://github.com/openshift/origin/issues/4187


### Deploy router
---
https://docs.openshift.org/latest/admin_guide/install/deploy_router.html
