vagrant_assets = File.dirname(__FILE__) + "/vagrant"

Vagrant.configure("2") do |config|
    config.vm.box = "ubuntu/bionic64"
    config.vm.hostname = "containerd-bermuda-gap"
    config.vm.provider "virtualbox" do |v|
        v.memory = 3072
        v.cpus = 2
    end

    config.vm.provision "file", source: "File.dirname(__FILE__)", destination: "$HOME/workdir"
    config.vm.provision "shell", path: path: "#{vagrant_assets}/provision.sh", privileged: true
end