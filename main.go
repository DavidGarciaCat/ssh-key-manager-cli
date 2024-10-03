package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const sshFolder = ".ssh"

func main() {
	homeDir, _ := os.UserHomeDir()
	sshDir := filepath.Join(homeDir, sshFolder)

	fmt.Println("\nSSH Key-Pair Manager CLI")
	fmt.Println("------------------------\n")

	currentKeyPair := getCurrentKeyPair(sshDir)

	if currentKeyPair == "" {
		fmt.Println("No matching key-pair is currently active.")
	} else {
		fmt.Printf("Current active key-pair: %s\n", currentKeyPair)
	}

	for {
		fmt.Println("\n1) Switch SSH key-pair for another system")
		fmt.Println("2) Generate new SSH key-pair")
		fmt.Println("3) Quit")

		choice := readInput("\nChoose an option: ")

		switch choice {
		case "1":
			switchKeyPair()
		case "2":
			generateNewKeyPair()
		case "3":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid option. Try again.")
		}
	}
}

func getCurrentKeyPair(sshDir string) string {
	// Step 1: Define a list of file names
	existingFiles := []string{"id_docker", "id_ed25519", "id_rsa"}

	// Step 2: Get a list of subfolders
	for _, folder := range getSubfolders(sshDir) {
		subDir := filepath.Join(sshDir, folder)

		// Step 3: Check for files in each subfolder
		for _, file := range existingFiles {
			currentPrivateKeyPath := filepath.Join(sshDir, file)
			currentPublicKeyPath := filepath.Join(sshDir, file+".pub")

			subPrivateKeyPath := filepath.Join(subDir, file)
			subPublicKeyPath := filepath.Join(subDir, file+".pub")

			// Step 4: If both private and public keys exist in both locations
			if fileExists(currentPrivateKeyPath) && fileExists(subPrivateKeyPath) &&
				fileExists(currentPublicKeyPath) && fileExists(subPublicKeyPath) {

				// Compare contents
				if compareFiles(currentPrivateKeyPath, subPrivateKeyPath) &&
					compareFiles(currentPublicKeyPath, subPublicKeyPath) {
					return folder // Found a match, return this folder
				}
			}
		}
	}

	// Step 5: No matches found, return <unknown>
	return "<unknown>"
}

// Helper function to check if a file exists
func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// Helper function to compare two files
func compareFiles(file1, file2 string) bool {
	content1, err1 := ioutil.ReadFile(file1)
	content2, err2 := ioutil.ReadFile(file2)

	if err1 != nil || err2 != nil {
		return false
	}
	return string(content1) == string(content2)
}

func listAvailableKeys() {
	homeDir, _ := os.UserHomeDir()
	sshDir := filepath.Join(homeDir, sshFolder)

	subfolders := getSubfolders(sshDir)

	fmt.Println("\nAvailable SSH key-pair folders:\n")
	for i, folder := range subfolders {
		fmt.Printf("%d) %s\n", i+1, folder)
	}
}

func getSubfolders(sshDir string) []string {
	files, err := ioutil.ReadDir(sshDir)
	if err != nil {
		fmt.Println("Failed to read SSH directory:", err)
		return nil
	}

	var subfolders []string
	for _, file := range files {
		if file.IsDir() && strings.HasPrefix(file.Name(), "_") {
			subfolders = append(subfolders, file.Name())
		}
	}
	return subfolders
}

func switchKeyPair() {
	homeDir, _ := os.UserHomeDir()
	sshDir := filepath.Join(homeDir, sshFolder)

	subfolders := getSubfolders(sshDir)

	// List available key-pairs with numeric prefix
	listAvailableKeys()

	// Get the user's choice by number
	choiceStr := readInput("\nEnter the number of the folder to switch to: ")
	choice, err := strconv.Atoi(choiceStr)
	if err != nil || choice < 1 || choice > len(subfolders) {
		fmt.Println("Invalid choice.")
		return
	}

	folder := subfolders[choice-1]

	// Verify the folder exists
	subDir := filepath.Join(sshDir, folder)
	if _, err := os.Stat(subDir); os.IsNotExist(err) {
		fmt.Println("\nFolder does not exist.\n")
		return
	}

	fmt.Println("\nRemoving SSH key-pair files...\n")

	// Remove existing key-pair files
	existingKeys := []string{"id_docker", "id_docker.pub", "id_ed25519", "id_ed25519.pub", "id_rsa", "id_rsa.pub"}
	for _, file := range existingKeys {
		filePath := filepath.Join(sshDir, file)
		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			os.Remove(filePath)
			fmt.Printf("Removed: %s\n", filePath)
		}
	}

	// Copy new key-pairs
	files, err := ioutil.ReadDir(subDir)
	if err != nil {
		fmt.Println("Failed to read folder:", err)
		return
	}

	fmt.Println("\nCopying SSH key-pair files...\n")

	for _, file := range files {
		srcFile := filepath.Join(subDir, file.Name())
		destFile := filepath.Join(sshDir, file.Name())
		copyFile(srcFile, destFile)
		fmt.Printf("Copied %s to %s\n", srcFile, destFile)
	}
	fmt.Println("\nSwitched to", folder)
}

func generateNewKeyPair() {
	homeDir, _ := os.UserHomeDir()
	sshDir := filepath.Join(homeDir, sshFolder)

	subFolder := readInput("\nEnter the subfolder name for the new key-pair (required): ")

	// Convert to snake_case
	subFolder = toSnakeCase(subFolder)

	// Ensure the subfolder name has the underscore prefix
	if !strings.HasPrefix(subFolder, "_") {
		subFolder = "_" + subFolder
	}

	subDir := filepath.Join(sshDir, subFolder)

	// Create subfolder if it doesn't exist
	if _, err := os.Stat(subDir); os.IsNotExist(err) {
		os.Mkdir(subDir, 0700)
		fmt.Printf("\nCreated folder: %s\n", subDir)
	}

	// Default algorithm to ed25519 if blank
	cipher := readInput("\nEnter the cipher (e.g., ed25519, rsa) [default: ed25519]: ")
	if cipher == "" {
		cipher = "ed25519"
	}

	// Default key name to id_<algorithm> if blank
	keyName := readInput("Enter the key name [default: id_" + cipher + "]: ")
	if keyName == "" {
		keyName = "id_" + cipher
	}

	// Get current username and hostname for default comment
	currentUser, _ := os.UserHomeDir()
	hostName, _ := os.Hostname()
	defaultComment := fmt.Sprintf("%s@%s", filepath.Base(currentUser), hostName)

	// Default comment to user@hostname if blank
	comment := readInput("Enter a comment for the key [default: " + defaultComment + "]: ")
	if comment == "" {
		comment = defaultComment
	}

	fmt.Println("\nGenerating new key pair...\n")

	keyPath := filepath.Join(subDir, keyName)
	cmd := exec.Command("ssh-keygen", "-t", cipher, "-C", comment, "-f", keyPath, "-N", "")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("\nFailed to generate key:", err)
		return
	}

	fmt.Printf("\nGenerated new key pair in %s\n", subDir)
}

// Helper function to convert a string to snake_case
func toSnakeCase(s string) string {
	// Replace spaces with underscores
	s = strings.ReplaceAll(s, " ", "_")

	// Use a regular expression to convert camel case and handle underscores
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	s = re.ReplaceAllString(s, "${1}_${2}")

	// Convert the entire string to lowercase
	return strings.ToLower(s)
}

func copyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Check if the source file is not a directory
	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	// Open the source file
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	// Create the destination file
	// destination, err := os.Create(dst)
	// Create the destination file with 0600 permissions (only owner can read/write)
	destination, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer destination.Close()

	// Copy the file content
	if _, err = io.Copy(destination, source); err != nil {
		return err
	}

	// Set permissions for the destination file (only readable/writable by owner)
	if err := os.Chmod(dst, 0600); err != nil {
		return err
	}

	return nil
}

func readInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
