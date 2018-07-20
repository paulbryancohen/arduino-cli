/*
 * This file is part of arduino-cli.
 *
 * arduino-cli is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
 *
 * As a special exception, you may use this file as part of a free software
 * library without restriction.  Specifically, if other files instantiate
 * templates or use macros or inline functions from this file, or you compile
 * this file and link it with other files to produce an executable, this
 * file does not by itself cause the resulting executable to be covered by
 * the GNU General Public License.  This exception does not however
 * invalidate any other reasons why the executable file might be covered by
 * the GNU General Public License.
 *
 * Copyright 2017 ARDUINO AG (http://www.arduino.cc/)
 */

package lib

import (
	"os"

	"github.com/bcmi-labs/arduino-cli/arduino/libraries/librariesindex"
	"github.com/bcmi-labs/arduino-cli/arduino/libraries/librariesmanager"
	"github.com/bcmi-labs/arduino-cli/commands"
	"github.com/bcmi-labs/arduino-cli/common/formatter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func initInstallCommand() *cobra.Command {
	installCommand := &cobra.Command{
		Use:   "install LIBRARY[@VERSION_NUMBER](S)",
		Short: "Installs one of more specified libraries into the system.",
		Long:  "Installs one or more specified libraries into the system.",
		Example: "" +
			"arduino lib install AudioZero       # for the latest version.\n" +
			"arduino lib install AudioZero@1.0.0 # for the specific version.",
		Args: cobra.MinimumNArgs(1),
		Run:  runInstallCommand,
	}
	return installCommand
}

func runInstallCommand(cmd *cobra.Command, args []string) {
	logrus.Info("Executing `arduino lib install`")
	lm := commands.InitLibraryManager(nil)

	refs := librariesindex.ParseArgs(args)
	downloadLibraries(lm, refs)
	installLibraries(lm, refs)
}

func installLibraries(lm *librariesmanager.LibrariesManager, refs []*librariesindex.Reference) {
	libReleasesToInstall := []*librariesindex.Release{}
	for _, ref := range refs {
		rel := lm.Index.FindRelease(ref)
		if rel == nil {
			formatter.PrintErrorMessage("Error: library " + ref.String() + " not found")
			os.Exit(commands.ErrBadCall)
		}
		libReleasesToInstall = append(libReleasesToInstall, rel)
	}

	for _, libRelease := range libReleasesToInstall {
		// FIXME: the library is installed again even if it's already installed

		logrus.WithField("library", libRelease).Info("Installing library")

		if _, err := librariesmanager.Install(libRelease); err != nil {
			logrus.WithError(err).Warn("Error installing library ", libRelease)
			formatter.PrintError(err, "Error installing library: "+libRelease.String())
			os.Exit(commands.ErrGeneric)
		}

		formatter.Print("Installed " + libRelease.String())
	}
}
