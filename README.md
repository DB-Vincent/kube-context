# kube-context
![kube-context demo](./demo.gif)
kube-context is a simple and easy-to-use CLI tool written in Go, which allows you to choose a Kubernetes config in a user-friendly way. It simplifies the process of switching between Kubernetes contexts by providing a menu-driven interface to list, select and switch between Kubernetes contexts.

## Installation
kube-context can be installed by downloading the binary and executing it:

```bash
wget -O ~/.local/bin/kube-context https://github.com/DB-Vincent/kube-context/releases/download/v0.0.1/kube-context
chmod +x ~/.local/bin/kube-context
```

## Usage
To use kube-context, simply run the kube-context command in your terminal:
```bash
kube-context
```

This will display a list of available Kubernetes contexts, which you can select from using the arrow keys and Enter key.

Once you have selected a context, kube-context will switch your current context to the one you selected.

## Contributing
If you want to contribute to kube-context, you can fork the repository and make your changes. Once you are done with your changes, create a pull request and we will review your changes.

## License
kube-context is licensed under the [GNU GPLv3](https://github.com/DB-Vincent/kube-context/blob/v0.0.1/LICENSE) license.