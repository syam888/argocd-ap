# argocd-ap

**A command-line wrapper for Argo CD application commands.**

`ap` lets you select an Argo CD application once and then run standard Argo CD app commands without retyping the application name.

## Installation

Download the latest release for your platform from [GitHub Releases](https://github.com/syam888/argocd-ap/releases/latest), then place the executable on your `PATH`.

## Usage

```bash
ap select                    # Select an application to work with
ap current                   # Display the currently selected application
ap clear                     # Clear the current application selection
ap get [args]                # Get application status
ap diff [args]               # Show application differences
ap sync [args]               # Synchronize application state
ap history [args]            # View application sync history
ap resources [args]          # List application resources
ap manifests [args]          # Display application manifests
```

## Requirements

* `argocd` CLI installed and authenticated
* `kubectl` configured for the same cluster as Argo CD
* Go 1.22+ only required for building from source
