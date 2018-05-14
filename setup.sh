#!/usr/bin/env bash
         
echo -e '----------------------------------------------\n| Running Setup For WebApp Backend (By Sten) |\n----------------------------------------------'

echo 'Installing MySQL Server...'
sudo debconf-set-selections <<< 'mysql-server mysql-server/root_password password root'
sudo debconf-set-selections <<< 'mysql-server mysql-server/root_password_again password root'

sudo apt-get update -q > /dev/null
sudo apt-get install -yq mysql-server > /dev/null && echo 'SUCCESS'


echo 'Installing Git...'
sudo apt-get install -yq git > /dev/null && echo 'SUCCESS'


echo 'Downloading and Setting up Golang...'
# Download Go package 
wget -q https://dl.google.com/go/go1.10.1.linux-amd64.tar.gz \
    && tar -C /usr/local -zxf go1.10.1.linux-amd64.tar.gz \
    && rm -rf go1.10.1.linux-amd64.tar.gz \
    && echo 'SUCCESS'

echo 'Setting PATH Variables...'
export PROJROOT=/vagrant
export GOPATH=/home/vagrant/go
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:$GOPATH/bin

echo 'Setting up paths for the vagrant ssh user'
echo -e 'export GOPATH=/home/vagrant/go' >> /home/vagrant/.bashrc
echo -e 'export PATH=$PATH:/usr/local/go/bin' >> /home/vagrant/.bashrc
echo -e 'export PATH=$PATH:$GOPATH/bin' >> /home/vagrant/.bashrc
echo -e 'export PROJROOT=/vagrant' >> /home/vagrant/.bashrc
echo 'DONE'

echo 'Downloading Golang Dependencies...'
make deps -C $PROJROOT/API && echo 'SUCCESS'


echo 'Setting up Database...'
mysql -u root --password=root <<< 'drop database if exists TestDB' \
    && mysql -u root --password=root <<< 'create database TestDB' \
    && mysql -u root --password=root TestDB < $PROJROOT/database/DDL_statements.sql \
    && mysql -u root --password=root TestDB < $PROJROOT/database/test_insert_statements.sql \
    && echo 'SUCCESS'

# OLD! Command for automatically starting the API
echo "function old_start_api() { echo -e \"root\nroot\nTestDB\n192.168.50.4:2002\" | make run --quiet -C $PROJROOT/API; }" >> /home/vagrant/.bashrc

# Command for completely reloading the Database
echo "function reload_db() { mysql -u root --password=root <<< \"drop database TestDB\" \\
    && mysql -u root --password=root <<< \"create database TestDB\" \\
    && mysql -u root --password=root TestDB < $PROJROOT/database/DDL_statements.sql \\
    && mysql -u root --password=root TestDB < $PROJROOT/database/test_insert_statements.sql; }" >> /home/vagrant/.bashrc

# Command for starting the API, where you specify the arguments:
# 1] = username for database user
# 2] = corresponding password for database user
# 3] = name of the database
# 4] = location 'ip:port' for api to listen
echo "function start_api_args() { echo -e \"\$1\n\$2\n\$3\n\$4\" | make run --quiet -C $PROJROOT/API; }" >> /home/vagrant/.bashrc

# Command for retrieving the ip
echo "function get_ip() { ifconfig eth1 | grep -o 'inet addr:[1-9.]*' | cut -d \":\" -f 2; }" >> /home/vagrant/.bashrc

# Command for starting the api with a specified ip address
echo "function start_api_ip() { echo -e \"root\nroot\nTestDB\n\${1}:2002\" | make run --quiet -C $PROJROOT/API; }" >> /home/vagrant/.bashrc

# New command for automatically starting the API
echo "function start_api { start_api_ip \$(get_ip); }" >> /home/vagrant/.bashrc

echo -e '------------------\n| Finished Setup |\n------------------'

