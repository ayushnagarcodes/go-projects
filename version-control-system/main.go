package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func displayHelp() {
	fmt.Println("These are SVCS commands:")
	fmt.Printf("%-10s Get and set a username.\n", "config")
	fmt.Printf("%-10s Add a file to the index.\n", "add")
	fmt.Printf("%-10s Show commit logs.\n", "log")
	fmt.Printf("%-10s Save changes.\n", "commit")
	fmt.Printf("%-10s Restore a file.\n", "checkout")
}

func setUsername(username string) {
	if err := os.WriteFile("./vcs/config.txt", []byte(username), 0644); err != nil {
		log.Fatal(err)
	}
}

func getUsername() []byte {
	data, err := os.ReadFile("./vcs/config.txt")
	if err != nil {
		if os.IsNotExist(err) {
			return []byte("")
		}
		log.Fatal(err)
	}
	return data
}

func config() {
	if len(os.Args) == 3 {
		setUsername(os.Args[2])
	}

	username := getUsername()
	if len(username) != 0 {
		fmt.Printf("The username is %s.\n", string(username))
	} else {
		fmt.Println("Please, tell me who you are.")
	}
}

func addFile(fileName string) {
	existingData, err := os.ReadFile("./vcs/index.txt")
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatal(err)
		}
	}

	// exit: if fileName already exists in index.txt
	names := strings.Split(string(existingData), "\n")
	for _, name := range names {
		if name == fileName {
			fmt.Println("File already tracked!")
			return
		}
	}

	var newData string
	if string(existingData) == "" {
		newData = fileName
	} else {
		newData = string(existingData) + "\n" + fileName
	}

	if err := os.WriteFile("./vcs/index.txt", []byte(newData), 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The file '%s' is tracked.\n", fileName)
}

func getFiles() {
	file, err := os.Open("./vcs/index.txt")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Add a file to the index.")
			return
		}
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	fmt.Println("Tracked files:")
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func add() {
	if len(os.Args) == 3 {
		fileName := os.Args[2]

		// checking if the given file exists or not
		_, err := os.Stat(fileName)
		if os.IsNotExist(err) {
			fmt.Printf("Can't find '%s'.\n", fileName)
			return
		}

		addFile(fileName)
	} else {
		getFiles()
	}
}

func showLog() {
	file, err := os.Open("./vcs/log.txt")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No commits yet.")
			return
		}
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func writeLog(commitHash string, commitMessage string) {
	author, err := os.ReadFile("./vcs/config.txt")
	if err != nil {
		log.Fatal(err)
	}

	existingData, err := os.ReadFile("./vcs/log.txt")
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatal(err)
		}
	}

	// constructing info string that contains info about the commit
	info := fmt.Sprintf("commit %s\nAuthor: %s\n%s\n", commitHash, author, commitMessage)
	if string(existingData) != "" {
		info += "\n"
	}

	// rewriting contents of log.txt so that the latest commit info is at the top
	newData := append([]byte(info), existingData...)
	if err := os.WriteFile("./vcs/log.txt", newData, 0644); err != nil {
		log.Fatal(err)
	}
}

func commitExists(commitHash string) bool {
	entries, err := os.ReadDir("./vcs/commits")
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		if entry.IsDir() && entry.Name() == commitHash {
			return true
		}
	}
	return false
}

func computeHash(filePaths []string) string {
	sha256Hash := sha256.New()
	for _, path := range filePaths {
		file, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}

		if _, err := io.Copy(sha256Hash, file); err != nil {
			log.Fatal(err)
		}
		file.Close()
	}
	return hex.EncodeToString(sha256Hash.Sum(nil))
}

func createCommit(commitHash string, filePaths []string) {
	commitHashDir := filepath.Join("./vcs/commits", commitHash)
	err := os.Mkdir(commitHashDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	// copying staged/tracked files inside the commitHashDir
	for _, path := range filePaths {
		newPath := filepath.Join(commitHashDir, path)

		// Open the source file for reading
		srcFile, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer srcFile.Close()

		// Create or open the destination file for writing
		destFile, err := os.Create(newPath)
		if err != nil {
			log.Fatal(err)
		}
		defer destFile.Close()

		// Copy the contents from the source to the destination
		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Changes are committed.")
}

func commit() {
	if len(os.Args) == 3 {
		// exit: if config command has not run yet
		_, err := os.ReadFile("./vcs/config.txt")
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("Config file not set up!")
				return
			}
		}

		// getting the list of staged/tracked files
		data, err := os.ReadFile("./vcs/index.txt")
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("No staged files.")
				return
			} else {
				log.Fatal(err)
			}
		}

		if len(data) == 0 {
			fmt.Println("No staged files.")
			return
		}

		filePaths := strings.Split(string(data), "\n")
		commitHash := computeHash(filePaths)
		// checking whether a directory with that commitHash exists inside "vcs/commits" directory
		if !commitExists(commitHash) {
			// if commitHash directory doesn't exist, then creating and logging the commit
			createCommit(commitHash, filePaths)
			writeLog(commitHash, os.Args[2])
		} else {
			fmt.Println("Nothing to commit.")
		}
	} else {
		fmt.Println("Message was not passed.")
	}
}

func doCheckout(commitHash string) {
	commitHashDir := filepath.Join("./vcs/commits", commitHash)

	// getting the list of files in commitHashDir
	entries, err := os.ReadDir(commitHashDir)
	if err != nil {
		log.Fatal(err)
	}

	// copying files inside the commitHashDir to the respective files in the project directory
	for _, entry := range entries {
		if !entry.IsDir() {
			path := filepath.Join(commitHashDir, entry.Name())

			// Open the source file for reading
			srcFile, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}

			// Create or open the destination file for writing
			destFile, err := os.Create(entry.Name())
			if err != nil {
				log.Fatal(err)
			}

			// Copy the contents from the source to the destination
			_, err = io.Copy(destFile, srcFile)
			if err != nil {
				log.Fatal(err)
			}

			srcFile.Close()
			destFile.Close()
		}
	}

	fmt.Printf("Switched to commit %s.\n", commitHash)
}

func checkout() {
	if len(os.Args) == 3 {
		commitHash := os.Args[2]
		if commitExists(commitHash) {
			doCheckout(commitHash)
		} else {
			fmt.Println("Commit does not exist.")
		}
	} else {
		fmt.Println("Commit id was not passed.")
	}
}

func main() {
	err := os.MkdirAll("./vcs/commits", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) == 1 {
		displayHelp()
	}

	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "--help":
			displayHelp()
		case "config":
			config()
		case "add":
			add()
		case "log":
			showLog()
		case "commit":
			commit()
		case "checkout":
			checkout()
		default:
			fmt.Printf("'%s' is not a SVCS command.\n", os.Args[1])
		}
	}
}
