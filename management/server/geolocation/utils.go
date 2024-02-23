package geolocation

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

// decompressTarGzFile decompresses a .tar.gz file.
func decompressTarGzFile(filepath, destDir string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			err := os.MkdirAll(destDir+"/"+header.Name, 0755)
			if err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err := os.Create(destDir + "/" + header.Name)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}

	return nil
}

// decompressZipFile decompresses a .zip file.
func decompressZipFile(filepath, destDir string) error {
	r, err := zip.OpenReader(filepath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := destDir + "/" + f.Name

		if f.FileInfo().IsDir() {
			err = os.MkdirAll(fpath, os.ModePerm)
			if err != nil {
				return err
			}
			continue
		} else {
			dir, _ := path.Split(fpath)
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return err
			}
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

// calculateSHA256 calculates the SHA256 checksum of a file.
func calculateSHA256(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

// loadChecksumFromFile loads the first checksum from a file.
func loadChecksumFromFile(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) > 0 {
			return parts[0], nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", nil
}

// verifyChecksum compares the calculated SHA256 checksum of a file against the expected checksum.
func verifyChecksum(filepath, expectedChecksum string) error {
	calculatedChecksum, err := calculateSHA256(filepath)
	if err != nil {
		return err
	}

	if calculatedChecksum != expectedChecksum {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedChecksum, calculatedChecksum)
	}

	return nil
}

// downloadFile downloads a file from a URL and saves it to a local file path.
func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
