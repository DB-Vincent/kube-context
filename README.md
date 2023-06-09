# kube-context
![kube-context demo](./demo/demo.gif)
kube-context is a simple and easy-to-use CLI tool written in Go, which allows you to choose a Kubernetes config in a user-friendly way. It simplifies the process of switching between Kubernetes contexts by providing a menu-driven interface to list, select and switch between Kubernetes contexts.

## Installation
### Homebrew
Please note that kube-context requires the gcc package to be installed, so please install it with your favorite package manager.

```shell
brew tap DB-Vincent/kube-context https://github.com/DB-Vincent/kube-context
brew install DB-Vincent/kube-context
```

### Manual
#### Linux
```shell
export LATEST_VERSION=$(curl --silent "https://api.github.com/repos/DB-Vincent/kube-context/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
mkdir /tmp/kube-context
wget -qO- https://github.com/DB-Vincent/kube-context/releases/download/$LATEST_VERSION/kube-context_Linux_x86_64.tar.gz | tar xvz -C /tmp/kube-context
mv /tmp/kube-context/kube-context ~/.local/bin/kube-context
rm -rf /tmp/kube-context
chmod +x ~/.local/bin/kube-context
```

You should now be able to execute the `kube-context` command and switch between contexts.

#### Windows
```powershell
$LATEST_VERSION = ((Invoke-WebRequest -Uri https://api.github.com/repos/DB-Vincent/go-aws-mfa/releases/latest).Content | ConvertFrom-JSON).tag_name
Invoke-WebRequest -Uri https://github.com/DB-Vincent/kube-context/releases/download/$LATEST_VERSION/kube-context_Windows_x86_64.zip -OutFile ~\Downloads\kube-context.zip
Expand-Archive ~\Downloads\kube-context.zip -DestinationPath ~\Downloads\kube-context
Move-Item -Path ~\Downloads\kube-context -Destination ~/Documents/
```

You can now execute `~\Documents\kube-context\kube-context.exe` to run the kube-context command.
In order to execute the command without `~\Documents\kube-context\`, add the `~\Documents\kube-context\` path to your environment variables.

## Usage
To use kube-context, simply run the kube-context command in your terminal: `kube-context`

This will display a list of available Kubernetes contexts, which you can select from using the arrow keys and Enter key.

Once you have selected a context, kube-context will switch your current context to the one you selected.

### Renaming a context

![kube-context-rename](./demo/demo-rename.gif)

### Setting a default namespace

![kube-context-rename](./demo/demo-default-namespace.gif)

## Contributing
If you want to contribute to kube-context, you can fork the repository and make your changes. Once you are done with your changes, create a pull request and we will review your changes.

## License
kube-context is licensed under the [GNU GPLv3](https://github.com/DB-Vincent/kube-context/blob/v0.0.1/LICENSE) license.