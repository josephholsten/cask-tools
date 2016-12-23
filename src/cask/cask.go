package cask

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"appcast"
	"general"
)

type Cask struct {
	Name     string
	Versions []Version
	Content  string
}

// New creates a new Cask instance based on provided parameters and returns its
// pointer. Support cask names with and without extension.
func New(name string, a ...interface{}) *Cask {
	name = strings.TrimSuffix(name, filepath.Ext(name)) // remove file extension

	c := new(Cask)
	c.Name = name

	if len(a) > 0 {
		path := a[0].(string)
		c.ReadCaskFile(path)
	} else {
		c.ReadCaskFile()
	}

	c.ExtractVersionsWithAppcasts()

	return c
}

// String returns a string representation of the Cask.
func (self Cask) String() string {
	return self.Name
}

// IsOutdated checks if one of the cask versions current checkpoints doesn't
// match the latest one. In this case the cask is considered to be outdated.
func (self Cask) IsOutdated() (result bool) {
	for _, version := range self.Versions {
		if version.Appcast.Checkpoint.Latest != "" && version.Appcast.Checkpoint.Current != version.Appcast.Checkpoint.Latest {
			return true
		}
	}

	return false
}

// AddVersion adds a provided Version to the cask Versions.
func (self *Cask) AddVersion(version Version) {
	self.Versions = append(self.Versions, version)
}

// ReadCaskFile reads cask file and save its content into Content struct.
// By default uses current working directory (pwd) for reading.
func (self *Cask) ReadCaskFile(a ...interface{}) {
	path := fmt.Sprintf("%s.rb", self.Name)
	if len(a) > 0 {
		path = a[0].(string)
	}

	self.Content = string(general.GetFileContent(path))
}

// GetStanzaValues gets cask stanza values as a string array. Also checks if the
// stanza is global (belongs to all versions in the cask). Returns string array
// of stanza values and bool array of global statuses.
func (self Cask) GetStanzaValues(stanza string) ([]string, []bool) {
	re := regexp.MustCompile(fmt.Sprintf(`.*%s ['|"](?P<value>.*)['|"]`, stanza))
	matches := re.FindAllStringSubmatch(self.Content, -1)

	result := make([]string, len(matches))
	global := make([]bool, len(matches))

	regexIfElseEnd := regexp.MustCompile("(?s)if.*(?:els(?:e|if)?).*?end")
	regexIfElseEndContent := ""
	if regexIfElseEnd.MatchString(self.Content) {
		regexIfElseEndContent = regexIfElseEnd.FindAllString(self.Content, -1)[0]
	}

	for i, match := range matches {
		result[i] = match[1]

		re := regexp.MustCompile(match[1])
		if regexIfElseEndContent == "" || !re.MatchString(regexIfElseEndContent) {
			global[i] = true
		} else {
			global[i] = false
		}
	}

	return result, global
}

// ExtractVersionsWithAppcasts extracts versions, appcasts and checkpoints from
// the cask content and gets new releases from appcasts. The releases then are
// added as a cask Versions.
func (self *Cask) ExtractVersionsWithAppcasts() {
	// remove previous versions
	self.Versions = nil

	// extract stanza values
	versions, _ := self.GetStanzaValues("version")
	appcasts, appcastsGlobal := self.GetStanzaValues("appcast")
	checkpoints, _ := self.GetStanzaValues("checkpoint:")

	var keys []string
	versionsAppcasts := map[string][]string{}

	if len(versions) > len(appcasts) {
		// when there are more versions than appcasts
		for i, appcast := range appcasts {
			versionsMatches := map[string]string{}

			for _, version := range versions {
				keys = append(keys, version)
				versionsAppcasts[version] = []string{}
				re := regexp.MustCompile(fmt.Sprintf(`(?s)(%s).*(?:%s)`, version, appcast))
				matches := re.FindAllStringSubmatch(self.Content, -1)

				versionsMatches[version] = matches[0][0]
			}

			encountered := map[string]bool{}

			for j, key := range keys {
				encountered[key] = false
				if !appcastsGlobal[i] {
					for l, k := range keys {
						if key != k && l > j {
							match := versionsMatches[k]
							re := regexp.MustCompile(match)
							if re.MatchString(match) {
								encountered[key] = true
								continue
							}
						}
					}
				}
			}

			for _, key := range keys {
				if encountered[key] == false {
					versionsAppcasts[key] = []string{appcasts[i], checkpoints[i]}
				}
			}
		}
	} else if len(versions) == len(appcasts) {
		// when number of versions are appcasts are equal
		for i, version := range versions {
			keys = append(keys, version)
			versionsAppcasts[version] = []string{appcasts[i], checkpoints[i]}
		}
	}

	for _, key := range keys {
		v := NewVersion(key)
		if len(versionsAppcasts[key]) > 0 {
			a := appcast.New(versionsAppcasts[key][0])
			a.Checkpoint.Current = versionsAppcasts[key][1]

			v = NewVersion(key, a)
		}

		self.AddVersion(*v)
	}
}

// LoadAppcasts gets new releases for each version from appcasts.
func (self *Cask) LoadAppcasts() {
	for i, version := range self.Versions {
		version.LoadAppcast()
		self.Versions[i] = version
	}
}