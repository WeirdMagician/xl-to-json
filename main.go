package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/tealeg/xlsx"
)

const (
	sheetno = 3
	row     = 3
	col     = 1
)

type sblock struct {
	Solution string
	Priority int
}

type record struct {
	CaseNumber, Product, ProductVersion, ProductFeature, AppsFeature string
	CaseComments                                                     []*sblock
}

var list []record

var tagdata = make(map[string][]string)

func init() {
	err := yaml.Unmarshal(func() ([]byte, interface{}) {
		b, err := ioutil.ReadFile("Tags.yml")
		if err != nil {
			log.Fatal("Error while reading the input YML file", err)
		}
		return b, tagdata
	}())
	if err != nil {
		log.Fatal("Error while parsing the input YML file", err)
	}
	fmt.Println(tagdata)
}

func main() {
	myslice, err := xlsx.FileToSlice("input.xlsx")
	if err != nil {
		log.Fatalln("Unable to Convert the input file to 3-dimentional slice", err)
	}
	fmt.Println(myslice[3][2])

	var pitem *record
	for r, v := range myslice[sheetno] {
		if r >= row {
			if strings.TrimSpace(v[0]) == "" {
				insert(pitem, v)
			} else {
				if pitem != nil {
					list = append(list, *pitem)
				}
				item := record{CaseNumber: v[0], Product: v[1], ProductVersion: v[2], ProductFeature: v[3], AppsFeature: v[4]}
				set := extract(v[5])
				sb := sblock{Priority: 1, Solution: set["#Solution#"]}
				item.CaseComments = append(item.CaseComments, &sb)
				pitem = &item
			}
		}
	}
	b, err := json.MarshalIndent(list, "", "\t")
	if err != nil {
		log.Fatalln("Json Conversion failed", err)
	}
	ioutil.WriteFile("output.json", b, os.FileMode(0777))

}

func insert(item *record, v []string) {
	set := extract(v[5])
	p := item.CaseComments[len(item.CaseComments)-1].Priority
	sb := sblock{Priority: p + 1, Solution: set["#Solution#"]}
	item.CaseComments = append(item.CaseComments, &sb)
}

func extract(s string) map[string]string {
	ls := strings.ToLower(s)
	var point = make(map[int]string)
	for _, v := range tagdata["rca"] {
		index := strings.Index(ls, strings.ToLower(v))
		if index != -1 {
			point[index] = v
		}
	}

	var keys []int
	for k := range point {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	var set = make(map[string]string)
	for index, key := range keys {
		matter := data(point, keys, s, index, key)
		set[point[key]] = matter
	}
	// for k, v := range set {
	// 	fmt.Println(k, "\t", v)
	// }
	return set
	// begin := strings.Index(strings.ToLower(s), "#solution#")
	// fmt.Println(begin)
}

func data(point map[int]string, keys []int, s string, index int, key int) string {
	start := key + len(point[key])
	stop := end(start, keys, index, s)

	fmt.Println(len(s), start, stop, point[key])
	matter := s[start:stop]
	if strings.TrimSpace(matter) == "" {
		if len(keys)-1 >= index+1 {
			return data(point, keys, s, index+1, keys[index+1])
		}
	}
	return matter
}

func end(start int, keys []int, index int, s string) int {
	if len(keys)-1 >= index+1 {
		if start >= keys[index+1] {
			return end(start, keys, index+1, s)
		}
		return keys[index+1]
	}
	return len(s)
}
