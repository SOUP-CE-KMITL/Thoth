alias k=kubectl

# create
alias kc='kubectl $NS create'
#
# describe
alias kdp='kubectl $NS describe pods'
alias kds='kubectl $NS describe service'
alias kdr='kubectl $NS describe rc'
alias kdl='kubectl $NS describe limits'
alias kdq='kubectl $NS describe quota'
#
# Get
alias kgp='kubectl $NS get pods'
alias kgs='kubectl $NS get service'
alias kgr='kubectl $NS get rc'
alias kgl='kubectl $NS get limits'
alias kgq='kubectl $NS get quota'
alias kga='kubectl $NS get all'
#
# delete
alias kxp='kubectl $NS delete pods'
alias kxs='kubectl $NS delete service'
alias kxr='kubectl $NS delete rc'
alias kxl='kubectl $NS delete limits'
alias kxq='kubectl $NS delete quota'

