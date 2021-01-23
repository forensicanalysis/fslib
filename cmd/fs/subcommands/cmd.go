// Copyright (c) 2019 Siemens AG
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//
// Author(s): Jonas Plum

package subcommands

import (
	"crypto/md5"  // #nosec
	"crypto/sha1" // #nosec
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"github.com/forensicanalysis/fslib"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"

	"github.com/spf13/cobra"
	"github.com/xlab/treeprint"

	"github.com/forensicanalysis/fslib/content"
	"github.com/forensicanalysis/fslib/filesystem/recursivefs"
	"github.com/forensicanalysis/fslib/filetype"
)

const (
	bashCompletionFunc = `__fs_ls_completion()
{
local fs_output out

if [[ -z $cur ]]; then
	return
fi

if fs_output=$(fs ls "$cur" 2>/dev/null); then
	if [[ $cur == */ ]]; then
		COMPREPLY=($( compgen -P "$cur" -W "${fs_output}" ))
		return
	fi
	COMPREPLY=( "$cur " ) # TODO add / for folder ...; add already in first completion...
	return
fi

parent=$(dirname $cur)
if [[ parent == "." ]]; then
	parent=""
fi
partialbase=$(basename $cur)
if fs_output=$(fs ls "$parent" 2>/dev/null); then
	COMPREPLY=($( compgen -P "$parent/" -W "${fs_output}" -- "$partialbase" ))
fi
}

__fs_get_resource()
{
__fs_ls_completion
if [[ $? -eq 0 ]]; then
	return 0
fi
}

__fs_custom_func() {
case ${last_command} in
	fs_cat | fs_ls)
		__fs_get_resource
		return
		;;
	*)
		;;
esac
}
`
)

// Execute runs the fs command.
func Execute() {
	exitOnError(executeCommands())
}

func executeCommands() error {
	var debug bool
	rootCmd := &cobra.Command{Use: "fs", Short: "recursive file, filesystem and archive commands",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log.SetFlags(log.LstdFlags | log.Llongfile)
			if !debug {
				log.SetOutput(ioutil.Discard)
			}
		},
		BashCompletionFunction: bashCompletionFunc,
	}

	cat := &cobra.Command{Use: "cat", Short: "print files", Run: catCmd}
	file := &cobra.Command{Use: "file", Short: "determine file type", Run: fileCmd}
	hashsum := &cobra.Command{Use: "hashsum", Short: "print hashsums", Run: hashsumCmd}
	ls := &cobra.Command{Use: "ls", Short: "list directory contents", Run: lsCmd}
	stat := &cobra.Command{Use: "stat", Short: "display file status", Run: statCmd}
	strings := &cobra.Command{Use: "strings", Short: "find the printable strings in an object, or other binary, file", Run: stringsCmd}
	tree := &cobra.Command{Use: "tree", Short: "list contents of directories in a tree-like format", Run: treeCmd}
	complete := &cobra.Command{Use: "complete", Run: func(cmd *cobra.Command, args []string) {
		if err := rootCmd.GenBashCompletionFile(".bash_completion.sh"); err == nil {
			log.Println("--")
			if _, err := os.Stat("/usr/local/etc/bash_completion.d"); !os.IsNotExist(err) {
				err = os.Rename(".bash_completion.sh", "/usr/local/etc/bash_completion.d/fs")
				log.Println(err)
			} else if _, err := os.Stat("/etc/bash_completion.d"); !os.IsNotExist(err) {
				err = os.Rename(".bash_completion.sh", "/etc/bash_completion.d/fs")
				log.Println(err)
			}
		}
	}}

	/*
		go install .
		fs complete
		. /usr/local/etc/bash_completion
		fs ls <tab>
	*/

	rootCmd.AddCommand(cat, file, hashsum, ls, stat, strings, tree, complete)
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug output")
	return rootCmd.Execute()
}

func exitOnError(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		// os.Exit(1)
		log.Fatal(err)
	}
}

func catCmd(_ *cobra.Command, args []string) {
	for _, arg := range args {
		func() {
			name, err := fslib.ToForensicPath(arg)
			exitOnError(err)
			r, err := recursivefs.New().Open(name)
			exitOnError(err)
			defer r.(*recursivefs.Item).Close()

			_, err = io.Copy(os.Stdout, r)
			exitOnError(err)
		}()
	}
}

func fileCmd(_ *cobra.Command, args []string) {
	b := make([]byte, 8192)
	for _, arg := range args {
		func() {
			name, err := fslib.ToForensicPath(arg)
			exitOnError(err)
			f, err := recursivefs.New().Open(name)
			exitOnError(err)
			defer f.(*recursivefs.Item).Close()
			_, err = f.Read(b)
			exitOnError(err)
			// f.Close()
			fmt.Printf("%s: %s\n", arg, filetype.Detect(b).Mimetype.Value)
		}()
	}
}

func hashsumCmd(_ *cobra.Command, args []string) {
	for _, arg := range args {
		md5hash := md5.New()   // #nosec
		sha1hash := sha1.New() // #nosec
		sha256hash := sha256.New()
		sha512hash := sha512.New()
		hash := io.MultiWriter(md5hash, sha1hash, sha256hash, sha512hash)

		name, err := fslib.ToForensicPath(arg)
		exitOnError(err)
		r, err := recursivefs.New().Open(name)
		exitOnError(err)
		_, err = io.Copy(hash, r)
		exitOnError(err)
		exitOnError(r.Close())
		fmt.Printf("MD5: %x\n", md5hash.Sum(nil))
		fmt.Printf("SHA1: %x\n", sha1hash.Sum(nil))
		fmt.Printf("SHA256: %x\n", sha256hash.Sum(nil))
		fmt.Printf("SHA512: %x\n", sha512hash.Sum(nil))
	}
}

func lsCmd(_ *cobra.Command, args []string) {
	if len(args) == 0 {
		args = append(args, ".")
	}

	fs := recursivefs.New()
	for _, arg := range args {
		name, err := fslib.ToForensicPath(arg)
		exitOnError(err)
		fi, err := recursivefs.New().Stat(name)
		exitOnError(err)
		if fi.IsDir() {
			f, err := recursivefs.New().Open(name)
			exitOnError(err)

			infos, err := fslib.ReadDir(f, 0)
			exitOnError(err)
			fileNames := fslib.InfosToNames(infos)

			sort.Strings(fileNames)
			for _, fileName := range fileNames {
				child, err := fs.Stat(path.Join(name, fileName))
				if err != nil {
					fmt.Println(fileName, err)
					continue
				}
				if child.IsDir() {
					fmt.Println(fileName + "/")
				} else {
					fmt.Println(fileName)
				}
			}
		} else {
			fmt.Println(arg)
		}
	}
}

func statCmd(_ *cobra.Command, args []string) {
	for _, arg := range args {
		name, err := fslib.ToForensicPath(arg)
		exitOnError(err)
		fi, err := recursivefs.New().Stat(name)
		exitOnError(err)
		fmt.Printf("Name: %v\n", fi.Name())
		fmt.Printf("Size: %v\n", fi.Size())
		fmt.Printf("IsDir: %v\n", fi.IsDir())
		fmt.Printf("Mode: %s\n", fi.Mode())
		fmt.Printf("Modified: %s\n", fi.ModTime())
	}
}

func stringsCmd(_ *cobra.Command, args []string) {
	for _, arg := range args {
		name, err := fslib.ToForensicPath(arg)
		exitOnError(err)
		fsys := recursivefs.New()
		r, err := fsys.Open(name)
		exitOnError(err)

		fx, err := fslib.FileX(r)
		exitOnError(err)

		ascii, err := content.Content(fx)
		exitOnError(err)
		fmt.Println(ascii)
	}
}

func treeCmd(_ *cobra.Command, args []string) {
	for _, arg := range args {
		tree := treeprint.New()
		tree.SetValue(arg)
		getChildren(tree, arg)
		fmt.Println(tree.String())
	}
}

func getChildren(tree treeprint.Tree, subpath string) {
	name, err := fslib.ToForensicPath(subpath)
	exitOnError(err)
	fi, err := recursivefs.New().Stat(name)
	if err != nil {
		fmt.Println("tree", err)
		return
	}
	exitOnError(err)
	if fi.IsDir() {
		f, err := recursivefs.New().Open(name)
		if err != nil {
			fmt.Println("tree", err)
			return
		}

		infos, err := fslib.ReadDir(f, 0)
		exitOnError(err)
		files := fslib.InfosToNames(infos)
		exitOnError(f.Close())
		sort.Strings(files)
		for _, item := range files {
			child := tree.AddBranch(item)
			getChildren(child, path.Join(subpath, item))
		}
	}
}
