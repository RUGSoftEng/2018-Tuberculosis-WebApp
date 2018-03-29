# 2018-Tuberculosis

## Vagrant
Waring: 
* Some commands may only availible for bash and not for the 'amazing ;)' command prompt from Windows. It might be better to use 'PowerShell' instead.
* I (Sten) am also new to Vagrant, so if there are any problems, suggestions or improvements, please tell!

### Prerequisites
Both the latest versions of:
* [Vagrant](https://www.vagrantup.com/downloads.html)
* [VirtualBox](https://www.virtualbox.org/wiki/Downloads)
Note that for Debian distributions (Ubuntu, Mint etc.) the apt repositories do not contain the latest version of Vagrant and does not work.

### Starting
Somewhere in the project folder:
run `vagrant up` to start the Virtual Machine (this may take a minute)

### Connecting
When the VM is started:
`vagrant ssh` to start a ssh connection. To close use `Ctrl-D` or type `exit` in your shell. This will only stop the ssh connection, not the complete VM.

The directory `/vagrant` is a shared directory. It is the same as the project root directory. Any changes in this directory also occur in the project root directory and vice versa.

### Closing
There are three options for closing the VM:
* `vagrant suspend` Stops, but remembers the current state of the VM. This means it will use all the memory it was using, but will start up fast. It is similar to 'suspending' your own system (or closing the lid of your laptop).
* `vagrant halt` Closes down the VM, similar as to when you would shut down your system. Takes up less memory than suspending, but also will start up slower.
* `vagrant destroy` This completely shuts down the VM and removes all the used resources. It will not leave any memory in use. This is generally what you want to do if you are done.

Starting the VM back up use `vagrant up`.
Restaring can be done by running `vagrant reload`. This is similar to running halt and up.


For a more complete introduction go to the [Vagrant introduction tab](https://www.vagrantup.com/intro/getting-started/index.html) or to the [documentation](https://www.vagrantup.com/docs/index.html)
