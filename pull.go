// auth token -> json with layers -> pull layers -> save to disk
package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Layer struct {
	Digest string `json:"digest"`
}

type Manifest struct {
	Layers []Layer `json:"layers"`
}

type ManifestEntry struct {
	Digest   string `json:"digest"`
	Platform struct {
		Architecture string `json:"architecture"`
		OS           string `json:"os"`
	} `json:"platform"`
}

type ManifestList struct {
	Manifests []ManifestEntry `json:"manifests"`
}

func getAuthToken(imageName string) (string, error) {
	resp, err := http.Get("https://auth.docker.io/token?service=registry.docker.io&scope=repository:library/" + imageName + ":pull")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var result struct {
		Token string `json:"token"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}
	return result.Token, nil
}

func getImageLayers(imageName, token string) ([]string, error) {
	req, err := http.NewRequest("GET", "https://registry-1.docker.io/v2/library/"+imageName+"/manifests/latest", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var manifest Manifest
	json.Unmarshal(body, &manifest)

	// if there is no "layers" field, it might be a manifest list (multi-arch image) ua
	if len(manifest.Layers) == 0 {
		var mList ManifestList
		json.Unmarshal(body, &mList)

		var digest string
		for _, m := range mList.Manifests {
			if m.Platform.Architecture == "amd64" && m.Platform.OS == "linux" {
				digest = m.Digest
				break
			}
		}
		if digest == "" {
			return nil, fmt.Errorf("no linux/amd64 manifest found")
		}

		req2, err := http.NewRequest("GET", "https://registry-1.docker.io/v2/library/"+imageName+"/manifests/"+digest, nil)
		if err != nil {
			return nil, err
		}
		req2.Header.Set("Authorization", "Bearer "+token)
		req2.Header.Set("Accept", "application/vnd.oci.image.manifest.v1+json")

		resp2, err := http.DefaultClient.Do(req2)
		if err != nil {
			return nil, err
		}
		defer resp2.Body.Close()
		body, _ = io.ReadAll(resp2.Body)
		json.Unmarshal(body, &manifest)
	}

	var layers []string
	for _, layer := range manifest.Layers {
		layers = append(layers, layer.Digest)
	}
	return layers, nil
}

// predpostavicu da radi ovaj kod
func downloadLayer(imageName, layer, token, dest string) error {
	req, err := http.NewRequest("GET", "https://registry-1.docker.io/v2/library/"+imageName+"/blobs/"+layer, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(dest, hdr.Name)

		switch hdr.Typeflag {
		case tar.TypeDir:
			os.MkdirAll(target, os.FileMode(hdr.Mode))

		case tar.TypeReg:
			os.MkdirAll(filepath.Dir(target), 0755)
			f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			io.Copy(f, tr)
			f.Close()

		case tar.TypeSymlink:
			os.Remove(target)
			os.Symlink(hdr.Linkname, target)
		}
	}
	return nil
}
func pullImage(imageName string) error {
	dest := filepath.Join(imagesDir, imageName, "rootfs")
	os.MkdirAll(dest, 0755)

	token, err := getAuthToken(imageName)
	if err != nil {
		return err
	}
	fmt.Println("Token:", token[:20]+"...")

	layers, err := getImageLayers(imageName, token)
	if err != nil {
		return err
	}
	fmt.Printf("Found %d layers\n", len(layers))
	for i, l := range layers {
		fmt.Printf("  Layer %d: %s\n", i, l[:30])
	}

	for _, layer := range layers {
		fmt.Printf("Downloading %s...\n", layer[:30])
		err = downloadLayer(imageName, layer, token, dest)
		if err != nil {
			return err
		}
	}
	return nil
}
