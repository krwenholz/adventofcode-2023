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
	X, Y, Z    float64
	DX, DY, DZ float64
}

func (h *Hailstone) String() string {
	return fmt.Sprintf("Hailstone{x=%f, y=%f, z=%f, dx=%f, dy=%f, dz=%f}", h.X, h.Y, h.Z, h.DX, h.DY, h.DZ)
}

func writeXYHailstonesOctaveFile(hailstones []*Hailstone) (string, error) {
	tmpl, err := template.New("hailstones").Parse(`
lines = {
{{- range .}}
  struct('start', [{{.X}}, {{.Y}}], 'velocity', [{{.DX}}, {{.DY}}]),
{{- end}}
};
function intersecting_pairs = find_intersections(lines)
    n = length(lines);
    
    % Extract start points and velocities
    starts = cell2mat(cellfun(@(x) x.start(:), lines, 'UniformOutput', false));
    velocities = cell2mat(cellfun(@(x) x.velocity(:), lines, 'UniformOutput', false));
    
    % Create all possible pairs
    [i, j] = meshgrid(1:n, 1:n);
    pairs = [i(:), j(:)];
    pairs = pairs(pairs(:,1) < pairs(:,2), :);
    
    % Set up the system of equations for all pairs at once
    P1 = starts(:, pairs(:,1));
    P2 = starts(:, pairs(:,2));
    v1 = velocities(:, pairs(:,1));
    v2 = velocities(:, pairs(:,2));
    
    % Solve the system for all pairs
    A = cat(3, v1, -v2);
    b = P2 - P1;
    params = zeros(2, size(pairs, 1));
    
    for k = 1:size(pairs, 1)
        params(:,k) = A(:,:,k) \ b(:,k);
    end
    
    % Check for valid intersections
    valid = all(isfinite(params), 1);
    
    % Return intersecting pairs and their parameters
    intersecting_pairs = [pairs(valid,:), params(:,valid)'];
end

% Assuming 'lines' is your array of line structures
intersecting_pairs = find_intersections(lines);

% Display results
for k = 1:size(intersecting_pairs, 1)
    i = intersecting_pairs(k, 1);
    j = intersecting_pairs(k, 2);
    t = intersecting_pairs(k, 3);
    s = intersecting_pairs(k, 4);
    intersection_point = lines{i}.start + t * lines{i}.velocity;
    printf('(%d, %d) intersect at point: [%.2f, %.2f, %.2f]\n', ...
           i, j, intersection_point(1), intersection_point(2), intersection_point(3));
end

printf('Found %d intersections\n', size(intersecting_pairs, 1));
`)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, hailstones)
	if err != nil {
		return "", err
	}

	tempFile, err := os.CreateTemp("/tmp", "hailstones_*.m")
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

func partOne(puzzleFile string) {
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
		slog.Info("Read line", "line", line)
		parts := strings.Split(line, "@")
		coords := strings.Split(parts[0], ",")
		velocities := strings.Split(parts[1], ",")

		x, _ := strconv.ParseFloat(strings.TrimSpace(coords[0]), 32)
		y, _ := strconv.ParseFloat(strings.TrimSpace(coords[1]), 32)
		z, _ := strconv.ParseFloat(strings.TrimSpace(coords[2]), 32)

		dx, _ := strconv.ParseFloat(strings.TrimSpace(velocities[0]), 32)
		dy, _ := strconv.ParseFloat(strings.TrimSpace(velocities[1]), 32)
		dz, _ := strconv.ParseFloat(strings.TrimSpace(velocities[2]), 32)

		hs = append(hs, &Hailstone{X: x, Y: y, Z: z, DX: dx, DY: dy, DZ: dz})
	}

	if err := scanner.Err(); err != nil {
		slog.Error("Error reading file", "error", err)
	}

	tempFile, err := writeXYHailstonesOctaveFile(hs)
	if err != nil {
		slog.Error("Error writing to temp file", "error", err)
		return
	}
	slog.Debug("Temp file", "file", tempFile)

	output, err := exec.Command("octave", tempFile).Output()
	if err != nil {
		slog.Error("Error executing command", "error", err)
		return
	}
	intersectionLines := strings.Split(string(output), "\n")
	slog.Info("Octave output fetched", "intersections", len(intersectionLines))
	for _, l := range intersectionLines {
		slog.Info("Octave output", "line", l)
	}

	slog.Info("Finished Day TwentyFour part one", "expected", expected)
}

func partTwo(puzzleFile string) {
	slog.Info("Day TwentyFour part two", "puzzle file", puzzleFile)
}

var Cmd = &cobra.Command{
	Use: "dayTwentyFour",
	Run: func(cmd *cobra.Command, args []string) {
		puzzleInput, _ := cmd.Flags().GetString("puzzle-input")
		if !cmd.Flag("part-two").Changed {
			partOne(puzzleInput)
		} else {
			partTwo(puzzleInput)
		}
	},
}

func init() {
	Cmd.Flags().Bool("part-two", false, "Whether to run part two of the day's challenge")
}
