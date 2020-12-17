package ec2

const userDataTemplate = `Content-Type: multipart/mixed; boundary="//"
MIME-Version: 1.0

--//
Content-Type: text/cloud-config; charset="us-ascii"
MIME-Version: 1.0
Content-Transfer-Encoding: 7bit
Content-Disposition: attachment; filename="cloud-config.txt"

#cloud-config
cloud_final_modules:
- [scripts-user, always]

--//
Content-Type: text/x-shellscript; charset="us-ascii"
MIME-Version: 1.0
Content-Transfer-Encoding: 7bit
Content-Disposition: attachment; filename="userdata.txt"

#!/bin/bash
sudo yum update -y
curl -sL https://rpm.nodesource.com/setup_14.x | sudo -E bash -
sudo yum install -y nodejs
sudo npm install pm2 -g
cd /home/ec2-user
wget -O defaultProject.zip "%s"
unzip defaultProject.zip
cd defaultProject
sudo pkill -f PM2
sudo PORT=80 pm2 start app.js --name server
sudo pm2 save
--//
`
