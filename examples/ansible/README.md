# Ansible agent deployment

You can use this Ansible playbook if you have are going to have, or already have a park of Linux servers.

## Prerequisites

* [ansible](https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html) needs to be installed on your PC
* Linux servers must be running *systemd* service manager
* Linux servers must already have your SSH key for authorization

## Deploy the agent

1. Fill the `hosts` file inventory with server IP addresses, info provided there can be used as an example.
2. Run the command `$ ansible-playbook -i hosts playbook.yaml`

## Verify if agent is running

You can connect to your server with SSH and verify if agent is working properly using: 

`$ journalctl -u mg -f`