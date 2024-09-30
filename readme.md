# SSH Key-Pair Manager CLI

This is a Go-based CLI tool that helps you manage multiple SSH key-pairs for different companies or Git repositories. It allows you to switch between SSH certificates, generate new key-pairs, and manage your SSH setup efficiently.

## Features

- List available SSH key-pair and switch between them.
- Generate new key-pairs with different algorithms.
- Manage folders for storing key-pairs.
- Uses default values for algorithms and comments if not provided.

## Installation

1. Clone the repository and navigate to the project folder:

```
git clone git@github.com:DavidGarciaCat/ssh-key-manager-cli.git
cd ssh-key-manager-cli
```

2. Build the binary:

Run the following command to build the Go binary:

```
go build -o ssh-manager main.go
```

3. Install the binary globally:

To make the ssh-manager available globally, move it to a directory in your $PATH (e.g., /usr/local/bin):

```
sudo mv ssh-manager /usr/local/bin/
```

Now you can run the ssh-manager from any location in your terminal.

## Usage

1. Launch the SSH Key Manager CLI:

```
ssh-manager
```

Example output:

```
SSH Key Manager CLI
-------------------

Current active key-pair: _personal

1) Switch SSH key-pair for another system
2) Generate new SSH key-pair
3) Quit

Choose an option:
```

2. Change the active SSH key-pair:

```
...
Choose an option: 1

Available SSH key-pair folders:

1) _company_a
2) _company_b
3) _personal

Enter the number of the folder to switch to:
```

Choose the new one to switch to it.

```
Enter the number of the folder to switch to: 1

Removing SSH key-pair files...

Removed: /home/<username>/.ssh/id_ed25519
Removed: /home/<username>/.ssh/id_ed25519.pub

Copying SSH key-pair files...

Copied /home/<username>/.ssh/_company_a/id_ed25519 to /home/<username>/.ssh/id_ed25519
Copied /home/<username>/.ssh/_company_a/id_ed25519.pub to /home/<username>/.ssh/id_ed25519.pub

Switched to _company_a
```

3. Generate a new SSH key-pair:

```
Choose an option: 2

Enter the subfolder name for the new key-pair (required): testing key pair

Created folder: /home/<username>/.ssh/_testing_key_pair

Enter the cipher (e.g., ed25519, rsa) [default: ed25519]:
Enter the key name [default: id_ed25519]:
Enter a comment for the key [default: username@hostname]:

Generating new key pair...

Generating public/private ed25519 key pair.
...

Generated new key pair in /home/<username>/.ssh/_testing_key_pair
```
