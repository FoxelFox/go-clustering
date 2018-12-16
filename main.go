package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/jinzhu/copier"
)

type dataStructure struct {
	id       int32
	edges    []Edge
	position mgl32.Vec3
}

type Edge struct {
	attraction string
	force      float32
}

// 974610|242760|1
// 974610|599140|1
// 974610|496300|1
// 974610|457140|1
// 974610|227300|1

func main() {

	start := time.Now()

	file := os.Args[1]
	data, _ := ioutil.ReadFile(file)

	lines := strings.Split(string(data), "\n")[1:]

	m := make(map[string]*dataStructure)
	var index int32
	var maxForce float32
	var maxEdges int

	t := time.Now()

	// create map
	for _, line := range lines {
		v := strings.Split(line, "|")

		var prem *dataStructure
		var ok bool
		if prem, ok = m[v[0]]; !ok {
			index++

			randomVector := mgl32.Vec3{
				rand.Float32()*2 - 1.0,
				rand.Float32()*2 - 1.0,
				rand.Float32()*2 - 1.0,
			}.Normalize().Mul(10)

			prem = &dataStructure{id: index, edges: []Edge{}, position: randomVector}
			m[v[0]] = prem
		}

		value, _ := strconv.ParseFloat(v[2], 32)

		force := float32(value)
		if maxForce < force && force < 500 {
			maxForce = force
		}

		prem.edges = append(prem.edges, Edge{attraction: v[1], force: force})

		if maxEdges < len(prem.edges) {
			maxEdges++
		}
	}

	// buffer 2

	m2 := make(map[string]*dataStructure)
	copier.Copy(&m2, &m)
	workersDone := make(chan bool)

	for i := 0; i <= 50; i += 2 {
		for k := range m {
			go funkyCluster(m, m2, k, workersDone)
		}

		for i := 0; i < len(m); i++ {
			<-workersDone
		}

		writeToFile(m)
		fmt.Println(i)

		for k := range m {
			go funkyCluster(m2, m, k, workersDone)
		}

		for i := 0; i < len(m); i++ {
			<-workersDone
		}

		writeToFile(m2)
		fmt.Println(i + 1)
	}

	t = time.Now()
	elapsed := t.Sub(start)
	fmt.Println(elapsed)
	fmt.Println()

}

func funkyCluster(wData map[string]*dataStructure, rData map[string]*dataStructure, id string, done chan bool) {

	self := rData[id]

	var velocity mgl32.Vec3
	dist := float32(0.0001)

	for _, edge := range rData[id].edges {
		attraction := rData[edge.attraction]
		if attraction == nil {
			continue
		}
		if self.id != attraction.id {
			distance := self.position.Sub(attraction.position).Len()

			if distance > dist {
				velocity = velocity.Sub(self.position.Sub(attraction.position).Normalize().Mul(mgl32.Clamp((distance-dist), 0.0, 100.0) * 0.5))
			}

		}
	}

	// for _, repulsion := range rData {
	// 	if self.id != repulsion.id {
	// 		distance := self.position.Sub(repulsion.position).Len()

	// 		if distance <= dist {
	// 			velocity = velocity.Add(self.position.Sub(repulsion.position).Normalize().Mul(mgl32.Clamp((dist-distance), 0.0, 10.0) * 0.5))
	// 		}
	// 	}
	// }

	self.position = self.position.Add(velocity)

	done <- true

}

func writeToFile(data map[string]*dataStructure) {
	f, _ := os.Create("test.bin")
	buf := new(bytes.Buffer)

	for _, item := range data {
		binary.Write(buf, binary.LittleEndian, item.position.X())
		binary.Write(buf, binary.LittleEndian, item.position.Y())
		binary.Write(buf, binary.LittleEndian, item.position.Z())
		binary.Write(buf, binary.LittleEndian, float32(1.0))

	}

	//fmt.Println(float32(1.0))
	f.Write(buf.Bytes())

	f.Close()
}
