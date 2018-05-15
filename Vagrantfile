# Configuration file for Vagrant
Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/trusty64"
  config.vm.provision :shell, path: "setup.sh"
  config.vm.network "public_network"
end
