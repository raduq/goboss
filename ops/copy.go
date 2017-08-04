package ops

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/hpcloud/tail"
)

var wg sync.WaitGroup

//CopyFile from a src to a dst
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

//RemoveContents of a folder
func RemoveContents(dir string) (err error) {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func Start(jb string) (err error) {
	log := "standalone/log/server.log"

	err = CleanLogs(jb)
	if err != nil {
		return err
	}

	sh1 := "bin/standalone.sh -b localhost --server-config=standalone.xml -Djboss.server.base.dir=" + jb + "standalone "
	sh2 := "-P=" + jb + "standalone/configuration/contaazul.properties"
	shPath := jb + sh1 + sh2

	cmd := exec.Command("/bin/sh", "-c", shPath)
	err = cmd.Start()
	if err != nil {
		return err
	}

	err = DoTail(jb + log)
	if err != nil {
		return err
	}
	return nil
}

func Execute(dir, comm string, args ...string) {
	cmd := exec.Command(comm, args...)
	cmd.Dir = dir
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalln(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalln(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatalln(err)
	}
	wg.Add(2)
	go printReader(stdout)
	go printReader(stderr)
	if err := cmd.Wait(); err != nil {
		log.Fatalln(err)
	}
	wg.Wait()
}

func printReader(rd io.ReadCloser) {
	r := bufio.NewReader(rd)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			defer wg.Done()
			if err == io.EOF {
				break
			}
			log.Fatalln("failed to read line: ", err)
		}
		fmt.Println(string(line))
	}
}

func DoTail(file string) (err error) {
	t, err := tail.TailFile(file, tail.Config{Follow: true})
	if err != nil {
		return err
	}
	for line := range t.Lines {
		fmt.Println(line.Text)
	}
	return nil
}

func CleanLogs(jb string) (err error) {
	logdir := jb + "standalone/log"

	exists, err := exists(logdir)
	if err != nil {
		return err
	}

	if exists {
		err = RemoveContents(logdir)
		if err != nil {
			return err
		}
	}
	return nil
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
