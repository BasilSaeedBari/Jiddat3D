package slicer

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"path/filepath"
	"strconv"
	"strings"
)

// Point3D represents a 3D coordinate in millimeters
type Point3D struct {
	X float64
	Y float64
	Z float64
}

// Triangle represents a 3D triangle facet
type Triangle struct {
	V1 Point3D
	V2 Point3D
	V3 Point3D
}

// MeshStats holds the calculated geometrical properties of the 3D model
type MeshStats struct {
	Filename       string  `json:"filename"`
	Format         string  `json:"format"`
	TriangleCount  int     `json:"triangle_count"`
	VertexCount    int     `json:"vertex_count"`
	MinX           float64 `json:"min_x"`
	MaxX           float64 `json:"max_x"`
	MinY           float64 `json:"min_y"`
	MaxY           float64 `json:"max_y"`
	MinZ           float64 `json:"min_z"`
	MaxZ           float64 `json:"max_z"`
	SizeX          float64 `json:"size_x"`          // width in mm
	SizeY          float64 `json:"size_y"`          // depth in mm
	SizeZ          float64 `json:"size_z"`          // height in mm
	SurfaceAreaCm2 float64 `json:"surface_area_cm2"` // in cm^2
	VolumeCm3      float64 `json:"volume_cm3"`      // solid volume in cm^3
}

// signedVolumeOfTriangle calculates signed volume of tetrahedron formed by triangle and origin
func signedVolumeOfTriangle(p1, p2, p3 Point3D) float64 {
	return (1.0 / 6.0) * (
		-p3.X*p2.Y*p1.Z + p2.X*p3.Y*p1.Z +
		p3.X*p1.Y*p2.Z - p1.X*p3.Y*p2.Z -
		p2.X*p1.Y*p3.Z + p1.X*p2.Y*p3.Z)
}

// areaOfTriangle calculates the surface area of a 3D triangle in mm^2
func areaOfTriangle(p1, p2, p3 Point3D) float64 {
	ax := p2.X - p1.X
	ay := p2.Y - p1.Y
	az := p2.Z - p1.Z

	bx := p3.X - p1.X
	by := p3.Y - p1.Y
	bz := p3.Z - p1.Z

	cx := ay*bz - az*by
	cy := az*bx - ax*bz
	cz := ax*by - ay*bx

	return 0.5 * math.Sqrt(cx*cx + cy*cy + cz*cz)
}

// ParseModel parses a 3D model file (.stl, .obj, .3mf) and extracts geometry
func ParseModel(filename string, r io.ReaderAt, size int64) (*MeshStats, error) {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".stl":
		return parseSTL(io.NewSectionReader(r, 0, size), size, filename)
	case ".obj":
		return parseOBJ(io.NewSectionReader(r, 0, size), filename)
	case ".3mf":
		return parse3MF(r, size, filename)
	default:
		return nil, fmt.Errorf("unsupported file extension '%s'. Please upload .3mf, .obj, or .stl", ext)
	}
}

func isFinitePoint(p Point3D) bool {
	return !math.IsNaN(p.X) && !math.IsInf(p.X, 0) &&
		!math.IsNaN(p.Y) && !math.IsInf(p.Y, 0) &&
		!math.IsNaN(p.Z) && !math.IsInf(p.Z, 0)
}

// parseSTL parses both binary and ASCII STL streams
func parseSTL(r io.Reader, size int64, filename string) (*MeshStats, error) {
	// Read first 84 bytes to inspect binary vs ascii
	header := make([]byte, 84)
	n, err := io.ReadFull(r, header)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return nil, fmt.Errorf("failed to read STL header: %w", err)
	}

	// A valid binary STL has exact byte length 84 + numTriangles*50.
	isBinary := false
	if n == 84 {
		numTriangles := uint64(binary.LittleEndian.Uint32(header[80:84]))
		if numTriangles <= 50000000 && uint64(size) == 84+numTriangles*50 {
			isBinary = true
		}
	}

	if !isBinary {
		trimmed := strings.TrimSpace(string(header[:n]))
		if strings.HasPrefix(strings.ToLower(trimmed), "solid") {
			combinedReader := io.MultiReader(bytes.NewReader(header[:n]), r)
			return parseAsciiSTL(combinedReader, filename)
		}
		return nil, fmt.Errorf("invalid or truncated STL file")
	}

	return parseBinarySTL(header, r, filename)
}

func parseBinarySTL(header []byte, r io.Reader, filename string) (*MeshStats, error) {
	if len(header) < 84 {
		return nil, fmt.Errorf("invalid binary STL file: truncated header")
	}

	numTriangles := int(binary.LittleEndian.Uint32(header[80:84]))
	if numTriangles < 0 || numTriangles > 50000000 {
		return nil, fmt.Errorf("invalid triangle count in STL: %d", numTriangles)
	}

	stats := &MeshStats{
		Filename:      filename,
		Format:        "Binary STL",
		TriangleCount: numTriangles,
		MinX: math.MaxFloat64, MinY: math.MaxFloat64, MinZ: math.MaxFloat64,
		MaxX: -math.MaxFloat64, MaxY: -math.MaxFloat64, MaxZ: -math.MaxFloat64,
	}

	var totalVol float64
	var totalArea float64
	br := bufio.NewReaderSize(r, 256*1024)
	buf := make([]byte, 50)

	for i := 0; i < numTriangles; i++ {
		if _, err := io.ReadFull(br, buf); err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				return nil, fmt.Errorf("truncated binary STL: expected %d triangles, read %d", numTriangles, i)
			}
			return nil, err
		}

		// 12 floats: Normal(3), V1(3), V2(3), V3(3)
		v1 := Point3D{
			X: float64(math.Float32frombits(binary.LittleEndian.Uint32(buf[12:16]))),
			Y: float64(math.Float32frombits(binary.LittleEndian.Uint32(buf[16:20]))),
			Z: float64(math.Float32frombits(binary.LittleEndian.Uint32(buf[20:24]))),
		}
		v2 := Point3D{
			X: float64(math.Float32frombits(binary.LittleEndian.Uint32(buf[24:28]))),
			Y: float64(math.Float32frombits(binary.LittleEndian.Uint32(buf[28:32]))),
			Z: float64(math.Float32frombits(binary.LittleEndian.Uint32(buf[32:36]))),
		}
		v3 := Point3D{
			X: float64(math.Float32frombits(binary.LittleEndian.Uint32(buf[36:40]))),
			Y: float64(math.Float32frombits(binary.LittleEndian.Uint32(buf[40:44]))),
			Z: float64(math.Float32frombits(binary.LittleEndian.Uint32(buf[44:48]))),
		}
		if !isFinitePoint(v1) || !isFinitePoint(v2) || !isFinitePoint(v3) {
			return nil, fmt.Errorf("invalid binary STL: non-finite vertex coordinate")
		}

		updateBounds(stats, v1)
		updateBounds(stats, v2)
		updateBounds(stats, v3)

		totalVol += signedVolumeOfTriangle(v1, v2, v3)
		totalArea += areaOfTriangle(v1, v2, v3)
	}

	finalizeMeshStats(stats, totalVol, totalArea)
	return stats, nil
}

func parseAsciiSTL(r io.Reader, filename string) (*MeshStats, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	stats := &MeshStats{
		Filename: filename,
		Format:   "ASCII STL",
		MinX:     math.MaxFloat64, MinY: math.MaxFloat64, MinZ: math.MaxFloat64,
		MaxX: -math.MaxFloat64, MaxY: -math.MaxFloat64, MaxZ: -math.MaxFloat64,
	}

	var vertices []Point3D
	var totalVol float64
	var totalArea float64
	triangleCount := 0

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 4 && strings.EqualFold(fields[0], "vertex") {
			x, errX := strconv.ParseFloat(fields[1], 64)
			y, errY := strconv.ParseFloat(fields[2], 64)
			z, errZ := strconv.ParseFloat(fields[3], 64)
			pt := Point3D{X: x, Y: y, Z: z}
			if errX != nil || errY != nil || errZ != nil || !isFinitePoint(pt) {
				return nil, fmt.Errorf("invalid ASCII STL vertex coordinates")
			}
			vertices = append(vertices, pt)
			updateBounds(stats, pt)

			if len(vertices) == 3 {
				totalVol += signedVolumeOfTriangle(vertices[0], vertices[1], vertices[2])
				totalArea += areaOfTriangle(vertices[0], vertices[1], vertices[2])
				triangleCount++
				vertices = vertices[:0]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read ASCII STL: %w", err)
	}

	stats.TriangleCount = triangleCount
	finalizeMeshStats(stats, totalVol, totalArea)
	return stats, nil
}

func parseOBJ(r io.Reader, filename string) (*MeshStats, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	stats := &MeshStats{
		Filename: filename,
		Format:   "Wavefront OBJ",
		MinX:     math.MaxFloat64, MinY: math.MaxFloat64, MinZ: math.MaxFloat64,
		MaxX: -math.MaxFloat64, MaxY: -math.MaxFloat64, MaxZ: -math.MaxFloat64,
	}

	var vertices []Point3D
	var totalVol float64
	var totalArea float64
	triangleCount := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) < 2 || line[0] == '#' {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		switch fields[0] {
		case "v":
			if len(fields) >= 4 {
				x, errX := strconv.ParseFloat(fields[1], 64)
				y, errY := strconv.ParseFloat(fields[2], 64)
				z, errZ := strconv.ParseFloat(fields[3], 64)
				pt := Point3D{X: x, Y: y, Z: z}
				if errX != nil || errY != nil || errZ != nil || !isFinitePoint(pt) {
					return nil, fmt.Errorf("invalid OBJ vertex coordinates")
				}
				vertices = append(vertices, pt)
				updateBounds(stats, pt)
			}
		case "f":
			// Faces can be polygon. Triangulate into fan: v0, v1, v2; v0, v2, v3; etc.
			faceIndices := make([]int, 0, len(fields)-1)
			for i := 1; i < len(fields); i++ {
				parts := strings.Split(fields[i], "/")
				idx, err := strconv.Atoi(parts[0])
				if err != nil {
					return nil, fmt.Errorf("invalid OBJ face vertex index '%s'", parts[0])
				}
				// OBJ 1-based indexing, or negative for relative
				var resolved int
				if idx > 0 {
					resolved = idx - 1
				} else if idx < 0 {
					resolved = len(vertices) + idx
				} else {
					return nil, fmt.Errorf("invalid OBJ face vertex index 0")
				}
				if resolved < 0 || resolved >= len(vertices) {
					return nil, fmt.Errorf("OBJ face vertex index %d out of range", idx)
				}
				faceIndices = append(faceIndices, resolved)
			}

			if len(faceIndices) >= 3 {
				v0 := vertices[faceIndices[0]]
				for i := 1; i < len(faceIndices)-1; i++ {
					v1 := vertices[faceIndices[i]]
					v2 := vertices[faceIndices[i+1]]
					totalVol += signedVolumeOfTriangle(v0, v1, v2)
					totalArea += areaOfTriangle(v0, v1, v2)
					triangleCount++
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read OBJ: %w", err)
	}

	stats.VertexCount = len(vertices)
	stats.TriangleCount = triangleCount
	finalizeMeshStats(stats, totalVol, totalArea)
	return stats, nil
}

// 3MF XML Structs
type Model3MF struct {
	XMLName   xml.Name    `xml:"model"`
	Resources Resources3MF `xml:"resources"`
}

type Resources3MF struct {
	Objects []Object3MF `xml:"object"`
}

type Object3MF struct {
	ID   string  `xml:"id,attr"`
	Type string  `xml:"type,attr"`
	Mesh Mesh3MF `xml:"mesh"`
}

type Mesh3MF struct {
	Vertices  Vertices3MF  `xml:"vertices"`
	Triangles Triangles3MF `xml:"triangles"`
}

type Vertices3MF struct {
	Vertices []Vertex3MF `xml:"vertex"`
}

type Vertex3MF struct {
	X float64 `xml:"x,attr"`
	Y float64 `xml:"y,attr"`
	Z float64 `xml:"z,attr"`
}

type Triangles3MF struct {
	Triangles []Triangle3MF `xml:"triangle"`
}

type Triangle3MF struct {
	V1 int `xml:"v1,attr"`
	V2 int `xml:"v2,attr"`
	V3 int `xml:"v3,attr"`
}

const max3MFModelBytes = 64 << 20 // 64 MB uncompressed XML cap against zip bombs

func parse3MF(r io.ReaderAt, size int64, filename string) (*MeshStats, error) {
	zipReader, err := zip.NewReader(r, size)
	if err != nil {
		return nil, fmt.Errorf("failed to open 3MF zip container: %w", err)
	}

	var modelFile *zip.File
	for _, f := range zipReader.File {
		if strings.HasSuffix(strings.ToLower(f.Name), ".model") {
			modelFile = f
			break
		}
	}

	if modelFile == nil {
		return nil, fmt.Errorf("invalid 3MF: missing .model XML geometry file")
	}
	if modelFile.UncompressedSize64 > max3MFModelBytes {
		return nil, fmt.Errorf("3MF model XML exceeds maximum uncompressed size of 64MB")
	}

	rc, err := modelFile.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to read 3MF model XML: %w", err)
	}
	defer rc.Close()

	limited := &io.LimitedReader{R: rc, N: max3MFModelBytes + 1}
	var model Model3MF
	decoder := xml.NewDecoder(limited)
	if err := decoder.Decode(&model); err != nil {
		if limited.N == 0 {
			return nil, fmt.Errorf("3MF model XML exceeds maximum uncompressed size of 64MB")
		}
		return nil, fmt.Errorf("failed to parse 3MF XML structure: %w", err)
	}
	if limited.N == 0 {
		return nil, fmt.Errorf("3MF model XML exceeds maximum uncompressed size of 64MB")
	}

	stats := &MeshStats{
		Filename: filename,
		Format:   "3D Manufacturing Format (3MF)",
		MinX:     math.MaxFloat64, MinY: math.MaxFloat64, MinZ: math.MaxFloat64,
		MaxX: -math.MaxFloat64, MaxY: -math.MaxFloat64, MaxZ: -math.MaxFloat64,
	}

	var totalVol float64
	var totalArea float64
	triangleCount := 0
	vertexCount := 0

	for _, obj := range model.Resources.Objects {
		verts := obj.Mesh.Vertices.Vertices
		vertexCount += len(verts)

		for _, v := range verts {
			pt := Point3D{X: v.X, Y: v.Y, Z: v.Z}
			if !isFinitePoint(pt) {
				return nil, fmt.Errorf("invalid 3MF: non-finite vertex coordinate")
			}
			updateBounds(stats, pt)
		}

		for _, t := range obj.Mesh.Triangles.Triangles {
			if t.V1 < 0 || t.V2 < 0 || t.V3 < 0 || t.V1 >= len(verts) || t.V2 >= len(verts) || t.V3 >= len(verts) {
				return nil, fmt.Errorf("invalid 3MF: triangle vertex index out of range")
			}
			p1 := Point3D{X: verts[t.V1].X, Y: verts[t.V1].Y, Z: verts[t.V1].Z}
			p2 := Point3D{X: verts[t.V2].X, Y: verts[t.V2].Y, Z: verts[t.V2].Z}
			p3 := Point3D{X: verts[t.V3].X, Y: verts[t.V3].Y, Z: verts[t.V3].Z}

			totalVol += signedVolumeOfTriangle(p1, p2, p3)
			totalArea += areaOfTriangle(p1, p2, p3)
			triangleCount++
		}
	}

	stats.VertexCount = vertexCount
	stats.TriangleCount = triangleCount
	finalizeMeshStats(stats, totalVol, totalArea)
	return stats, nil
}

func updateBounds(stats *MeshStats, p Point3D) {
	if p.X < stats.MinX { stats.MinX = p.X }
	if p.X > stats.MaxX { stats.MaxX = p.X }
	if p.Y < stats.MinY { stats.MinY = p.Y }
	if p.Y > stats.MaxY { stats.MaxY = p.Y }
	if p.Z < stats.MinZ { stats.MinZ = p.Z }
	if p.Z > stats.MaxZ { stats.MaxZ = p.Z }
}

func finalizeMeshStats(stats *MeshStats, rawVol float64, rawArea float64) {
	if stats.TriangleCount == 0 {
		stats.MinX, stats.MinY, stats.MinZ = 0, 0, 0
		stats.MaxX, stats.MaxY, stats.MaxZ = 0, 0, 0
	}

	stats.SizeX = math.Max(0, stats.MaxX-stats.MinX)
	stats.SizeY = math.Max(0, stats.MaxY-stats.MinY)
	stats.SizeZ = math.Max(0, stats.MaxZ-stats.MinZ)

	// Volume in mm^3 converted to cm^3 (1 cm^3 = 1000 mm^3)
	stats.VolumeCm3 = math.Abs(rawVol) / 1000.0
	// Surface area in mm^2 converted to cm^2 (1 cm^2 = 100 mm^2)
	stats.SurfaceAreaCm2 = rawArea / 100.0

	// Sanity safeguard: if mesh is non-manifold or planar (near-zero volume relative to its bbox),
	// approximate minimum volume from bounding box without corrupting tiny/meter-exported watertight meshes.
	if stats.SizeX > 0 && stats.SizeY > 0 && stats.SizeZ > 0 {
		bboxVolCm3 := (stats.SizeX * stats.SizeY * stats.SizeZ) / 1000.0
		if stats.VolumeCm3 < bboxVolCm3*0.005 {
			stats.VolumeCm3 = bboxVolCm3 * 0.35 // 35% bounding envelope approximation
		}
	}
}
