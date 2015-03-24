package main

import "github.com/savaki/go.hue" 
import "time"
import "encoding/json"
import "os"
import "strings"
import "strconv"
import "log"

var Colors map[string]Color;

type Color struct {
    Hue string
    Sat string
    Bri string
}



func main() {
    Colors = make(map[string]Color)
    Colors["green"] = Color{"120", "100", "100"}
    Colors["red"] = Color{"0", "100", "100"}
    Colors["magenta"] = Color{"300", "100", "100"}
    Colors["purple"] = Color{"280", "100", "100"}
    Colors["yellow"] = Color{"60", "100", "100"}

    for {
        rerun()
        time.Sleep(1000 * 60 * time.Millisecond)
    }
}

func rerun(){
    r, err := os.Open("hue.json")
    if err != nil {
        log.Fatal(err)
    }
    reader := json.NewDecoder(r);
    var data map[string]map[string]map[string]string;
    if err := reader.Decode(&data); err != nil {
        log.Fatal("Error with decoding data: %s\n", err)
    }
    hour, min, sec := time.Now().Clock()
    
    val := sec + (min * 60) + (hour * 3600)
    for k, v := range data {
         max := val
         min := 0
         var maxdata map[string]string
         var mindata map[string]string
         for l, w := range v {
            __ := strings.Split(l, ":")
            hr, err := strconv.Atoi(__[0])
            if err != nil {
                log.Fatal(err)
            }
            mn, err := strconv.Atoi(__[1])
            if err != nil {
                log.Fatal(err)
            }
            sum := (hr * 3600) + (mn * 60)
            if sum >= min && sum <= val {
                min = sum
                mindata = w
            }
            if sum >= max {
                max = sum
                maxdata = w
            }
        }
        var set map[string]string
        if mindata != nil {
            set = mindata
        } else if maxdata != nil {
            set = maxdata
        }
        light(k, set)
    }
}

func light(name string, data map[string]string) {
    var ip = "0.1.2.3"
    var username = "puffrfish"
    brg := hue.NewBridge(ip, username)
    light, err := brg.FindLightByName(name)
    if err != nil {
        log.Fatal(err)
    }
    var color Color = Colors[data["color"]]
    state := hue.SetLightState{On:"true", Alert:data["blink"], Hue:color.Hue, Sat:color.Sat, Bri:color.Bri}
    light.SetState(state)
}

