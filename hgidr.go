package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type Record struct {
	Season  int
	Episode int
}

type DataFile struct {
	filename        string
	records         map[string]Record
	defaultFileMode os.FileMode
}

func (self *DataFile) initRecords() {
	self.records = make(map[string]Record, 0)
}

/*
read reads the csv file `datafile` into `records`.
*/
func (self *DataFile) read() (err error) {
	os.MkdirAll(path.Dir(self.filename), 0755)
	b, err := ioutil.ReadFile(self.filename)

	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("The data file doesn't exist.")
			// The file does not exist - create and initialize it.
			f, err := os.Create(self.filename)
			f.Close()
			self.initRecords()
			return err
		} else {
			return
		}
	}

	if len(b) == 0 {
		fmt.Println("The data file is empty.")
		self.initRecords()
		return
	}

	err = json.Unmarshal(b, &self.records)
	return
}

/*
write writes `records` into `filename` as JSON
*/
func (self DataFile) write() (err error) {
	data, err := json.MarshalIndent(self.records, "", "  ")

	if err != nil {
		return
	}

	err = ioutil.WriteFile(self.filename, data, self.defaultFileMode)
	return
}

/*
createNewSeries creates a new `Record` for a series called `name` with the
Season and Episode initialized to 1.
*/
func (self *DataFile) createNewSeries(name string) {
	fmt.Println("Creating ", name)
	newRecord := Record{Season: 1, Episode: 1}
	self.records[name] = newRecord
}

/*
incEpisode increments the last watched episode of `name`.
*/
func (self *DataFile) incEpisode(name string) {
	record := self.records[name]
	record.Episode++
	self.records[name] = record
}

/*
incSeason increments the season of `name`.
*/
func (self *DataFile) incSeason(name string) {
	record := self.records[name]
	record.Season++
	record.Episode = 1
	self.records[name] = record
}

/*
setEpisode sets the last watched episode of `name` to `episode`.
*/
func (self *DataFile) setEpisode(name string, episode int) {
	record := self.records[name]
	record.Episode = episode
	self.records[name] = record
}

/*
setSeason sets the season of `name` to `season`.
*/
func (self *DataFile) setSeason(name string, season int) {
	record := self.records[name]
	record.Season = season
	self.records[name] = record
}

/*
stats displays information about series `name`.
*/
func (self DataFile) stats(name string) {
	record := self.records[name]
	fmt.Printf("Season %d Episode %d", record.Season, record.Episode)
}

/*
get_data_path builds and returns the path to the data file.
*/
func get_data_path() (p string) {
	p = os.Getenv("XDG_DATA_HOME")

	if p == "" {
		p = path.Join(os.Getenv("HOME"), ".local/share")
	}

	p = path.Join(p, "hgidr/data.json")
	return
}

/*
read_datafile reads and returns the data file.
*/
func read_datafile() (datafile DataFile) {
	datafile = DataFile{filename: get_data_path(), defaultFileMode: 0644}
	err := datafile.read()
	if err != nil {
		panic(err)
	}
	return
}

func main() {
	var newSeries = flag.Bool("newseries", false, "Create a new series")
	var episode = flag.Bool("ep", false, "Increment the episode counter of the series")
	var season = flag.Bool("season", false, "Increment the season counter of the series")
	var setEpisode = flag.Int("set-ep", 0, "Set the episode counter of the series")
	var setSeason = flag.Int("set-season", 0, "Set the Episode counter of the series")

	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("You need to specify the name of the series")
		return
	}

	name := args[0]

	datafile := read_datafile()

	if *newSeries {
		datafile.createNewSeries(name)
	}
	if *episode {
		datafile.incEpisode(name)
	}
	if *season {
		datafile.incSeason(name)
	}
	if *setEpisode > 0 {
		datafile.setEpisode(name, *setEpisode)
	}
	if *setSeason > 0 {
		datafile.setSeason(name, *setSeason)
	}

	datafile.stats(name)

	datafile.write()
}
