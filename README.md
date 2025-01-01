# gothic

This is a fork of [github.com/nsf/gothic](https://github.com/nsf/gothic).

The main changes in this fork are:

* Adds Go modules support (`go.mod`)
* Adds support for using a Tcl/Tk [KitDLL](https://kitcreator.rkeene.org/fossil/wiki?name=KitDLL)
* Uses [pkg-config](https://en.wikipedia.org/wiki/pkg-config) to locate Tcl/Tk libraries instead of hard-coded paths

The original [README](README) file is included, without the `VERSION NOTICE` section since it doesn't apply due to the use of `pkg-config`.

## Quickstart Guide

There are a number of options for obtaining Tcl/Tk lib and header files.  One of the easiest is to download the KitDLL SDK nighlty build from the [KitCreator](https://kitcreator.rkeene.org/) website.  This will also allow the distribution of the Go executable with a single DLL file and does not require Tcl/Tk to be installed on the target systems.

On the KitCreator [nightly builds](https://www.rkeene.org/devel/kitcreator/kitbuild/nightly/) page, scroll down to the section corresponding to your desired Tcl version, OS, and arch.  Note that there are downloads for both Tcl versions 8.5 and 8.6, as well as separate builds for 32-bit 64-bit.

In the "**Built as a Library (sdk)**" row for the desired version, click the "**sdk**" link to download the SDK.  The "**Built as a Library**" link is just the DLL, but the SDK is required to build the go binary, and it also includes the required DLL.

Unpack the downloaded `.tar.gz` file to a local directory, which will be referenced as `$LIBDIR` below, and copy the `assets/libtclkit.pc` file from this repo to the same directory, so it looks like the following:

```
$LIBDIR
├── bin
│   └── ...
├── doc
│   └── ...
├── include
│   └── ...
├── lib
│   └── ...
└── libtclkit.pc
```

Edit the `libtclkit.pc` file so the **Libs:** line references the same version that was downloaded.  It should match the version number in the filenames of the versioned files in the `lib` directory.

From the `lib` directory, copy the `.dll` or `.so` file to your project directory where you will run `go run` or the built executable so it can be found during execution.  This file will also need to be included along with the final executable.

Set the `PKG_CONFIG_PATH` env variable to `$LIBDIR` mentioned above, to tell `pkg-config` where `libtclkit.pc` is located.

**NOTE**: Basic widgets, such as `label` and `button` don't work with the cross-compiled Windows KitDLL from KitCreator, but `ttk` widgets do work, which also look better.  To support basic widgets, a KitDLL of the same version will need to be obtained from elsewhere, or built locally from source.  Refer to the KitCreator [documentation](https://kitcreator.rkeene.org/fossil/doc/trunk/README?mimetype=text/plain) for building a fully working Windows KitDLL (using MinGW).

Run `go build` to build your Go application.  See below for further details when using Windows.

### Obtaining a smaller KitDLL

The KitDLL incluced in the SDK has extra packages that may not be required.  THe KitCreator website has a [web interface](https://kitcreator.rkeene.org/kitcreator) for building just the KitDLL with various options.  On this page, select the Tcl version and Platform at the top and following options below, and click the Create button to create a much smaller KitDLL:
* Package: Tk
* Kit: Build Library (KitDLL)
* Kit: "Minimal" build (remove extra packages shipped as part of Tcl and reduce encodings)


## Adding extra packages

The Tcl/Tk KitDLL allows addition of extra packages by appending a zip file on the end of the DLL.

Some packages may have only binary components, some may have only Tcl components, and some may be a combination of binary and Tcl components.  All of these are supported with this method.

Start by creating a `tcl-libs.zip` file with the required contents.  The name of the zip file can be anything, but the contents must follow this structure.  Binary components (`.dll`/`.so` files) go in the root directory of the zip, and Tcl components (`.tcl` files) go in the directory named for the package under the `lib` directory.

For example, to add [TkTreeCtrl](https://tktreectrl.sourceforge.net/) and [tooltip](https://core.tcl-lang.org/tklib/doc/trunk/embedded/md/tklib/files/modules/tooltip/tooltip.md) (from Tklib), download the packages, and create a `zip` file with the following contents:

```
tcl-libs.zip
├── lib
│   ├── tooltip
│   │   ├── pkgIndex.tcl
│   │   ├── tipstack.tcl
│   │   └── tooltip.tcl
│   └── treectrl2.4.1
│       ├── filelist-bindings.tcl
│       ├── pkgIndex.tcl
│       └── treectrl.tcl
└── treectrl24.dll
```

Note that the `tooltip` package uses basic widgets, so it will not work with the cross-compiled Windows KitDLL from KitCreator.  The package is included in this repo in `assets/tooltip` with the basic widgets replaced with Ttk widgets to allow it to work with the cross-compiled Windows KitDLL from KitCreator.  Compare this version with the official version for an example of how to make this modification and apply it to other packages.

To create the final KitDLL with the additional packages, make a backup copy of the original KitDLL, and concatenate it and the zip file to make the zip file contents available to Tcl.

```
C:> copy /b libtclkit860.dll libtclkit860.dll.bak
C:> copy /b libtclkit860.dll.bak + tcl-libs.zip libtclkit860.dll
```
or
```
$ cp libtclkit860.so libtclkit860.so.bak
$ cat libtclkit860.so.bak tcl-libs.zip > libtclkit860.so
```

With this updated KitDLL, the included packages can be loaded with the standard method in Tcl.  The location of the `.dll`/`.so` files that were placed in the root directory of the zip file are available in the `/.KITDLL_USER/` directory, and the `lib` directory is automatically added to the list of directories Tcl looks in for packages:

```
load /.KITDLL_USER/treectrl24.dll
package require treectrl
package require tooltip
```

## Adding an icon to the Windows executable

Use [goversioninfo](https://github.com/josephspurrier/goversioninfo) to add version information and an icon to the `.exe` built with `go build`.

Add an icon file (using `icon.ico` here) and a `versioninfo.json` file (which could contain as little as just `{}`) in the directory where you run `go build`, and a Go directive in the Go source code:

```go
//go:generate goversioninfo -icon=icon.ico -64=true
```

Then run `go generate`, which should generate a `resource.syso` file which will be used by the compiler when running `go build` to add the version information and icon to the file.

See the [goversioninfo](https://github.com/josephspurrier/goversioninfo) documentation for details.

## Adding an icon to the Windows application's system menu / taskbar button

With an icon added to the executable using the method above, the same icon can be used in the system menu and task bar with the addition of the following.  Note that `os.Args[0]` is outside the Tcl init string.

```go
func main() {
	ir := tk.NewInterpreter(`
		wm iconbitmap . -default {` + os.Args[0] + `}
	`)
	<-ir.Done
}
```

## Removing the Windows terminal window when running the Windows executables

When running the Windows executable, a terminal window is shown automatically.  To hide this window, build it with:

```
go build -ldflags -H=windowsgui .
```
