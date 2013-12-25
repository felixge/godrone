# -*- mode: ruby -*-
# vi: set ft=ruby :

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

home="/home/vagrant"
share_dst="#{home}/godrone"
provision_script="scripts/vagrant-provision.bash"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.box = "ubuntu-12.04-64bit"
  # config.vm.box_url = "http://domain.com/path/to/above.box"
  config.vm.synced_folder ".", "#{share_dst}"

  # vagrant seems to execute our script as root by default. let's use the same
  # user as vagrant ssh to avoid permission issues.
  config.vm.provision :shell, :inline => "cd #{share_dst} && ./#{provision_script}", :privileged => false, :keep_color => true
end
