# Configuration file for Vagrant
Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/trusty64"
  config.vm.provision :shell, path: "setup.sh"
  config.vm.network "private_network", ip: "192.168.50.4"
  #config.vm.network "forwarded_port", guest: 2002, host: 8080
end
