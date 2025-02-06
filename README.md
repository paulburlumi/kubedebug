# kubedebug

`kubedebug` is a terminal user interface (TUI) wrapper for `kubectl debug`. It provides a set of interactive forms using [`huh?`](https://github.com/charmbracelet/huh) to select the resources you want to debug. Specifically, it allows you to choose the:

* Context
* Namespace
* Pod
* Container
* Image

## Usage

To get started with `kubedebug`, follow these steps:

1. **Install `kubedebug`**:
    The easiest way to install `kubedebug` is by using the `go install` command:

    ```shell
    go install github.com/paulburlumi/kubedebug/cmd/kubedebug@latest
    ```

1. **Run `kubedebug`**:
    To debug your container, simply run the command and complete the interactive forms:

    ```shell
    kubedebug
    ```

1. **Specify a Custom Image**:
    By default, `kubedebug` uses the `busybox` image for debugging. If you prefer a different image, you can specify it using the `-i` flag:

    ```shell
    kubedebug -i alpine
    ```

1. **Set Up an Alias**:
    For convenience, you can set up an alias. On Linux, using `bash`, add the following line to your `~/.bashrc` file:

    ```shell
    alias kdebug='kubedebug -i alpine'
    ```

1. **Navigate the Forms**:
    Follow the instructions at the bottom of each form for help on navigating the options:

    ![image showing navigation options](https://camo.githubusercontent.com/fbdf5bc8b45878a6959083f9840d1743af8f8ad87f8e21678595e37651dcd4ed/68747470733a2f2f7668732e636861726d2e73682f7668732d377746715a6c784d576762576d4f49704271584a54692e676966)
