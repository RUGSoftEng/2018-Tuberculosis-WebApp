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
# Adding go binaries to PATH
export PATH=$PATH:/usr/local/go/bin && echo 'SUCCESS'
export GOPATH=/home/vagrant/go && echo 'SUCCESS'
export PROJROOT=/vagrant

# Reads lines from 'go_dependencies' file and installs them
echo 'Installing Go Dependencies...'
while IFS='' read -r line || [[ -n "$line" ]]; do
    echo "Downloading $line..."
    go get -u $line && echo "SUCCESS" || echo "FAILED"
done < "$PROJROOT/go_dependencies"
echo "DONE"

echo 'Setting up paths for the vagrant ssh user'
echo -e 'export GOPATH=/home/vagrant/go' >> /home/vagrant/.bashrc
echo -e 'export PATH=$PATH:/usr/local/go/bin' >> /home/vagrant/.bashrc
echo -e 'export PATH=$PATH:$GOPATH/bin' >> /home/vagrant/.bashrc
echo -e 'export PROJROOT=/vagrant' >> /home/vagrant/.bashrc
echo 'DONE'

echo 'Setting up Database...'
mysql -u root --password=root <<< 'drop database if exists TestDB' \
    && mysql -u root --password=root <<< 'create database TestDB' \
    && mysql -u root --password=root TestDB < $PROJROOT/database/DDL_statements.sql \
    && mysql -u root --password=root TestDB < $PROJROOT/database/test_insert_statements.sql \
    && echo 'SUCCESS'

# Command for automatically starting the API
echo "alias start_api='echo -e \"root\nTestDB\n192.168.50.4:2002\" | make run --quiet -C $PROJROOT/API'" >> /home/vagrant/.bashrc

# Command for completely reloading the Database
echo "alias reload_db='mysql -u root --password=root <<< \"drop database TestDB\" \\
    && mysql -u root --password=root <<< \"create database TestDB\" \\
    && mysql -u root --password=root TestDB < $PROJROOT/database/DDL_statements.sql \\
    && mysql -u root --password=root TestDB < $PROJROOT/database/test_insert_statements.sql'" >> /home/vagrant/.bashrc

echo -e '------------------\n| Finished Setup |\n------------------'
