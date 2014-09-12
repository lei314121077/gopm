// Copyright 2014 Unknwon
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/gpmgo/gopm/modules/base"
	"github.com/gpmgo/gopm/modules/cli"
	"github.com/gpmgo/gopm/modules/errors"
	"github.com/gpmgo/gopm/modules/log"
	"github.com/gpmgo/gopm/modules/setting"
)

var CmdBuild = cli.Command{
	Name:  "build",
	Usage: "link dependencies and go build",
	Description: `Command build links dependencies according to gopmfile,
and execute 'go build'

gopm build <go build commands>`,
	Action: runBuild,
	Flags: []cli.Flag{
		cli.BoolFlag{"update, u", "update pakcage(s) and dependencies if any", ""},
		cli.BoolFlag{"remote, r", "build with pakcages in gopm local repository only", ""},
		cli.BoolFlag{"verbose, v", "show process details", ""},
	},
}

func buildBinary(ctx *cli.Context, args ...string) error {
	_, target, err := parseGopmfile(setting.GOPMFILE)
	if err != nil {
		return err
	}

	if err := linkVendors(ctx); err != nil {
		return err
	}

	log.Info("Building...")

	cmdArgs := []string{"go", "build"}
	cmdArgs = append(cmdArgs, ctx.Args()...)
	if err := execCmd(setting.DefaultVendor, setting.WorkDir, cmdArgs...); err != nil {
		return fmt.Errorf("fail to build program: %v", err)
	}

	if setting.IsWindowsXP {
		fName := path.Base(target)
		binName := fName + ".exe"
		os.Remove(binName)
		exePath := path.Join(setting.DefaultVendorSrc, target, binName)
		if base.IsFile(exePath) {
			if err := os.Rename(exePath, path.Join(setting.WorkDir, binName)); err != nil {
				return fmt.Errorf("fail to move binary: %v", err)
			}
		} else {
			log.Warn("No binary generated")
		}
	}
	return nil
}

func runBuild(ctx *cli.Context) {
	if err := setup(ctx); err != nil {
		errors.SetError(err)
		return
	}

	os.RemoveAll(setting.DefaultVendor)
	if !setting.Debug {
		defer os.RemoveAll(setting.DefaultVendor)
	}
	if err := buildBinary(ctx, ctx.Args()...); err != nil {
		errors.SetError(err)
		return
	}

	log.Info("Command executed successfully!")
}
