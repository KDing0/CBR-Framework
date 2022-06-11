package main

import (
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

/*---FILE DESCRIPTION---
This file contains the code used to save and load the replay data necessary for the CBR ai to function.
CBR data is stored in whatever directory desired, but the name of the file should be based on the character
for which it contains data.
---FILE DESCRIPTION---*/

func (x *CBRData) insertReplaytoCaseData(replay *CBRData_ReplayFile) bool {
	x.ReplayFile = append(x.ReplayFile, replay)
	return true
}

//saves the CBR Data of a specific character, directory should be the place to save in and saveFileName is the filename based on the characters name
func saveCBRData(cbr *CBRData, directory string, saveFileName string) bool {
	//saveFile := cbr.ReplayFile[0].CharName[cbr.ReplayFile[0].CbrFocusCharNr]

	data, err := proto.Marshal(cbr)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	_, err = create(strings.ToValidUTF8(directory+saveFileName+".cbr", ""))
	if err != nil {
		log.Fatal("Savedata Creation error: ", err)
	}

	ioutil.WriteFile(directory+saveFileName+".cbr", data, 0666)

	return true
}

// loads the CBR data from a directory, saveFileName is the name of the file based on the character to load
func loadCBRData(directory string, saveFileName string) *CBRData {
	cbr := CBRData{}
	data, err := ioutil.ReadFile(strings.ToValidUTF8(directory+saveFileName+".cbr", ""))
	if err == nil {
		err := proto.Unmarshal(data, &cbr)
		if err != nil {
			log.Fatal("unmarshaling error: ", err)
		}
	}
	return &cbr
}

//saves the CBR Data of a specific character, directory should be the place to save in and saveFileName is the filename based on the characters name
func saveDebugData(text string, directory string, saveFileName string) bool {
	//saveFile := cbr.ReplayFile[0].CharName[cbr.ReplayFile[0].CbrFocusCharNr]

	file, err := create(strings.ToValidUTF8(directory+saveFileName+".log", ""))
	if err != nil {
		log.Fatal("Savedata Creation error: ", err)
	}
	file.WriteString(text)

	return true
}

// create is used to make a directory and file before saving if it does not exist
func create(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}

func deleteAllData(directory string) bool {
	err := removeContents(directory)
	if err != nil {
		return false
	}
	return true
}
func deleteCharData(directory string, saveFileName string) bool {
	err := removeFile(directory, strings.ToValidUTF8(saveFileName+".cbr", ""))
	if err != nil {
		return false
	}
	return true
}
func removeContents(dir string) error {
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

func removeFile(dir string, filename string) error {
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
		if name == filename {
			err = os.Remove(filepath.Join(dir, name))
		}
		if err != nil {
			return err
		}
	}
	return nil
}
