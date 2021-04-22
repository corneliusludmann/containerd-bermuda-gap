vagrant_assets = File.dirname(__FILE__) + "/vagrant"
client_dir = File.dirname(__FILE__) + "/client"
registry_data = File.dirname(__FILE__) + "/registry"
logs_dir = File.dirname(__FILE__) + "/logs"

Vagrant.configure("2") do |config|
    config.vm.box = "ubuntu/bionic64"
    config.vm.hostname = "containerd-bermuda-gap"
    config.vm.provider "virtualbox" do |v|
        v.memory = 1024
        v.cpus = 1
    end
    config.disksize.size = "80GB"

    config.vm.synced_folder "#{logs_dir}", "/home/vagrant/logs"

    config.vm.provision "file", source: "#{client_dir}", destination: "/home/vagrant/workdir/client"
    config.vm.provision "file", source: "#{registry_data}", destination: "/home/vagrant/registry"
    config.vm.provision "shell", path: "#{vagrant_assets}/provision.sh", privileged: true
end
