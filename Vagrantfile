# Configuration file for Vagrant
Vagrant.configure("2") do |config|
  config.vm.box = "hashicorp/precise64"
  config.vm.provision :shell, path: "mysql_setup.sh"
end
