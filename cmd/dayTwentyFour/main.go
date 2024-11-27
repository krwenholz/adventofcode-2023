package dayTwentyFour

import (
	"bufio"
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

/*
It's a linear algebra problem! I'll want to
- remove z
- build the system of linear equations
- construct an all pairs version of the matrices (probably a slick way to do this?)
- inspect the intersection points for being between the desired x and y values

19, 13, 30 @ -2,  1, -2
x = 19t - 2 -> t = (x+2)/19
y = 13t + 1 -> t = (y-1)/13
z = 30t - 2 -> t = (z+2)/30
3t = (x+2)/19 + (y-1)/13 + (z+2)/30

18, 19, 22 @ -1, -1, -2
3t = (x+1)/18 + (y+1)/19 + (z+2)/22
*/

type Hailstone struct {
	X, Y, Z    int64
	DX, DY, DZ int64
}

func (h *Hailstone) String() string {
	return fmt.Sprintf("Hailstone{x=%d, y=%d, z=%d, dx=%d, dy=%d, dz=%d}", h.X, h.Y, h.Z, h.DX, h.DY, h.DZ)
}

func writeXYHailstonesPythonFile(hailstones []*Hailstone, testAreaStart, testAreaEnd int) (string, error) {
	tmpl, err := template.New("hailstones").Parse(`
from sympy import Point, Line

lines = [
{{- range .Hs}}
  (Point({{.X}}, {{.Y}}), Point({{.X}} + ({{.DX}}), {{.Y}} + ({{.DY}}))),
{{- end}}
]

area_start = {{.TestAreaStart}}
area_end = {{.TestAreaEnd}}

def is_future(p1, p2, intersection):
	v = p2 - p1
	if v.x > 0 and intersection.x < p1.x:
		return False
	if v.x < 0 and intersection.x > p1.x:
		return False
	if v.y > 0 and intersection.y < p1.y:
		return False
	if v.y < 0 and intersection.y > p1.y:
		return False
	return True

def is_between(p, start, end):
	return p.x >= start and p.x <= end and p.y >= start and p.y <= end

intersections = 0
for i in range(len(lines)):
	pi1, pi2 = lines[i]
	line1 = Line(pi1, pi2)
	for j in range(i+1, len(lines)):
		pj1, pj2 = lines[j]
		line2 = Line(pj1, pj2)

		# Find the intersection
		intersection = line1.intersection(line2)

		if intersection and is_future(pi1, pi2, intersection[0]) and is_future(pj1, pj2, intersection[0]):
			if is_between(intersection[0], area_start, area_end):
				print(f"Lines ({pi1},{pj1}) intersect at: {intersection[0]}")
				intersections += 1
			else:
				print(f"Lines ({pi1},{pj1}) intersect at: {intersection[0]} but outside area")
		else:
			print(f"Lines ({pi1}, {pj2}) do not intersect")
	
print(intersections)
`)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, struct {
		Hs            []*Hailstone
		TestAreaStart int
		TestAreaEnd   int
	}{
		Hs:            hailstones,
		TestAreaStart: testAreaStart,
		TestAreaEnd:   testAreaEnd,
	})
	if err != nil {
		return "", err
	}

	tempFile, err := os.CreateTemp("/tmp", "hailstones_*.py")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	_, err = tempFile.Write(buf.Bytes())
	if err != nil {
		return "", err
	}

	return tempFile.Name(), nil
}

func partOne(puzzleFile string, testAreaStart, testAreaEnd int) {
	slog.Info("Day TwentyFour part one", "puzzle file", puzzleFile)
	file, err := os.Open(puzzleFile)
	if err != nil {
		slog.Error("Error opening file", "error", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	expected := scanner.Text()

	hs := []*Hailstone{}
	for scanner.Scan() {
		line := scanner.Text()
		// Process the line here
		parts := strings.Split(line, "@")
		coords := strings.Split(parts[0], ",")
		velocities := strings.Split(parts[1], ",")

		x, _ := strconv.ParseInt(strings.TrimSpace(coords[0]), 10, 64)
		y, _ := strconv.ParseInt(strings.TrimSpace(coords[1]), 10, 64)
		z, _ := strconv.ParseInt(strings.TrimSpace(coords[2]), 10, 64)

		dx, _ := strconv.ParseInt(strings.TrimSpace(velocities[0]), 10, 64)
		dy, _ := strconv.ParseInt(strings.TrimSpace(velocities[1]), 10, 64)
		dz, _ := strconv.ParseInt(strings.TrimSpace(velocities[2]), 10, 64)

		hs = append(hs, &Hailstone{X: x, Y: y, Z: z, DX: dx, DY: dy, DZ: dz})
	}

	if err := scanner.Err(); err != nil {
		slog.Error("Error reading file", "error", err)
	}

	slog.Debug("hs", "hs", hs)

	tempFile, err := writeXYHailstonesPythonFile(hs, testAreaStart, testAreaEnd)
	if err != nil {
		slog.Error("Error writing to temp file", "error", err)
		return
	}

	slog.Info("Processing with Python", "file", tempFile)
	output, err := exec.Command("python", tempFile).Output()
	if err != nil {
		slog.Error("Error executing command", "error", err)
		return
	}
	pout := strings.Split(string(output), "\n")
	for _, l := range pout {
		slog.Debug("pout", "line", l)
	}

	intersections, _ := strconv.ParseInt(pout[len(pout)-2], 10, 32)

	slog.Info("Finished Day TwentyFour part one", "intersections", intersections, "expected", expected)
}

func partTwo(puzzleFile string) {
	slog.Info("Day TwentyFour part two", "puzzle file", puzzleFile)
}

var Cmd = &cobra.Command{
	Use: "dayTwentyFour",
	Run: func(cmd *cobra.Command, args []string) {
		puzzleInput, _ := cmd.Flags().GetString("puzzle-input")
		testAreaStart, _ := cmd.Flags().GetInt("test-area-start")
		testAreaEnd, _ := cmd.Flags().GetInt("test-area-end")
		if !cmd.Flag("part-two").Changed {
			partOne(puzzleInput, testAreaStart, testAreaEnd)
		} else {
			partTwo(puzzleInput)
		}
	},
}

func init() {
	Cmd.Flags().Bool("part-two", false, "Whether to run part two of the day's challenge")
	Cmd.Flags().Int("test-area-start", 7, "")
	Cmd.Flags().Int("test-area-end", 27, "")
}
