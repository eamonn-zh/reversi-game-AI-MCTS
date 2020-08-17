# ubuntu_mirror = 'http://mirror.its.sfu.ca/mirror/ubuntu/'
ubuntu_mirror = 'http://mirror.csclub.uwaterloo.ca/ubuntu/'
ubuntu_release = 'bionic'
ubuntu_version = '18.04'
username = 'vagrant'
user_home = '/home/' + username
project_home = user_home + '/project' # you may need to change the working directory to match your project


python3_packages = '/usr/local/lib/python3.6/dist-packages'
ruby_gems = '/var/lib/gems/2.5.0/gems/'


# Get Ubuntu sources set up and packages up to date.

template '/etc/apt/sources.list' do
  variables(
    :mirror => ubuntu_mirror,
    :release => ubuntu_release
  )
  notifies :run, 'execute[apt-get update]', :immediately
end
execute 'apt-get update' do
  action :nothing
end
execute 'apt-get upgrade' do
  command 'apt-get dist-upgrade -y'
  only_if 'apt list --upgradeable | grep -q upgradable'
end
directory '/opt'
directory '/opt/installers'


# Basic packages many of us probably want. Includes gcc C and C++ compilers.
package ['build-essential', 'cmake']


# Other core language tools you might want
package ['python3', 'python3-pip', 'python3-dev']  # Python


# Go (more modern than Ubuntu golang-go package)
execute 'snap install --classic go' do
end


# FastAPI
execute 'pip3 install fastapi' do
end


# ASGI server
execute 'pip3 install uvicorn' do
end


# NodeJS (more modern than Ubuntu nodejs package) and NPM
remote_file '/opt/installers/node-setup.sh' do
  source 'https://deb.nodesource.com/setup_14.x'
  mode '0755'
end
execute '/opt/installers/node-setup.sh' do
  creates '/etc/apt/sources.list.d/nodesource.list'
  notifies :run, 'execute[apt-get update]', :immediately
end
package ['nodejs']


# Compile Go share library
execute 'export GOPATH=`pwd` && go build -buildmode=c-shared -o MonteCarloTreeSearch.so main' do
  cwd project_home + '/Go_MCTS'
  user username
  environment 'HOME' => user_home
end


# Copy Go share library to FastAPI_Server folder
execute 'cp MonteCarloTreeSearch.so ' + project_home + '/FastAPI_Server/lib' do
  cwd project_home + '/Go_MCTS'
  user username
  environment 'HOME' => user_home
end

execute 'cp MonteCarloTreeSearch.h ' + project_home + '/FastAPI_Server/lib' do
  cwd project_home + '/Go_MCTS'
  user username
  environment 'HOME' => user_home
end


# Run FastAPI Server
execute 'python3 main.py &' do
  cwd project_home + '/FastAPI_Server'
  user username
  environment 'HOME' => user_home
end


# Run Node.js Server
execute 'npm install && npm run serve &' do
  cwd project_home + '/Web_Reversi'
  user username
  environment 'HOME' => user_home
end

