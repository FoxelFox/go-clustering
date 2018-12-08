package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/image/math/f32"
)

type DataStructure struct {
	id       int32
	edges    []Edge
	position f32.Vec3
}

type Edge struct {
	conc  string
	force float32
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

	m := make(map[string]DataStructure)
	var index int32
	var maxForce float32
	var maxEdges int

	t := time.Now()
	elapsed1 := t.Sub(start)

	// create map
	for _, line := range lines {
		v := strings.Split(line, "|")
		prem := DataStructure{}
		if value, ok := m[v[0]]; ok {
			prem = value
		} else {
			index++
			m[v[0]] = DataStructure{id: index, edges: []Edge{}}
			prem = m[v[0]]
		}

		value, _ := strconv.ParseFloat(v[2], 32)

		force := float32(value)
		if maxForce < force && force < 500 {
			maxForce = force
		}

		prem.edges = append(prem.edges, Edge{conc: v[1], force: force})

		if maxEdges < len(prem.edges) {
			maxEdges++
		}
	}

	// create dataset
	t = time.Now()
	elapsed2 := t.Sub(start)
	fmt.Println(elapsed1)
	fmt.Println(elapsed2)

	fmt.Print("Max Force: ")
	fmt.Printf("%d\n", maxEdges)
}

func cluster(data DataStructure) {
	// vec4 o = texture(image, v_texCoord);
	// vec3 position = vec3(o.x, o.y, o.z);
	// vec3 velocity = vec3(0.0, 0.0, 0.0);

	// if (o.w > 0.5) {
	//     if (forceActive < 0.5) {
	//         for (float x = 0.0; x < maxEdges; x++) {

	//             vec3 ref = texelFetch(edges, ivec3(x, v_texCoord.x * size,v_texCoord.y * size), 0).xyz;

	//             if (ref.x >= 0.0) {
	//                 vec3 p = texelFetch(image, ivec2(ref.x, ref.y), 0).xyz;

	//                 if (!(p.x == position.x && p.y == position.y)) {
	//                     float radius = distance(p, position);
	//                     float force = ref.z * radius * radius;
	//                     velocity -= normalize(position - p) * clamp(0.0125 * force, 0.0, 0.25);
	//                 }
	//             }

	//         }

	// 			for (float x = 0.0; x < size; x++) {
	// 			for (float y = 0.0; y < size; y++) {

	// 				vec4 pp = texelFetch(image, ivec2(x, y), 0);

	// 				if(pp.w > 0.55) {
	// 					vec3 p = pp.xyz;
	// 					if (!(p.x == position.x && p.y == position.y)) {
	// 						float radius = distance(p, position);
	// 						float force = 16.0 / (radius * radius);
	// 						velocity += normalize(position - p) * clamp(0.0001 * force, 0.0, 0.0025);

	// 					}
	// 				}

	// 			}
	// 		}

	// 		position += velocity;
	// 	}
	// }

	// outColor = vec4(position, o.w);
}
