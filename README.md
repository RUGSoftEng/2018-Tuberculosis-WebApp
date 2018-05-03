# 2018-Tuberculosis

master [![Build Status](https://travis-ci.org/RUGSoftEng/2018-Tuberculosis-WebApp.svg?branch=master)](https://travis-ci.org/RUGSoftEng/2018-Tuberculosis-WebApp)

development [![Build Status](https://travis-ci.org/RUGSoftEng/2018-Tuberculosis-WebApp.svg?branch=development)](https://travis-ci.org/RUGSoftEng/2018-Tuberculosis-WebApp)

dev-web [![Build Status](https://travis-ci.org/RUGSoftEng/2018-Tuberculosis-WebApp.svg?branch=dev-web)](https://travis-ci.org/RUGSoftEng/2018-Tuberculosis-WebApp)

---
* [Vagrant](#vagrant)
* [API](#api)
* [Database](#database)
---


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
* `vagrant destroy` This completely shuts down the VM and removes all the used resources. It will not leave any memory in use.

Starting the VM back up use `vagrant up`.
Restaring can be done by running `vagrant reload`. This is similar to running halt and up.


For a more complete introduction go to the Vagrant [introduction](https://www.vagrantup.com/intro/getting-started/index.html)  tab or to the [documentation](https://www.vagrantup.com/docs/index.html)

## API
### Starting
Starting the API is currently done by ssh-ing to the VM and run `start_api` in the shell. The API will be started and listening on location (ip:port) `192.168.50.4:2002`.

Note that the API is definitely not fully finished, so don't expect the smoothest behaviour. Please make an Issue or something if you want to report an error.
### Documentation
Can be found in the [API Documentation](https://documenter.getpostman.com/view/4186694/tubuddy-api-documentation-for-android-and-ios/RW1aJKy3) file.

## Database
The database TestDB automatically gets created when the Vagrant VM setup is ran. To avoid having to rerun the complete setup every time the database needs to be reset, it is possible to reset the database manually:

In the VM run `reload_db`. This essentially destroys the current database, recreates a new one and runs all the DDL and Insert statements.
