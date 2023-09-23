package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	pb "addressbook/addressbookpb"

	"google.golang.org/protobuf/proto"
)

func promptForDelete(r io.Reader, book *pb.AddressBook) error {
	var (
		id             int32
		selectedIndex  int
		selectedPerson *pb.Person
	)

	reader := bufio.NewReader(r)

	fmt.Print("Enter person ID number: ")

	n, err := fmt.Fscanf(reader, "%d\n", &id)
	if err != nil {
		return err
	}
	if n != 1 {
		return fmt.Errorf("You can only delete one address at a time.")
	}

	for index, address := range book.People {
		if address.Id == id {
			selectedPerson = address
			selectedIndex = index
			break
		}
	}

	if selectedPerson == nil {
		return fmt.Errorf("No person is found.")
	}

	fmt.Print("Are you sure you want to delete it (y/n): ")
	confirm, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	confirm = strings.TrimSpace(confirm)
	if confirm == "y" {
		book.People = append(book.People[0:selectedIndex], book.People[selectedIndex+1:]...)
	}
	if confirm == "n" {
		return nil
	}

	return nil
}

func getFileName() (string, error) {
	// Only accept filename as the argument
	// program name + filename = 2 arguments
	if len(os.Args) != 2 {
		return "", errors.New("Invalid arguments, only one filename is accepted ")
	}

	return os.Args[1], nil
}

func readFile(fname string) ([]byte, error) {
	in, err := os.ReadFile(fname)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("%s: File not found. Creating new file after address is input.\n", fname)
			return []byte{}, nil
		} else {
			return nil, fmt.Errorf("Error reading file: %s", fname)
		}
	}
	return in, nil
}

func run() error {
	fname, err := getFileName()
	if err != nil {
		return err
	}

	in, err := readFile(fname)
	if err != nil {
		return err
	}

	book := &pb.AddressBook{}
	if err := proto.Unmarshal(in, book); err != nil {
		return fmt.Errorf("Failed to parse address book: %s", err)
	}

	err = promptForDelete(os.Stdin, book)
	if err != nil {
		return fmt.Errorf("Error with address: %s", err)
	}

	out, err := proto.Marshal(book)
	if err != nil {
		return fmt.Errorf("Failed to encode address book: %s", err)
	}

	if err := os.WriteFile(fname, out, 0644); err != nil {
		return fmt.Errorf("Failed to write address book: %s", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}
