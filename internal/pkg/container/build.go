package container

import (
	"fmt"
	"github.com/hpcng/warewulf/internal/pkg/errors"
	"github.com/hpcng/warewulf/internal/pkg/util"
	"github.com/hpcng/warewulf/internal/pkg/wwlog"
	"os"
	"os/exec"
	"path"
)

func Build(name string, buildForce bool) (string, error) {

	rootfsPath := RootFsDir(name)
	imagePath := ImageFile(name)

	if ValidSource(name) == false {
		return "", errors.New("Container does not exist")
	}

	if buildForce == false {
		wwlog.Printf(wwlog.DEBUG, "Checking if there have been any updates to the VNFS directory\n")
		if util.PathIsNewer(rootfsPath, imagePath) {
			return "Skipping (VNFS is current)", nil
		}
	}

	wwlog.Printf(wwlog.DEBUG, "Making parent directory for: %s\n", name)
	err := os.MkdirAll(path.Dir(imagePath), 0755)
	if err != nil {
		return "Failed creating directory", err
	}

	wwlog.Printf(wwlog.DEBUG, "Making parent directory for: %s\n", rootfsPath)
	err = os.MkdirAll(path.Dir(rootfsPath), 0755)
	if err != nil {
		return "Failed creating directory", err
	}

	wwlog.Printf(wwlog.DEBUG, "Building VNFS image: '%s' -> '%s'\n", rootfsPath, imagePath)
	cmd := fmt.Sprintf("cd %s; find . | cpio --quiet -o -H newc | gzip -c > \"%s\"", rootfsPath, imagePath)
	err = exec.Command("/bin/sh", "-c", cmd).Run()
	if err != nil {
		return "Failed building VNFS", err
	}

	return "Done", nil
}
