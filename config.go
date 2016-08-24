package shelflib

import (
    "encoding/json"
    "os"
    "runtime"
    "path"
)

// Config is the struct representing shelflib config.
type Config struct {
    ShelfHost string
    ShelfPathConst string
}

// Load loads the default config in the repository base directory and file name config.json.
// This config can be overriden with an override config file.
// It returns type Config
func LoadConfig() (Config, error){
    _, filename, _, _ := runtime.Caller(1)
    file, err := os.Open(path.Join(path.Dir(filename), "config.json"))

    if err != nil {
        return nil, err
    }

    decoder, err := json.NewDecoder(file)

    if err != nil {
        return nil, err
    }

    config := Config{}
    err := decoder.Decode(&config)

    return config, err
}
